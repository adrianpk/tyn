package capture

import (
	"testing"
	"time"
)

func TestNewNode(t *testing.T) {
	baseTime := time.Now()

	overrideDate, _ := time.Parse("2006-01-02", "2025-06-17")

	tests := []struct {
		name    string
		input   string
		want    Node
		wantErr bool
	}{
		{
			name:  "simple note",
			input: "this is a simple note",
			want: Node{
				Type:    "note",
				Content: "this is a simple note",
				Date:    baseTime,
			},
			wantErr: false,
		},
		{
			name:  "note with tags",
			input: "note with #tag1 #tag2",
			want: Node{
				Type:    "note",
				Content: "note with",
				Tags:    []string{"tag1", "tag2"},
				Date:    baseTime,
			},
			wantErr: false,
		},
		{
			name:  "task with status",
			input: "task !todo",
			want: Node{
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
			want: Node{
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
			want: Node{
				Type:    "link",
				Content: "check this",
				Link:    "https://example.com",
				Date:    baseTime,
			},
			wantErr: false,
		},
		{
			name:  "complete note",
			input: "complete note #tag @place !todo ^2025-06-17 https://example.com",
			want: Node{
				Type:         "task",
				Content:      "complete note",
				Tags:         []string{"tag"},
				Places:       []string{"place"},
				Status:       "todo",
				Link:         "https://example.com",
				Date:         baseTime,
				OverrideDate: &overrideDate,
			},
			wantErr: false,
		},
		{
			name:    "invalid date",
			input:   "note with ^invalid-date",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNode(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// Since Date is set to time.Now(), we only check if it's not zero
			if got.Date.IsZero() {
				t.Error("NewNode() Date is zero")
			}

			// Replace the actual Date with the expected one for comparison
			got.Date = tt.want.Date

			if got.Type != tt.want.Type {
				t.Errorf("NewNode() Type = %v, want %v", got.Type, tt.want.Type)
			}
			if got.Content != tt.want.Content {
				t.Errorf("NewNode() Content = %v, want %v", got.Content, tt.want.Content)
			}
			if got.Link != tt.want.Link {
				t.Errorf("NewNode() Link = %v, want %v", got.Link, tt.want.Link)
			}
			if !sliceEqual(got.Tags, tt.want.Tags) {
				t.Errorf("NewNode() Tags = %v, want %v", got.Tags, tt.want.Tags)
			}
			if !sliceEqual(got.Places, tt.want.Places) {
				t.Errorf("NewNode() Places = %v, want %v", got.Places, tt.want.Places)
			}
			if got.Status != tt.want.Status {
				t.Errorf("NewNode() Status = %v, want %v", got.Status, tt.want.Status)
			}
			if !timePointersEqual(got.OverrideDate, tt.want.OverrideDate) {
				t.Errorf("NewNode() OverrideDate = %v, want %v", got.OverrideDate, tt.want.OverrideDate)
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
	return a.Equal(*b)
}
