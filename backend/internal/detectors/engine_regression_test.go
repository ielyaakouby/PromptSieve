package detectors

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/promptsieve/promptsieve/internal/risk"
	"github.com/promptsieve/promptsieve/internal/sanitizer"
	"github.com/promptsieve/promptsieve/internal/scanner"
	"github.com/promptsieve/promptsieve/pkg/models"
)

const (
	fakeOpenAIKey = "sk-abcdefghijklmnopqrstuvwxyzABCDEF12"
	fakeGitHubPAT = "ghp_abcdefghijklmnopqrstuvwxyz012345"
	fakeAWSAccess = "AKIAIOSFODNN7EXAMPLE"
	fakeAWSSecret = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	fakeJWT       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	fakeCard      = "4242424242424242"
	fakePassword  = "FakePassword123!"
	fakeDBPass    = "FakeDatabasePassword123!"
	fakeEmail     = "jean.dupont@example.com"
	fakePhone     = "+33 6 12 34 56 78"
	fakeDOB       = "14/05/1987"
	fakeAddress   = "25 rue des Fleurs, 75015 Paris"
	fakeName      = "Jean Dupont"
	fakeCVV       = "123"
	fakeExpiry    = "12/28"
	fakeBIC       = "BNPAFRPPXXX"
)

func scanAll(t *testing.T, input string) []models.Finding {
	t.Helper()
	sc := scanner.NewScanner(scanner.Config{Timeout: 5 * time.Second})
	for _, detector := range Default() {
		sc.Register(detector)
	}
	findings, err := sc.Scan(context.Background(), input)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	return findings
}

func typesOf(findings []models.Finding) map[string]bool {
	out := map[string]bool{}
	for _, f := range findings {
		out[f.Type] = true
	}
	return out
}

func mustHave(t *testing.T, found map[string]bool, types ...string) {
	t.Helper()
	for _, typ := range types {
		if !found[typ] {
			t.Errorf("missing finding type %q", typ)
		}
	}
}

func mustNotHaveRaw(t *testing.T, sanitized string, secrets ...string) {
	t.Helper()
	for _, secret := range secrets {
		if secret != "" && strings.Contains(sanitized, secret) {
			t.Errorf("sanitized prompt still contains %q", secret)
		}
	}
}

func mustNotExposeSecrets(t *testing.T, findings []models.Finding, secrets ...string) {
	t.Helper()
	payload, err := json.Marshal(findings)
	if err != nil {
		t.Fatal(err)
	}
	body := string(payload)
	for _, secret := range secrets {
		if secret != "" && strings.Contains(body, secret) {
			t.Errorf("API findings leaked secret %q", secret)
		}
	}
	for _, f := range findings {
		if f.Match != "" && strings.Contains(body, f.Match) && f.Match != f.MaskedValue {
			t.Errorf("finding match leaked for type %s", f.Type)
		}
	}
}

func TestFrenchPersonalData(t *testing.T) {
	input := fmt.Sprintf(`Client: %s
né le %s
Email: %s
Téléphone: %s
Adresse: %s
mot de passe : %s
`, fakeName, fakeDOB, fakeEmail, fakePhone, fakeAddress, fakePassword)

	findings := scanAll(t, input)
	found := typesOf(findings)
	mustHave(t, found, "email", "phone", "date_of_birth", "address", "password")

	score := risk.CalculateRisk(findings)
	if score.Status == "SAFE" {
		t.Fatalf("expected non-SAFE status, got %s", score.Status)
	}

	sanitized := sanitizer.Sanitize(input, findings)
	mustNotHaveRaw(t, sanitized, fakeEmail, fakePhone, fakeDOB, fakeAddress, fakePassword)
	mustNotExposeSecrets(t, findings, fakePassword)
	if !strings.Contains(sanitized, "<REDACTED_EMAIL>") {
		t.Error("expected email redaction token")
	}
	if !strings.Contains(sanitized, "<REDACTED_PASSWORD>") {
		t.Error("expected password redaction token")
	}
}

func TestDeveloperSecrets(t *testing.T) {
	input := fmt.Sprintf("OPENAI_API_KEY=%s\nDATABASE_PASSWORD=%s\nGITHUB_TOKEN=%s\n", fakeOpenAIKey, fakeDBPass, fakeGitHubPAT)
	findings := scanAll(t, input)
	found := typesOf(findings)
	mustHave(t, found, "openai_api_key", "github_token")
	if !found["password"] && !found["database_password"] {
		t.Error("expected database password detection")
	}

	sanitized := sanitizer.Sanitize(input, findings)
	mustNotHaveRaw(t, sanitized, fakeOpenAIKey, fakeDBPass, fakeGitHubPAT)
	mustNotExposeSecrets(t, findings, fakeOpenAIKey, fakeDBPass, fakeGitHubPAT)
	if strings.Contains(sanitized, "DATABASE_PASSWORD="+fakeDBPass) {
		t.Error("database password value was not sanitized")
	}
	if risk.CalculateRisk(findings).Status == "SAFE" {
		t.Fatal("developer secrets must not be SAFE")
	}
}

func TestAWSAndJWT(t *testing.T) {
	input := fmt.Sprintf("AWS_ACCESS_KEY_ID=%s\nAWS_SECRET_ACCESS_KEY=%s\ntoken=%s\n", fakeAWSAccess, fakeAWSSecret, fakeJWT)
	findings := scanAll(t, input)
	found := typesOf(findings)
	mustHave(t, found, "aws_access_key", "aws_secret_access_key", "jwt")

	sanitized := sanitizer.Sanitize(input, findings)
	mustNotHaveRaw(t, sanitized, fakeAWSAccess, fakeAWSSecret, fakeJWT)
	if strings.Contains(sanitized, fakeAWSSecret) {
		t.Fatal("AWS secret access key remained in sanitized prompt")
	}
	if !strings.Contains(sanitized, "<REDACTED_AWS_SECRET_KEY>") {
		t.Error("expected AWS secret redaction token")
	}
	mustNotExposeSecrets(t, findings, fakeAWSSecret, fakeJWT)
	if risk.CalculateRisk(findings).Status != "CRITICAL" {
		t.Fatalf("expected CRITICAL, got %s", risk.CalculateRisk(findings).Status)
	}
}

func TestFinancialInformation(t *testing.T) {
	iban := mustIBAN("FR", "30006000011234567890189")
	input := fmt.Sprintf(`Customer: Marie Martin
IBAN: %s
BIC: %s
card: %s
expiry: %s
CVV: %s
`, iban, fakeBIC, fakeCard, fakeExpiry, fakeCVV)

	findings := scanAll(t, input)
	found := typesOf(findings)
	mustHave(t, found, "iban", "bic", "payment_card", "card_expiry", "cvv")

	sanitized := sanitizer.Sanitize(input, findings)
	mustNotHaveRaw(t, sanitized, iban, fakeBIC, fakeCard, fakeExpiry, fakeCVV)
	mustNotExposeSecrets(t, findings, fakeCard, fakeCVV, iban)
	if risk.CalculateRisk(findings).Status == "SAFE" {
		t.Fatal("financial data must not be SAFE")
	}
}

func TestFindingLinesAreAccurate(t *testing.T) {
	input := "safe line\nEmail: jean.dupont@example.com\n"
	findings := scanAll(t, input)
	if len(findings) == 0 {
		t.Fatal("expected email finding")
	}
	for _, f := range findings {
		if f.Type == "email" && f.Line != 2 {
			t.Fatalf("email line=%d want 2", f.Line)
		}
	}
}

func TestDedupePrefersSpecificGitHubToken(t *testing.T) {
	input := "GITHUB_TOKEN=" + fakeGitHubPAT
	findings := scanAll(t, input)
	count := 0
	for _, f := range findings {
		if f.Start < len("GITHUB_TOKEN=") {
			continue
		}
		count++
		if f.Type != "github_token" {
			t.Errorf("overlapping secret should be github_token, got %s", f.Type)
		}
	}
	if count == 0 {
		t.Fatal("expected github token")
	}
}

func mustIBAN(country, bban string) string {
	for i := 0; i < 100; i++ {
		cand := fmt.Sprintf("%s%02d%s", country, i, bban)
		if validIBAN(cand) {
			return cand
		}
	}
	panic("could not construct test IBAN")
}
