package detectors

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"unicode"
)

var ibanLengths = map[string]int{
	"AD": 24, "AE": 23, "AL": 28, "AT": 20, "AZ": 28, "BA": 20, "BE": 16,
	"BG": 22, "BH": 22, "BR": 29, "BY": 28, "CH": 21, "CR": 22, "CY": 28,
	"CZ": 24, "DE": 22, "DK": 18, "DO": 28, "EE": 20, "EG": 29, "ES": 24,
	"FI": 18, "FO": 18, "FR": 27, "GB": 22, "GE": 22, "GI": 23, "GL": 18,
	"GR": 27, "GT": 28, "HR": 21, "HU": 28, "IE": 22, "IL": 23, "IQ": 23,
	"IS": 26, "IT": 27, "JO": 30, "KW": 30, "KZ": 20, "LB": 28, "LC": 32,
	"LI": 21, "LT": 20, "LU": 20, "LV": 21, "LY": 25, "MC": 27, "MD": 24,
	"ME": 22, "MK": 19, "MR": 27, "MT": 31, "MU": 30, "NL": 18, "NO": 15,
	"PK": 24, "PL": 28, "PS": 29, "PT": 25, "QA": 29, "RO": 24, "RS": 22,
	"SA": 24, "SE": 24, "SI": 19, "SK": 24, "SM": 27, "TN": 24, "TR": 26,
	"UA": 29, "VA": 22, "VG": 24, "XK": 20,
}

func normalizeIBAN(value string) string {
	var b strings.Builder
	for _, r := range strings.ToUpper(value) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func validIBAN(value string) bool {
	iban := normalizeIBAN(value)
	if len(iban) < 15 || len(iban) > 34 {
		return false
	}
	country := iban[:2]
	if !unicode.IsLetter(rune(country[0])) || !unicode.IsLetter(rune(country[1])) {
		return false
	}
	if expected, ok := ibanLengths[country]; ok && len(iban) != expected {
		return false
	}
	rearranged := iban[4:] + iban[:4]
	return mod97(rearranged) == 1
}

func mod97(value string) int {
	total := 0
	for _, r := range value {
		var n int
		switch {
		case r >= '0' && r <= '9':
			n = int(r - '0')
			total = (total*10 + n) % 97
		case r >= 'A' && r <= 'Z':
			n = int(r-'A') + 10
			if n >= 10 {
				total = (total*10 + n/10) % 97
				total = (total*10 + n%10) % 97
			}
		default:
			return 0
		}
	}
	return total
}

func validLuhn(digits string) bool {
	if len(digits) < 13 || len(digits) > 19 {
		return false
	}
	sum := 0
	alt := false
	for i := len(digits) - 1; i >= 0; i-- {
		n := int(digits[i] - '0')
		if n < 0 || n > 9 {
			return false
		}
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}

func validBIC(value string) bool {
	bic := strings.ToUpper(strings.TrimSpace(value))
	if len(bic) != 8 && len(bic) != 11 {
		return false
	}
	for i, r := range bic {
		switch {
		case i < 6:
			if r < 'A' || r > 'Z' {
				return false
			}
		default:
			if (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
				return false
			}
		}
	}
	return true
}

func validJWT(token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}
	if len(parts[0]) < 8 || len(parts[1]) < 8 || len(parts[2]) < 10 {
		return false
	}
	header, err := decodeB64URL(parts[0])
	if err != nil {
		return false
	}
	var hdr map[string]any
	if json.Unmarshal(header, &hdr) != nil {
		return false
	}
	if _, ok := hdr["alg"]; !ok {
		return false
	}
	payload, err := decodeB64URL(parts[1])
	if err != nil {
		return false
	}
	var body any
	if json.Unmarshal(payload, &body) != nil {
		return false
	}
	if _, err := decodeB64URL(parts[2]); err != nil {
		for _, r := range parts[2] {
			if !isBase64URLRune(r) {
				return false
			}
		}
	}
	return true
}

func decodeB64URL(value string) ([]byte, error) {
	if m := len(value) % 4; m != 0 {
		value += strings.Repeat("=", 4-m)
	}
	if data, err := base64.URLEncoding.DecodeString(value); err == nil {
		return data, nil
	}
	return base64.RawURLEncoding.DecodeString(strings.TrimRight(value, "="))
}

func isBase64URLRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_'
}

func validFrenchPhone(digits string) bool {
	switch {
	case strings.HasPrefix(digits, "33") && len(digits) == 11:
		return digits[2] >= '1' && digits[2] <= '9'
	case strings.HasPrefix(digits, "0") && len(digits) == 10:
		return digits[1] >= '1' && digits[1] <= '9'
	default:
		return false
	}
}

func validE164(digits string) bool {
	if len(digits) < 8 || len(digits) > 15 {
		return false
	}
	if digits[0] == '0' {
		return false
	}
	return true
}

func isSafeAssignmentValue(value string) bool {
	v := strings.Trim(strings.TrimSpace(value), `"'`)
	if v == "" {
		return true
	}
	switch strings.ToLower(v) {
	case "true", "false", "yes", "no", "on", "off", "null", "none", "undefined",
		"production", "development", "test", "dev", "prod", "staging", "local":
		return true
	}
	if _, err := strconv.Atoi(v); err == nil && len(v) <= 6 {
		return true
	}
	if strings.HasPrefix(v, "/") || strings.HasPrefix(v, "./") || strings.HasPrefix(v, "~") {
		return true
	}
	return false
}
