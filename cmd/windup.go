package main

import (
	"bytes"
	"github.com/konveyor/tackle-hub/api"
	"os/exec"
)

//
// Windup application analyzer.
type Windup struct {
	*Data
	repository Repository
	bucket     *api.Bucket
}

//
// run windup.
func (r *Windup) Run() (err error) {
	_ = addon.Activity("Running windup.")
	options := r.options()
	cmd := exec.Command("/opt/mta-cli/bin/mta-cli", options...)
	cmd.Dir = "/opt/mta-cli"
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
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
	options.add("--target", r.Windup.Targets...)
	options.add("--input", r.repository.Path())
	if !r.repository.Binary() {
		options.add("--sourceMode")
	}
	if len(r.Windup.Packages) > 0 {
		options.add("--packages", r.Windup.Packages...)
	}
	return
}

//
// Options are CLI options.
type Options []string

//
// add
func (a *Options) add(option string, s ...string) {
	*a = append(*a, option)
	*a = append(*a, s...)
}
