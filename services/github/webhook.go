package github_service

import (
	"errors"
	"github.com/google/go-github/github"
)

type HookService struct {
	Client    *github.Client
	RepoOwner string
	RepoName  string
	Hooks     *[]github.Hook
}

func NewHookService(client *github.Client, owner, name string) (*HookService, error) {
	hs := HookService{client, owner, name, nil}

	hooks, _, err := hs.Client.Repositories.ListHooks(hs.RepoOwner, hs.RepoName, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	hs.Hooks = &hooks
	return &hs, nil
}

func (hs *HookService) byUrl(url string) (*github.Hook, error) {
	if hs.Hooks == nil {
		return nil, errors.New("Fetch hooks from github first")
	}

	for i, h := range *hs.Hooks {
		if h.Config["url"] == url {
			return &(*hs.Hooks)[i], nil
		}
	}
	return nil, nil
}

func (hs *HookService) CreateOrUpdate(params *github.Hook) error {
	hookUrl := params.Config["url"].(string)
	hook, err := hs.byUrl(hookUrl)
	if err != nil {
		return err
	}

	if hook == nil {
		_, _, err := hs.Client.Repositories.CreateHook(hs.RepoOwner, hs.RepoName, params)
		return err
	}

	_, _, updateErr := hs.Client.Repositories.EditHook(hs.RepoOwner, hs.RepoName, *hook.ID, params)
	return updateErr
}

func (hs *HookService) Remove(params *github.Hook) error {
	hook, err := hs.byUrl(params.Config["url"].(string))
	if err != nil {
		return err
	}
	if hook == nil {
		return nil
	}

	_, deleteErr := hs.Client.Repositories.DeleteHook(hs.RepoOwner, hs.RepoName, *hook.ID)
	return deleteErr
}
