package main

import (
	"errors"
	"github.com/konveyor/tackle-hub/api"
)

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
		err = errors.New("")
	}

	return
}

type Repository interface {
	Fetch(path string) (err error)
	Path() string
}

type Git struct {
	*api.Repository
	path string
}

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

func (r *Git) Path() string {
	return r.path
}
