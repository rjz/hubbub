package travis_service

import (
	"encoding/json"
	"errors"
	"github.com/rjz/go-travis/travis"
	hubbub "github.com/rjz/hubbub/common"
)

const PRO = "travis.pro_token"
const ORG = "travis.org_token"

// TravisService handles policy goals related to travis-ci
type TravisService struct {
	Client        *travis.Client
	RepoID        int
	EnvVarService *EnvVarService
}

// repositorySettingsParams describe the state of repository settings in travis
type repositorySettingsParams travis.RepositorySettings

// parseRepositorySettingsParams parses a JSON goal into repositorySettingsParams
func parseRepositorySettingsParams(rawGoal *json.RawMessage) (*repositorySettingsParams, error) {
	params := repositorySettingsParams{}
	if err := json.Unmarshal([]byte(*rawGoal), &params); err != nil {
		return nil, err
	}
	return &params, nil
}

// envVarParams describes the state of an environment variable in travis
type envVarParams struct {
	State string `json:"state,omitempty"`
	*travis.EnvironmentVariable
}

// parseEnvVarParams parses a JSON goal into envVarParams
func parseEnvVarParams(rawGoal *json.RawMessage) (*envVarParams, error) {
	params := envVarParams{}
	if err := json.Unmarshal([]byte(*rawGoal), &params); err != nil {
		return nil, err
	}
	return &params, nil
}

// configureRepositoryId configures the repo's travis-ci ID for the service
func (ts *TravisService) configureRepositoryId(owner, name string) error {
	travisRepo, err := ts.Client.GetRepository(owner, name)
	if err != nil {
		return err
	}
	ts.RepoID = travisRepo.ID
	return nil
}

// repositorySettings updates settings to match the provided goal
func (ts *TravisService) repositorySettings(rawGoal *json.RawMessage) error {
	settings, err := parseRepositorySettingsParams(rawGoal)
	if err != nil {
		return err
	}

	travisSettings := travis.RepositorySettings(*settings)
	_, updateErr := ts.Client.UpdateRepositorySettings(ts.RepoID, &travisSettings)
	return updateErr
}

func (ts *TravisService) envVar(rawGoal *json.RawMessage) error {
	params, err := parseEnvVarParams(rawGoal)
	if err != nil {
		return err
	}

	// lazily configure the var service, allowing other travis-related tasks to
	// be completed without fetching environment variables
	if ts.EnvVarService == nil {
		evs, err := NewEnvVarService(ts.Client, ts.RepoID)
		if err != nil {
			return err
		}
		ts.EnvVarService = evs
	}

	switch params.State {
	case "present":
		return ts.EnvVarService.CreateOrUpdate(params.EnvironmentVariable)
	case "absent":
		return ts.EnvVarService.RemoveByName(*params.EnvironmentVariable.Name)
	default:
		return errors.New("unknown state.")
	}
}

// Do executes a single policy goal
func (ts *TravisService) Do(name string, rawGoal *json.RawMessage) error {
	switch name {
	case "travis_env_var":
		return ts.envVar(rawGoal)
	case "travis_repository_settings":
		return ts.repositorySettings(rawGoal)
	default:
		return errors.New("unknown goal (this shouldn't happen..)")
	}
	return nil
}

// configureClient attempts to access the travis API with the specified token.
var configureClient = func(ts *TravisService, client *travis.Client, owner, name string) bool {
	if ts.Client != nil {
		return true
	}

	ts.Client = client

	// Fetching the repo ID is a useful 'hello world'--most requests to the
	// travis API require a valid ID anyway!
	if err := ts.configureRepositoryId(owner, name); err == nil {
		return true
	}

	ts.Client = nil
	return false
}

// Construct an instance of TravisService configured for the given facts
//
// Attempts to configure for either travis.org, or--failing that--for travis
// pro / travis.com.
func TravisServiceFactory(facts *hubbub.Facts) (*hubbub.Service, error) {

	ts := TravisService{}
	owner := facts.GetString("repo.owner")
	name := facts.GetString("repo.name")

	// try travis.com (travis pro)
	if facts.IsAvailable(PRO) && configureClient(&ts, travis.NewProClient(hubbub.String(facts.GetString(PRO))), owner, name) {
		svc := hubbub.Service(&ts)
		return &svc, nil
	}

	// try travis.org
	if facts.IsAvailable(ORG) && configureClient(&ts, travis.NewClient(hubbub.String(facts.GetString(ORG))), owner, name) {
		svc := hubbub.Service(&ts)
		return &svc, nil
	}

	return nil, errors.New("Failed to configure Travis client")
}

func init() {
	hubbub.RegisterService([]string{
		"travis_repository_settings",
		"travis_env_var",
	}, TravisServiceFactory)
}
