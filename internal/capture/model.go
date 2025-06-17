package capture

import "time"

type Node struct {
	Type         string
	Content      string
	Link         string
	Tags         []string
	Places       []string
	Status       string
	Date         time.Time
	OverrideDate *time.Time
}

// NewNode creates a new Node from the given input string using the parser
func NewNode(input string) (Node, error) {
	return Parse(input)
}
