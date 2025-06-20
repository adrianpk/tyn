package list

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adrianpk/tyn/internal/bkg"
	"github.com/adrianpk/tyn/internal/command/common"
	"github.com/adrianpk/tyn/internal/model"
	"github.com/adrianpk/tyn/internal/svc"
	"github.com/spf13/cobra"
)

type ListCommand struct {
	common.BaseCommand
	tagFilter    string
	placeFilter  string
	statusFilter string
}

func NewCommand(svc *svc.Svc) *cobra.Command {
	cmd := &ListCommand{
		BaseCommand: common.BaseCommand{
			Svc:         svc,
			CommandName: "list",
		},
	}

	cobraCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "list all nodes or filter by type (note, task, link), tag, or place",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cobra *cobra.Command, args []string) error {
			flags := common.ExtractFlagsFromCommand(cobra)
			return cmd.Execute(cobra.Context(), args, flags, cmd.ExecuteDirect, cmd.ExecuteViaIPC)
		},
	}

	cobraCmd.Flags().StringVarP(&cmd.tagFilter, "tag", "t", "", "filter by tag")
	cobraCmd.Flags().StringVarP(&cmd.placeFilter, "place", "p", "", "filter by place")
	cobraCmd.Flags().StringVarP(&cmd.statusFilter, "status", "s", "", "filter by status")

	cmd.CobraCmd = cobraCmd
	return cobraCmd
}

func (c *ListCommand) ExecuteDirect(ctx context.Context, args []string, flags map[string]interface{}) error {
	var filterType string
	if len(args) == 1 {
		filterType = args[0]
	}

	nodes, err := c.Svc.Repo.List(ctx)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		if matchesFilters(node, filterType, c.tagFilter, c.placeFilter, c.statusFilter) {
			fmt.Println(printNode(node))
		}
	}

	return nil
}

func (c *ListCommand) ExecuteViaIPC(args []string, flags map[string]interface{}) error {
	var filterType string
	if len(args) == 1 {
		filterType = args[0]
	}

	params := bkg.ListParams{
		Type: filterType,
	}

	if c.tagFilter != "" {
		params.Tags = strings.Split(c.tagFilter, ",")
	}

	if c.placeFilter != "" {
		params.Places = strings.Split(c.placeFilter, ",")
	}

	if c.statusFilter != "" {
		params.Status = c.statusFilter
	}

	resp, err := bkg.SendCommand("list", params)
	if err != nil {
		return fmt.Errorf("error communicating with daemon: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("daemon returned error: %s", resp.Error)
	}

	var nodes []model.Node
	err = json.Unmarshal(resp.Data, &nodes)
	if err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	for _, node := range nodes {
		fmt.Println(printNode(node))
	}

	return nil
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
