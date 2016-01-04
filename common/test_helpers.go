package common

import (
	"encoding/json"
	"errors"
)

type FooService struct{}

func (m *FooService) Do(name string, msg *json.RawMessage) error {
	if name != "foo_do" && name != "foo_echo" {
		return errors.New("unknown foo action")
	}

	dst := make(map[string]interface{})
	json.Unmarshal(*msg, &dst)
	if dst["bar"].(string) != "baz" {
		return errors.New("unexpected value")
	}
	return nil
}

func FooServiceFactory(pc *Facts) (*Service, error) {
	svc := Service(&FooService{})
	return &svc, nil
}
