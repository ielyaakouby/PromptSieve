package detectors

import (
	"context"
	"strings"

	"github.com/promptsieve/promptsieve/pkg/models"
)

type privateKeyDetector struct{}

func NewPrivateKeyDetector() models.Detector { return privateKeyDetector{} }

func (privateKeyDetector) Name() string { return "private-key" }

func (privateKeyDetector) Detect(ctx context.Context, input string) ([]models.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var findings []models.Finding
	search := 0
	for i := 0; i < maxFindingsPerDetector; i++ {
		begin := strings.Index(input[search:], "-----BEGIN ")
		if begin < 0 {
			break
		}
		begin += search
		headerEnd := strings.Index(input[begin:], "PRIVATE KEY-----")
		if headerEnd < 0 {
			break
		}
		headerEnd += begin + len("PRIVATE KEY-----")
		end := strings.Index(input[headerEnd:], "-----END ")
		if end < 0 {
			break
		}
		end += headerEnd
		footer := strings.Index(input[end:], "PRIVATE KEY-----")
		if footer < 0 {
			search = headerEnd
			continue
		}
		end += footer + len("PRIVATE KEY-----")
		if end-begin > 64*1024 {
			search = headerEnd
			continue
		}
		value, start, stop, ok := validUTF8Slice(input, begin, end)
		if ok {
			findings = append(findings, newFinding(
				"private-key", "cryptographic_material", "private_key", models.SeverityCritical,
				start, stop, value, "<REDACTED_PRIVATE_KEY>",
				"Private cryptographic key material detected.",
				"Never share private keys.",
				0.99,
			))
		}
		search = end
	}
	return findings, nil
}
