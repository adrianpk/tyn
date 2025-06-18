package svc

import (
	"testing"
	"time"

	"github.com/adrianpk/tyn/internal/model"
)

func TestParseAndNodeConstruction(t *testing.T) {
	baseTime := time.Now()
	dueDate, _ := time.Parse("2006-01-02", "2025-06-17")
	tests := []struct {
		name    string
		input   string
		want    model.Node
		wantErr bool
	}{
		{
			name:  "simple note",
			input: "this is a simple note",
			want: model.Node{
				Type:    "note",
				Content: "this is a simple note",
				Date:    baseTime,
			},
			wantErr: false,
		},
		{
			name:  "note with tags",
			input: "note with #tag1 #tag2",
			want: model.Node{
				Type:    "note",
				Content: "note with",
				Tags:    []string{"tag1", "tag2"},
				Date:    baseTime,
			},
			wantErr: false,
		},
		{
			name:  "task with status",
			input: "task :todo",
			want: model.Node{
				Type:    "task",
				Content: "task",
				Status:  "todo",
				Date:    baseTime,
			},
			wantErr: false,
		},
		{
			name:  "note with places",
			input: "note from @home @office",
			want: model.Node{
				Type:    "note",
				Content: "note from",
				Places:  []string{"home", "office"},
				Date:    baseTime,
			},
			wantErr: false,
		},
		{
			name:  "link with url",
			input: "check this https://example.com",
			want: model.Node{
				Type:    "link",
				Content: "check this",
				Link:    "https://example.com",
				Date:    baseTime,
			},
			wantErr: false,
		},
		{
			name:  "complete note",
			input: "complete note #tag @place :todo ^2025-06-17 https://example.com",
			want: model.Node{
				Type:    "task",
				Content: "complete note",
				Tags:    []string{"tag"},
				Places:  []string{"place"},
				Status:  "todo",
				Link:    "https://example.com",
				Date:    baseTime,
				DueDate: &dueDate,
			},
			wantErr: false,
		},
		{
			name:    "invalid date",
			input:   "note with ^invalid-date",
			want:    model.Node{Type: "note", Content: "note with ^invalid-date", Date: baseTime},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if node.Date.IsZero() {
				t.Error("Parse() Date is zero")
			}
			node.Date = tt.want.Date
			if node.Type != tt.want.Type {
				t.Errorf("Parse() Type = %v, want %v", node.Type, tt.want.Type)
			}
			if node.Content != tt.want.Content {
				t.Errorf("Parse() Content = %v, want %v", node.Content, tt.want.Content)
			}
			if node.Link != tt.want.Link {
				t.Errorf("Parse() Link = %v, want %v", node.Link, tt.want.Link)
			}
			if !sliceEqual(node.Tags, tt.want.Tags) {
				t.Errorf("Parse() Tags = %v, want %v", node.Tags, tt.want.Tags)
			}
			if !sliceEqual(node.Places, tt.want.Places) {
				t.Errorf("Parse() Places = %v, want %v", node.Places, tt.want.Places)
			}
			if node.Status != tt.want.Status {
				t.Errorf("Parse() Status = %v, want %v", node.Status, tt.want.Status)
			}
			if node.DueDate != nil && tt.want.DueDate != nil {
				d, e := node.DueDate.Local(), tt.want.DueDate.Local()
				if d.Year() != e.Year() || d.Month() != e.Month() || d.Day() != e.Day() {
					t.Errorf("Parse() DueDate = %v, want %v", d, e)
				}
			} else if node.DueDate != nil || tt.want.DueDate != nil {
				t.Errorf("Parse() DueDate = %v, want %v", node.DueDate, tt.want.DueDate)
			}
		})
	}
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func timePointersEqual(a, b *time.Time) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if a == nil {
		return true
	}
	return a.UTC().Equal(b.UTC())
}
