package github_service

import (
	"encoding/json"
	"errors"
	"github.com/google/go-github/github"
	hubbub "github.com/rjz/hubbub/common"
	"golang.org/x/oauth2"
	"io/ioutil"
)

type sha string

// Serves policy goals related to github
type GithubService struct {
	Client      *github.Client
	HookService *HookService
	FileService *FileService
	RepoOwner   string
	RepoName    string
}

// fileParams describe a "github_file" goal
type fileParams struct {
	State    *string `json:"state,omitempty"`
	Content  *string `json:"content,omitempty"`
	Filename *string `json:"filename,omitempty"`
	Name     *string `json:"name,omitempty"`
	Ref      *string `json:"ref,omitempty"`
}

func (params *fileParams) loadContent() error {
	if params.Content != nil {
		return errors.New("Ambiguous argument: cannot specify both content and filename")
	}

	data, err := ioutil.ReadFile(*params.Filename)
	if err != nil {
		return err
	}

	strData := string(data)
	params.Content = &strData
	return nil
}

func parseFileParams(attrs *json.RawMessage) (*fileParams, error) {
	params := fileParams{}
	if err := json.Unmarshal([]byte(*attrs), &params); err != nil {
		return nil, err
	}

	if params.Filename != nil {
		if err := params.loadContent(); err != nil {
			return nil, err
		}
	}

	return &params, nil
}

// fileParams describe a "github_webhook" goal
type hookParams struct {
	State string `json:"state,omitempty"`
	*github.Hook
}

func parseHookParams(attrs *json.RawMessage) (*hookParams, error) {
	params := hookParams{}
	if err := json.Unmarshal([]byte(*attrs), &params); err != nil {
		return nil, err
	}

	if params.Name == nil {
		// default to webhook, users can override with service hooks if needed
		// https://developer.github.com/webhooks/#service-hooks
		params.Name = hubbub.String("web")
	}

	return &params, nil
}

func (s *GithubService) doWebhook(msg *json.RawMessage) error {
	if s.HookService == nil {
		hs, err := NewHookService(s.Client, s.RepoOwner, s.RepoName)
		if err != nil {
			return err
		}
		s.HookService = hs
	}

	params, err := parseHookParams(msg)
	if err != nil {
		return err
	}
	switch params.State {
	case "present":
		return s.HookService.CreateOrUpdate(params.Hook)
	case "absent":
		return s.HookService.Remove(params.Hook)
	default:
		return errors.New("unknown state.")
	}
}

func (s *GithubService) doFile(msg *json.RawMessage) error {
	params, err := parseFileParams(msg)
	if err != nil {
		return err
	}

	// find current SHA for ref
	SHA, err := s.refSHA(*params.Ref)
	if err != nil {
		return err
	}

	if s.FileService == nil {
		s.FileService = NewFileService(s.Client, s.RepoOwner, s.RepoName)
	}

	if err := s.FileService.TreeFacts(*SHA); err != nil {
		return err
	}

	switch *params.State {
	case "present":
		return s.FileService.CreateOrUpdate(*SHA, *params)
	case "absent":
		return s.FileService.Remove(*SHA, *params)
	default:
		return errors.New("unknown state.")
	}
}

// RefFacts fetches the current state (SHA, tree) of the reference
func (s *GithubService) refSHA(refName string) (*sha, error) {
	ref, _, refErr := s.Client.Git.GetRef(s.RepoOwner, s.RepoName, refName)
	if refErr != nil {
		return nil, refErr
	}
	refSHA := sha(*ref.Object.SHA)
	return &refSHA, nil
}

func (s *GithubService) Do(goal string, msg *json.RawMessage) error {
	switch goal {
	case "github_webhook":
		return s.doWebhook(msg)
	case "github_file":
		return s.doFile(msg)
	}
	return nil
}

func GithubServiceFactory(facts *hubbub.Facts) (*hubbub.Service, error) {
	if !facts.IsAvailable("github.access_token") {
		return nil, errors.New("no github access token available")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: facts.GetString("github.access_token")})
	oc := oauth2.NewClient(oauth2.NoContext, ts)
	gs := GithubService{github.NewClient(oc), nil, nil, facts.GetString("repo.owner"), facts.GetString("repo.name")}

	svc := hubbub.Service(&gs)
	return &svc, nil
}

func init() {
	hubbub.RegisterService([]string{
		"github_webhook",
		"github_file",
	}, GithubServiceFactory)
}
