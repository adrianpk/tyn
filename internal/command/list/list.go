package list

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adrianpk/tyn/internal/bkg"
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
		Aliases: []string{"ls", "l"},
		Short:   "list all nodes or filter by type (note, task, link), tag, or place",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var nodes []model.Node
			var err error

			var filterType string
			if len(args) == 1 {
				filterType = args[0]
			}

			// Direct svc use if available
			if svc != nil {
				nodes, err = svc.Repo.List(cmd.Context())
				if err != nil {
					return err
				}

				for _, node := range nodes {
					if matchesFilters(node, filterType, tagFilter, placeFilter, statusFilter) {
						fmt.Println(printNode(node))
					}
				}

				return nil
			}

			// IPC Otherwise
			params := bkg.ListParams{
				Type: filterType,
			}

			if tagFilter != "" {
				params.Tags = strings.Split(tagFilter, ",")
			}

			if placeFilter != "" {
				params.Places = strings.Split(placeFilter, ",")
			}

			if statusFilter != "" {
				params.Status = statusFilter
			}

			resp, err := bkg.SendCommand("list", params)
			if err != nil {
				return fmt.Errorf("error communicating with daemon: %w", err)
			}

			if !resp.Success {
				return fmt.Errorf("daemon returned error: %s", resp.Error)
			}

			err = json.Unmarshal(resp.Data, &nodes)
			if err != nil {
				return fmt.Errorf("error parsing response: %w", err)
			}

			for _, node := range nodes {
				fmt.Println(printNode(node))
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
	dueDateStr := "nil"
	if node.DueDate != nil {
		dueDateStr = node.DueDate.Format("2006-01-02 15:04:05 -0700 MST")
	}

	dateStr := node.Date.Format("2006-01-02 15:04:05 -0700 MST")

	return fmt.Sprintf(
		"ID: %s\nType: %s\nContent: %s\nTags: %v\nPlaces: %v\nStatus: %s\nLink: %s\nDate: %s\nDueDate: %s\n",
		node.ID,
		node.Type,
		node.Content,
		node.Tags,
		node.Places,
		node.Status,
		node.Link,
		dateStr,
		dueDateStr,
	)
}
