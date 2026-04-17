package path

import "testing"

func TestMatcher(t *testing.T) {
	pattern := "/**"
	target := "/hello"
	match := Match(pattern, target)
	if !match {
		t.Errorf("Match fail: pattern: %s,target: %s", pattern, target)
	}

	pattern = "!/**"
	target = "/hello"
	match = Match(pattern, target)
	if match {
		t.Errorf("Match fail: pattern: %s,target: %s", pattern, target)
	}
}

func TestAnyMatcher(t *testing.T) {
	pattern := "/hello"
	pattern2 := "/world"
	target := "/hello"
	match := AnyMatch([]string{pattern, pattern2}, target)
	if !match {
		t.Errorf("Match fail: pattern: %s,target: %s", pattern, target)
	}

	allMatch := AllMatch([]string{pattern, pattern2}, target)
	if allMatch {
		t.Errorf("Match fail: pattern: %s,target: %s", pattern, target)
	}
}
