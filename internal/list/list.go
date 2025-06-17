package list

import (
	"context"
	"fmt"

	"github.com/adrianpk/tyn/internal/model"
	"github.com/adrianpk/tyn/internal/svc"
	"github.com/spf13/cobra"
)

type Repo interface {
	List(ctx context.Context) ([]model.Node, error)
}

func NewCommand(svc *svc.Svc) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list all nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			nodes, err := svc.Repo.List(cmd.Context())
			if err != nil {
				return err
			}

			for _, node := range nodes {
				fmt.Println(printNode(node))
			}
			return nil
		},
	}
	return cmd
}

func printNode(node model.Node) string {
	return fmt.Sprintf(
		"ID: %s\nType: %s\nContent: %s\nTags: %v\nPlaces: %v\nStatus: %s\nLink: %s\nDate: %s\nOverrideDate: %v\n",
		node.ID,
		node.Type,
		node.Content,
		node.Tags,
		node.Places,
		node.Status,
		node.Link,
		node.Date.Format("2006-01-02 15:04:05"),
		node.OverrideDate,
	)
}
