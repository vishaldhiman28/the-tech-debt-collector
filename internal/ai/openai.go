package ai

import (
	"context"
	"fmt"
	"time"

	"tech-debt-collector/internal/models"

	"github.com/sashabaranov/go-openai"
)

type OpenAIBackend struct {
	client    *openai.Client
	model     string
	maxTokens int
	costPer1K float64
}

// NewOpenAIBackend creates OpenAI backend
func NewOpenAIBackend(apiKey, model string) *OpenAIBackend {
	return &OpenAIBackend{
		client:    openai.NewClient(apiKey),
		model:     model,
		maxTokens: 4096,
		costPer1K: 0.03,
	}
}

func (ob *OpenAIBackend) Name() string {
	return fmt.Sprintf("openai-%s", ob.model)
}

func (ob *OpenAIBackend) Cost() float64 {
	return ob.costPer1K
}

func (ob *OpenAIBackend) IsAvailable() bool {
	return ob.client != nil
}

func (ob *OpenAIBackend) MaxTokens() int {
	return ob.maxTokens
}

// Analyze performs analysis
func (ob *OpenAIBackend) Analyze(ctx context.Context, item *models.DebtItem, context string) (*AnalysisResult, error) {
	start := time.Now()

	prompt := buildAnalysisPrompt(item, context)

	resp, err := ob.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: ob.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt(),
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.7,
		MaxTokens:   2000,
	})

	if err != nil {
		return nil, fmt.Errorf("openai error: %w", err)
	}

	result := &AnalysisResult{
		Explanation: resp.Choices[0].Message.Content,
		Model:       ob.Name(),
		Latency:     time.Since(start).Milliseconds(),
		Cost:        (float64(resp.Usage.TotalTokens) / 1000) * ob.costPer1K,
		Severity:    item.Severity,
		Priority:    "MEDIUM",
	}

	return result, nil
}

func systemPrompt() string {
	return `You are an expert software architect analyzing technical debt.
Be concise and actionable. Provide analysis in this format:
EXPLANATION: [why it's risky]
SEVERITY: [1-5]
PRIORITY: [HIGH/MEDIUM/LOW]
IMPACT: [business consequences]
FIX: [time estimate]`
}

func buildAnalysisPrompt(item *models.DebtItem, context string) string {
	return fmt.Sprintf(`Analyze this technical debt:

File: %s (Line %d)
Type: %s
Comment: %s
Risk: %.1f/100

Context:
%s

Provide analysis with explanation, severity, priority, business impact, and fix estimate.`,
		item.FilePath,
		item.LineNumber,
		item.Type,
		item.Comment,
		item.Risk,
		context,
	)
}
