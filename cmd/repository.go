package main

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/konveyor/tackle-hub/api"
)

func newRepository(d *Data) (rp Repository, err error) {
	repository, err := addon.Repository.Get(d.Repository)
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
	Init(path string) (err error)
	Path() string
	Binary() bool
}

type Git struct {
	*api.Repository
	path string
}

func (r *Git) Init(path string) (err error) {
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

func (r *Git) Binary() bool {
	return false
}
