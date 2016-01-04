package common

import (
	"encoding/json"
	"io/ioutil"
)

type PolicyGoal struct {
	Goal       *string
	RawMessage json.RawMessage
}

type Policy []PolicyGoal

// Goals lists all goals included in the policy
func (p Policy) Goals() []string {
	var goals []string
	for _, goal := range p {
		goals = append(goals, *goal.Goal)
	}
	return goals
}

// LoadPolicy reads a raw (JSON) policy from filename
func LoadPolicy(filename string) (policy Policy, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var parsed []map[string]json.RawMessage
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	for _, s := range parsed {
		pg := PolicyGoal{}
		for k, v := range s {
			pg.RawMessage = v
			pg.Goal = &k
			break
		}

		policy = append(policy, pg)
	}
	return
}
