# Tech Debt Collector

AI-powered CLI to find and prioritize technical debt.

## Quick Start

- Build: `go build -o bin/tech-debt-collector ./cmd/tech-debt-collector`
- Run: `./bin/tech-debt-collector -path .`
- With LLM: `export OPENAI_API_KEY=sk-... && ./bin/tech-debt-collector -path . -llm`

## Key Features

- Scan for TODO/FIXME/HACK and other debt markers
- Risk scoring and JSON/text reports
- Optional AI analysis and web dashboard

## Requirements

- Go 1.21+
- OpenAI API key (optional, only for LLM features)

## License

MIT

