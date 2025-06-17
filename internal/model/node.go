package model

import (
	"strings"
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

func EncodeCSV(s []string) string {
	return strings.Join(s, ",")
}

func DecodeCSV(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func (n *Node) GenID() {
	n.ID = uuid.NewString()
}
