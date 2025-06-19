package model

import (
	"time"

	"github.com/google/uuid"
)

type Node struct {
	ID      string
	Type    string
	Content string
	Link    string
	Tags    []string
	Places  []string
	Status  string
	Date    time.Time
	DueDate *time.Time
}

func (n *Node) GenID() {
	n.ID = uuid.NewString()
}

func (n *Node) IsOverdue() bool {
	if n.Type != "task" || n.DueDate == nil {
		return false
	}

	if n.Status == Status.Done || n.Status == Status.Canceled {
		return false
	}

	return time.Now().After(*n.DueDate)
}

type Filter struct {
	Type   string
	Tags   []string
	Places []string
	Status string
}
