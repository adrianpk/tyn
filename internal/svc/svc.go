package svc

import (
	"context"

	"github.com/adrianpk/tyn/internal/model"
)

type ParserFunc func(input string) (model.Node, error)

type Svc struct {
	Repo   Repo
	Parser ParserFunc
}

func New(repo Repo) *Svc {
	return &Svc{
		Repo:   repo,
		Parser: Parse,
	}
}

// Capture parses the input text and stores the resulting node
func (s *Svc) Capture(text string) (model.Node, error) {
	// Parse the input text
	node, err := s.Parser(text)
	if err != nil {
		return model.Node{}, err
	}

	node.GenID()

	err = s.Repo.Create(context.Background(), node)
	if err != nil {
		return model.Node{}, err
	}

	return node, nil
}

// List retrieves nodes from the repository based on the provided filter
func (s *Svc) List(filter model.Filter) ([]model.Node, error) {
	nodes, err := s.Repo.List(context.Background())
	if err != nil {
		return nil, err
	}

	if isEmptyFilter(filter) {
		return nodes, nil
	}

	var filteredNodes []model.Node
	for _, node := range nodes {
		if matches(node, filter) {
			filteredNodes = append(filteredNodes, node)
		}
	}

	return filteredNodes, nil
}

func isEmptyFilter(filter model.Filter) bool {
	return filter.Type == "" && filter.Status == "" && len(filter.Tags) == 0 && len(filter.Places) == 0
}

func matches(node model.Node, filter model.Filter) bool {
	if filter.Type != "" && node.Type != filter.Type {
		return false
	}

	if filter.Status != "" && node.Status != filter.Status {
		return false
	}

	if len(filter.Tags) > 0 {
		tagMatch := false
		for _, tag := range filter.Tags {
			if hasTag(node, tag) {
				tagMatch = true
				break
			}
		}
		if !tagMatch {
			return false
		}
	}

	if len(filter.Places) > 0 {
		placeMatch := false
		for _, place := range filter.Places {
			if hasPlace(node, place) {
				placeMatch = true
				break
			}
		}
		if !placeMatch {
			return false
		}
	}

	return true
}

func hasTag(node model.Node, tag string) bool {
	for _, t := range node.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func hasPlace(node model.Node, place string) bool {
	for _, p := range node.Places {
		if p == place {
			return true
		}
	}
	return false
}
