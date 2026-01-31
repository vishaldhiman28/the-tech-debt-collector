package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// Scanner scans a repository for source files
type Scanner struct {
	RootPath          string
	ExcludeDirs       map[string]bool
	IncludeExtensions map[string]bool
	SkipHiddenDirs    bool
}

// NewScanner creates a new repository scanner
func NewScanner(rootPath string, excludeDirs, includeExtensions []string, skipHidden bool) *Scanner {
	excludeMap := make(map[string]bool)
	for _, dir := range excludeDirs {
		excludeMap[dir] = true
	}

	extMap := make(map[string]bool)
	for _, ext := range includeExtensions {
		extMap[ext] = true
	}

	// Default extensions if none specified
	if len(extMap) == 0 {
		extMap[".go"] = true
		extMap[".py"] = true
		extMap[".js"] = true
		extMap[".ts"] = true
		extMap[".java"] = true
		extMap[".c"] = true
		extMap[".cpp"] = true
		extMap[".h"] = true
		extMap[".hpp"] = true
		extMap[".rs"] = true
		extMap[".rb"] = true
		extMap[".php"] = true
		extMap[".sh"] = true
	}

	return &Scanner{
		RootPath:          rootPath,
		ExcludeDirs:       excludeMap,
		IncludeExtensions: extMap,
		SkipHiddenDirs:    skipHidden,
	}
}

// ScanFiles recursively scans for source files
func (s *Scanner) ScanFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(s.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip on error
		}

		// Skip hidden files/dirs
		if s.SkipHiddenDirs && strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip excluded directories
		if info.IsDir() {
			if s.ExcludeDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file extension is in include list
		ext := filepath.Ext(path)
		if s.IncludeExtensions[ext] {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// GetFileImportance scores file importance (1-5)
// Higher scores for core files
func (s *Scanner) GetFileImportance(filePath string) int {
	criticalPaths := []string{
		"main.go", "main.py", "index.js",
		"/core/", "/kernel/", "/engine/",
		"/models/", "/handlers/", "/routers/",
		"config", "setup", "init",
	}

	lower := strings.ToLower(filePath)

	for _, critical := range criticalPaths {
		if strings.Contains(lower, critical) {
			return 5
		}
	}

	// Check file depth - deeper = less critical
	depth := strings.Count(filePath, string(os.PathSeparator))
	if depth > 5 {
		return 1
	} else if depth > 3 {
		return 2
	} else if depth > 1 {
		return 3
	}

	return 3 // Default
}
