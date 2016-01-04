package common

import (
	"fmt"
	"log"
	"os"
)

type Session struct {
	*Policy
	*Facts
	*log.Logger
	*ServiceFactoryRegistry
}

// NewSession creates a new session configured with the policy, facts, and globally-registered services
func NewSession(rp *Policy, f *Facts) *Session {
	logger := log.New(os.Stdout, fmt.Sprintf("%s - ", f.GetString("repo.url")), log.LstdFlags)
	factories := ServiceFactories()
	s := Session{rp, f, logger, &factories}
	return &s
}

// Run the session
func (s *Session) Run() error {

	s.Logger.Println("BEGIN")

	services, err := s.ServiceFactoryRegistry.CreateServices(s.Policy.Goals(), s.Facts)
	if err != nil {
		return err
	}

	for _, pg := range *s.Policy {

		goalName := *pg.Goal
		s.Logger.Println(" --", goalName)

		svc := (*services)[goalName]
		if err := (*svc).Do(goalName, &pg.RawMessage); err != nil {
			s.Logger.Println("FAILED", err)
			return err
		}
	}

	s.Logger.Println("END")
	return nil
}
