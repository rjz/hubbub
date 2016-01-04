package common

import (
	"errors"
	"testing"
)

func TestPolicyCreateServicesUnavailableGoal(t *testing.T) {
	r := NewServiceFactoryRegistry()
	r.Register([]string{"foo_do", "foo_echo"}, FooServiceFactory)
	if _, err := r.CreateServices([]string{"ixnay"}, &Facts{}); err == nil {
		t.Error("expected error for invalid service; didn't get it.")
	}
}

func TestPolicyCreateServicesInstanceFails(t *testing.T) {
	r := NewServiceFactoryRegistry()
	r.Register([]string{"broken_task"}, func(pc *Facts) (*Service, error) {
		return nil, errors.New("service failed to start")
	})

	if _, err := r.CreateServices([]string{"broken_task"}, &Facts{}); err == nil {
		t.Error("expected error when service failed to start; didn't get it.")
	}
}

func TestPolicyCreateServicesSetAllServiceGoals(t *testing.T) {
	r := NewServiceFactoryRegistry()
	r.Register([]string{"foo_do", "foo_echo"}, FooServiceFactory)

	services, _ := r.CreateServices([]string{"foo_do"}, &Facts{})
	if (*services)["foo_do"] == nil || (*services)["foo_echo"] == nil {
		t.Error("expected all foo_* goals to be registered; they weren't.")
	}
}

func TestPolicyCreateServicesShareInstance(t *testing.T) {
	r := NewServiceFactoryRegistry()
	r.Register([]string{"foo_do", "foo_echo"}, FooServiceFactory)

	services, _ := r.CreateServices([]string{"foo_do", "foo_echo"}, &Facts{})
	if (*services)["foo_do"] != (*services)["foo_echo"] {
		t.Error("expected all foo_* goals to share single service instance; they didn't.")
	}
}
