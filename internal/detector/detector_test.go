package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"tech-debt-collector/internal/models"
)

func TestDetectorFindsTODO(t *testing.T) {
	d := NewDetector()
	
	tests := []struct {
		input    string
		expected string
	}{
		{"// TODO: fix this", "TODO"},
		{"// FIXME: bug here", "FIXME"},
		{"# HACK: temporary solution", "HACK"},
		{"// XXX: security issue", "XXX"},
	}

	for _, tt := range tests {
		// Write test file
		tmpFile := "/tmp/test_" + tt.expected + ".go"
		content := tt.input
		err := writeTestFile(tmpFile, content)
		assert.NoError(t, err)

		// Scan file
		items, err := d.Detect(tmpFile)
		assert.NoError(t, err)
		assert.Greater(t, len(items), 0)
		assert.Equal(t, tt.expected, items[0].Type)
	}
}

func TestDetectorExtracts Message(t *testing.T) {
	d := NewDetector()
	tmpFile := "/tmp/test_message.go"
	content := "// TODO: implement authentication logic"
	
	writeTestFile(tmpFile, content)
	items, _ := d.Detect(tmpFile)

	assert.NotEmpty(t, items[0].Comment)
	assert.Contains(t, items[0].Comment, "authentication")
}

func TestDetectorSeverityDetection(t *testing.T) {
	d := NewDetector()
	
	tests := []struct {
		comment  string
		minSev   int
	}{
		{"TODO: nice to have", 1},
		{"FIXME: critical security issue", 5},
		{"HACK: memory leak", 4},
	}

	for _, tt := range tests {
		tmpFile := "/tmp/test_sev.go"
		content := "// " + tt.comment
		writeTestFile(tmpFile, content)
		
		items, _ := d.Detect(tmpFile)
		assert.GreaterOrEqual(t, items[0].Severity, tt.minSev)
	}
}

func TestDetectorFrequency(t *testing.T) {
	d := NewDetector()
	tmpFile := "/tmp/test_freq.go"
	content := `// TODO: fix 1
	// TODO: fix 2
	// TODO: fix 3
	// FIXME: bug`
	
	writeTestFile(tmpFile, content)
	items, _ := d.Detect(tmpFile)

	d.CalculateFrequency(items, tmpFile)

	// Check frequency assigned
	todoCount := 0
	for _, item := range items {
		if item.Type == "TODO" {
			todoCount++
			assert.Greater(t, item.Frequency, 0)
		}
	}
	assert.Equal(t, 3, todoCount)
}

func writeTestFile(path, content string) error {
	return writeToFile(path, content)
}
