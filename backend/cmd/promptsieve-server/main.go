package main

import (
	"log"
	"net/http"
	"time"

	"github.com/promptsieve/promptsieve/internal/api"
	"github.com/promptsieve/promptsieve/internal/detectors"
	"github.com/promptsieve/promptsieve/internal/scanner"
)

func main() {
	sc := scanner.NewScanner(scanner.Config{
		Timeout:       5 * time.Second,
		MaxConcurrent: 10,
	})
	for _, detector := range detectors.Default() {
		sc.Register(detector)
	}

	log.Println("Starting PromptSieve server on 0.0.0.0:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", api.NewHandler(sc)); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
