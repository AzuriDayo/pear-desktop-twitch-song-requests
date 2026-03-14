//go:build !windows

package main

import "log"

func checkUpdatesResult(newVersion string) {
	log.Println("Download the latest release here: https://github.com/" + githubRepositoryOwner + "/" + githubRepositoryName + "/releases/tag/" + newVersion)
}
