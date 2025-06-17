package svc

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
			wantType:    "note",
			wantContent: "just a note",
		},
		{
			name:        "with tags and places",
			input:       "task #tag1 #tag2 @home @office",
			wantType:    "note", // no status → still a note
			wantContent: "task",
			wantTags:    []string{"tag1", "tag2"},
			wantPlaces:  []string{"home", "office"},
		},
		{
			name:         "with status and override date",
			input:        "meeting :done ^2025-06-20",
			wantType:     "task",
			wantContent:  "meeting",
			wantStatus:   "done",
			wantOverride: &overrideDate,
		},
		{
			name:        "link detection",
			input:       "check this https://example.com/page",
			wantType:    "link",
			wantContent: "check this",
			wantLink:    "https://example.com/page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if node.Type != tt.wantType {
				t.Errorf("got type %q, want %q", node.Type, tt.wantType)
			}
			if node.Content != tt.wantContent {
				t.Errorf("got content %q, want %q", node.Content, tt.wantContent)
			}
			if !reflect.DeepEqual(node.Tags, tt.wantTags) {
				t.Errorf("got tags %v, want %v", node.Tags, tt.wantTags)
			}
			if !reflect.DeepEqual(node.Places, tt.wantPlaces) {
				t.Errorf("got places %v, want %v", node.Places, tt.wantPlaces)
			}
			if node.Status != tt.wantStatus {
				t.Errorf("got status %q, want %q", node.Status, tt.wantStatus)
			}
			if tt.wantOverride != nil {
				if node.OverrideDate == nil || !node.OverrideDate.Equal(*tt.wantOverride) {
					t.Errorf("got override %v, want %v", node.OverrideDate, tt.wantOverride)
				}
			} else if node.OverrideDate != nil {
				t.Errorf("got override %v, want nil", node.OverrideDate)
			}
			if node.Link != tt.wantLink {
				t.Errorf("got link %q, want %q", node.Link, tt.wantLink)
			}
		})
	}
}
