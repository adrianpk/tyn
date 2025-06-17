package capture

import (
	"context"
)

type Repo interface {
	Create(ctx context.Context, node Node) error
	Get(ctx context.Context, id string) (Node, error)
	Update(ctx context.Context, node Node) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]Node, error)
}
