package agent

import (
	"testing"

	"tech-debt-collector/internal/ai"
	"tech-debt-collector/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestAgentCreation(t *testing.T) {
	router := ai.NewBackendRouter("test")
	agent := NewAgent(router, nil)

	assert.NotNil(t, agent)
	assert.Equal(t, 5, agent.maxSteps)
}

func TestAgentStepStructure(t *testing.T) {
	step := AgentStep{
		StepNum:   1,
		Action:    "test_action",
		Reasoning: "test reason",
		Result:    "test result",
	}

	assert.Equal(t, 1, step.StepNum)
	assert.NotEmpty(t, step.Action)
}

func TestAgentAnalysisStructure(t *testing.T) {
	item := &models.DebtItem{
		FilePath:   "test.go",
		LineNumber: 10,
		Type:       "TODO",
		Comment:    "test comment",
	}

	analysis := &AgentAnalysis{
		DebtItem: item,
		Steps:    []AgentStep{},
	}

	assert.Equal(t, item, analysis.DebtItem)
	assert.Empty(t, analysis.Steps)
}

func TestGetFinalAnalysis(t *testing.T) {
	initial := &ai.AnalysisResult{
		Explanation: "initial",
		Severity:    3,
	}

	final := &ai.AnalysisResult{
		Explanation: "final",
		Severity:    5,
	}

	analysis := &AgentAnalysis{
		InitialAnalysis: initial,
		FinalAnalysis:   final,
	}

	result := analysis.GetFinalAnalysis()
	assert.Equal(t, final, result)
	assert.Equal(t, "final", result.Explanation)
}
