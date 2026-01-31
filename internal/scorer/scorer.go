package scorer

import (
	"tech-debt-collector/internal/models"
)

// Scorer calculates risk scores for debt items
type Scorer struct {
	SeverityWeight    float64
	CriticalityWeight float64
	FrequencyWeight   float64
}

// NewScorer creates a new risk scorer with default weights
func NewScorer() *Scorer {
	return &Scorer{
		SeverityWeight:    0.5,  // 50% from severity
		CriticalityWeight: 0.35, // 35% from file criticality
		FrequencyWeight:   0.15, // 15% from frequency
	}
}

// ScoreItem calculates risk score for a single item
func (s *Scorer) ScoreItem(item *models.DebtItem) float64 {
	// Normalize scores to 0-1 range
	severityScore := float64(item.Severity) / 5.0
	criticalityScore := float64(item.FileImportance) / 5.0
	frequencyScore := float64(item.Frequency) / 5.0

	// Calculate weighted risk
	risk := (severityScore * s.SeverityWeight) +
		(criticalityScore * s.CriticalityWeight) +
		(frequencyScore * s.FrequencyWeight)

	// Scale to 0-100
	return risk * 100
}

// ScoreAll calculates risk scores for all items
func (s *Scorer) ScoreAll(items []models.DebtItem) []models.DebtItem {
	for i := range items {
		items[i].Risk = s.ScoreItem(&items[i])
	}
	return items
}

// CategorizeRisk assigns priority level based on risk score
func (s *Scorer) CategorizeRisk(score float64) string {
	switch {
	case score >= 75:
		return "HIGH"
	case score >= 50:
		return "MEDIUM"
	default:
		return "LOW"
	}
}

// SortByRisk sorts items by risk score (descending)
func (s *Scorer) SortByRisk(items []models.DebtItem) []models.DebtItem {
	// Simple bubble sort for small datasets (can optimize with quicksort)
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].Risk > items[i].Risk {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	return items
}

// GetStats returns statistical summary
func (s *Scorer) GetStats(items []models.DebtItem) (critical, high, medium, low int) {
	for _, item := range items {
		priority := s.CategorizeRisk(item.Risk)
		switch priority {
		case "HIGH":
			high++
		case "MEDIUM":
			medium++
		case "LOW":
			low++
		}

		if item.Risk >= 80 {
			critical++
		}
	}
	return
}
