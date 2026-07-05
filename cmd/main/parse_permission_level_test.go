package main

import "testing"

func TestParsePermissionLevel(t *testing.T) {
	tests := []struct {
		input   string
		want    int
		wantErr bool
	}{
		// ── Valid inputs ───────────────────────────────────────────────────────
		{"0", 0, false},
		{"1", 1, false},
		{"2", 2, false},
		{"3", 3, false},
		{"4", 4, false},

		// ── Out-of-range ──────────────────────────────────────────────────────
		{"-1", 0, true},
		{"5", 0, true},
		{"99", 0, true},
		{"100", 0, true},

		// ── Non-numeric ───────────────────────────────────────────────────────
		{"", 0, true},
		{"mod", 0, true},
		{"broadcaster", 0, true},
		{"4.0", 0, true},
		{" 3", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parsePermissionLevel(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePermissionLevel(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parsePermissionLevel(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
