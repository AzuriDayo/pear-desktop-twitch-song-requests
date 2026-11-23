package songrequests

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

var apiPayloadDumpEnabled atomic.Bool

func SetAPIPayloadDumpEnabled(enabled bool) {
	apiPayloadDumpEnabled.Store(enabled)
}

func APIPayloadDumpEnabled() bool {
	return apiPayloadDumpEnabled.Load()
}

func searchPayloadDumpFilename(timestamp time.Time) string {
	name := timestamp.UTC().Format("2006-01-02T15:04:05.000Z") + ".json"
	if runtime.GOOS == "windows" {
		name = strings.ReplaceAll(name, ":", "-")
	}
	return name
}

func dumpSearchPayload(query string, payload []byte) error {
	if !APIPayloadDumpEnabled() {
		return nil
	}

	dump, err := json.Marshal(struct {
		SearchQuery string         `json:"searchQuery"`
		Payload     jsontext.Value `json:"payload"`
	}{SearchQuery: query, Payload: payload},
		jsontext.WithIndent("  "),
		jsontext.PreserveRawStrings(true),
		jsontext.AllowDuplicateNames(true),
	)
	if err != nil {
		return err
	}
	dir := filepath.Join("debug", "search")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Reserve each filename exclusively so searches in the same millisecond
	// cannot overwrite an earlier dump.
	for timestamp := time.Now(); ; timestamp = timestamp.Add(time.Millisecond) {
		path := filepath.Join(dir, searchPayloadDumpFilename(timestamp))
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if errors.Is(err, os.ErrExist) {
			continue
		}
		if err != nil {
			return err
		}
		_, writeErr := file.Write(dump)
		closeErr := file.Close()
		return errors.Join(writeErr, closeErr)
	}
}
