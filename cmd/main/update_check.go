package main

import (
	"context"
	"log"
	"net/http"

	"github.com/google/go-github/v84/github"
	"golang.org/x/mod/semver"
)

var githubRepositoryOwner = "AzuriDayo"
var githubRepositoryName = "pear-desktop-twitch-song-requests"

func checkForUpdates() {
	if !semver.IsValid(version) || semver.Canonical(version) == "v0.0.0" {
		// invalid version possibly means development
		return
	}
	canonicalVersion := semver.Canonical(version)
	client := github.NewClient(nil)
	release, repsonse, err := client.Repositories.GetLatestRelease(context.Background(), githubRepositoryOwner, githubRepositoryName)
	if err != nil || repsonse.StatusCode != http.StatusOK {
		log.Println("Failed to check for updates. Is https://github.com alive?")
		return
	}
	if latestVersion := release.GetTagName(); latestVersion != canonicalVersion {
		log.Println("New version of " + githubRepositoryName + " found! " + latestVersion)
		checkUpdatesResult(latestVersion)
	} else {
		log.Println("You are running the latest version of " + githubRepositoryName)
	}
}
