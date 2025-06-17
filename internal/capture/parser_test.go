// internal/capture/parser_test.go
package capture

import (
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	overrideDate := time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		input        string
		wantType     string
		wantContent  string
		wantTags     []string
		wantPlaces   []string
		wantStatus   string
		wantOverride *time.Time
		wantLink     string
	}{
		{
			name:        "note only",
			input:       "just a note",
			wantType:    Type.Note,
			wantContent: "just a note",
		},
		{
			name:        "with tags and places",
			input:       "task #tag1 #tag2 @home @office",
			wantType:    Type.Note, // no !status → still a note
			wantContent: "task",
			wantTags:    []string{"tag1", "tag2"},
			wantPlaces:  []string{"home", "office"},
		},
		{
			name:         "with status and override date",
			input:        "meeting !done ^2025-06-20",
			wantType:     Type.Task,
			wantContent:  "meeting",
			wantStatus:   "done",
			wantOverride: &overrideDate,
		},
		{
			name:        "link detection",
			input:       "check this https://example.com/page",
			wantType:    Type.Link,
			wantContent: "check this",
			wantLink:    "https://example.com/page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}
			if got.Type != tt.wantType {
				t.Errorf("Type = %q; want %q", got.Type, tt.wantType)
			}
			if got.Content != tt.wantContent {
				t.Errorf("Content = %q; want %q", got.Content, tt.wantContent)
			}
			if !reflect.DeepEqual(got.Tags, tt.wantTags) {
				t.Errorf("Tags = %v; want %v", got.Tags, tt.wantTags)
			}
			if !reflect.DeepEqual(got.Places, tt.wantPlaces) {
				t.Errorf("Places = %v; want %v", got.Places, tt.wantPlaces)
			}
			if got.Status != tt.wantStatus {
				t.Errorf("Status = %q; want %q", got.Status, tt.wantStatus)
			}
			if tt.wantOverride != nil {
				if got.OverrideDate == nil || !got.OverrideDate.Equal(*tt.wantOverride) {
					t.Errorf("OverrideDate = %v; want %v", got.OverrideDate, tt.wantOverride)
				}
			} else if got.OverrideDate != nil {
				t.Errorf("OverrideDate = %v; want nil", got.OverrideDate)
			}
			if got.Link != tt.wantLink {
				t.Errorf("Link = %q; want %q", got.Link, tt.wantLink)
			}
		})
	}
}
