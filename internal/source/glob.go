package source

import (
	"fmt"
	"path/filepath"
	"sort"
)

// ExpandGlobs expands any glob patterns in paths and returns the
// deduplicated, sorted list of matching file paths.
// Paths without wildcards are returned as-is if they are valid patterns.
func ExpandGlobs(patterns []string) ([]string, error) {
	seen := make(map[string]struct{})
	var result []string

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("source: invalid glob pattern %q: %w", pattern, err)
		}
		// filepath.Glob returns nil when there are no matches.
		if len(matches) == 0 {
			// Preserve the original path so the caller gets a meaningful error.
			matches = []string{pattern}
		}
		for _, m := range matches {
			if _, ok := seen[m]; !ok {
				seen[m] = struct{}{}
				result = append(result, m)
			}
		}
	}
	sort.Strings(result)
	return result, nil
}
