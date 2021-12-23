package main

import (
	"bytes"
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
	Application uint   `json:"application"`
	GitURL      string `json:"git_url"`
	GitBranch   string `json:"git_branch"`
	GitPath     string `json:"git_path"`
}

//
// main
func main() {
	var err error
	fmt.Printf("--- Tackle Addon - Discovery - Languages ---\n")

	// Get the addon data associated with the task.
	d := &Data{}
	_ = addon.DataWith(d)

	fmt.Printf("Data passed to the addon:\n")
	fmt.Printf("  - Application ID: %d\n", d.Application)
	fmt.Printf("  - Git URL: %s\n", d.GitURL)
	fmt.Printf("  - Git Branch: %s\n", d.GitBranch)
	fmt.Printf("  - Git Path: %s\n", d.GitPath)

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
	if d.Application == 0 {
		fmt.Printf("Field 'application' is missing in addon data\n")
		_ = addon.Failed("field 'application' is missing in addon data")
		os.Exit(1)
	}
	if d.GitURL == "" {
		fmt.Printf("Field 'git_url' is missing in addon data\n")
		_ = addon.Failed("field 'git_url' is missing in addon data")
		os.Exit(1)
	}

	// Clone Git repository
	cloneGitRepository(d)

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
// Clone Git repository
// TODO: Add support for non anonymous Git operations
// TODO: Add support for fetching credentials from Hub
func cloneGitRepository(d *Data) (err error) {
	fmt.Printf("Cloning Git repository\n")
	_ = addon.Activity("cloning Git repository")

	gitCloneOptions := &git.CloneOptions{
		URL:               d.GitURL,
		ReferenceName:     plumbing.ReferenceName(d.GitBranch),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}
	fmt.Printf("Git clone options\n")
	fmt.Printf("  - URL: %s\n", gitCloneOptions.URL)
	fmt.Printf("  - ReferenceName: %s\n", gitCloneOptions.ReferenceName)
	//fmt.Printf("  - RecurseSubmodules: %s", git.DefaultSubmoduleRecursionDepth)

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

	appPath := "/tmp/app/" + d.GitPath
	fmt.Printf("Application path to analyze: %s\n", appPath)

	cmd := exec.Command("/opt/mta-cli/bin/mta-cli",
		"--batchMode", "--sourceMode", "--overwrite",
		"--target", "cloud-readiness",
		"--input", appPath, "--output", bucket.Path,
	)
	cmd.Dir = "/opt/mta-cli"
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	fmt.Printf("Calling Windup...\n")
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
