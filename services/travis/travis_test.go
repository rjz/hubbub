package travis_service

import (
	hubbub "github.com/rjz/hubbub/common"
	"testing"
)

var unstubbedConfigureClient func(ts *TravisService, token string) bool

func storeConfigureClient() {
	unstubbedConfigureClient = configureClient
}

func restoreConfigureClient() {
	configureClient = unstubbedConfigureClient
}

func TestTravisServiceFactoryNoToken(t *testing.T) {
	pc := &hubbub.Facts{}
	pc.SetString("repo.owner", "rjz")
	pc.SetString("repo.name", "dingus")

	if _, err := TravisServiceFactory(pc); err == nil {
		t.Fatal("expected a failing, didn't get it.")
	}
}

func TestTravisServiceFactoryOrgToken(t *testing.T) {
	storeConfigureClient()
	defer restoreConfigureClient()

	pc := &hubbub.Facts{}
	pc.SetString("repo.owner", "rjz")
	pc.SetString("repo.name", "dingus")
	pc.SetString("travis.org_token", "xyz")

	configureClient = func(ts *TravisService, token string) bool {
		if token != "xyz" {
			t.Error("expected org configuration, didn't get it.", token)
		}
		return true
	}

	if _, err := TravisServiceFactory(pc); err != nil {
		t.Fatal("expected pass, didn't get it.")
	}
}

func TestTravisServiceFactoryProToken(t *testing.T) {
	storeConfigureClient()
	defer restoreConfigureClient()

	pc := &hubbub.Facts{}
	pc.SetString("repo.owner", "rjz")
	pc.SetString("repo.name", "dingus")
	pc.SetString("travis.org_token", "xyz")
	pc.SetString("travis.pro_token", "zyx")

	configureClient = func(ts *TravisService, token string) bool {
		if token == "xyz" {
			return false
		}
		return true
	}

	if _, err := TravisServiceFactory(pc); err != nil {
		t.Fatal("expected pass, didn't get it.")
	}
}
