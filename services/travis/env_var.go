package travis_service

import (
	"errors"
	"fmt"
	"github.com/rjz/go-travis/travis"
)

// Wraps environment variable list for a repo
type EnvVarService struct {
	client *travis.Client
	repoID int
	vars   *[]travis.EnvironmentVariable
}

// NewEnvVarService configures a new EnvVarService using the specified client
// and travis-ci repoId
func NewEnvVarService(client *travis.Client, repoId int) (*EnvVarService, error) {
	vars, err := client.ListEnvironmentVariables(repoId)
	if err != nil {
		return nil, err
	}
	return &EnvVarService{client, repoId, &vars}, nil
}

// byName returns environment variables matching (case-sensitive) name
func (evs *EnvVarService) byName(name string) []*travis.EnvironmentVariable {
	var matches []*travis.EnvironmentVariable
	all := *evs.vars
	for i, v := range all {
		if *v.Name == name {
			matches = append(matches, &all[i])
		}
	}
	return matches
}

// Create a new environment variable
func (evs *EnvVarService) create(ev *travis.EnvironmentVariable) error {
	_, err := evs.client.CreateEnvironmentVariable(evs.repoID, ev)
	if err == nil {
		// Add new var to internal list
		newVars := append(*evs.vars, *ev)
		evs.vars = &newVars
	}
	return err
}

// CreateOrUpdate sets an environment variable
//
// If the environment variable has multiple definitions, the update will
// overwrite the most recent entry and all other entries will be removed
func (evs *EnvVarService) CreateOrUpdate(ev *travis.EnvironmentVariable) error {
	existingVars := evs.byName(*ev.Name)
	if len(existingVars) == 0 {
		return evs.create(ev)
	}

	// Travis doesn't enforce uniqueness on variable names: pop latest and
	// prepare list of dups for deletion.
	existingVar, dups := existingVars[len(existingVars)-1], existingVars[:len(existingVars)-1]

	// Update internal var
	existingVar.Value = ev.Value
	existingVar.Public = ev.Public

	// Update remote var
	if _, err := evs.client.UpdateEnvironmentVariable(evs.repoID, *existingVar.ID, existingVar); err != nil {
		return err
	}

	// TODO: We know IDs in advance and can run requests in parallel. We should.
	return evs.removeAll(dups)
}

// removeInternalById omits a var from the internal list
func (evs *EnvVarService) removeInternalById(id string) error {
	oldVars := *evs.vars
	for i, v := range *evs.vars {
		if *v.ID == id {
			newVars := append(oldVars[:i], oldVars[i+1:]...)
			evs.vars = &newVars
			return nil
		}
	}
	return errors.New(fmt.Sprintf("var '%d' does not exist", id))
}

// removeAll removes all vars from the list
func (evs *EnvVarService) removeAll(vars []*travis.EnvironmentVariable) error {
	for _, v := range vars {
		// Remove internal var
		if err := evs.removeInternalById(*v.ID); err != nil {
			return err
		}

		// Destroy remote var
		if err := evs.client.DestroyEnvironmentVariable(evs.repoID, *v.ID); err != nil {
			return err
		}
	}
	return nil
}

// RemoveByName deletes one or more environment variables
func (evs *EnvVarService) RemoveByName(name string) error {
	return evs.removeAll(evs.byName(name))
}
