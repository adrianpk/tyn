package model

import (
	"time"

	"github.com/google/uuid"
)

var NodeType = struct {
	Note  string
	Task  string
	Link  string
	Draft string
}{
	Note:  "note",
	Task:  "task",
	Link:  "link",
	Draft: "draft",
}

type Node struct {
	ID      string
	Type    string
	Content string
	Link    string
	Tags    []string
	Places  []string
	Status  string
	Draft   string
	Date    time.Time
	DueDate *time.Time
}

func (n *Node) GenID() {
	n.ID = uuid.NewString()
}

func (n *Node) IsOverdue() bool {
	if n.Type != NodeType.Task || n.DueDate == nil {
		return false
	}

	if n.Status == Status.Done || n.Status == Status.Canceled {
		return false
	}

	return time.Now().After(*n.DueDate)
}

func (n *Node) ShortID() string {
	if len(n.ID) < 4 {
		return n.ID
	}
	return n.ID[0:4]
}

type Filter struct {
	Type   string
	Tags   []string
	Places []string
	Status string
}
