package common

import (
	"encoding/json"
	"reflect"
	"testing"
)

func setup() {
	serviceFactories.Register([]string{"foo_do", "foo_echo", "foo_invalid"}, FooServiceFactory)
}

func teardown() {
	serviceFactories = ServiceFactoryRegistry{}
	serviceFactories.goals = make(map[string]*int)
}

func expectJsonObject(t *testing.T, msg json.RawMessage, expected map[string]interface{}) {
	dest := map[string]interface{}{}
	if err := json.Unmarshal([]byte(msg), &dest); err != nil {
		t.Error("Invalid Json", string(msg))
		return
	}

	if !reflect.DeepEqual(expected, dest) {
		t.Error("expected", expected, "got", dest)
	}
}

func expectGoal(t *testing.T, pg PolicyGoal, goal string, body map[string]interface{}) {
	if *pg.Goal != goal {
		t.Error("expected goal", goal, "got", *pg.Goal)
		return
	}

	expectJsonObject(t, pg.RawMessage, body)
}

func TestLoadPolicy(t *testing.T) {
	setup()
	defer teardown()

	fixture, err := LoadPolicy("__fixtures/test.json")
	if err != nil {
		t.Error(err)
	}

	count := len(fixture)
	if count != 2 {
		t.Error("expected 2, saw", count)
	}

	expectGoal(t, fixture[0], "foo_do", map[string]interface{}{
		"state": "ambivalent",
		"bar":   "baz",
	})

	expectGoal(t, fixture[1], "foo_echo", map[string]interface{}{
		"bar": "baz",
	})
}

func TestPolicyGoals(t *testing.T) {
	p := Policy{
		PolicyGoal{Goal: String("goal_one")},
		PolicyGoal{Goal: String("goal_two")},
	}

	expected := []string{"goal_one", "goal_two"}
	if !reflect.DeepEqual(p.Goals(), expected) {
		t.Error("expected", expected, "got", p.Goals())
	}
}
