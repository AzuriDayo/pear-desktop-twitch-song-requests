package main

import "testing"

func TestUserMeetsPermission(t *testing.T) {
	tests := []struct {
		name          string
		level         int
		isBroadcaster bool
		isModerator   bool
		isVip         bool
		isSub         bool
		want          bool
	}{
		// ── Level 0: broadcaster only ──────────────────────────────────────────
		{"broadcaster at level 0", 0, true, false, false, false, true},
		{"mod at level 0", 0, false, true, false, false, false},
		{"vip at level 0", 0, false, false, true, false, false},
		{"sub at level 0", 0, false, false, false, true, false},
		{"viewer at level 0", 0, false, false, false, false, false},

		// ── Level 1: moderator or above ───────────────────────────────────────
		{"broadcaster at level 1", 1, true, false, false, false, true},
		{"mod at level 1", 1, false, true, false, false, true},
		{"vip at level 1", 1, false, false, true, false, false},
		{"sub at level 1", 1, false, false, false, true, false},
		{"viewer at level 1", 1, false, false, false, false, false},

		// ── Level 2: VIP or above ─────────────────────────────────────────────
		{"broadcaster at level 2", 2, true, false, false, false, true},
		{"mod at level 2", 2, false, true, false, false, true},
		{"vip at level 2", 2, false, false, true, false, true},
		{"sub at level 2", 2, false, false, false, true, false},
		{"viewer at level 2", 2, false, false, false, false, false},

		// ── Level 3: subscriber or above ──────────────────────────────────────
		{"broadcaster at level 3", 3, true, false, false, false, true},
		{"mod at level 3", 3, false, true, false, false, true},
		{"vip at level 3", 3, false, false, true, false, true},
		{"sub at level 3", 3, false, false, false, true, true},
		{"viewer at level 3", 3, false, false, false, false, false},

		// ── Level 4: viewer / everyone ────────────────────────────────────────
		{"broadcaster at level 4", 4, true, false, false, false, true},
		{"mod at level 4", 4, false, true, false, false, true},
		{"vip at level 4", 4, false, false, true, false, true},
		{"sub at level 4", 4, false, false, false, true, true},
		{"viewer at level 4", 4, false, false, false, false, true},

		// ── Boundary / invalid levels ──────────────────────────────────────────
		{"level -1 always false", -1, true, true, true, true, false},
		{"level 5 always false", 5, true, true, true, true, false},
		{"level 99 always false", 99, true, true, true, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserMeetsPermission(tt.level, tt.isBroadcaster, tt.isModerator, tt.isVip, tt.isSub)
			if got != tt.want {
				t.Errorf("UserMeetsPermission(%d, broadcaster=%v, mod=%v, vip=%v, sub=%v) = %v, want %v",
					tt.level, tt.isBroadcaster, tt.isModerator, tt.isVip, tt.isSub, got, tt.want)
			}
		})
	}
}
