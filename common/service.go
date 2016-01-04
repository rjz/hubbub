package common

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Service represents a hubbub service implementation
type Service interface {

	// Do applies a single policy goal using a string (the goal name) and a raw,
	// JSON configuration, returning an error if the goal cannot be achieved.
	Do(string, *json.RawMessage) error
}

// ServiceRegistry organizes Service implementations by goal name
type ServiceRegistry map[string]*Service

// ServiceFactory returns a Service implementation configured using the
// specified Facts
type ServiceFactory func(*Facts) (*Service, error)

// ServiceFactoryRegistry enumerates goals by the services that implement them
type ServiceFactoryRegistry struct {
	goals     map[string]*int
	factories []ServiceFactory
}

// NewServiceFactoryRegistry initializes an empty registry
func NewServiceFactoryRegistry() *ServiceFactoryRegistry {
	return &ServiceFactoryRegistry{
		goals: make(map[string]*int),
	}
}

// Register associates `goals` with the provided `factory`
func (r *ServiceFactoryRegistry) Register(goals []string, factory ServiceFactory) {
	// add factory
	factoryIndex := len(r.factories)
	r.factories = append(r.factories, factory)

	// point goals to new factory
	for _, goal := range goals {
		if r.goals[goal] != nil {
			panic(fmt.Sprintf("goal '%s' was previously defined", goal))
		}
		r.goals[goal] = &factoryIndex
	}
}

func (r *ServiceFactoryRegistry) factoryId(name string) int {
	return *r.goals[name]
}

func (r *ServiceFactoryRegistry) createService(factoryId int, facts *Facts) (*Service, error) {
	return r.factories[factoryId](facts)
}

// Goals returns the list of all goals registered by available services
func (r *ServiceFactoryRegistry) Goals() []string {
	var goals []string
	for name, _ := range r.goals {
		goals = append(goals, name)
	}
	return goals
}

func (r *ServiceFactoryRegistry) goalsByFactoryId(factoryId int) []string {
	var relatives []string
	for _, goal := range r.Goals() {
		if r.factoryId(goal) == factoryId {
			relatives = append(relatives, goal)
		}
	}
	return relatives
}

// CreateServices constructs a `ServiceRegistry` containing service
// instances for the requested goals.
func (r *ServiceFactoryRegistry) CreateServices(goals []string, facts *Facts) (*ServiceRegistry, error) {
	services := ServiceRegistry{}
	for _, goal := range goals {
		if r.goals[goal] == nil {
			return nil, errors.New(fmt.Sprintf("no service available for '%s'", goal))
		}

		if services[goal] == nil {
			factoryIndex := r.factoryId(goal)
			svc, err := r.createService(factoryIndex, facts)
			if err != nil {
				return nil, err
			}

			for _, alias := range r.goalsByFactoryId(factoryIndex) {
				services[alias] = svc
			}
		}
	}

	return &services, nil
}

var serviceFactories ServiceFactoryRegistry

// RegisterService adds a factory to the global registry for the specified goals
func RegisterService(goals []string, factory ServiceFactory) {
	serviceFactories.Register(goals, factory)
}

// ServiceFactories returns the global service registry
func ServiceFactories() ServiceFactoryRegistry {
	return serviceFactories
}

func init() {
	serviceFactories = *NewServiceFactoryRegistry()
}
