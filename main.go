package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	hub "github.com/konveyor/tackle-hub/addon"
	"github.com/konveyor/tackle-hub/api"
	"os"
	"os/exec"
)

const (
	Kind = "Folder:application/html"
)

var (
	// addon adapter.
	addon = hub.Addon
)

//
// Data Addon data passed in the secret.
// TODO: Replace Git* fields with fetching the info from the hub
type Data struct {
	Application uint `json:"application"`
	Git         struct {
		URL    string `json:"url"`
		Branch string `json:"branch"`
		Path   string `json:"path"`
	} `json:"git"`
	Maven struct {
		File string `json:"file"`
	} `json:"maven"`
	Windup struct {
		Targets  []string `json:"targets"`
		Packages []string `json:"packages"`
	} `json:"windup"`
}

//
// main
func main() {
	var err error
	fmt.Printf("--- Tackle Addon - Discovery - Languages ---\n")

	// Get the addon data associated with the task.
	d := &Data{}
	_ = addon.DataWith(d)

	// Error handler
	defer func() {
		if err != nil {
			fmt.Printf("Addon failed: %s\n", err.Error())
			_ = addon.Failed(err.Error())
			os.Exit(1)
		}
	}()

	// Signal that addon has started
	fmt.Printf("Addon started\n")
	_ = addon.Started()

	// Validate the addon data and enforce defaults
	err = validateData(d)
	if err != nil {
		return
	}

	// Clone Git repository
	if d.Git.URL != "" {
		err = cloneGitRepository(d)
		if err != nil {
			return
		}
	}

	// Create the bucket to store the result
	bucket, err := createBucket(d)
	if err != nil {
		return
	}

	// Get the main language
	err = execWindup(d, bucket)
	if err != nil {
		return
	}

	// Task update: The addon has succeeded
	_ = addon.Succeeded()
}

//
// Validate the addon data and enforce defaults
func validateData(d *Data) (err error) {
	if d.Application == 0 {
		err = errors.New("Field 'application' is missing in addon data")
		return
	}
	if d.Git.URL == "" && d.Maven.File == "" {
		err = errors.New("neither Git URL, nor Maven File was provided")
		return
	}
	if d.Git.URL != "" && d.Maven.File != "" {
		err = errors.New("Git URL and Maven File have both been provided")
		return
	}
	if len(d.Windup.Targets) == 0 {
		fmt.Printf("No Windup target set. Using 'cloud-readiness'.\n")
		d.Windup.Targets = append(d.Windup.Targets, "cloud-readiness")
	}

	dJSON, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return
	}
	fmt.Printf("Data used by the addon:\n%s\n", string(dJSON))

	return
}

//
// Clone Git repository
// TODO: Add support for non anonymous Git operations
// TODO: Add support for fetching credentials from Hub
func cloneGitRepository(d *Data) (err error) {
	fmt.Printf("Cloning Git repository\n")
	_ = addon.Activity("cloning Git repository")

	gitCloneOptions := &git.CloneOptions{
		URL:               d.Git.URL,
		ReferenceName:     plumbing.ReferenceName(d.Git.Branch),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}
	fmt.Printf("Git clone options\n")
	fmt.Printf("  - URL: %s\n", gitCloneOptions.URL)
	fmt.Printf("  - ReferenceName: %s\n", gitCloneOptions.ReferenceName)

	_, err = git.PlainClone("/tmp/app", false, gitCloneOptions)
	if err != nil {
		return
	}

	fmt.Printf("Git clone completed\n")

	return
}

//
// Returns the most represented language in the repository
func execWindup(d *Data, bucket *api.Bucket) (err error) {
	fmt.Printf("Analyzing the repository\n")
	_ = addon.Activity("analyzing the repository")

	// Build mta-cli command arguments
	args := []string{
		"--batchMode",
		"--output", bucket.Path,
	}

	target := []string{"--target"}
	target = append(target, d.Windup.Targets...)
	args = append(args, target...)

	input := []string{"--input"}
	if d.Git.URL != "" {
		srcInput := []string{"/tmp/app" + d.Git.Path, "--sourceMode"}
		input = append(input, srcInput...)
	}
	if d.Maven.File != "" {
		input = append(input, d.Maven.File)
	}
	args = append(args, input...)

	if len(d.Windup.Packages) > 0 {
		packages := []string{"--packages"}
		packages = append(packages, d.Windup.Packages...)
		args = append(args, packages...)
	}

	// Invoke the mta-cli command
	cmd := exec.Command("/opt/mta-cli/bin/mta-cli", args...)
	cmd.Dir = "/opt/mta-cli"
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	fmt.Printf("Calling mta-cli...\n")
	err = cmd.Run()
	if err != nil {
		return
	}
	fmt.Printf("Analysis completed\n")
	fmt.Printf("Windup stdout\n%s\n", stdout.String())
	fmt.Printf("Windup stderr\n%s\n", stderr.String())

	return
}

//
// Upload full languages list as an artifact
func createBucket(d *Data) (bucket *api.Bucket, err error) {
	fmt.Printf("Creating the bucket to store the result\n")
	_ = addon.Activity("creating the bucket to store the result\n")

	bucket = &api.Bucket{}
	bucket.CreateUser = "addon"
	bucket.Name = "AnalysisWindup"
	bucket.ApplicationID = d.Application
	err = addon.Bucket.Create(bucket)
	if err != nil {
		_ = addon.Bucket.Delete(bucket)
	}

	return
}
