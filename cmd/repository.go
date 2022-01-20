package main

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
	gitCloneOptions := &git.CloneOptions{
		URL:               r.URL,
		ReferenceName:     plumbing.ReferenceName(r.Branch),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}
	_, err = git.PlainClone(path, false, gitCloneOptions)
	if err != nil {
		return
	}

	return
}

func (r *Git) Path() string {
	return r.path
}
