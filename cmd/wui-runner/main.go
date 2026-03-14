//go:build windows

package main

import (
	"os/exec"

	"github.com/ChromeTemp/Popup"
)

func main() {
	b := Popup.LazyDialog("Open github", "Do you want to open github?")
	if b {
		exec.Command("cmd", "/c", "start", "https://github.com/AzuriDayo/pear-desktop-twitch-song-requests").Start()
	}
}
