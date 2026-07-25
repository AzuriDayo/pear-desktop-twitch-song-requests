package main

import (
	"strings"
	"testing"
)

func TestNormalizeAlias(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"q", "!q"},
		{"!Q", "!q"},
		{"  !Np  ", "!np"},
		{"", ""},
		{"   ", ""},
	}
	for _, tt := range tests {
		if got := normalizeAlias(tt.in); got != tt.want {
			t.Errorf("normalizeAlias(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestValidateCmdAliases(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		got, err := validateCmdAliases(map[string]string{
			"!q":  "queue",
			"np":  "!song",
			"!REQ": "sr",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got["!q"] != "queue" || got["!np"] != "song" || got["!req"] != "sr" {
			t.Fatalf("unexpected map: %#v", got)
		}
	})

	t.Run("rejects built-in collision", func(t *testing.T) {
		_, err := validateCmdAliases(map[string]string{"!queue": "song"})
		if err == nil || !strings.Contains(err.Error(), "conflicts with a built-in") {
			t.Fatalf("expected built-in conflict error, got %v", err)
		}
	})

	t.Run("rejects invalid target", func(t *testing.T) {
		_, err := validateCmdAliases(map[string]string{"!q": "srversion"})
		if err == nil || !strings.Contains(err.Error(), "invalid alias target") {
			t.Fatalf("expected invalid target error, got %v", err)
		}
	})

	t.Run("rejects bad alias pattern", func(t *testing.T) {
		_, err := validateCmdAliases(map[string]string{"!bad alias": "queue"})
		if err == nil || !strings.Contains(err.Error(), "invalid alias") {
			t.Fatalf("expected invalid alias error, got %v", err)
		}
	})

	t.Run("nil becomes empty", func(t *testing.T) {
		got, err := validateCmdAliases(nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Fatalf("want empty map, got %#v", got)
		}
	})
}

func TestValidateCmdDisabledBuiltins(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		got, err := validateCmdDisabledBuiltins(map[string]bool{
			"queue": true,
			"!Song": true,
			"sr":    false, // ignored
		})
		if err != nil {
			t.Fatal(err)
		}
		if !got["queue"] || !got["song"] || got["sr"] {
			t.Fatalf("unexpected map: %#v", got)
		}
	})

	t.Run("rejects unknown", func(t *testing.T) {
		_, err := validateCmdDisabledBuiltins(map[string]bool{"srversion": true})
		if err == nil || !strings.Contains(err.Error(), "invalid disabled built-in") {
			t.Fatalf("expected invalid disabled error, got %v", err)
		}
	})
}

func TestParseAndMarshalCmdAliasesJSON(t *testing.T) {
	raw := `{"!q":"queue","NP":"song"}`
	parsed, err := parseCmdAliasesJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	validated, err := validateCmdAliases(parsed)
	if err != nil {
		t.Fatal(err)
	}
	out, err := marshalCmdAliasesJSON(validated)
	if err != nil {
		t.Fatal(err)
	}
	roundTrip, err := parseCmdAliasesJSON(out)
	if err != nil {
		t.Fatal(err)
	}
	if roundTrip["!q"] != "queue" || roundTrip["!np"] != "song" {
		t.Fatalf("round trip failed: %#v", roundTrip)
	}
}

func TestParseAndMarshalCmdDisabledBuiltinsJSON(t *testing.T) {
	raw := `["queue","!song"]`
	parsed, err := parseCmdDisabledBuiltinsJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	validated, err := validateCmdDisabledBuiltins(parsed)
	if err != nil {
		t.Fatal(err)
	}
	out, err := marshalCmdDisabledBuiltinsJSON(validated)
	if err != nil {
		t.Fatal(err)
	}
	if out != `["queue","song"]` {
		t.Fatalf("marshal = %s, want [\"queue\",\"song\"]", out)
	}
}

func TestResolveCmdAlias(t *testing.T) {
	a := &App{
		cmdAliases: map[string]string{
			"!q":   "queue",
			"!np":  "song",
			"!req": "sr",
			"!ds":  "delsong",
		},
	}

	tests := []struct {
		in         string
		wantText   string
		wantAlias  bool
	}{
		{"!q", "!queue", true},
		{"!Q", "!queue", true},
		{"!np", "!song", true},
		{"!req yena nemonemo", "!sr yena nemonemo", true},
		{"!REQ  https://youtu.be/abc", "!sr https://youtu.be/abc", true},
		{"!ds 2", "!delsong 2", true},
		{"!ds", "!delsong", true},
		{"!queue", "!queue", false},
		{"!q extra", "!q extra", false}, // exact-match aliases ignore trailing args
		{"hello", "hello", false},
		{"", "", false},
	}

	for _, tt := range tests {
		got, fromAlias := a.resolveCmdAlias(tt.in)
		if got != tt.wantText || fromAlias != tt.wantAlias {
			t.Errorf("resolveCmdAlias(%q) = (%q, %v), want (%q, %v)",
				tt.in, got, fromAlias, tt.wantText, tt.wantAlias)
		}
	}

	t.Run("empty aliases", func(t *testing.T) {
		empty := &App{cmdAliases: map[string]string{}}
		got, fromAlias := empty.resolveCmdAlias("!q")
		if got != "!q" || fromAlias {
			t.Fatalf("got (%q, %v)", got, fromAlias)
		}
	})
}

func TestBuiltinAllowed(t *testing.T) {
	a := &App{
		cmdDisabledBuiltins: map[string]bool{
			"queue": true,
			"song":  true,
		},
	}

	if a.builtinAllowed("queue", false) {
		t.Error("disabled built-in should not be allowed without alias")
	}
	if !a.builtinAllowed("queue", true) {
		t.Error("disabled built-in should still be allowed via alias")
	}
	if !a.builtinAllowed("sr", false) {
		t.Error("enabled built-in should be allowed")
	}
	if !a.builtinAllowed("skip", true) {
		t.Error("enabled built-in via alias should be allowed")
	}

	t.Run("nil app / nil map", func(t *testing.T) {
		var nilApp *App
		if !nilApp.builtinAllowed("queue", false) {
			t.Error("nil app should allow")
		}
		empty := &App{}
		if !empty.builtinAllowed("queue", false) {
			t.Error("nil disabled map should allow")
		}
	})
}

func TestResolveAliasThenBuiltinAllowed(t *testing.T) {
	// Simulates the chat-handler path: alias rewrite + disable original !queue.
	a := &App{
		cmdAliases:          map[string]string{"!q": "queue"},
		cmdDisabledBuiltins: map[string]bool{"queue": true},
	}

	resolved, fromAlias := a.resolveCmdAlias("!q")
	if resolved != "!queue" || !fromAlias {
		t.Fatalf("alias resolve failed: %q %v", resolved, fromAlias)
	}
	if !a.builtinAllowed("queue", fromAlias) {
		t.Fatal("!q should still run queue when built-in is disabled")
	}

	resolved, fromAlias = a.resolveCmdAlias("!queue")
	if resolved != "!queue" || fromAlias {
		t.Fatalf("built-in resolve unexpected: %q %v", resolved, fromAlias)
	}
	if a.builtinAllowed("queue", fromAlias) {
		t.Fatal("!queue should be blocked when disabled")
	}
}
