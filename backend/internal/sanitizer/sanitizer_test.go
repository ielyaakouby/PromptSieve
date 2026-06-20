package sanitizer

import (
	"strings"
	"testing"

	"github.com/promptsieve/promptsieve/pkg/models"
)

func TestSanitizeUsesFindingRanges(t *testing.T) {
	input := "user=alpha password=secret token=keep"
	got := Sanitize(input, []models.Finding{
		{Start: 18, End: 24, ReplacementToken: "<REDACTED_PASSWORD>"},
	})
	if strings.Contains(got, "secret") {
		t.Fatalf("got %q", got)
	}
	if !strings.Contains(got, "token=keep") {
		t.Fatalf("unrelated text changed: %q", got)
	}
}

func TestSanitizeOverlapKeepsLonger(t *testing.T) {
	input := "ABCDEFGH"
	got := Sanitize(input, []models.Finding{
		{Start: 0, End: 4, ReplacementToken: "<A>"},
		{Start: 0, End: 8, ReplacementToken: "<ALL>"},
	})
	if got != "<ALL>" {
		t.Fatalf("got %q", got)
	}
}

func TestSanitizeFallsBackToMaskedValue(t *testing.T) {
	input := "abc"
	got := Sanitize(input, []models.Finding{
		{Start: 0, End: 3, MaskedValue: "<X>"},
	})
	if got != "<X>" {
		t.Fatalf("got %q", got)
	}
}
