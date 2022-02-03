package main

import (
	"errors"
	"github.com/konveyor/tackle-hub/api"
)

//
// Factory.
func newRepository(id uint) (rp Repository, err error) {
	repository, err := addon.Repository.Get(id)
	if err != nil {
		return
	}
	switch repository.Kind {
	case "git":
		rp = &Git{Repository: repository}
	case "svn":
	case "mvn":
	default:
		err = errors.New("unknown kind")
	}

	return
}

//
// Repository interface.
type Repository interface {
	Fetch(path string) (err error)
	Path() string
}

//
// Git repository.
type Git struct {
	*api.Repository
	path string
}

//
// Fetch the repository.
func (r *Git) Fetch(path string) (err error) {
	r.path = path
	cmd := Command{Path: "/usr/bin/git"}
	cmd.Options.add("clone", r.URL)
	err = cmd.Run()
	if err != nil {
		return
	}

	return
}

//
// Path to the fetched repository.
func (r *Git) Path() string {
	return r.path
}
