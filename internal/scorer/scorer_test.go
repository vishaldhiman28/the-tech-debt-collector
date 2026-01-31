package scorer

import (
	"testing"

	"tech-debt-collector/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestScorerCalculatesRisk(t *testing.T) {
	s := NewScorer()

	item := &models.DebtItem{
		Severity:       5,
		FileImportance: 5,
		Frequency:      1,
	}

	risk := s.ScoreItem(item)
	assert.Greater(t, risk, 0.0)
	assert.LessOrEqual(t, risk, 100.0)
}

func TestScorerHighRisk(t *testing.T) {
	s := NewScorer()

	highRisk := &models.DebtItem{
		Severity:       5,
		FileImportance: 5,
		Frequency:      5,
	}

	lowRisk := &models.DebtItem{
		Severity:       1,
		FileImportance: 1,
		Frequency:      1,
	}

	highScore := s.ScoreItem(highRisk)
	lowScore := s.ScoreItem(lowRisk)

	assert.Greater(t, highScore, lowScore)
}

func TestScorerCategorizes(t *testing.T) {
	s := NewScorer()

	tests := []struct {
		score    float64
		expected string
	}{
		{85.0, "HIGH"},
		{60.0, "MEDIUM"},
		{30.0, "LOW"},
	}

	for _, tt := range tests {
		result := s.CategorizeRisk(tt.score)
		assert.Equal(t, tt.expected, result)
	}
}

func TestScorerSorts(t *testing.T) {
	s := NewScorer()

	items := []models.DebtItem{
		{Severity: 1, FileImportance: 1, Frequency: 1},
		{Severity: 5, FileImportance: 5, Frequency: 5},
		{Severity: 3, FileImportance: 3, Frequency: 3},
	}

	for i := range items {
		items[i].Risk = s.ScoreItem(&items[i])
	}

	sorted := s.SortByRisk(items)

	// Should be descending
	for i := 0; i < len(sorted)-1; i++ {
		assert.GreaterOrEqual(t, sorted[i].Risk, sorted[i+1].Risk)
	}
}

func TestScorerStats(t *testing.T) {
	s := NewScorer()

	items := []models.DebtItem{
		{Severity: 5, FileImportance: 5, Frequency: 5}, // High
		{Severity: 5, FileImportance: 5, Frequency: 4}, // High
		{Severity: 3, FileImportance: 3, Frequency: 3}, // Medium
		{Severity: 1, FileImportance: 1, Frequency: 1}, // Low
	}

	for i := range items {
		items[i].Risk = s.ScoreItem(&items[i])
	}

	critical, high, medium, low := s.GetStats(items)

	assert.Greater(t, high, 0)
	assert.Greater(t, medium, 0)
	assert.Greater(t, low, 0)
}
