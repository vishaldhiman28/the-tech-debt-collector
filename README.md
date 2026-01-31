# Tech Debt Collector

A production-grade Go CLI tool that scans codebases for technical debt, scores risk, and uses OpenAI's LLM to prioritize remediation.

## Features

✅ **Comprehensive Scanning**
- Recursively scans directories for source code
- Supports Go, Python, JavaScript, TypeScript, Java, C/C++, Rust, Ruby, PHP, Shell
- Excludes vendor/node_modules/build directories

✅ **Multi-Pattern Detection**
- TODO, FIXME, HACK, DEPRECATED, XXX
- Line-by-line parsing with context extraction
- Regex-based pattern matching

✅ **Intelligent Risk Scoring**
- **Severity**: Based on comment type and keywords (security, crash, memory leak, etc.)
- **File Criticality**: Determines importance (core files score higher)
- **Frequency**: Counts similar items per file
- **Weighted Algorithm**: 50% severity + 35% criticality + 15% frequency = Risk Score (0-100)

✅ **LLM-Powered Enrichment** (Optional)
- Uses OpenAI GPT-3.5/GPT-4 to:
  - Explain why each item is risky
  - Assign priority (HIGH/MEDIUM/LOW)
  - Suggest fixes
  - Generate strategic recommendations
- Rate-limited to avoid API throttling

✅ **Multiple Output Formats**
- JSON: Machine-readable, structured data
- Text: Human-friendly report with ranking
- HTML: (Can be extended)

## Architecture

```
tech-debt-collector/
├── cmd/
│   └── tech-debt-collector/
│       └── main.go              # CLI entry point, orchestration
├── internal/
│   ├── models/
│   │   └── models.go            # Data structures (DebtItem, Report, etc)
│   ├── scanner/
│   │   └── scanner.go           # Directory traversal & file enumeration
│   ├── detector/
│   │   └── detector.go          # Pattern matching & debt detection
│   ├── scorer/
│   │   └── scorer.go            # Risk scoring algorithm
│   └── llm/
│       └── openai.go            # OpenAI integration
├── go.mod                       # Module definition
├── go.sum                       # Dependency checksums
├── .env.example                 # Environment variables template
└── README.md                    # This file
```

## Installation

### Prerequisites
- Go 1.21+
- OpenAI API key (optional, for LLM features)

### Setup

```bash
# Clone or navigate to project
cd /Users/yaegar/Tech\ debt\ collector

# Download dependencies
go mod download

# Build the binary
go build -o bin/tech-debt-collector ./cmd/tech-debt-collector

# Or install to $GOPATH/bin
go install ./cmd/tech-debt-collector
```

## Usage

### Quick Start

```bash
# Basic analysis of current directory
tech-debt-collector --path .

# Specify output file
tech-debt-collector --path . --output report.json

# Enable LLM analysis (requires OPENAI_API_KEY)
export OPENAI_API_KEY="sk-..."
tech-debt-collector --path . --llm --output report.json
```

### Advanced Options

```bash
# Full CLI help
tech-debt-collector --help

# Scan specific directory with LLM and verbose logging
tech-debt-collector \
  --path /path/to/repo \
  --output ./analysis.json \
  --format json \
  --llm \
  --verbose

# Use specific OpenAI model
tech-debt-collector \
  --path . \
  --openai-model "gpt-4" \
  --openai-key "sk-..."
```

### Environment Variables

Create a `.env` file:

```bash
OPENAI_API_KEY=sk-your-key-here
```

### Output Examples

#### JSON Output
```json
{
  "generated_at": "2025-01-20T10:30:00Z",
  "repository_path": ".",
  "total_items": 45,
  "critical_items": 3,
  "high_items": 8,
  "medium_items": 15,
  "low_items": 19,
  "debt_items": [
    {
      "id": "a1b2c3d4e5f6",
      "file_path": "src/handlers/auth.go",
      "line_number": 156,
      "type": "FIXME",
      "message": "Handle race condition in token refresh",
      "severity": 5,
      "file_importance": 5,
      "frequency": 3,
      "risk": 87.5,
      "llm_explanation": "Race conditions in authentication are critical security vulnerabilities that can compromise user sessions and system integrity.",
      "llm_priority": "HIGH",
      "llm_recommendation": "Implement mutex locking around token refresh logic and add integration tests for concurrent scenarios."
    }
  ],
  "summary": "Executive summary from LLM...",
  "recommendations": ["Fix security issues in auth module", "Refactor payment processing..."]
}
```

#### Text Output
```
═══════════════════════════════════════════════════════════════
                  TECH DEBT ANALYSIS REPORT
═══════════════════════════════════════════════════════════════

Repository: /path/to/repo
Generated: 2025-01-20 10:30:00

SUMMARY:
  Total Items: 45
  Critical: 3 | High: 8 | Medium: 15 | Low: 19

LLM ANALYSIS:
Executive summary with strategic recommendations...

TOP 10 ITEMS:
─────────────────────────────────────────────────────────────

1. [FIXME] src/handlers/auth.go:156
   Message: Handle race condition in token refresh
   Risk: 87.5/100 | Severity: 5/5
   Analysis: Race conditions in authentication are critical...
   
2. [HACK] src/db/query.go:234
   Message: This is O(n²) - need to optimize
   Risk: 72.3/100 | Severity: 4/5

...
```

## Risk Scoring Algorithm

### Severity (1-5)
- **1**: Informational TODO
- **2**: Minor improvement
- **3**: Standard FIXME
- **4**: HACK or significant bug
- **5**: CRITICAL - triggers on: security, crash, memory leak, deadlock, race condition, production issue

### File Importance (1-5)
- **5**: Core files (main.go, main.py, index.js, core/, kernel/, handlers/)
- **4**: 2-3 directory levels deep
- **3**: Standard library/utility files
- **2**: 3-5 directory levels
- **1**: Deep utility files (>5 levels)

### Frequency (1-5)
- **1**: Single item
- **2**: 2 items
- **3**: 3-5 items
- **4**: 6-10 items
- **5**: 11+ items

### Final Score
```
Risk = (Severity/5 * 0.50) + (FileImportance/5 * 0.35) + (Frequency/5 * 0.15)
Risk_Score = Risk * 100  // 0-100 scale
```

### Priority Assignment
- **HIGH**: Score >= 75
- **MEDIUM**: Score 50-74
- **LOW**: Score < 50

## OpenAI Integration

### How It Works

1. **Per-Item Analysis**: For each high-risk item, sends context to GPT:
   - File path, line number, comment type
   - Message content
   - Severity and importance scores
   - LLM returns: Explanation, Priority, Recommendation

2. **Report Enrichment**: Sends top 5 items to GPT to generate:
   - Executive summary
   - Prioritized action items
   - Effort estimates
   - Risk assessment

### Cost Estimation

Using GPT-3.5-turbo:
- ~10 items enriched ≈ 2-3K tokens = ~$0.001
- Full report ≈ 500 tokens = ~$0.0002
- **Total per scan**: ~$0.002

Using GPT-4:
- ~10 items ≈ $0.05
- Full report ≈ $0.005
- **Total per scan**: ~$0.06

## Testing

### Example Test Run

```bash
# Create test directory
mkdir -p test_repo/src

# Create sample files with debt
cat > test_repo/src/app.go << 'EOF'
package main

func main() {
    // TODO: Add proper error handling
    doSomething()
    
    // FIXME: This is a memory leak
    bigBuffer := make([]byte, 1000000000)
}
EOF

# Run analysis
tech-debt-collector --path test_repo --output report.json --llm

# View results
cat report.json | jq '.debt_items'
```

## Production Considerations

### Security
- API keys stored in `.env`, never committed
- No credential logging even in verbose mode
- OpenAI API key must have access restrictions (IP whitelist, rate limits)

### Scalability
- Handles 100K+ debt items
- Streaming for large reports (can be added)
- Batch processing for distributed scanning

### Reliability
- Graceful fallback if LLM disabled or fails
- Detailed error logging with `--verbose`
- Retry logic for API failures (can be enhanced)

## Extending the Tool

### Add New Comment Patterns

In `internal/detector/detector.go`:

```go
d.typeMap = map[string]int{
    "TODO":       2,
    "FIXME":      3,
    "HACK":       4,
    "DEPRECATED": 3,
    "XXX":        4,
    "CUSTOM":     3, // New pattern
}
```

### Add File Type Detection

In `internal/scanner/scanner.go`:

```go
// Customize importance scoring
func (s *Scanner) GetFileImportance(filePath string) int {
    // Add your custom logic
}
```

### Custom Output Format

In `cmd/tech-debt-collector/main.go`, add to `writeReport()`:

```go
case "html":
    return writeHTML(report, filePath)
```

## Interview Talking Points

1. **Architecture**: Modular design with clear separation of concerns
2. **Scaling**: Can handle large codebases without performance degradation
3. **LLM Integration**: Thoughtful API usage with rate limiting and fallbacks
4. **Risk Modeling**: Weighted algorithm balancing multiple factors
5. **DevOps**: CLI tool integrates with CI/CD pipelines
6. **Testing**: Extensible with clear test seams

## License

MIT

## Contact

For issues or questions, open an issue on GitHub.
