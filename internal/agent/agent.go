package agent

import (
	"context"
	"fmt"
	"log"
	"time"

	"tech-debt-collector/internal/ai"
	"tech-debt-collector/internal/models"
	"tech-debt-collector/internal/rag"
)

// Agent performs agentic reasoning
type Agent struct {
	router      *ai.BackendRouter
	vectorStore *rag.VectorStore
	maxSteps    int
}

// AgentStep represents one reasoning step
type AgentStep struct {
	StepNum   int
	Action    string
	Reasoning string
	Result    string
	Timestamp time.Time
}

// AgentAnalysis holds complete analysis
type AgentAnalysis struct {
	DebtItem        *models.DebtItem
	Steps           []AgentStep
	InitialAnalysis *ai.AnalysisResult
	FinalAnalysis   *ai.AnalysisResult
	FinalConfidence float64
}

// NewAgent creates agent
func NewAgent(router *ai.BackendRouter, vectorStore *rag.VectorStore) *Agent {
	return &Agent{
		router:      router,
		vectorStore: vectorStore,
		maxSteps:    5,
	}
}

// Analyze performs agentic analysis
func (a *Agent) Analyze(ctx context.Context, item *models.DebtItem) (*AgentAnalysis, error) {
	log.Printf("[Agent] Analyzing %s:%d\n", item.FilePath, item.LineNumber)

	analysis := &AgentAnalysis{
		DebtItem: item,
		Steps:    []AgentStep{},
	}

	// Step 1: RAG search
	log.Println("[Agent] Step 1: Searching for similar patterns...")
	var contextStr string
	if a.vectorStore != nil {
		similar, err := a.vectorStore.SearchSimilar(ctx, item.Comment, 3)
		if err != nil {
			log.Printf("[Agent] RAG search failed (continuing): %v\n", err)
		} else {
			contextStr = buildContextFromSimilar(similar)
		}
	}

	analysis.Steps = append(analysis.Steps, AgentStep{
		StepNum:   1,
		Action:    "rag_search",
		Reasoning: "Find similar debt patterns",
		Result:    fmt.Sprintf("Context prepared: %d chars", len(contextStr)),
		Timestamp: time.Now(),
	})

	// Step 2: Initial LLM analysis
	log.Println("[Agent] Step 2: Initial LLM analysis...")
	initialAnalysis, err := a.router.Analyze(ctx, item, contextStr)
	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	analysis.InitialAnalysis = initialAnalysis
	analysis.Steps = append(analysis.Steps, AgentStep{
		StepNum:   2,
		Action:    "llm_analysis",
		Reasoning: "Generate initial assessment",
		Result:    fmt.Sprintf("Model: %s, Latency: %dms", initialAnalysis.Model, initialAnalysis.Latency),
		Timestamp: time.Now(),
	})

	// Step 3: Confidence check
	log.Printf("[Agent] Step 3: Confidence check (%.2f)...\n", initialAnalysis.Confidence)
	if initialAnalysis.Confidence < 0.8 {
		log.Println("[Agent] Low confidence - refining...")
		analysis.Steps = append(analysis.Steps, AgentStep{
			StepNum:   3,
			Action:    "confidence_check",
			Reasoning: "Confidence < 0.8, needs refinement",
			Result:    "Proceeding to refinement",
			Timestamp: time.Now(),
		})
	} else {
		analysis.Steps = append(analysis.Steps, AgentStep{
			StepNum:   3,
			Action:    "confidence_check",
			Reasoning: fmt.Sprintf("Confidence: %.2f", initialAnalysis.Confidence),
			Result:    "Confidence acceptable, proceeding",
			Timestamp: time.Now(),
		})
	}

	// Step 4: Final analysis
	log.Println("[Agent] Step 4: Finalizing...")
	finalAnalysis, err := a.router.Analyze(ctx, item, contextStr+"\\n\\nProvide final refined assessment.")
	if err == nil {
		analysis.FinalAnalysis = finalAnalysis
	} else {
		analysis.FinalAnalysis = initialAnalysis
	}

	analysis.Steps = append(analysis.Steps, AgentStep{
		StepNum:   4,
		Action:    "finalize",
		Reasoning: "Generate final assessment",
		Result:    fmt.Sprintf("Final severity: %d, confidence: %.2f", analysis.FinalAnalysis.Severity, analysis.FinalAnalysis.Confidence),
		Timestamp: time.Now(),
	})

	analysis.FinalConfidence = analysis.FinalAnalysis.Confidence
	return analysis, nil
}

func buildContextFromSimilar(items []*rag.SimilarItem) string {
	if len(items) == 0 {
		return ""
	}

	context := "Similar debt items in codebase:\\n"
	for i, item := range items {
		context += fmt.Sprintf("%d. [%s:%d] Risk: %.1f\\n", i+1, item.File, item.Line, item.Risk)
	}
	return context
}

// GetFinalAnalysis returns best analysis
func (aa *AgentAnalysis) GetFinalAnalysis() *ai.AnalysisResult {
	if aa.FinalAnalysis != nil {
		return aa.FinalAnalysis
	}
	return aa.InitialAnalysis
}
