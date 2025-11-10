package gitignore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParsePattern(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantOk    bool
		wantPat   pattern
	}{
		{
			name:    "simple pattern",
			line:    "*.log",
			wantOk:  true,
			wantPat: pattern{pattern: "*.log", negate: false, dirOnly: false},
		},
		{
			name:    "negation pattern",
			line:    "!important.log",
			wantOk:  true,
			wantPat: pattern{pattern: "important.log", negate: true, dirOnly: false},
		},
		{
			name:    "directory only pattern",
			line:    "node_modules/",
			wantOk:  true,
			wantPat: pattern{pattern: "node_modules", negate: false, dirOnly: true},
		},
		{
			name:    "anchored pattern",
			line:    "/build",
			wantOk:  true,
			wantPat: pattern{pattern: "build", negate: false, dirOnly: false},
		},
		{
			name:   "comment",
			line:   "# this is a comment",
			wantOk: false,
		},
		{
			name:   "empty line",
			line:   "",
			wantOk: false,
		},
		{
			name:   "whitespace only",
			line:   "   ",
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parsePattern(tt.line)
			if ok != tt.wantOk {
				t.Errorf("parsePattern() ok = %v, want %v", ok, tt.wantOk)
				return
			}
			if ok && got != tt.wantPat {
				t.Errorf("parsePattern() = %+v, want %+v", got, tt.wantPat)
			}
		})
	}
}

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		pattern string
		want    bool
	}{
		{
			name:    "exact match",
			path:    "test.log",
			pattern: "test.log",
			want:    true,
		},
		{
			name:    "wildcard extension",
			path:    "test.log",
			pattern: "*.log",
			want:    true,
		},
		{
			name:    "wildcard name",
			path:    "test.log",
			pattern: "test.*",
			want:    true,
		},
		{
			name:    "directory pattern",
			path:    "node_modules/package/file.js",
			pattern: "node_modules",
			want:    true,
		},
		{
			name:    "nested file with wildcard",
			path:    "src/test.log",
			pattern: "*.log",
			want:    true,
		},
		{
			name:    "full path pattern",
			path:    "src/test/file.js",
			pattern: "src/test/*.js",
			want:    true,
		},
		{
			name:    "double wildcard",
			path:    "src/deep/nested/file.js",
			pattern: "src/**/*.js",
			want:    true,
		},
		{
			name:    "no match",
			path:    "test.txt",
			pattern: "*.log",
			want:    false,
		},
		{
			name:    "directory trailing slash",
			path:    "build/output.js",
			pattern: "build/*",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchPattern(tt.path, tt.pattern); got != tt.want {
				t.Errorf("matchPattern(%q, %q) = %v, want %v", tt.path, tt.pattern, got, tt.want)
			}
		})
	}
}

func TestLoadIgnoreRules(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create a test .gitignore file
	gitignoreContent := `# Test .gitignore
*.log
node_modules/
!important.log
/build
src/**/*.tmp
`
	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .gitignore: %v", err)
	}

	rules, err := LoadIgnoreRules(tmpDir)
	if err != nil {
		t.Fatalf("LoadIgnoreRules() error = %v", err)
	}

	if rules == nil {
		t.Fatal("LoadIgnoreRules() returned nil rules")
	}

	// Should have: .mini-git (default) + 5 patterns from file = 6 total
	// (comment and empty line should be skipped)
	expectedPatterns := 6
	if len(rules.patterns) != expectedPatterns {
		t.Errorf("LoadIgnoreRules() loaded %d patterns, want %d", len(rules.patterns), expectedPatterns)
	}
}

func TestLoadIgnoreRulesNoFile(t *testing.T) {
	// Create a temporary directory without .gitignore
	tmpDir := t.TempDir()

	rules, err := LoadIgnoreRules(tmpDir)
	if err != nil {
		t.Fatalf("LoadIgnoreRules() error = %v, want nil", err)
	}

	if rules == nil {
		t.Fatal("LoadIgnoreRules() returned nil rules")
	}

	// Should only have the default .mini-git pattern
	if len(rules.patterns) != 1 {
		t.Errorf("LoadIgnoreRules() loaded %d patterns, want 1 (default)", len(rules.patterns))
	}
}

func TestIsIgnored(t *testing.T) {
	tmpDir := t.TempDir()

	gitignoreContent := `*.log
node_modules/
!important.log
build/
*.tmp
src/cache/
`
	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .gitignore: %v", err)
	}

	rules, err := LoadIgnoreRules(tmpDir)
	if err != nil {
		t.Fatalf("LoadIgnoreRules() error = %v", err)
	}

	tests := []struct {
		name       string
		path       string
		isDir      bool
		wantIgnore bool
	}{
		{
			name:       "ignore .log file",
			path:       filepath.Join(tmpDir, "test.log"),
			isDir:      false,
			wantIgnore: true,
		},
		{
			name:       "don't ignore negated file",
			path:       filepath.Join(tmpDir, "important.log"),
			isDir:      false,
			wantIgnore: false,
		},
		{
			name:       "ignore node_modules directory",
			path:       filepath.Join(tmpDir, "node_modules"),
			isDir:      true,
			wantIgnore: true,
		},
		{
			name:       "ignore build directory",
			path:       filepath.Join(tmpDir, "build"),
			isDir:      true,
			wantIgnore: true,
		},
		{
			name:       "ignore .mini-git (default)",
			path:       filepath.Join(tmpDir, ".mini-git"),
			isDir:      true,
			wantIgnore: true,
		},
		{
			name:       "don't ignore regular file",
			path:       filepath.Join(tmpDir, "src", "main.go"),
			isDir:      false,
			wantIgnore: false,
		},
		{
			name:       "ignore .tmp files",
			path:       filepath.Join(tmpDir, "data.tmp"),
			isDir:      false,
			wantIgnore: true,
		},
		{
			name:       "ignore nested .log file",
			path:       filepath.Join(tmpDir, "src", "debug.log"),
			isDir:      false,
			wantIgnore: true,
		},
		{
			name:       "ignore cache directory",
			path:       filepath.Join(tmpDir, "src", "cache"),
			isDir:      true,
			wantIgnore: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules.IsIgnored(tt.path, tt.isDir)
			if got != tt.wantIgnore {
				t.Errorf("IsIgnored(%q, %v) = %v, want %v", tt.path, tt.isDir, got, tt.wantIgnore)
			}
		})
	}
}

func TestShouldSkipDir(t *testing.T) {
	tmpDir := t.TempDir()

	gitignoreContent := `node_modules/
.cache/
build/
`
	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .gitignore: %v", err)
	}

	rules, err := LoadIgnoreRules(tmpDir)
	if err != nil {
		t.Fatalf("LoadIgnoreRules() error = %v", err)
	}

	tests := []struct {
		name     string
		path     string
		wantSkip bool
	}{
		{
			name:     "skip node_modules",
			path:     filepath.Join(tmpDir, "node_modules"),
			wantSkip: true,
		},
		{
			name:     "skip .cache",
			path:     filepath.Join(tmpDir, ".cache"),
			wantSkip: true,
		},
		{
			name:     "skip .mini-git (default)",
			path:     filepath.Join(tmpDir, ".mini-git"),
			wantSkip: true,
		},
		{
			name:     "don't skip src",
			path:     filepath.Join(tmpDir, "src"),
			wantSkip: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rules.ShouldSkipDir(tt.path)
			if got != tt.wantSkip {
				t.Errorf("ShouldSkipDir(%q) = %v, want %v", tt.path, got, tt.wantSkip)
			}
		})
	}
}
