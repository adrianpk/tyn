package model

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID               string
	NodeID           string
	NotificationType string
	LastNotifiedAt   time.Time
	TimesNotified    int
}

func (n *Notification) GenID() {
	n.ID = uuid.NewString()
}

// NotificationType constants
var NotificationType = struct {
	DueDate string
}{
	DueDate: "due_date",
}
