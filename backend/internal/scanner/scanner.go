package scanner

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/promptsieve/promptsieve/pkg/models"
)

type Config struct {
	Timeout       time.Duration
	MaxConcurrent int
}

type Scanner struct {
	detectors []models.Detector
	config    Config
}

func NewScanner(cfg Config) *Scanner {
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}
	if cfg.MaxConcurrent == 0 {
		cfg.MaxConcurrent = 10
	}
	return &Scanner{
		detectors: make([]models.Detector, 0),
		config:    cfg,
	}
}

func (s *Scanner) Register(d models.Detector) {
	s.detectors = append(s.detectors, d)
}

func (s *Scanner) Detectors() []models.Detector {
	return s.detectors
}

func (s *Scanner) Scan(parentCtx context.Context, input string) ([]models.Finding, error) {
	ctx, cancel := context.WithTimeout(parentCtx, s.config.Timeout)
	defer cancel()

	var wg sync.WaitGroup
	findingsCh := make(chan []models.Finding, len(s.detectors))
	errCh := make(chan error, len(s.detectors))
	sem := make(chan struct{}, s.config.MaxConcurrent)

	for _, d := range s.detectors {
		wg.Add(1)
		go func(det models.Detector) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }()

			res, err := det.Detect(ctx, input)
			if err != nil {
				errCh <- err
				return
			}
			if len(res) > 0 {
				findingsCh <- res
			}
		}(d)
	}

	go func() {
		wg.Wait()
		close(findingsCh)
		close(errCh)
	}()

	var allFindings []models.Finding
	for res := range findingsCh {
		allFindings = append(allFindings, res...)
	}

	select {
	case err := <-errCh:
		if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
			return Deduplicate(assignLines(input, allFindings)), err
		}
	default:
	}

	return Deduplicate(assignLines(input, allFindings)), nil
}

func assignLines(input string, findings []models.Finding) []models.Finding {
	for i := range findings {
		idx := findings[i].Start
		if idx < 0 {
			idx = 0
		}
		if idx > len(input) {
			idx = len(input)
		}
		findings[i].Line = strings.Count(input[:idx], "\n") + 1
	}
	return findings
}
