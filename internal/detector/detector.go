package detector

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"tech-debt-collector/internal/models"
)

// Detector detects technical debt items in source files
type Detector struct {
	patterns map[string]*regexp.Regexp
	typeMap  map[string]int // Type to default severity
}

// NewDetector creates a new tech debt detector
func NewDetector() *Detector {
	d := &Detector{
		patterns: make(map[string]*regexp.Regexp),
		typeMap: map[string]int{
			"TODO":       2,
			"FIXME":      3,
			"HACK":       4,
			"DEPRECATED": 3,
			"XXX":        4,
		},
	}

	// Compile regex patterns
	for typeStr := range d.typeMap {
		// Match: TODO, TODO:, TODO: message, # TODO: message, etc.
		pattern := regexp.MustCompile(fmt.Sprintf(`(?i)(%s)[\s:]*(.*)$`, typeStr))
		d.patterns[typeStr] = pattern
	}

	return d
}

// DetectInFile scans a file for technical debt items
func (d *Detector) DetectInFile(filePath string, fileImportance int) ([]models.DebtItem, error) {
	var items []models.DebtItem

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Check against each pattern
		for typeStr, pattern := range d.patterns {
			matches := pattern.FindStringSubmatch(line)
			if len(matches) > 0 {
				message := ""
				if len(matches) > 2 {
					message = strings.TrimSpace(matches[2])
				}

				// Detect severity from message context
				severity := d.detectSeverity(typeStr, message)

				item := models.DebtItem{
					ID:             d.generateID(filePath, lineNumber),
					FilePath:       filePath,
					LineNumber:     lineNumber,
					Type:           typeStr,
					Message:        message,
					Severity:       severity,
					FileImportance: fileImportance,
					DetectedAt:     time.Now(),
				}

				items = append(items, item)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// detectSeverity determines severity level based on type and message
func (d *Detector) detectSeverity(typeStr, message string) int {
	severity := d.typeMap[typeStr]

	// Escalate severity for critical keywords
	criticalKeywords := []string{
		"security", "crash", "memory", "leak", "deadlock",
		"race", "critical", "production", "urgent", "asap",
	}

	lower := strings.ToLower(message)
	for _, keyword := range criticalKeywords {
		if strings.Contains(lower, keyword) {
			severity = 5
			break
		}
	}

	// High severity keywords
	highKeywords := []string{
		"error", "bug", "fix", "broken", "broken", "severe",
	}
	for _, keyword := range highKeywords {
		if strings.Contains(lower, keyword) && severity < 4 {
			severity = 4
			break
		}
	}

	return severity
}

// generateID creates a unique ID for a debt item
func (d *Detector) generateID(filePath string, lineNumber int) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s:%d:%d", filePath, lineNumber, time.Now().Unix())))
	return fmt.Sprintf("%x", hash)[:12]
}

// CalculateFrequency counts similar items in a list
func (d *Detector) CalculateFrequency(items []models.DebtItem, filePath string) map[string]int {
	frequencies := make(map[string]int)

	for _, item := range items {
		if item.FilePath == filePath {
			frequencies[item.Type]++
		}
	}

	// Assign frequency scores (1-5)
	for i := range items {
		if items[i].FilePath == filePath {
			count := frequencies[items[i].Type]
			items[i].Frequency = d.frequencyToScore(count)
		}
	}

	return frequencies
}

// frequencyToScore converts count to score (1-5)
func (d *Detector) frequencyToScore(count int) int {
	switch {
	case count <= 1:
		return 1
	case count <= 2:
		return 2
	case count <= 5:
		return 3
	case count <= 10:
		return 4
	default:
		return 5
	}
}
