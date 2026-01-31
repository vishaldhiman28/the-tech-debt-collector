#!/bin/bash
# Quick Start Script for Tech Debt Collector

set -e

echo "🚀 Tech Debt Collector - Quick Setup"
echo "===================================="
echo ""

# Check prerequisites
echo "📋 Checking prerequisites..."
command -v go >/dev/null 2>&1 || { echo "❌ Go not installed"; exit 1; }
command -v docker >/dev/null 2>&1 || { echo "❌ Docker not installed"; exit 1; }
echo "✅ Prerequisites OK"
echo ""

# Verify API key
echo "🔑 Checking OpenAI API key..."
if [ -z "$OPENAI_API_KEY" ]; then
    echo "❌ OPENAI_API_KEY not set"
    echo "Run: export OPENAI_API_KEY=sk-..."
    exit 1
fi
echo "✅ API key found"
echo ""

# Build binary
echo "🔨 Building binary..."
mkdir -p bin
go build -o bin/tech-debt-collector ./cmd/tech-debt-collector
echo "✅ Binary built: bin/tech-debt-collector"
echo ""

# Option 1: Local development
echo "🎯 Choose setup mode:"
echo "1. Local development (requires manual service startup)"
echo "2. Docker Compose (all services included)"
echo ""
read -p "Select [1-2]: " choice

if [ "$choice" = "1" ]; then
    echo ""
    echo "📝 Local Development Setup:"
    echo ""
    echo "1️⃣  Start Qdrant vector database:"
    echo "   docker run -p 6333:6333 qdrant/qdrant:latest"
    echo ""
    echo "2️⃣  Start Prometheus:"
    echo "   prometheus --config.file=observability/prometheus.yml"
    echo ""
    echo "3️⃣  Run tests:"
    echo "   go test ./..."
    echo ""
    echo "4️⃣  Run analysis:"
    echo "   ./bin/tech-debt-collector -path ./test_repo -agent -rag -verbose"
    echo ""
    
elif [ "$choice" = "2" ]; then
    echo ""
    echo "🐳 Starting Docker Compose stack..."
    docker-compose up -d
    
    echo ""
    echo "⏳ Waiting for services to start..."
    sleep 10
    
    echo "✅ Services started!"
    echo ""
    echo "📊 Access points:"
    echo "   • Feedback Dashboard: http://localhost:8080"
    echo "   • Prometheus:         http://localhost:9090"
    echo "   • Grafana:            http://localhost:3000 (admin/admin)"
    echo "   • Qdrant API:         http://localhost:6333"
    echo ""
    echo "▶️  Run analysis inside container:"
    echo "   docker-compose exec collector ./tech-debt-collector -path /workspace -agent -rag"
    echo ""
    
else
    echo "❌ Invalid choice"
    exit 1
fi

echo ""
echo "📚 Next steps:"
echo "   1. Read SETUP.md for detailed instructions"
echo "   2. Check PROJECT_SUMMARY.md for architecture overview"
echo "   3. Run tests: go test -v ./..."
echo ""
echo "✨ Ready to analyze technical debt!"
