//go:build windows

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/windows"
)

const (
	// Local\ scopes the mutex to the current login session, which is the right
	// granularity for a per-user desktop app and avoids needing the
	// SeCreateGlobalPrivilege that the Global\ namespace requires.
	singleInstanceMutexName = "Local\\PearDesktopTwitchSongRequests"
	singleInstancePIDFile   = "pear-desktop-twitch-song-requests.pid"
)

func singleInstancePIDPath() string {
	return filepath.Join(os.TempDir(), singleInstancePIDFile)
}

func writeSingleInstancePID() {
	if err := os.WriteFile(singleInstancePIDPath(), []byte(strconv.Itoa(os.Getpid())), 0o644); err != nil {
		log.Printf("Warning: could not write single-instance PID file: %v", err)
	}
}

// releaseSingleInstanceLock removes the PID file on a clean shutdown. The named
// mutex is released automatically by the OS when the process exits, so this only
// tidies up the informational PID file.
func releaseSingleInstanceLock() {
	if err := os.Remove(singleInstancePIDPath()); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: could not remove single-instance PID file: %v", err)
	}
}

// tryCreateMutex attempts to create (or open) the named mutex.
// Returns (acquired, handle). When acquired is true the caller owns the mutex
// and must NOT close the handle — the OS releases it on process exit.
// When acquired is false the handle is already closed.
func tryCreateMutex() (bool, windows.Handle) {
	h, err := windows.CreateMutex(nil, true, windows.StringToUTF16Ptr(singleInstanceMutexName))
	if err == windows.ERROR_ALREADY_EXISTS {
		windows.CloseHandle(h)
		return false, 0
	}
	if err != nil {
		log.Printf("Warning: could not create single-instance mutex: %v", err)
	}
	return true, h
}

// acquireSingleInstanceLock ensures only one instance runs at a time using a
// named Windows mutex. If another instance already owns the mutex the user is
// prompted to force-quit it. Stale PID files from a clean or dirty previous
// exit are handled automatically because the OS always releases the mutex when
// the owning process dies.
func acquireSingleInstanceLock() {
	acquired, _ := tryCreateMutex()
	if acquired {
		writeSingleInstancePID()
		return
	}

	// A live instance owns the mutex. Read its PID so we can kill it.
	pid := 0
	if raw, err := os.ReadFile(singleInstancePIDPath()); err == nil {
		pid, _ = strconv.Atoi(strings.TrimSpace(string(raw)))
	}

	if pid != 0 {
		fmt.Printf("Pear Desktop is already open (PID %d).\n", pid)
	} else {
		fmt.Println("Pear Desktop is already open.")
	}
	fmt.Print("Close the existing instance and relaunch? [y/N]: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if answer != "y" && answer != "yes" {
		os.Exit(1)
	}

	if pid != 0 {
		proc, err := os.FindProcess(pid)
		if err == nil {
			if killErr := proc.Kill(); killErr != nil {
				log.Printf("Warning: could not kill old instance (PID %d): %v", pid, killErr)
			} else {
				log.Printf("Killed old instance (PID %d), waiting for mutex to be released...", pid)
			}
		}
	}

	// Poll until the OS releases the mutex from the dead process (up to ~3 s).
	for range 10 {
		time.Sleep(300 * time.Millisecond)
		if ok, _ := tryCreateMutex(); ok {
			writeSingleInstancePID()
			return
		}
	}

	log.Println("Could not re-acquire single-instance lock after force-quit. Exiting.")
	os.Exit(1)
}
