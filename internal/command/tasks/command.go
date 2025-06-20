package tasks

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adrianpk/tyn/internal/bkg"
	"github.com/adrianpk/tyn/internal/model"
	"github.com/adrianpk/tyn/internal/svc"
	"github.com/spf13/cobra"
)

func NewCommand(svc *svc.Svc) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tasks",
		Aliases: []string{"t"},
		Short:   "Manage tasks",
		Long:    "Manage tasks (list, filter, change status) in the system",
		RunE: func(cmd *cobra.Command, args []string) error {
			listCmd := newListCommand(svc)
			return listCmd.RunE(listCmd, args)
		},
	}

	cmd.AddCommand(newListCommand(svc))
	cmd.AddCommand(newStatusCommand(svc))

	return cmd
}

func newStatusCommand(svc *svc.Svc) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"st", "s"},
		Short:   "Change task status",
		Long:    "Change the status of a task using its ID",
	}

	cmd.AddCommand(newStatusSetCommand(svc))
	cmd.AddCommand(newStatusNextCommand(svc))
	cmd.AddCommand(newStatusPrevCommand(svc))

	return cmd
}

func newStatusSetCommand(svc *svc.Svc) *cobra.Command {
	return &cobra.Command{
		Use:   "set <id> <status>",
		Short: "Set specific status",
		Long:  "Set a specific status for a task using its ID",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			status := args[1]

			return changeTaskStatus(svc, id, status, "set")
		},
	}
}

func newStatusNextCommand(svc *svc.Svc) *cobra.Command {
	return &cobra.Command{
		Use:   "next <id>",
		Short: "Move to next status",
		Long:  "Advance a task to the next status in the cycle",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			return changeTaskStatus(svc, id, "", "next")
		},
	}
}

func newStatusPrevCommand(svc *svc.Svc) *cobra.Command {
	return &cobra.Command{
		Use:   "prev <id>",
		Short: "Move to previous status",
		Long:  "Move a task to the previous status in the cycle",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			return changeTaskStatus(svc, id, "", "prev")
		},
	}
}

func changeTaskStatus(svc *svc.Svc, id, targetStatus, operation string) error {
	if svc != nil {
		nodes, err := svc.Repo.List(nil)
		if err != nil {
			return err
		}

		targetNode, found := findNodeByShortID(nodes, id)
		if !found {
			return fmt.Errorf("task with ID '%s' not found", id)
		}

		originalStatus := targetNode.Status
		var newStatus string

		switch operation {
		case "set":
			if !model.ValidStatus(targetStatus) {
				return fmt.Errorf("invalid status: %s", targetStatus)
			}
			newStatus = targetStatus
		case "next":
			newStatus = model.NextStatus(targetNode.Status)
		case "prev":
			newStatus = model.PreviousStatus(targetNode.Status)
		}

		targetNode.Status = newStatus

		err = svc.Repo.Update(nil, targetNode)
		if err != nil {
			return err
		}

		fmt.Printf("Task status updated: '%s' → '%s'\n", originalStatus, newStatus)
		displayStatusCycle(originalStatus, newStatus)

		return nil
	}

	params := bkg.StatusParams{
		ID:        id,
		Status:    targetStatus,
		Operation: operation,
	}

	resp, err := bkg.SendCommand("status", params)
	if err != nil {
		return fmt.Errorf("error communicating with daemon: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("daemon returned error: %s", resp.Error)
	}

	var result struct {
		OriginalStatus string `json:"original_status"`
		NewStatus      string `json:"new_status"`
	}

	err = json.Unmarshal(resp.Data, &result)
	if err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Printf("Task status updated: '%s' → '%s'\n", result.OriginalStatus, result.NewStatus)
	displayStatusCycle(result.OriginalStatus, result.NewStatus)

	return nil
}

func findNodeByShortID(nodes []model.Node, shortID string) (model.Node, bool) {
	for _, node := range nodes {
		if strings.HasPrefix(node.ID, shortID) || node.ShortID() == shortID {
			return node, true
		}
	}
	return model.Node{}, false
}

func displayStatusCycle(originalStatus, newStatus string) {
	var statusDisplay string

	for i, status := range model.StatusCycle {
		if i > 0 {
			statusDisplay += " → "
		}

		if status == newStatus {
			statusDisplay += fmt.Sprintf("[%s]", status)
		} else if status == originalStatus {
			statusDisplay += fmt.Sprintf("<%s>", status)
		} else {
			statusDisplay += status
		}
	}

	fmt.Println(statusDisplay)
}

func newListCommand(svc *svc.Svc) *cobra.Command {
	var tagFilter string
	var placeFilter string
	var statusFilter string

	cmd := &cobra.Command{
		Use:     "list [status] [#tag] [@place]",
		Aliases: []string{"ls", "l"},
		Short:   "List tasks with optional filtering",
		Long:    "List tasks with optional filtering by status, tags, and places",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if strings.HasPrefix(arg, ":") {
					statusFilter = strings.TrimPrefix(arg, ":")
				} else if strings.HasPrefix(arg, "#") {
					tagFilter = strings.TrimPrefix(arg, "#")
				} else if strings.HasPrefix(arg, "@") {
					placeFilter = strings.TrimPrefix(arg, "@")
				}
			}

			if svc != nil {
				nodes, err := svc.Repo.List(cmd.Context())
				if err != nil {
					return err
				}

				tasks := filterTasks(nodes, statusFilter, tagFilter, placeFilter)
				printTasks(tasks)
				return nil
			}

			params := bkg.ListParams{
				Type:   "task",
				Status: statusFilter,
			}

			if tagFilter != "" {
				params.Tags = strings.Split(tagFilter, ",")
			}

			if placeFilter != "" {
				params.Places = strings.Split(placeFilter, ",")
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

			tasks := filterTasks(nodes, statusFilter, tagFilter, placeFilter)
			printTasks(tasks)
			return nil
		},
	}

	cmd.Flags().StringVarP(&tagFilter, "tag", "t", "", "filter by tag")
	cmd.Flags().StringVarP(&placeFilter, "place", "p", "", "filter by place")
	cmd.Flags().StringVarP(&statusFilter, "status", "s", "", "filter by status")

	return cmd
}

func filterTasks(nodes []model.Node, statusFilter, tagFilter, placeFilter string) []model.Node {
	var tasks []model.Node

	for _, node := range nodes {
		if node.Type != "task" {
			continue
		}

		if statusFilter != "" && node.Status != statusFilter {
			continue
		}

		if tagFilter != "" && !hasTag(node, tagFilter) {
			continue
		}

		if placeFilter != "" && !hasPlace(node, placeFilter) {
			continue
		}

		tasks = append(tasks, node)
	}

	return tasks
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

func printTasks(tasks []model.Node) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	contentWidth := 45

	fmt.Printf("%-6s %-10s %-45s %-20s %s\n", "ID", "STATUS", "CONTENT", "TAGS/PLACES", "!")
	fmt.Println(strings.Repeat("-", 90))

	for _, task := range tasks {
		var metadata []string

		for _, tag := range task.Tags {
			metadata = append(metadata, fmt.Sprintf("#%s", tag))
		}

		for _, place := range task.Places {
			metadata = append(metadata, fmt.Sprintf("@%s", place))
		}

		metadataStr := strings.Join(metadata, " ")

		if len(metadataStr) > 20 {
			metadataStr = metadataStr[:17] + "..."
		}

		statusDisplay := fmt.Sprintf("[%s]", task.Status)

		overdueIndicator := " "
		if task.IsOverdue() {
			overdueIndicator = "⌛"
		}

		content := strings.TrimSpace(task.Content)
		content = strings.Join(strings.Fields(content), " ")

		if len(content) > contentWidth {
			content = content[:contentWidth-3] + "..."
		}

		fmt.Printf("%-6s %-10s %-45s %-20s %s\n",
			task.ShortID(),
			statusDisplay,
			content,
			metadataStr,
			overdueIndicator)
	}
}
