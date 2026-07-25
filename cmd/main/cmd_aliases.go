package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Canonical command short names that may be targeted by aliases or disabled.
var aliasableCommands = map[string]bool{
	"sr":      true,
	"queue":   true,
	"song":    true,
	"delsong": true,
	"skip":    true,
}

// Built-in command tokens (with leading !) that aliases must not collide with.
var builtinCommandTokens = map[string]bool{
	"!sr":        true,
	"!queue":     true,
	"!song":      true,
	"!delsong":   true,
	"!skip":      true,
	"!srversion": true,
}

// Commands that take arguments after the command token.
var aliasArgCommands = map[string]bool{
	"sr":      true,
	"delsong": true,
}

var aliasPattern = regexp.MustCompile(`^![a-z0-9_]{1,31}$`)

// normalizeAlias ensures the alias starts with "!" and is lowercase.
func normalizeAlias(alias string) string {
	alias = strings.TrimSpace(strings.ToLower(alias))
	if alias == "" {
		return ""
	}
	if !strings.HasPrefix(alias, "!") {
		alias = "!" + alias
	}
	return alias
}

// normalizeCmdName lowercases and strips a leading "!" from a command short name.
func normalizeCmdName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	return strings.TrimPrefix(name, "!")
}

// parseCmdAliasesJSON unmarshals a cmd_aliases JSON object into a normalized map.
// Invalid individual entries are skipped; malformed JSON returns an error.
func parseCmdAliasesJSON(raw string) (map[string]string, error) {
	if strings.TrimSpace(raw) == "" {
		return map[string]string{}, nil
	}
	var incoming map[string]string
	if err := json.Unmarshal([]byte(raw), &incoming); err != nil {
		return nil, fmt.Errorf("cmd_aliases must be a JSON object: %w", err)
	}
	out := make(map[string]string, len(incoming))
	for k, v := range incoming {
		alias := normalizeAlias(k)
		target := normalizeCmdName(v)
		if alias == "" || target == "" {
			continue
		}
		out[alias] = target
	}
	return out, nil
}

// validateCmdAliases checks that every alias/target pair is valid.
// Returns a normalized copy on success.
func validateCmdAliases(aliases map[string]string) (map[string]string, error) {
	if aliases == nil {
		return map[string]string{}, nil
	}
	out := make(map[string]string, len(aliases))
	for k, v := range aliases {
		alias := normalizeAlias(k)
		target := normalizeCmdName(v)

		if alias == "" {
			return nil, fmt.Errorf("alias must not be empty")
		}
		if !aliasPattern.MatchString(alias) {
			return nil, fmt.Errorf("invalid alias %q: must match ![a-z0-9_] (2–32 chars including !)", alias)
		}
		if builtinCommandTokens[alias] {
			return nil, fmt.Errorf("alias %q conflicts with a built-in command", alias)
		}
		if !aliasableCommands[target] {
			return nil, fmt.Errorf("invalid alias target %q for %q: must be one of sr, queue, song, delsong, skip", v, alias)
		}
		if _, exists := out[alias]; exists {
			return nil, fmt.Errorf("duplicate alias %q", alias)
		}
		out[alias] = target
	}
	return out, nil
}

// marshalCmdAliasesJSON serializes aliases to a JSON object string.
func marshalCmdAliasesJSON(aliases map[string]string) (string, error) {
	if aliases == nil {
		aliases = map[string]string{}
	}
	b, err := json.Marshal(aliases)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// parseCmdDisabledBuiltinsJSON unmarshals a JSON array of command short names.
func parseCmdDisabledBuiltinsJSON(raw string) (map[string]bool, error) {
	if strings.TrimSpace(raw) == "" {
		return map[string]bool{}, nil
	}
	var incoming []string
	if err := json.Unmarshal([]byte(raw), &incoming); err != nil {
		return nil, fmt.Errorf("cmd_disabled_builtins must be a JSON array: %w", err)
	}
	out := make(map[string]bool, len(incoming))
	for _, name := range incoming {
		n := normalizeCmdName(name)
		if n == "" {
			continue
		}
		out[n] = true
	}
	return out, nil
}

// validateCmdDisabledBuiltins checks that every name is a known disableable command.
func validateCmdDisabledBuiltins(disabled map[string]bool) (map[string]bool, error) {
	if disabled == nil {
		return map[string]bool{}, nil
	}
	out := make(map[string]bool, len(disabled))
	for name, isDisabled := range disabled {
		if !isDisabled {
			continue
		}
		n := normalizeCmdName(name)
		if !aliasableCommands[n] {
			return nil, fmt.Errorf("invalid disabled built-in %q: must be one of sr, queue, song, delsong, skip", name)
		}
		out[n] = true
	}
	return out, nil
}

// marshalCmdDisabledBuiltinsJSON serializes disabled built-ins to a JSON array string.
func marshalCmdDisabledBuiltinsJSON(disabled map[string]bool) (string, error) {
	names := make([]string, 0, len(disabled))
	for name, isDisabled := range disabled {
		if isDisabled {
			names = append(names, name)
		}
	}
	// Stable order for predictable API/UI diffs.
	order := []string{"sr", "queue", "song", "delsong", "skip"}
	sorted := make([]string, 0, len(names))
	seen := map[string]bool{}
	for _, n := range order {
		if disabled[n] {
			sorted = append(sorted, n)
			seen[n] = true
		}
	}
	for _, n := range names {
		if !seen[n] {
			sorted = append(sorted, n)
		}
	}
	b, err := json.Marshal(sorted)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// resolveCmdAlias rewrites aliased command text to the canonical built-in form.
// fromAlias is true when an alias was applied.
// Examples: "!q" -> "!queue", "!np" -> "!song", "!req yena" -> "!sr yena".
func (a *App) resolveCmdAlias(text string) (resolved string, fromAlias bool) {
	if a == nil || len(a.cmdAliases) == 0 || text == "" {
		return text, false
	}

	fields := strings.Fields(text)
	if len(fields) == 0 {
		return text, false
	}

	token := strings.ToLower(fields[0])
	target, ok := a.cmdAliases[token]
	if !ok {
		return text, false
	}

	canonical := "!" + target
	if aliasArgCommands[target] {
		if len(fields) == 1 {
			return canonical, true
		}
		rest := strings.TrimSpace(text[len(fields[0]):])
		return canonical + " " + rest, true
	}

	// Exact-match commands: only rewrite when the whole message is the alias.
	if len(fields) != 1 {
		return text, false
	}
	return canonical, true
}

// builtinAllowed reports whether the built-in command name may run.
// Aliases targeting a disabled built-in are still allowed (fromAlias=true).
func (a *App) builtinAllowed(name string, fromAlias bool) bool {
	if fromAlias {
		return true
	}
	if a == nil || a.cmdDisabledBuiltins == nil {
		return true
	}
	return !a.cmdDisabledBuiltins[normalizeCmdName(name)]
}
