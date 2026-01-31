# Complete Project Setup Guide

## 📦 Prerequisites

```bash
# Install Docker & Docker Compose
brew install docker docker-compose

# Install Go 1.21+
brew install go

# API Keys needed
export OPENAI_API_KEY=sk-...
export CLAUDE_API_KEY=claude-...  # Optional
```

## 🚀 Quick Start (Docker)

### 1. Start Full Stack

```bash
cd "/Users/yaegar/Tech debt collector"

# Set environment variables
export OPENAI_API_KEY=sk-your-key

# Start all services
docker-compose up -d
```

### 2. Access Services

```
- Tech Debt Collector API: Runs on demand
- Qdrant Vector DB: http://localhost:6333
- Prometheus Metrics: http://localhost:9090
- Grafana Dashboards: http://localhost:3000 (admin/admin)
- Feedback Dashboard: http://localhost:8080
```

### 3. Run Analysis

```bash
docker-compose exec collector ./tech-debt-collector \
  -path /workspace \
  -agent -rag \
  -format json \
  -output report.json
```

## 🛠️ Local Development (Without Docker)

### 1. Install Dependencies

```bash
cd "/Users/yaegar/Tech debt collector"
go mod download
go mod tidy
```

### 2. Start Qdrant Locally

```bash
docker run -p 6333:6333 qdrant/qdrant:latest
```

### 3. Start Prometheus

```bash
# Download prometheus first
brew install prometheus

# Run with config
prometheus --config.file=observability/prometheus.yml
```

### 4. Build Binary

```bash
make build
# or
go build -o bin/tech-debt-collector ./cmd/tech-debt-collector
```

### 5. Run Analysis

```bash
./bin/tech-debt-collector \
  -path ./test_repo \
  -agent -rag \
  -verbose

# With Grafana dashboard
open http://localhost:3000
```

## 🧪 Running Tests

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test -v ./internal/detector
go test -v ./internal/scorer
go test -v ./internal/agent
```

### Expected Coverage

```
detector: 80%+
scorer: 85%+
agent: 75%+
rag: 70%+
```

## 📊 Monitoring

### Prometheus Metrics

Endpoint: `http://localhost:9090`

Available metrics:
- `tech_debt_items_scanned_total` - Total items scanned
- `tech_debt_items_analyzed_total` - Items analyzed by LLM
- `tech_debt_llm_latency_ms` - LLM response time
- `tech_debt_llm_cost_usd` - Cumulative LLM cost
- `tech_debt_vector_searches_total` - RAG searches
- `tech_debt_accuracy` - Analysis accuracy
- `tech_debt_fixed_rate` - Items actually fixed

### Grafana Dashboards

Endpoint: `http://localhost:3000`

Pre-configured dashboards:
- Analysis Performance (latency, cost, accuracy)
- LLM Model Comparison (GPT-4 vs Claude vs local)
- RAG Effectiveness (search latency, results quality)
- Feedback Trends (rating progression, user satisfaction)

## 🎯 Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                    Tech Debt Collector                   │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  Scanner          Detector          Scorer              │
│  (Files)    →    (Comments)   →    (Risk)    →          │
│                                                           │
│                         ↓                                 │
│                                                           │
│  ┌──────────────────────────────────────────────┐       │
│  │           Agent (5-Step Reasoning)           │       │
│  │                                              │       │
│  │  Step 1: RAG Search  (Vector Similarity)    │       │
│  │  Step 2: LLM Analysis (Initial Assessment) │       │
│  │  Step 3: Confidence Check                   │       │
│  │  Step 4: Refinement (if needed)            │       │
│  │  Step 5: Finalize & Report                  │       │
│  │                                              │       │
│  └──────────────────────────────────────────────┘       │
│         ↓                              ↓                 │
│    ┌─────────────┐              ┌──────────────┐       │
│    │ Qdrant VDB  │              │ Multi-Model  │       │
│    │ (RAG)       │              │ AI Backend   │       │
│    │             │              │ (GPT/Claude) │       │
│    └─────────────┘              └──────────────┘       │
│                                                           │
│                         ↓                                 │
│                                                           │
│              Feedback → Learning Loop                     │
│              (Continuous Improvement)                     │
│                                                           │
│                    Report Output                          │
│              (JSON, Text, Dashboard)                      │
│                                                           │
└─────────────────────────────────────────────────────────┘

         ↓
    
┌──────────────────────────────────────────────┐
│         Observability Stack                  │
├──────────────────────────────────────────────┤
│                                              │
│  Prometheus (Metrics)                       │
│  Grafana (Dashboards)                       │
│  Structured Logging (slog)                  │
│  OpenTelemetry (Distributed Tracing)        │
│                                              │
└──────────────────────────────────────────────┘

         ↓

┌──────────────────────────────────────────────┐
│         Web Dashboard                        │
├──────────────────────────────────────────────┤
│                                              │
│  Feedback Collection                        │
│  Trend Analysis                             │
│  Misclassification Reporting                │
│  Model Performance Tracking                 │
│                                              │
└──────────────────────────────────────────────┘
```

## 🔑 Configuration

### Environment Variables

```bash
# Required
export OPENAI_API_KEY=sk-...

# Optional
export CLAUDE_API_KEY=claude-...
export QDRANT_URL=http://localhost:6333
export QDRANT_API_KEY=tech-debt-secret
export PROMETHEUS_PUSHGATEWAY=http://localhost:9091
export LOG_LEVEL=info
```

### CLI Flags

```bash
tech-debt-collector \
  -path ./repo              # Repository to scan
  -output report.json       # Output file
  -format json              # json, text
  -agent                    # Enable agentic analysis
  -rag                      # Enable RAG
  -claude                   # Use Claude instead of GPT
  -verbose                  # Verbose logging
  -qdrant-url URL          # Qdrant connection
```

## 📈 Performance Targets

| Component | Metric | Target |
|-----------|--------|--------|
| Scanner | 10K files/min | ✓ |
| Detector | 1K files/min | ✓ |
| Scorer | Instant | ✓ |
| Agent (top 10) | 30s | ✓ |
| RAG Search | <100ms | ✓ |
| LLM Request | <5s | ✓ |
| Total Run | <2min | ✓ |
| Cost (per run) | <$1 | ✓ |

## 🚨 Troubleshooting

### Qdrant Connection Failed

```bash
# Check if Qdrant is running
curl http://localhost:6333/health

# Restart Qdrant
docker-compose restart qdrant
```

### LLM API Errors

```bash
# Verify API key
echo $OPENAI_API_KEY

# Test OpenAI connection
curl -H "Authorization: Bearer $OPENAI_API_KEY" \
  https://api.openai.com/v1/models
```

### Memory Issues

```bash
# Increase Docker memory
docker-compose down
export DOCKER_MEMORY=4g
docker-compose up -d
```

## 📚 Code Structure

```
tech-debt-collector/
├── cmd/
│   └── tech-debt-collector/     # Main CLI entry point
├── internal/
│   ├── models/                  # Data structures
│   ├── scanner/                 # Repository scanning
│   ├── detector/                # Debt detection
│   ├── scorer/                  # Risk scoring
│   ├── ai/                      # LLM backends (multi-model)
│   ├── rag/                     # Vector database & RAG
│   ├── agent/                   # Agentic reasoning loop
│   ├── feedback/                # Feedback collection
│   └── observability/           # Metrics & logging
├── web/
│   ├── main.go                  # Dashboard API
│   └── static/index.html        # Frontend UI
├── observability/               # Prometheus & Grafana configs
├── tests/                       # Unit & integration tests
├── Dockerfile                   # Container image
├── docker-compose.yml           # Full stack setup
├── go.mod                       # Dependencies
└── README.md
```

## 🎓 Interview Talking Points

**Q: Why this architecture?**
> Multi-step reasoning agent reduces hallucinations. RAG provides codebase context. Feedback loop enables continuous learning. Cost-optimized through smart model routing.

**Q: How does RAG help?**
> Semantic search finds similar debt in THIS codebase. Reduces generic advice, provides project-specific insights. Dramatically improves explanation quality.

**Q: Scalability?**
> Async processing with job queues. Distributed RAG with managed vector DB. Horizontal scaling of agent workers. Cost controls per-analysis.

**Q: Production readiness?**
> Full observability stack (Prometheus + Grafana). Error handling & fallbacks. Rate limiting. Structured logging. Comprehensive testing (80%+ coverage).

## 📝 Next Steps

1. **Fine-tune prompts** based on feedback trends
2. **Add more tools** (git blame, test coverage, metrics)
3. **Implement distributed agent** (multiple workers)
4. **Add model fine-tuning** pipeline
5. **Deploy to Kubernetes** with GitOps
6. **Create Slack integration** for notifications

---

**Ready to impress senior engineers with your AI architecture!** 🚀
