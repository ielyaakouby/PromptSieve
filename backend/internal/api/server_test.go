package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/promptsieve/promptsieve/internal/detectors"
	"github.com/promptsieve/promptsieve/internal/scanner"
)

func TestScanDoesNotReturnRawPassword(t *testing.T) {
	sc := scanner.NewScanner(scanner.Config{Timeout: 5 * time.Second})
	for _, detector := range detectors.Default() {
		sc.Register(detector)
	}
	srv := httptest.NewServer(NewHandler(sc))
	defer srv.Close()

	body := `{"text":"mot de passe : FakePassword123!"}`
	res, err := http.Post(srv.URL+"/api/v1/scan", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var payload ScanResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(payload)
	if strings.Contains(string(raw), "FakePassword123!") {
		t.Fatal("raw password leaked in API response")
	}
	if payload.Status == "SAFE" {
		t.Fatal("password prompt marked SAFE")
	}
	if !strings.Contains(payload.SanitizedText, "<REDACTED_PASSWORD>") {
		t.Fatalf("sanitized=%q", payload.SanitizedText)
	}
}

func TestScanRejectsInvalidJSON(t *testing.T) {
	sc := scanner.NewScanner(scanner.Config{Timeout: time.Second})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/scan", bytes.NewBufferString("{"))
	rec := httptest.NewRecorder()
	NewHandler(sc).ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestScanContextDoesNotNeedExternalServices(t *testing.T) {
	sc := scanner.NewScanner(scanner.Config{Timeout: time.Second})
	_, err := sc.Scan(context.Background(), "hello")
	if err != nil {
		t.Fatal(err)
	}
}
