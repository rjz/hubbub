package travis_service

import (
	"github.com/rjz/go-travis/travis"
	util "github.com/rjz/hubbub/common"
	"reflect"
	"testing"
)

func evsFixture() EnvVarService {
	return EnvVarService{
		vars: &[]travis.EnvironmentVariable{
			travis.EnvironmentVariable{Name: util.String("xyz"), ID: util.String("123")},
			travis.EnvironmentVariable{Name: util.String("xyz"), ID: util.String("456")},
			travis.EnvironmentVariable{Name: util.String("abc"), ID: util.String("789")},
		},
	}
}

func TestEnvVarByNameMultiple(t *testing.T) {
	evs := evsFixture()
	if count := len(evs.byName("xyz")); count != 2 {
		t.Error("expected 2, got", count)
	}
}

func TestEnvVarByNameSingle(t *testing.T) {
	evs := evsFixture()
	if count := len(evs.byName("abc")); count != 1 {
		t.Error("expected 1, got", count)
	}
}

func TestEnvVarByNameMissing(t *testing.T) {
	evs := evsFixture()
	if count := len(evs.byName("def")); count != 0 {
		t.Error("expected 0, got", count)
	}
}

func TestRemoveById(t *testing.T) {
	evs := evsFixture()
	evs.removeInternalById("456")

	if !reflect.DeepEqual(*evs.vars, []travis.EnvironmentVariable{
		travis.EnvironmentVariable{Name: util.String("xyz"), ID: util.String("123")},
		travis.EnvironmentVariable{Name: util.String("abc"), ID: util.String("789")},
	}) {
		t.Fatal("failed removing by Id")
	}
}

func TestRemoveByIdFirst(t *testing.T) {
	evs := evsFixture()
	evs.removeInternalById("123")

	if !reflect.DeepEqual(*evs.vars, []travis.EnvironmentVariable{
		travis.EnvironmentVariable{Name: util.String("xyz"), ID: util.String("456")},
		travis.EnvironmentVariable{Name: util.String("abc"), ID: util.String("789")},
	}) {
		t.Fatal("failed removing by Id")
	}
}

func TestRemoveByInvalidId(t *testing.T) {
	evs := evsFixture()
	if err := evs.removeInternalById("xyz"); err == nil {
		t.Error("Expected err, didn't get it.")
	}
}
