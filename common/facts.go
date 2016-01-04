package common

import (
	"errors"
	"fmt"
)

// Facts are a write-once key:value store
type Facts map[string]interface{}

// NewFacts initializes an empty set of facts
func NewFacts(defaults map[string]interface{}) *Facts {
	f := Facts{}
	f.SetMap(defaults)
	return &f
}

func (f *Facts) GetString(k string) string {
	return f.Get(k).(string)
}

func (f *Facts) GetInt(k string) int {
	return f.Get(k).(int)
}

func (f *Facts) Get(k string) interface{} {
	return (*f)[k]
}

func (f *Facts) set(k string, v interface{}) error {
	if f.IsAvailable(k) {
		panic(fmt.Sprintf("fact is known and cannot be reset: '%s'", k))
	}

	(*f)[k] = v
	return nil
}

// SetRepository assigns defaults for the repository
func (f *Facts) SetRepository(r *Repository) {
	f.SetString("repo.host", *r.Host())
	f.SetString("repo.owner", *r.Owner())
	f.SetString("repo.name", *r.Name())
	f.SetString("repo.url", r.URL)
}

// Set assigns k to the interface-type v, returning an error on failed assignment
func (f *Facts) Set(k string, v interface{}) error {
	switch v.(type) {
	case string:
		return f.SetString(k, v.(string))
	case int:
		return f.SetInt(k, v.(int))
	case bool:
		return f.SetBool(k, v.(bool))
	default:
		return errors.New(fmt.Sprintf("Invalid type for '%s'", k))
	}
}

// SetMap attempts merging the facts in m, returning an error if any
// assignments fail
//
// TODO: rollback on failed assignment instead of leaving the Facts in an
// indeterminate state
func (f *Facts) SetMap(m map[string]interface{}) error {
	for k, v := range m {
		if err := f.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (f *Facts) SetString(k, v string) error {
	return f.set(k, v)
}

func (f *Facts) SetInt(k string, v int) error {
	return f.set(k, v)
}

func (f *Facts) SetBool(k string, v bool) error {
	return f.set(k, v)
}

func (f *Facts) IsAvailable(k string) bool {
	return (*f)[k] != nil
}
