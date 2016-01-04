package common

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type Repository struct {
	URL       string `json:"url,omitempty"`
	urlPieces []string
}

func (r *Repository) urlFragment(n int) *string {
	if r.urlPieces == nil {
		r.urlPieces = strings.Split(r.URL, "/")
	}
	return &r.urlPieces[n]
}

func (r *Repository) Host() *string {
	return r.urlFragment(0)
}

func (r *Repository) Owner() *string {
	return r.urlFragment(1)
}

func (r *Repository) Name() *string {
	return r.urlFragment(2)
}

func LoadRepositories(filename string) (*[]Repository, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	rs := []Repository{}
	if err := json.Unmarshal(data, &rs); err != nil {
		return nil, err
	}

	return &rs, nil
}
