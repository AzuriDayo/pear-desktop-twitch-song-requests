package songrequests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestSearchPayloadDumpFilename(t *testing.T) {
	timestamp := time.Date(2026, 11, 23, 13, 4, 5, 123456789, time.FixedZone("EST", -5*60*60))
	want := "2026-11-23T18:04:05.123Z.json"
	if runtime.GOOS == "windows" {
		want = strings.ReplaceAll(want, ":", "-")
	}
	if got := searchPayloadDumpFilename(timestamp); got != want {
		t.Fatalf("filename = %q, want %q", got, want)
	}
	if got := searchPayloadDumpFilename(timestamp.Add(-123456789)); !strings.HasSuffix(got, ".000Z.json") {
		t.Fatalf("zero milliseconds must be retained: %q", got)
	}
}

func TestSearchSongPayloadDump(t *testing.T) {
	t.Chdir(t.TempDir())
	SetAPIPayloadDumpEnabled(false)
	t.Cleanup(func() { SetAPIPayloadDumpEnabled(false) })

	const payload = `{"contents":{},"unknown":{"nested":[1,true,null,"\u0065xtra"],"large":9007199254740993,"duplicate":1,"duplicate":2}}`
	const query = `artist "song"`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/search" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body struct {
			Query string `json:"query"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Query != query {
			t.Errorf("search request query = %q, error = %v", body.Query, err)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, payload)
	}))
	defer server.Close()
	originalHost := pearDesktopHost
	pearDesktopHost = strings.TrimPrefix(server.URL, "http://")
	t.Cleanup(func() { pearDesktopHost = originalHost })

	search := func() {
		t.Helper()
		if _, err := SearchSong("  "+query+" \n", 0, 600); !errors.Is(err, ErrNoResults) {
			t.Fatalf("search behavior changed: %v", err)
		}
	}
	search()
	if _, err := os.Stat("debug"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("disabled dump should not create a directory: %v", err)
	}

	SetAPIPayloadDumpEnabled(true)
	for range 10 {
		search()
	}
	SetAPIPayloadDumpEnabled(false)
	search()

	files, err := filepath.Glob(filepath.Join("debug", "search", "*.json"))
	if err != nil || len(files) != 10 {
		t.Fatalf("want 10 distinct dumps, got %d, error = %v", len(files), err)
	}
	for _, path := range files {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var dump map[string]json.RawMessage
		if err := json.Unmarshal(content, &dump); err != nil {
			t.Fatal(err)
		}
		var dumpedQuery string
		if err := json.Unmarshal(dump["searchQuery"], &dumpedQuery); err != nil || dumpedQuery != query || len(dump) != 2 {
			t.Fatalf("unexpected dump envelope: %s", content)
		}
		var compact bytes.Buffer
		if err := json.Compact(&compact, dump["payload"]); err != nil || compact.String() != payload {
			t.Fatalf("raw response was not preserved: %s (error = %v)", content, err)
		}
	}

	// Filesystem failures must not turn a successful API exchange into a failed search.
	t.Chdir(t.TempDir())
	if err := os.WriteFile("debug", []byte("blocks directory creation"), 0644); err != nil {
		t.Fatal(err)
	}
	SetAPIPayloadDumpEnabled(true)
	search()
}
