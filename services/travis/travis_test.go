package travis_service

import (
	"github.com/rjz/go-travis/travis"
	hubbub "github.com/rjz/hubbub/common"
	"testing"
)

var unstubbedConfigureClient func(ts *TravisService, client *travis.Client, owner, name string) bool

func stubConfigureClient() {
	unstubbedConfigureClient = configureClient
	configureClient = func(ts *TravisService, client *travis.Client, owner, name string) bool {
		return true
	}
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
	stubConfigureClient()
	defer restoreConfigureClient()

	pc := &hubbub.Facts{}
	pc.SetString("repo.owner", "rjz")
	pc.SetString("repo.name", "dingus")
	pc.SetString("travis.org_token", "xyz")

	if _, err := TravisServiceFactory(pc); err != nil {
		t.Fatal("expected pass, didn't get it.")
	}
}

func TestTravisServiceFactoryProToken(t *testing.T) {
	stubConfigureClient()
	defer restoreConfigureClient()

	pc := &hubbub.Facts{}
	pc.SetString("repo.owner", "rjz")
	pc.SetString("repo.name", "dingus")
	pc.SetString("travis.pro_token", "zyx")

	if _, err := TravisServiceFactory(pc); err != nil {
		t.Fatal("expected pass, didn't get it.")
	}
}
