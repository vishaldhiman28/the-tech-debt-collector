package models

import "time"

// DebtItem represents a single piece of technical debt detected in the codebase
type DebtItem struct {
	ID                string    `json:"id"`
	FilePath          string    `json:"file_path"`
	LineNumber        int       `json:"line_number"`
	Type              string    `json:"type"` // TODO, FIXME, HACK, DEPRECATED, XXX
	Message           string    `json:"message"`
	Severity          int       `json:"severity"`        // 1-5: low to critical
	FileImportance    int       `json:"file_importance"` // 1-5: low to critical
	Frequency         int       `json:"frequency"`       // How many similar items in file
	Risk              float64   `json:"risk"`            // Computed risk score (0-100)
	DetectedAt        time.Time `json:"detected_at"`
	LLMExplanation    string    `json:"llm_explanation"`
	LLMPriority       string    `json:"llm_priority"` // HIGH, MEDIUM, LOW
	LLMRecommendation string    `json:"llm_recommendation"`
}

// RiskScore holds the risk assessment
type RiskScore struct {
	Item              *DebtItem
	SeverityWeight    float64
	CriticalityWeight float64
	FrequencyWeight   float64
	FinalRisk         float64
}

// Report represents the final analysis report
type Report struct {
	GeneratedAt     time.Time  `json:"generated_at"`
	RepositoryPath  string     `json:"repository_path"`
	TotalItems      int        `json:"total_items"`
	CriticalItems   int        `json:"critical_items"`
	HighItems       int        `json:"high_items"`
	MediumItems     int        `json:"medium_items"`
	LowItems        int        `json:"low_items"`
	DebtItems       []DebtItem `json:"debt_items"`
	Summary         string     `json:"summary"`
	Recommendations []string   `json:"recommendations"`
}

// ScannerConfig holds scanner configuration
type ScannerConfig struct {
	RootPath          string
	ExcludeDirs       []string
	IncludeExtensions []string
	SkipHiddenDirs    bool
}

// Config holds application configuration
type Config struct {
	OpenAIAPIKey  string
	OpenAIModel   string
	ScannerConfig ScannerConfig
	EnableLLM     bool
	OutputFormat  string // json, text, html
	OutputPath    string
	Verbose       bool
}
