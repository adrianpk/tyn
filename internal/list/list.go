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
	var tagFilter string
	var placeFilter string
	var statusFilter string
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list all nodes or filter by type (note, task, link), tag, or place",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodes, err := svc.Repo.List(cmd.Context())
			if err != nil {
				return err
			}

			var filterType string
			if len(args) == 1 {
				filterType = args[0]
			}

			// WIP: Printing nodes is of limited utility, but this is just an initial
			// approach to define the command logic.
			for _, node := range nodes {
				if matchesFilters(node, filterType, tagFilter, placeFilter, statusFilter) {
					fmt.Println(printNode(node))
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&tagFilter, "tag", "t", "", "filter by tag")
	cmd.Flags().StringVarP(&placeFilter, "place", "p", "", "filter by place")
	cmd.Flags().StringVarP(&statusFilter, "status", "s", "", "filter by status")

	return cmd
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

func matchesFilters(node model.Node, filterType, tagFilter, placeFilter, statusFilter string) bool {
	if filterType != "" && node.Type != filterType {
		return false
	}

	if tagFilter != "" && !hasTag(node, tagFilter) {
		return false
	}

	if placeFilter != "" && !hasPlace(node, placeFilter) {
		return false
	}

	if statusFilter != "" && node.Status != statusFilter {
		return false
	}

	return true
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
