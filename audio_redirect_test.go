package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLegacyAudioRedirect(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/static/audio/2026-09-28-the-empty-chair-is-still-information.mp3", nil)
	rec := httptest.NewRecorder()
	legacyAudioRedirectHandler(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("status=%d", rec.Code)
	}
	want := "https://static.strangerloops.com/audio/2026-09-28-the-empty-chair-is-still-information.mp3?v=20261004-r2"
	if got := rec.Header().Get("Location"); got != want {
		t.Fatalf("Location=%q, want %q", got, want)
	}
}

func TestLegacyAudioRedirectRejectsInvalidNames(t *testing.T) {
	for _, path := range []string{"/static/audio/", "/static/audio/private.txt", "/static/audio/subdir/file.mp3"} {
		rec := httptest.NewRecorder()
		legacyAudioRedirectHandler(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusNotFound {
			t.Errorf("%s: status=%d", path, rec.Code)
		}
	}
}
