package llm

import (
	"context"
	"fmt"
	"strings"

	"tech-debt-collector/internal/models"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient wraps OpenAI API calls
type OpenAIClient struct {
	client    *openai.Client
	modelName string
	verbose   bool
}

// NewOpenAIClient creates a new OpenAI LLM client
func NewOpenAIClient(apiKey, modelName string, verbose bool) *OpenAIClient {
	return &OpenAIClient{
		client:    openai.NewClient(apiKey),
		modelName: modelName,
		verbose:   verbose,
	}
}

// EnrichItem uses LLM to explain and prioritize a debt item
func (c *OpenAIClient) EnrichItem(ctx context.Context, item *models.DebtItem) error {
	prompt := c.buildPrompt(item)

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.modelName,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an expert code reviewer analyzing technical debt. Be concise and actionable.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.3,
		MaxTokens:   300,
	})

	if err != nil {
		return fmt.Errorf("openai error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return fmt.Errorf("no response from OpenAI")
	}

	// Parse response
	response := resp.Choices[0].Message.Content
	c.parseResponse(item, response)

	return nil
}

// EnrichReport uses LLM to generate summary and recommendations
func (c *OpenAIClient) EnrichReport(ctx context.Context, report *models.Report) error {
	prompt := c.buildReportPrompt(report)

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.modelName,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a technical lead creating a strategic tech debt remediation plan. Be specific and prioritized.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.3,
		MaxTokens:   500,
	})

	if err != nil {
		return fmt.Errorf("openai error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return fmt.Errorf("no response from OpenAI")
	}

	response := resp.Choices[0].Message.Content
	report.Summary = response
	report.Recommendations = c.extractRecommendations(response)

	return nil
}

// buildPrompt creates a prompt for analyzing a single debt item
func (c *OpenAIClient) buildPrompt(item *models.DebtItem) string {
	return fmt.Sprintf(`
	Analyze this technical debt item and provide:
	1. Brief explanation of why it's risky (1-2 sentences)
    2. Priority level (HIGH, MEDIUM, or LOW)
    3. Recommended fix (1-2 sentences)

	File: %s (Line %d)
Type: %s
Message: %s
Severity Score: %d/5
File Importance: %d/5
Frequency in File: %d

Format your response as:
EXPLANATION: [your explanation]
PRIORITY: [HIGH/MEDIUM/LOW]
RECOMMENDATION: [your recommendation]
`,
		item.FilePath, item.LineNumber, item.Type, item.Message,
		item.Severity, item.FileImportance, item.Frequency,
	)
}

// buildReportPrompt creates a prompt for analyzing the overall report
func (c *OpenAIClient) buildReportPrompt(report *models.Report) string {
	topItems := ""
	for i, item := range report.DebtItems {
		if i >= 5 { // Top 5 items
			break
		}
		topItems += fmt.Sprintf("- %s (%s, Risk: %.1f): %s\n",
			item.Type, item.FilePath, item.Risk, item.Message)
	}

	return fmt.Sprintf(`
Create a strategic technical debt remediation plan.

Repository: %s
Total Items: %d
Critical: %d | High: %d | Medium: %d | Low: %d

Top Issues:
%s

Provide:
1. Executive summary (2-3 sentences)
2. Top 3 action items with effort estimates
3. Risk if not addressed
4. Quick wins and long-term improvements

Format clearly with headers.
`,
		report.RepositoryPath, report.TotalItems,
		report.CriticalItems, report.HighItems, report.MediumItems, report.LowItems,
		topItems,
	)
}

// parseResponse extracts structured data from LLM response
func (c *OpenAIClient) parseResponse(item *models.DebtItem, response string) {
	lines := strings.Split(response, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "EXPLANATION:") {
			item.LLMExplanation = strings.TrimSpace(strings.TrimPrefix(line, "EXPLANATION:"))
		} else if strings.HasPrefix(line, "PRIORITY:") {
			priority := strings.TrimSpace(strings.TrimPrefix(line, "PRIORITY:"))
			item.LLMPriority = strings.ToUpper(priority)
			if !isValidPriority(item.LLMPriority) {
				item.LLMPriority = "MEDIUM"
			}
		} else if strings.HasPrefix(line, "RECOMMENDATION:") {
			item.LLMRecommendation = strings.TrimSpace(strings.TrimPrefix(line, "RECOMMENDATION:"))
		}
	}

	// Fallback if parsing didn't work
	if item.LLMExplanation == "" {
		item.LLMExplanation = response[:min(len(response), 200)]
	}
	if item.LLMPriority == "" {
		item.LLMPriority = "MEDIUM"
	}
}

// extractRecommendations extracts action items from response
func (c *OpenAIClient) extractRecommendations(response string) []string {
	var recommendations []string
	lines := strings.Split(response, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "-") || strings.HasPrefix(trimmed, "•") {
			rec := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(trimmed, "-"), "•"))
			if len(rec) > 10 {
				recommendations = append(recommendations, rec)
			}
		}
	}

	return recommendations
}

// isValidPriority checks if priority is valid
func isValidPriority(p string) bool {
	return p == "HIGH" || p == "MEDIUM" || p == "LOW"
}

// min returns minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
