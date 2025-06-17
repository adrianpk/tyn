package svc

import "github.com/adrianpk/tyn/internal/model"

type ParserFunc func(input string) (model.Node, error)

type Svc struct {
	Repo   Repo
	Parser ParserFunc
}

// New creates a new Svc instance
func New(repo Repo) *Svc {
	return &Svc{
		Repo:   repo,
		Parser: Parse,
	}
}
