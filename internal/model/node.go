package model

import (
	"time"

	"github.com/google/uuid"
)

type Node struct {
	ID           string
	Type         string
	Content      string
	Link         string
	Tags         []string
	Places       []string
	Status       string
	Date         time.Time
	OverrideDate *time.Time
}

func (n *Node) GenID() {
	n.ID = uuid.NewString()
}
