package capture

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

func NewNode(input string) (Node, error) {
	node, err := Parse(input)
	if err != nil {
		return node, err
	}
	node.ID = uuid.NewString()

	return node, nil
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
