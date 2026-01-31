package feedback

import (
	"context"
	"sync"
	"time"
)

// UserFeedback represents user evaluation
type UserFeedback struct {
	ID                  string
	ItemID              string
	FilePath            string
	LineNumber          int
	ExplanationRating   int    // 1-5
	SeverityAccuracy    int    // 1-5
	PriorityAccuracy    int    // 1-5
	UserComment         string
	ActualSeverity      int
	ActualPriority      string
	IsFixed             bool
	TimeToFix           string
	CreatedAt           time.Time
}

// FeedbackStorage interface
type FeedbackStorage interface {
	Save(ctx context.Context, fb *UserFeedback) error
	GetByItemID(ctx context.Context, itemID string) ([]*UserFeedback, error)
	GetTrends(ctx context.Context, days int) (*Trends, error)
	GetLowRated(ctx context.Context, threshold float64) ([]*UserFeedback, error)
}

// Trends holds feedback analysis
type Trends struct {
	AvgExplanation     float64
	AvgSeverity        float64
	AvgPriority        float64
	FixedRate          float64
	TotalItems         int
	LowRatedCount      int
}

// Collector manages feedback
type Collector struct {
	storage FeedbackStorage
}

// NewCollector creates collector
func NewCollector(storage FeedbackStorage) *Collector {
	return &Collector{storage: storage}
}

// Record saves feedback
func (c *Collector) Record(ctx context.Context, fb *UserFeedback) error {
	fb.CreatedAt = time.Now()
	return c.storage.Save(ctx, fb)
}

// Analyze gets trends
func (c *Collector) Analyze(ctx context.Context, days int) (*Trends, error) {
	return c.storage.GetTrends(ctx, days)
}

// GetMisclassified finds low-rated items
func (c *Collector) GetMisclassified(ctx context.Context) ([]*UserFeedback, error) {
	return c.storage.GetLowRated(ctx, 3.0)
}

// MemoryStorage implements in-memory storage
type MemoryStorage struct {
	mu mu        sync.RWMutex
	items map[string]*UserFeedback
}

// NewMemoryStorage creates storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		items: make(map[string]*UserFeedback),
	}
}

func (ms *MemoryStorage) Save(ctx context.Context, fb *UserFeedback) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.items[fb.ID] = fb
	return nil
}

func (ms *MemoryStorage) GetByItemID(ctx context.Context, itemID string) ([]*UserFeedback, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var results []*UserFeedback
	for _, fb := range ms.items {
		if fb.ItemID == itemID {
			results = append(results, fb)
		}
	}
	return results, nil
}

func (ms *MemoryStorage) GetTrends(ctx context.Context, days int) (*Trends, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	cutoff := time.Now().AddDate(0, 0, -days)
	var recent []*UserFeedback

	for _, fb := range ms.items {
		if fb.CreatedAt.After(cutoff) {
			recent = append(recent, fb)
		}
	}

	if len(recent) == 0 {
		return &Trends{}, nil
	}

	var totalExp, totalSev, totalPri int
	fixedCount := 0

	for _, fb := range recent {
		totalExp += fb.ExplanationRating
		totalSev += fb.SeverityAccuracy
		totalPri += fb.PriorityAccuracy
		if fb.IsFixed {
			fixedCount++
		}
	}

	return &Trends{
		AvgExplanation: float64(totalExp) / float64(len(recent)),
		AvgSeverity:    float64(totalSev) / float64(len(recent)),
		AvgPriority:    float64(totalPri) / float64(len(recent)),
		FixedRate:      float64(fixedCount) / float64(len(recent)),
		TotalItems:     len(recent),
	}, nil
}

func (ms *MemoryStorage) GetLowRated(ctx context.Context, threshold float64) ([]*UserFeedback, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var results []*UserFeedback
	for _, fb := range ms.items {
		avg := float64(fb.ExplanationRating+fb.SeverityAccuracy+fb.PriorityAccuracy) / 3.0
		if avg < threshold {
			results = append(results, fb)
		}
	}
	return results, nil
}
