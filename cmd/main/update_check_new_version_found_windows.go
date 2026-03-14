//go:build windows

package main

import (
	"os/exec"

	"github.com/ChromeTemp/Popup"
)

func checkUpdatesResult(newVersion string) {
	b := Popup.LazyDialog("New version found", "Do you want to open the download URL for version "+newVersion+"?")
	if b {
		exec.Command("cmd", "/c", "start", "https://github.com/"+githubRepositoryOwner+"/"+githubRepositoryName+"/releases/tag/"+newVersion).Start()
	}
}
