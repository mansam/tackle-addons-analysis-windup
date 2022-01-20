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
// Data Addon data passed in the secret.
type Data struct {
	Application  uint     `json:"application"`
	Binary       bool     `json:"binary"`
	Dependencies bool     `json:"dependencies"`
	Targets      []string `json:"targets"`
	Packages     []string `json:"packages"`
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
// CreateBucket to store windup report.
func createBucket(d *Data) (bucket *api.Bucket, err error) {
	bucket = &api.Bucket{}
	bucket.CreateUser = "addon"
	bucket.Name = "AnalysisWindup"
	bucket.ApplicationID = d.Application
	err = addon.Bucket.Create(bucket)
	return
}

//
// main
func main() {
	var err error
	// Get the addon data associated with the task.
	d := &Data{}
	_ = addon.DataWith(d)
	// Error handling.
	defer func() {
		if err != nil {
			log.Error(err, "Addon failed.")
			_ = addon.Failed(err.Error())
			os.Exit(1)
		}
	}()
	// Report addon has started
	_ = addon.Started()
	// Validate the addon data.
	err = d.validate()
	if err != nil {
		return
	}
	// Run windup.
	windup := Windup{}
	// Fetch repository.
	if !d.Binary {
		_ = addon.Total(2)
		appRepository, err := addon.Repository.ByApplication(d.Application)
		if err != nil {
			return
		}
		repository, err := newRepository(appRepository.ID)
		if err != nil {
			return
		}
		err = repository.Fetch("/tmp")
		if err == nil {
			_ = addon.Increment()
			windup.repository = repository
		} else {
			return
		}
	}
	// Create the bucket.
	bucket, err := createBucket(d)
	if err == nil {
		windup.bucket = bucket
	} else {
		return
	}
	// Run windup.
	err = windup.Run()
	if err == nil {
		_ = addon.Increment()
	} else {
		return
	}

	// Task update: The addon has succeeded
	_ = addon.Succeeded()
}
