package gitignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// IgnoreRules represents the compiled .gitignore rules
type IgnoreRules struct {
	patterns []pattern
	repoRoot string
}

// pattern represents a single .gitignore pattern
type pattern struct {
	pattern string
	negate  bool
	dirOnly bool
}

// LoadIgnoreRules reads and parses the .gitignore file from the repository root
func LoadIgnoreRules(repoRoot string) (*IgnoreRules, error) {
	rules := &IgnoreRules{
		patterns: make([]pattern, 0),
		repoRoot: repoRoot,
	}

	// Always ignore .mini-git directory
	rules.patterns = append(rules.patterns, pattern{
		pattern: ".mini-git",
		negate:  false,
		dirOnly: true,
	})

	gitignorePath := filepath.Join(repoRoot, ".gitignore")
	file, err := os.Open(gitignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			// .gitignore doesn't exist, just return the default rules
			return rules, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if p, ok := parsePattern(line); ok {
			rules.patterns = append(rules.patterns, p)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

// parsePattern parses a single line from .gitignore file
func parsePattern(line string) (pattern, bool) {
	// Trim whitespace
	line = strings.TrimSpace(line)

	// Skip empty lines and comments
	if line == "" || strings.HasPrefix(line, "#") {
		return pattern{}, false
	}

	p := pattern{
		negate:  false,
		dirOnly: false,
	}

	// Check for negation pattern
	if strings.HasPrefix(line, "!") {
		p.negate = true
		line = line[1:]
	}

	// Check if pattern applies only to directories
	if strings.HasSuffix(line, "/") {
		p.dirOnly = true
		line = strings.TrimSuffix(line, "/")
	}

	// Handle leading slash (anchored to root)
	if strings.HasPrefix(line, "/") {
		line = line[1:]
	}

	p.pattern = line
	return p, true
}

// IsIgnored checks if a given path should be ignored according to the rules
func (r *IgnoreRules) IsIgnored(path string, isDir bool) bool {
	// Get relative path from repo root
	relPath, err := filepath.Rel(r.repoRoot, path)
	if err != nil {
		// If we can't get relative path, don't ignore it
		return false
	}

	// Normalize path separators to forward slashes for matching
	relPath = filepath.ToSlash(relPath)

	ignored := false

	// Apply patterns in order (later patterns can override earlier ones)
	for _, p := range r.patterns {
		// Skip directory-only patterns for files
		if p.dirOnly && !isDir {
			continue
		}

		if matchPattern(relPath, p.pattern) {
			ignored = !p.negate
		}
	}

	return ignored
}

// ShouldSkipDir checks if a directory should be skipped during traversal
func (r *IgnoreRules) ShouldSkipDir(path string) bool {
	return r.IsIgnored(path, true)
}

// matchPattern performs glob-style pattern matching
func matchPattern(path, pattern string) bool {
	// Handle exact matches
	if path == pattern {
		return true
	}

	// Check if pattern contains directory separators
	if strings.Contains(pattern, "/") {
		// Pattern with slashes matches full path
		matched, _ := filepath.Match(pattern, path)
		if matched {
			return true
		}
		// Also check with ** wildcard support
		return matchWildcard(path, pattern)
	}

	// Pattern without slashes matches basename
	basename := filepath.Base(path)
	matched, _ := filepath.Match(pattern, basename)
	if matched {
		return true
	}

	// Also check if pattern matches any component in the path
	parts := strings.Split(path, "/")
	for _, part := range parts {
		matched, _ := filepath.Match(pattern, part)
		if matched {
			return true
		}
	}

	// Check if pattern matches as a prefix with wildcard
	return matchWildcard(path, pattern)
}

// matchWildcard handles ** wildcards and more complex patterns
func matchWildcard(path, pattern string) bool {
	// Handle ** wildcard (matches any number of directories)
	if strings.Contains(pattern, "**") {
		parts := strings.Split(pattern, "**")
		if len(parts) == 2 {
			prefix := strings.TrimSuffix(parts[0], "/")
			suffix := strings.TrimPrefix(parts[1], "/")

			hasPrefix := prefix == "" || strings.HasPrefix(path, prefix)
			hasSuffix := suffix == "" || strings.HasSuffix(path, suffix)

			if hasPrefix && hasSuffix {
				return true
			}
		}
	}

	// Handle trailing * (matches everything under a directory)
	if strings.HasSuffix(pattern, "/*") {
		dir := strings.TrimSuffix(pattern, "/*")
		if strings.HasPrefix(path, dir+"/") {
			return true
		}
	}

	// Simple prefix match for directories
	if strings.HasPrefix(path, pattern+"/") {
		return true
	}

	return false
}
