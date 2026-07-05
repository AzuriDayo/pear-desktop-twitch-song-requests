package main

import (
	"testing"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
)

func TestDefaultCmdPermissions(t *testing.T) {
	defaults := defaultCmdPermissions()

	expected := map[string]int{
		data.DB_KEY_CMD_PERMISSION_SR:      data.PermissionLevelSubscriber, // 3
		data.DB_KEY_CMD_PERMISSION_QUEUE:   data.PermissionLevelViewer,     // 4
		data.DB_KEY_CMD_PERMISSION_SONG:    data.PermissionLevelViewer,     // 4
		data.DB_KEY_CMD_PERMISSION_DELSONG: data.PermissionLevelModerator,  // 1
	}

	for key, wantLevel := range expected {
		t.Run(key, func(t *testing.T) {
			got, ok := defaults[key]
			if !ok {
				t.Errorf("defaultCmdPermissions() missing key %q", key)
				return
			}
			if got != wantLevel {
				t.Errorf("defaultCmdPermissions()[%q] = %d, want %d", key, got, wantLevel)
			}
		})
	}

	// Ensure no unexpected keys are present
	if len(defaults) != len(expected) {
		t.Errorf("defaultCmdPermissions() has %d keys, want %d", len(defaults), len(expected))
	}
}
