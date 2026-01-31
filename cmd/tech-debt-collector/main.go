package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"tech-debt-collector/internal/detector"
	"tech-debt-collector/internal/llm"
	"tech-debt-collector/internal/models"
	"tech-debt-collector/internal/scanner"
	"tech-debt-collector/internal/scorer"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// Define flags
	repoPath := flag.String("path", ".", "Repository path to scan")
	outputPath := flag.String("output", "report.json", "Output file path")
	outputFormat := flag.String("format", "json", "Output format: json or text")
	enableLLM := flag.Bool("llm", true, "Enable LLM enrichment")
	verbose := flag.Bool("verbose", false, "Verbose output")
	openAIKey := flag.String("openai-key", os.Getenv("OPENAI_API_KEY"), "OpenAI API key")
	openAIModel := flag.String("openai-model", "gpt-3.5-turbo", "OpenAI model")
	help := flag.Bool("help", false, "Show help")

	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	err := runAnalysis(*repoPath, *outputPath, *outputFormat, *enableLLM, *verbose, *openAIKey, *openAIModel)

	if err != nil {
		log.Fatalf("❌ Error: %v", err)
	}
}

func runAnalysis(repoPath, outputPath, outputFormat string, enableLLM, verbose bool, openAIKey, openAIModel string) error {
	log.Println("🔍 Tech Debt Collector Analysis Starting...")

	// Validate inputs
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return fmt.Errorf("repository path does not exist: %s", repoPath)
	}

	// Step 1: Scan for files
	log.Printf("📁 Scanning repository: %s\n", repoPath)
	s := scanner.NewScanner(
		repoPath,
		[]string{".git", "node_modules", "vendor", "build", ".venv", ".next"},
		[]string{},
		true,
	)

	files, err := s.ScanFiles()
	if err != nil {
		return fmt.Errorf("scan error: %w", err)
	}
	log.Printf("   Found %d source files\n", len(files))

	// Step 2: Detect technical debt
	log.Println("🔎 Detecting technical debt items...")

	det := detector.NewDetector()
	var allItems []models.DebtItem

	for _, filePath := range files {
		items, err := det.DetectInFile(filePath, s.GetFileImportance(filePath))
		if err != nil {
			if verbose {
				log.Printf("   Warning: Could not scan %s: %v\n", filePath, err)
			}
			continue
		}
		allItems = append(allItems, items...)
	}
	log.Printf("   Found %d debt items\n", len(allItems))

	// Step 3: Calculate frequency
	det.CalculateFrequency(allItems, repoPath)

	// Step 4: Score items
	log.Println("📊 Scoring risk...")
	sc := scorer.NewScorer()
	allItems = sc.ScoreAll(allItems)
	allItems = sc.SortByRisk(allItems)

	critical, high, medium, low := sc.GetStats(allItems)
	log.Printf("   Risk Distribution: Critical:%d, High:%d, Medium:%d, Low:%d\n",
		critical, high, medium, low)

	// Step 5: Enrich with LLM (optional)
	if enableLLM && openAIKey != "" {
		log.Println("🤖 Enriching with LLM analysis...")
		client := llm.NewOpenAIClient(openAIKey, openAIModel, verbose)
		ctx := context.Background()

		// Enrich top 10 items
		limit := 10
		if len(allItems) < limit {
			limit = len(allItems)
		}

		for i := 0; i < limit; i++ {
			if err := client.EnrichItem(ctx, &allItems[i]); err != nil {
				if verbose {
					log.Printf("   Warning: Could not enrich item %d: %v\n", i, err)
				}
			} else if verbose {
				log.Printf("   ✓ Enriched: %s\n", allItems[i].Type)
			}
			// Rate limit
			time.Sleep(100 * time.Millisecond)
		}

		// Generate report summary
		report := createReport(allItems, repoPath, critical, high, medium, low)
		if err := client.EnrichReport(ctx, &report); err != nil {
			if verbose {
				log.Printf("   Warning: Could not enrich report: %v\n", err)
			}
		} else {
			log.Println("   ✓ Generated LLM summary and recommendations")
		}
	}

	// Step 6: Output results
	log.Printf("💾 Writing report to: %s\n", outputPath)
	report := createReport(allItems, repoPath, critical, high, medium, low)

	if err := writeReport(&report, outputPath, outputFormat); err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	// Print summary
	printSummary(&report, outputPath)

	return nil
}

// createReport builds the final report
func createReport(items []models.DebtItem, repoPath string, critical, high, medium, low int) models.Report {
	return models.Report{
		GeneratedAt:    time.Now(),
		RepositoryPath: repoPath,
		TotalItems:     len(items),
		CriticalItems:  critical,
		HighItems:      high,
		MediumItems:    medium,
		LowItems:       low,
		DebtItems:      items,
		Summary:        "Technical Debt Analysis Report",
	}
}

// writeReport saves the report in specified format
func writeReport(report *models.Report, filePath, format string) error {
	switch format {
	case "json":
		return writeJSON(report, filePath)
	case "text":
		return writeText(report, filePath)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// writeJSON writes report as JSON
func writeJSON(report *models.Report, filePath string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// writeText writes report as human-readable text
func writeText(report *models.Report, filePath string) error {
	var content string

	content += "═══════════════════════════════════════════════════════════════\n"
	content += "                  TECH DEBT ANALYSIS REPORT\n"
	content += "═══════════════════════════════════════════════════════════════\n\n"

	content += fmt.Sprintf("Repository: %s\n", report.RepositoryPath)
	content += fmt.Sprintf("Generated: %s\n\n", report.GeneratedAt.Format("2006-01-02 15:04:05"))

	content += fmt.Sprintf("SUMMARY:\n")
	content += fmt.Sprintf("  Total Items: %d\n", report.TotalItems)
	content += fmt.Sprintf("  Critical: %d | High: %d | Medium: %d | Low: %d\n\n",
		report.CriticalItems, report.HighItems, report.MediumItems, report.LowItems)

	if report.Summary != "" && report.Summary != "Technical Debt Analysis Report" {
		content += fmt.Sprintf("LLM ANALYSIS:\n%s\n\n", report.Summary)
	}

	content += "TOP 10 ITEMS:\n"
	content += "─────────────────────────────────────────────────────────────\n"

	for i, item := range report.DebtItems {
		if i >= 10 {
			break
		}
		content += fmt.Sprintf("\n%d. [%s] %s:%d\n", i+1, item.Type, item.FilePath, item.LineNumber)
		content += fmt.Sprintf("   Message: %s\n", item.Message)
		content += fmt.Sprintf("   Risk: %.1f/100 | Severity: %d/5\n", item.Risk, item.Severity)
		if item.LLMExplanation != "" {
			content += fmt.Sprintf("   Analysis: %s\n", item.LLMExplanation)
		}
	}

	if len(report.Recommendations) > 0 {
		content += "\n\nRECOMMENDATIONS:\n"
		content += "─────────────────────────────────────────────────────────────\n"
		for i, rec := range report.Recommendations {
			content += fmt.Sprintf("%d. %s\n", i+1, rec)
		}
	}

	content += "\n═══════════════════════════════════════════════════════════════\n"

	return os.WriteFile(filePath, []byte(content), 0644)
}

// printSummary prints summary to stdout
func printSummary(report *models.Report, outputPath string) {
	fmt.Println("\n✅ Analysis Complete!")
	fmt.Printf("   Total Items: %d\n", report.TotalItems)
	fmt.Printf("   Critical: %d\n", report.CriticalItems)
	fmt.Printf("   High: %d\n", report.HighItems)
	fmt.Printf("   Medium: %d\n", report.MediumItems)
	fmt.Printf("   Low: %d\n", report.LowItems)
	fmt.Printf("\n   Report saved to: %s\n", outputPath)
}

// printHelp displays help information
func printHelp() {
	fmt.Print(`
Tech Debt Collector - Scan and analyze technical debt in your codebase

USAGE:
  tech-debt-collector [flags]

FLAGS:
  -path string              Repository path to scan (default ".")
  -output string            Output file path (default "report.json")
  -format string            Output format: json or text (default "json")
  -llm                      Enable LLM enrichment (default true)
  -verbose                  Verbose output (default false)
  -openai-key string        OpenAI API key (from OPENAI_API_KEY env var)
  -openai-model string      OpenAI model to use (default "gpt-3.5-turbo")
  -help                     Show this help message

EXAMPLES:
  # Scan current directory
  tech-debt-collector

  # Scan specific path with LLM
  tech-debt-collector -path /path/to/repo -llm

  # Generate text report
  tech-debt-collector -format text -output report.txt

  # Verbose analysis with GPT-4
  tech-debt-collector -llm -verbose -openai-model gpt-4

ENVIRONMENT:
  OPENAI_API_KEY            Your OpenAI API key (optional)

For more information, visit: https://github.com/vishaldhiman28/tech-debt-collector
`)
}
