package ai

import (
	"context"
	"fmt"

	"tech-debt-collector/internal/models"
)

// AnalysisResult contains AI analysis output
type AnalysisResult struct {
	Explanation    string
	Severity       int
	Priority       string
	BusinessImpact string
	FixEstimate    string
	Confidence     float64
	ReasoningSteps []string
	Model          string
	Cost           float64
	Latency        int64
}

// AIBackend defines interface for LLM providers
type AIBackend interface {
	Analyze(ctx context.Context, item *models.DebtItem, context string) (*AnalysisResult, error)
	Name() string
	Cost() float64
	IsAvailable() bool
	MaxTokens() int
}

// BackendRouter routes requests to optimal backend
type BackendRouter struct {
	backends  map[string]AIBackend
	primary   string
	fallbacks []string
	costLimit float64
	speedMode bool
}

// NewBackendRouter creates router
func NewBackendRouter(primary string) *BackendRouter {
	return &BackendRouter{
		backends:  make(map[string]AIBackend),
		primary:   primary,
		fallbacks: []string{},
		costLimit: 0.10,
		speedMode: false,
	}
}

// Register adds backend
func (br *BackendRouter) Register(backend AIBackend) {
	br.backends[backend.Name()] = backend
}

// AddFallback adds fallback
func (br *BackendRouter) AddFallback(name string) {
	br.fallbacks = append(br.fallbacks, name)
}

// Route selects best backend
func (br *BackendRouter) Route(ctx context.Context) AIBackend {
	primary, ok := br.backends[br.primary]
	if ok && primary.IsAvailable() {
		return primary
	}

	for _, name := range br.fallbacks {
		if backend, ok := br.backends[name]; ok && backend.IsAvailable() {
			return backend
		}
	}

	return nil
}

// Analyze routes to best backend
func (br *BackendRouter) Analyze(ctx context.Context, item *models.DebtItem, context string) (*AnalysisResult, error) {
	backend := br.Route(ctx)
	if backend == nil {
		return nil, fmt.Errorf("no AI backend available")
	}

	return backend.Analyze(ctx, item, context)
}
