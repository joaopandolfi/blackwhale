package models

import (
	"strings"
	"testing"
)

func TestPermissionUpsert(t *testing.T) {
	p1 := []string{"a", "b"}
	p2 := []string{"b", "c"}
	p3 := UpsertPermissions(p1, p2)

	parsed := GeneratePermissions(p3...)
	if !strings.Contains(parsed, "a") {
		t.Errorf("expected a inside got: %s", parsed)
	}

	if !strings.Contains(parsed, "b") {
		t.Errorf("expected b inside got: %s", parsed)
	}

	if !strings.Contains(parsed, "c") {
		t.Errorf("expected c inside got: %s", parsed)
	}
}
