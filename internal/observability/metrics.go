package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics
type Metrics struct {
	// Analysis metrics
	ItemsScanned     prometheus.Counter
	ItemsAnalyzed    prometheus.Counter
	AnalysisErrors   prometheus.Counter
	AnalysisDuration prometheus.Histogram

	// LLM metrics
	LLMRequests prometheus.Counter
	LLMErrors   prometheus.Counter
	LLMLatency  prometheus.Histogram
	LLMCost     prometheus.Gauge

	// RAG metrics
	VectorIndexed  prometheus.Counter
	VectorSearches prometheus.Counter
	VectorLatency  prometheus.Histogram

	// Feedback metrics
	FeedbackRecorded prometheus.Counter
	AvgAccuracy      prometheus.Gauge
	FixedRate        prometheus.Gauge
}

// NewMetrics creates metrics registry
func NewMetrics() *Metrics {
	return &Metrics{
		ItemsScanned: promauto.NewCounter(prometheus.CounterOpts{
			Name: "tech_debt_items_scanned_total",
			Help: "Total items scanned",
		}),
		ItemsAnalyzed: promauto.NewCounter(prometheus.CounterOpts{
			Name: "tech_debt_items_analyzed_total",
			Help: "Total items analyzed by LLM",
		}),
		AnalysisErrors: promauto.NewCounter(prometheus.CounterOpts{
			Name: "tech_debt_analysis_errors_total",
			Help: "Total analysis errors",
		}),
		AnalysisDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "tech_debt_analysis_duration_seconds",
			Help:    "Analysis duration in seconds",
			Buckets: prometheus.DefBuckets,
		}),
		LLMRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: "tech_debt_llm_requests_total",
			Help: "Total LLM requests",
		}),
		LLMErrors: promauto.NewCounter(prometheus.CounterOpts{
			Name: "tech_debt_llm_errors_total",
			Help: "Total LLM errors",
		}),
		LLMLatency: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "tech_debt_llm_latency_ms",
			Help:    "LLM latency in milliseconds",
			Buckets: prometheus.DefBuckets,
		}),
		LLMCost: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "tech_debt_llm_cost_usd",
			Help: "Total LLM cost in USD",
		}),
		VectorIndexed: promauto.NewCounter(prometheus.CounterOpts{
			Name: "tech_debt_vectors_indexed_total",
			Help: "Total vectors indexed",
		}),
		VectorSearches: promauto.NewCounter(prometheus.CounterOpts{
			Name: "tech_debt_vector_searches_total",
			Help: "Total vector searches",
		}),
		VectorLatency: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "tech_debt_vector_latency_ms",
			Help:    "Vector search latency",
			Buckets: prometheus.DefBuckets,
		}),
		FeedbackRecorded: promauto.NewCounter(prometheus.CounterOpts{
			Name: "tech_debt_feedback_recorded_total",
			Help: "Total feedback entries",
		}),
		AvgAccuracy: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "tech_debt_accuracy",
			Help: "Average analysis accuracy (0-1)",
		}),
		FixedRate: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "tech_debt_fixed_rate",
			Help: "Rate of fixed items (0-1)",
		}),
	}
}

// RecordItemScanned records scanned item
func (m *Metrics) RecordItemScanned() {
	m.ItemsScanned.Inc()
}

// RecordAnalysis records analysis metrics
func (m *Metrics) RecordAnalysis(durationMs int64, cost float64) {
	m.ItemsAnalyzed.Inc()
	m.LLMRequests.Inc()
	m.LLMLatency.Observe(float64(durationMs))
	m.LLMCost.Add(cost)
}

// RecordError records error
func (m *Metrics) RecordError() {
	m.AnalysisErrors.Inc()
	m.LLMErrors.Inc()
}

// RecordVectorSearch records vector search
func (m *Metrics) RecordVectorSearch(latencyMs int64) {
	m.VectorSearches.Inc()
	m.VectorLatency.Observe(float64(latencyMs))
}

// RecordFeedback records feedback
func (m *Metrics) RecordFeedback(accuracy, fixedRate float64) {
	m.FeedbackRecorded.Inc()
	m.AvgAccuracy.Set(accuracy)
	m.FixedRate.Set(fixedRate)
}
