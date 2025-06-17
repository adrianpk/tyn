package svc

import (
	"context"

	"github.com/adrianpk/tyn/internal/model"
)

// Repo is the interface for repository operations
// (moved from svc.go)
type Repo interface {
	Create(ctx context.Context, node model.Node) error
	Get(ctx context.Context, id string) (model.Node, error)
	Update(ctx context.Context, node model.Node) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]model.Node, error)
}
