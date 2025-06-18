package svc

import (
	"testing"
	"time"

	"github.com/adrianpk/tyn/internal/model"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType string
		wantTags []string
	}{
		{
			name:     "note with tag",
			input:    "A simple note #tag1",
			wantType: model.Type.Note,
			wantTags: []string{"tag1"},
		},
		{
			name:     "task with status",
			input:    "A task :todo",
			wantType: model.Type.Task,
		},
		{
			name:     "link",
			input:    "https://example.com",
			wantType: model.Type.Link,
		},
		{
			name:     "multiple tags",
			input:    "Note with #tag1 #tag2",
			wantType: model.Type.Note,
			wantTags: []string{"tag1", "tag2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}
			if got.Type != tt.wantType {
				t.Errorf("Parse() type = %v, want %v", got.Type, tt.wantType)
			}

			for i, tag := range tt.wantTags {
				if i < len(got.Tags) && got.Tags[i] != tag {
					t.Errorf("Parse() tag = %v, want %v", got.Tags[i], tag)
				}
			}
		})
	}
}

func TestParseDueDate(t *testing.T) {
	input := "A task with due date ^2025-06-17 #tag1"
	node, err := Parse(input)
	if err != nil {
		t.Errorf("Parse() error = %v", err)
		return
	}

	expectedDate, _ := time.Parse("2006-01-02", "2025-06-17")

	if node.DueDate == nil {
		t.Error("Expected DueDate to be set, but it was nil")
	} else {
		d, e := node.DueDate.Local(), expectedDate.Local()
		if d.Year() != e.Year() || d.Month() != e.Month() || d.Day() != e.Day() {
			t.Errorf("Parse() DueDate = %v, want %v", d, e)
		}
	}

	if len(node.Tags) != 1 || node.Tags[0] != "tag1" {
		t.Errorf("Parse() Tags = %v, want %v", node.Tags, []string{"tag1"})
	}

	if node.Content != "A task with due date" {
		t.Errorf("Parse() Content = %v, want %v", node.Content, "A task with due date")
	}
}
