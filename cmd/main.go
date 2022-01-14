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
	Application uint `json:"application"`
	Repository  uint `json:"repository"`
	Identity    uint `json:"identity"`
	Windup      struct {
		Targets  []string `json:"targets"`
		Packages []string `json:"packages"`
	} `json:"windup"`
}

// validate settings.
// Default settings not specified.
func (d *Data) validate() (err error) {
	if d.Application == 0 {
		err = errors.New("Application not specified.")
		return
	}
	if d.Repository == 0 {
		err = errors.New("Repository not specified.")
		return
	}
	if len(d.Windup.Targets) == 0 {
		d.Windup.Targets = append(
			d.Windup.Targets,
			DefaultTarget)
		return
	}

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
	_ = addon.Total(2)
	// Validate the addon data.
	err = d.validate()
	if err != nil {
		return
	}
	//
	repository, err := newRepository(d)
	if err != nil {
		return
	}
	err = repository.Init("/tmp")
	if err == nil {
		_ = addon.Increment()
	} else {
		return
	}
	// Create the bucket.
	bucket, err := createBucket(d)
	if err != nil {
		return
	}
	// Run windup.
	windup := Windup{
		Data:       d,
		repository: repository,
		bucket:     bucket,
	}
	err = windup.Run()
	if err == nil {
		_ = addon.Increment()
	} else {
		return
	}

	// Task update: The addon has succeeded
	_ = addon.Succeeded()
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
