package model

import "testing"

func TestNextStatus(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus string
		expected      string
	}{
		{
			name:          "Todo to Ready",
			currentStatus: Status.Todo,
			expected:      Status.Ready,
		},
		{
			name:          "Ready to InProgress",
			currentStatus: Status.Ready,
			expected:      Status.InProgress,
		},
		{
			name:          "InProgress to Blocked",
			currentStatus: Status.InProgress,
			expected:      Status.Blocked,
		},
		{
			name:          "Blocked to OnHold",
			currentStatus: Status.Blocked,
			expected:      Status.OnHold,
		},
		{
			name:          "OnHold to Review",
			currentStatus: Status.OnHold,
			expected:      Status.Review,
		},
		{
			name:          "Review to Done",
			currentStatus: Status.Review,
			expected:      Status.Done,
		},
		{
			name:          "Done to Canceled",
			currentStatus: Status.Done,
			expected:      Status.Canceled,
		},
		{
			name:          "Canceled to Waiting",
			currentStatus: Status.Canceled,
			expected:      Status.Waiting,
		},
		{
			name:          "Waiting cycles back to Todo",
			currentStatus: Status.Waiting,
			expected:      Status.Todo,
		},
		{
			name:          "Unknown status defaults to Todo",
			currentStatus: "invalid-status",
			expected:      Status.Todo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NextStatus(tt.currentStatus)
			if result != tt.expected {
				t.Errorf("NextStatus(%q) = %q; want %q", tt.currentStatus, result, tt.expected)
			}
		})
	}
}

func TestPreviousStatus(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus string
		expected      string
	}{
		{
			name:          "Todo cycles back to Waiting",
			currentStatus: Status.Todo,
			expected:      Status.Waiting,
		},
		{
			name:          "Ready to Todo",
			currentStatus: Status.Ready,
			expected:      Status.Todo,
		},
		{
			name:          "InProgress to Ready",
			currentStatus: Status.InProgress,
			expected:      Status.Ready,
		},
		{
			name:          "Blocked to InProgress",
			currentStatus: Status.Blocked,
			expected:      Status.InProgress,
		},
		{
			name:          "OnHold to Blocked",
			currentStatus: Status.OnHold,
			expected:      Status.Blocked,
		},
		{
			name:          "Review to OnHold",
			currentStatus: Status.Review,
			expected:      Status.OnHold,
		},
		{
			name:          "Done to Review",
			currentStatus: Status.Done,
			expected:      Status.Review,
		},
		{
			name:          "Canceled to Done",
			currentStatus: Status.Canceled,
			expected:      Status.Done,
		},
		{
			name:          "Waiting to Canceled",
			currentStatus: Status.Waiting,
			expected:      Status.Canceled,
		},
		{
			name:          "Unknown status defaults to Todo",
			currentStatus: "invalid-status",
			expected:      Status.Todo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PreviousStatus(tt.currentStatus)
			if result != tt.expected {
				t.Errorf("PreviousStatus(%q) = %q; want %q", tt.currentStatus, result, tt.expected)
			}
		})
	}
}
