package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/promptsieve/promptsieve/internal/risk"
	"github.com/promptsieve/promptsieve/internal/sanitizer"
	"github.com/promptsieve/promptsieve/internal/scanner"
	"github.com/promptsieve/promptsieve/pkg/models"
)

const maxScanBytes = 2 << 20

type ScanRequest struct {
	Text string `json:"text"`
}

type ScanResponse struct {
	Status        string           `json:"status"`
	Score         int              `json:"score"`
	Findings      []models.Finding `json:"findings"`
	SanitizedText string           `json:"sanitizedText"`
}

func NewHandler(sc *scanner.Scanner) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/scan", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxScanBytes)
		var req ScanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		findings, err := sc.Scan(context.Background(), req.Text)
		if err != nil {
			log.Printf("scan completed with error findings=%d", len(findings))
			http.Error(w, "scan failed", http.StatusInternalServerError)
			return
		}
		if findings == nil {
			findings = []models.Finding{}
		}

		score := risk.CalculateRisk(findings)
		sanitized := sanitizer.Sanitize(req.Text, findings)
		log.Printf("scan completed findings=%d risk=%s score=%d", len(findings), score.Status, score.Score)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ScanResponse{
			Status:        score.Status,
			Score:         score.Score,
			Findings:      findings,
			SanitizedText: sanitized,
		})
	})

	return mux
}
