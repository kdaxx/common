package path

import (
	"github.com/bmatcuk/doublestar/v4"
	"strings"
)

// Match checks if the target matches the ant like pattern.
// and "!" supports.
func Match(pattern string, target string) bool {
	negate := strings.HasPrefix(pattern, "!")
	if negate {
		pattern = pattern[1:]
	}

	ok := doublestar.MatchUnvalidated(pattern, target)
	if negate {
		return !ok
	}
	return ok
}

// AnyMatch checks if the target matches any ant pattern.
func AnyMatch(patterns []string, target string) bool {
	for _, pattern := range patterns {
		if Match(pattern, target) {
			return true
		}
	}
	return false
}

// AllMatch checks if the target matches all ant pattern.
func AllMatch(patterns []string, target string) bool {
	for _, pattern := range patterns {
		if !Match(pattern, target) {
			return false
		}
	}
	return true
}
