package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/promptsieve/promptsieve/internal/detectors"
	"github.com/promptsieve/promptsieve/internal/risk"
	"github.com/promptsieve/promptsieve/internal/scanner"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: promptsieve scan <file>")
		os.Exit(4)
	}

	if os.Args[1] != "scan" {
		fmt.Println("Unknown command")
		os.Exit(4)
	}

	fileName := os.Args[2]
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	sc := scanner.NewScanner(scanner.Config{Timeout: 5 * time.Second})
	for _, detector := range detectors.Default() {
		sc.Register(detector)
	}

	findings, err := sc.Scan(context.Background(), string(data))
	if err != nil {
		log.Fatalf("Scan failed: %v", err)
	}

	score := risk.CalculateRisk(findings)

	fmt.Printf("Risk: %s\n", score.Status)
	fmt.Printf("Score: %d\n\n", score.Score)

	if len(findings) > 0 {
		fmt.Printf("Findings:\n")
		counts := make(map[string]int)
		for _, f := range findings {
			counts[f.Type]++
		}
		for k, v := range counts {
			fmt.Printf("%d %s\n", v, k)
		}
	}

	switch score.Status {
	case "CRITICAL":
		os.Exit(3)
	case "SENSITIVE":
		os.Exit(2)
	case "REVIEW":
		os.Exit(1)
	default:
		os.Exit(0)
	}
}
