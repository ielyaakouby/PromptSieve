package detectors

import (
	"strings"
	"testing"

	"github.com/promptsieve/promptsieve/internal/risk"
	"github.com/promptsieve/promptsieve/internal/sanitizer"
)

func TestNegativeFalsePositives(t *testing.T) {
	input := strings.Join([]string{
		"HTTP status code 401",
		"port 8080",
		"Kubernetes version 1.30.2",
		"not an IP 999.999.999.999",
		"The release was shipped on 14/05/1987",
		"random number 404",
		"uuid 550e8400-e29b-41d4-a716-446655440000",
		"sha256 e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"commit 1a2b3c4d5e6f7890aabbccddeeff001122334455",
		"PORT=3000",
		"DEBUG=true",
		"NODE_ENV=production",
		"REPLICAS=3",
	}, "\n")

	findings := scanAll(t, input)
	if len(findings) != 0 {
		for _, f := range findings {
			t.Errorf("unexpected finding type=%s value=%q", f.Type, f.MaskedValue)
		}
	}
	if risk.CalculateRisk(findings).Status != "SAFE" {
		t.Fatalf("expected SAFE, got %s", risk.CalculateRisk(findings).Status)
	}
	if sanitized := sanitizer.Sanitize(input, findings); sanitized != input {
		t.Error("negative input should not be rewritten")
	}
}

func TestCVVRequiresContext(t *testing.T) {
	findings := scanAll(t, "The room number is 404 and the code is unused.")
	for _, f := range findings {
		if f.Type == "cvv" {
			t.Fatal("unlabeled 3-digit number must not be a CVV")
		}
	}
}

func TestAWSSecretRequiresContext(t *testing.T) {
	findings := scanAll(t, "checksum "+fakeAWSSecret)
	for _, f := range findings {
		if f.Type == "aws_secret_access_key" {
			t.Fatal("40-char string without AWS context must not be an AWS secret")
		}
	}
}

func TestOrdinaryDateIsNotDOB(t *testing.T) {
	findings := scanAll(t, "Documentation updated on 14/05/1987 for the public release.")
	for _, f := range findings {
		if f.Type == "date_of_birth" {
			t.Fatal("ordinary documentation date must not be DOB")
		}
	}
}
