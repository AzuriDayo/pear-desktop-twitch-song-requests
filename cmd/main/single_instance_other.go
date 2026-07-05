//go:build !windows

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const singleInstanceLockFile = "pear-desktop-twitch-song-requests.lock"

func singleInstanceLockPath() string {
	return filepath.Join(os.TempDir(), singleInstanceLockFile)
}

// acquireSingleInstanceLock writes the current PID to a lock file in the OS
// temp directory. If a lock file already exists and its PID belongs to a live
// process, the user is prompted to force-quit it. Stale lock files left by a
// crashed previous run are silently overwritten (the dead PID won't respond to
// Signal(0)).
func acquireSingleInstanceLock() {
	lockPath := singleInstanceLockPath()

	if raw, err := os.ReadFile(lockPath); err == nil {
		pid, err := strconv.Atoi(strings.TrimSpace(string(raw)))
		if err == nil {
			if proc, err := os.FindProcess(pid); err == nil && proc.Signal(syscall.Signal(0)) == nil {
				fmt.Printf("Pear Desktop is already open (PID %d).\n", pid)
				fmt.Print("Close the existing instance and relaunch? [y/N]: ")

				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
				if answer != "y" && answer != "yes" {
					os.Exit(1)
				}

				if killErr := proc.Kill(); killErr != nil {
					log.Printf("Warning: could not kill old instance (PID %d): %v", pid, killErr)
				} else {
					log.Printf("Killed old instance (PID %d).", pid)
				}
			}
		}
	}

	if err := os.WriteFile(lockPath, []byte(fmt.Sprintf("%d", os.Getpid())), 0o644); err != nil {
		log.Printf("Warning: could not write single-instance lock file: %v", err)
	}
}

// releaseSingleInstanceLock removes the lock file on a clean shutdown so the
// next launch doesn't have to fall back to the stale-PID liveness check.
func releaseSingleInstanceLock() {
	if err := os.Remove(singleInstanceLockPath()); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: could not remove single-instance lock file: %v", err)
	}
}
