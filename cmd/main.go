package main

import (
	"errors"
	hub "github.com/konveyor/tackle-hub/addon"
	"github.com/konveyor/tackle-hub/api"
	"os"
)

var (
	// addon adapter.
	addon = hub.Addon
	// Logger
	log = hub.Log
)

const (
	DefaultTarget = "cloud-readiness"
)

//
// Artifact uploaded.
type Artifact struct {
	Bucket uint   `json:"bucket" binding:"required"`
	Path   string `json:"path" binding:"required"`
}

//
// Data Addon data passed in the secret.
type Data struct {
	Application  uint      `json:"application" binding:"required"`
	Binary       bool      `json:"binary"`
	Dependencies bool      `json:"dependencies"`
	Targets      []string  `json:"targets"`
	Packages     []string  `json:"packages"`
	Artifact     *Artifact `json:"artifact"`
}

// validate settings.
// Default settings not specified.
func (d *Data) validate() (err error) {
	if d.Application == 0 {
		err = errors.New("Application not specified.")
		return
	}
	if len(d.Targets) == 0 {
		d.Targets = []string{DefaultTarget}
	}

	return
}

//
// ensureBucket to store windup report.
func ensureBucket(d *Data) (bucket *api.Bucket, err error) {
	bucket, err = addon.Bucket.Ensure(d.Application, "Windup")
	if err != nil {
		return
	}
	err = addon.Bucket.Purge(bucket)
	return
}

//
// main
func main() {
	var err error
	// Error handling.
	defer func() {
		if err != nil {
			log.Error(err, "Addon failed.")
			_ = addon.Failed(err.Error())
			os.Exit(1)
		}
	}()
	// Get the addon data associated with the task.
	d := &Data{}
	err = addon.DataWith(d)
	if err != nil {
		return
	}
	// Report addon has started
	err = addon.Started()
	if err != nil {
		return
	}
	// Validate the addon data.
	err = d.validate()
	if err != nil {
		return
	}
	application, err := addon.Application.Get(d.Application)
	if err != nil {
		return
	}
	// Run windup.
	windup := Windup{}
	// Fetch repository.
	if !d.Binary {
		err = addon.Total(2)
		if err != nil {
			return
		}
		if application.Repository == nil {
			err = errors.New("Application repository not defined.")
			return
		}
		repository, err := newRepository(application.Repository)
		if err != nil {
			return
		}
		err = repository.Fetch("/tmp/git")
		if err == nil {
			err = addon.Increment()
			if err != nil {
				return
			}
			windup.repository = repository
		} else {
			return
		}
	}
	// Create the bucket.
	bucket, err := ensureBucket(d)
	if err == nil {
		windup.bucket = bucket
	} else {
		return
	}
	// Run windup.
	err = windup.Run()
	if err == nil {
		err = addon.Increment()
		if err != nil {
			return
		}
	} else {
		return
	}

	// Task update: The addon has succeeded
	_ = addon.Succeeded()
}
