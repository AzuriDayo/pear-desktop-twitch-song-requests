package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/databaseconn"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
	"github.com/labstack/echo/v4"
)

func TestAPIPayloadDumpSettings(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := databaseconn.Migrate(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { songrequests.SetAPIPayloadDumpEnabled(false) })
	a := NewApp()
	t.Cleanup(a.cancel)
	a.clientsBroadcast = make(chan string, 10)
	if err := a.loadSqliteSettings(); err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	assertSetting := func(want bool) {
		t.Helper()
		rec := httptest.NewRecorder()
		if err := a.handleApiV1SettingsGET(e.NewContext(httptest.NewRequest(http.MethodGet, "/api/v1/settings", nil), rec)); err != nil {
			t.Fatal(err)
		}
		var result map[string]json.RawMessage
		if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		var enabled bool
		if err := json.Unmarshal(result[data.DB_KEY_PEAR_DESKTOP_API_PAYLOAD_DUMP_ENABLED], &enabled); err != nil || enabled != want {
			t.Fatalf("GET setting = %s, want %v, error = %v", result[data.DB_KEY_PEAR_DESKTOP_API_PAYLOAD_DUMP_ENABLED], want, err)
		}
		if songrequests.APIPayloadDumpEnabled() != want {
			t.Fatalf("runtime setting = %v, want %v", songrequests.APIPayloadDumpEnabled(), want)
		}
	}
	assertSetting(false)
	for _, value := range []string{"true", "false"} {
		rec := httptest.NewRecorder()
		body := `{"pear_desktop_api_payload_dump_enabled":"` + value + `"}`
		if err := a.handleApiV1SettingsPATCH(e.NewContext(httptest.NewRequest(http.MethodPatch, "/api/v1/settings", strings.NewReader(body)), rec)); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("PATCH status = %d: %s", rec.Code, rec.Body.String())
		}
		want := value == "true"
		assertSetting(want)
		// Reset runtime state to prove the saved database value survives a reload.
		songrequests.SetAPIPayloadDumpEnabled(!want)
		if err := a.loadSqliteSettings(); err != nil {
			t.Fatal(err)
		}
		assertSetting(want)
	}
	for _, value := range []string{`"invalid"`, `true`, `null`} {
		rec := httptest.NewRecorder()
		body := `{"pear_desktop_api_payload_dump_enabled":` + value + `}`
		if err := a.handleApiV1SettingsPATCH(e.NewContext(httptest.NewRequest(http.MethodPatch, "/api/v1/settings", strings.NewReader(body)), rec)); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("invalid setting %s returned status %d", value, rec.Code)
		}
		assertSetting(false)
	}
}
