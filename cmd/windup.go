package main

import (
	"github.com/konveyor/tackle-hub/api"
)

//
// Windup application analyzer.
type Windup struct {
	repository Repository
	bucket     *api.Bucket
	packages   []string
	targets    []string
}

//
// Run windup.
func (r *Windup) Run() (err error) {
	_ = addon.Activity("Running windup.")
	cmd := Command{Path: "/opt/windup"}
	cmd.Options = r.options()
	cmd.Dir = "/opt/mta-cli"
	err = cmd.Run()
	if err != nil {
		return
	}

	return
}

//
// options builds CLL options.
func (r *Windup) options() (opt Options) {
	options := Options{
		"--batchMode",
		"--output",
		r.bucket.Path,
	}
	options.add("--output", r.bucket.Path)
	options.add("--target", r.targets...)
	options.add("--input", r.repository.Path())
	if r.repository == nil {
		options.add("--sourceMode")
	}
	if len(r.packages) > 0 {
		options.add("--packages", r.packages...)
	}
	return
}
