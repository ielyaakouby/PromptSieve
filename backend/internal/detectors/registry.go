package detectors

import "github.com/promptsieve/promptsieve/pkg/models"

func Default() []models.Detector {
	return []models.Detector{
		NewEmailDetector(),
		NewPhoneDetector(),
		NewDOBDetector(),
		NewAddressDetector(),
		NewPersonDetector(),
		NewPasswordDetector(),
		NewOpenAIDetector(),
		NewGitHubDetector(),
		NewAWSDetector(),
		NewJWTDetector(),
		NewPrivateKeyDetector(),
		NewIBANDetector(),
		NewBICDetector(),
		NewCardDetector(),
		NewCVVDetector(),
		NewCardExpiryDetector(),
		NewGenericSecretDetector(),
		NewBearerDetector(),
	}
}
