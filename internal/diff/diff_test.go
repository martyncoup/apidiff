package diff_test

import (
	"testing"

	"github.com/martyn/apidiff/internal/diff"
	"github.com/martyn/apidiff/internal/model"
	"github.com/martyn/apidiff/internal/parser"
)

func TestCompare(t *testing.T) {
	oldDoc, err := parser.Parse("../../testdata/swagger-v1.yaml")
	if err != nil {
		t.Fatalf("failed to parse old spec: %v", err)
	}

	newDoc, err := parser.Parse("../../testdata/swagger-v2.yaml")
	if err != nil {
		t.Fatalf("failed to parse new spec: %v", err)
	}

	changes := diff.Compare(oldDoc, newDoc)

	if len(changes) == 0 {
		t.Fatal("expected changes, got none")
	}

	summary := diff.Summarize(changes)

	if summary.Added == 0 {
		t.Error("expected at least one added endpoint")
	}
	if summary.Removed == 0 {
		t.Error("expected at least one removed endpoint")
	}
	if summary.Breaking == 0 {
		t.Error("expected at least one breaking change")
	}

	// Verify DELETE /users/{id} is reported as removed
	found := false
	for _, c := range changes {
		if c.Type == model.EndpointRemoved && c.Path == "DELETE /users/{id}" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected DELETE /users/{id} to be reported as removed")
	}

	// Verify User.lastName is reported as removed
	found = false
	for _, c := range changes {
		if c.Type == model.PropertyRemoved && c.Property == "User.lastName" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected User.lastName to be reported as removed")
	}
}
