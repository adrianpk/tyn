package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adrianpk/tyn/internal/bkg"
	"github.com/adrianpk/tyn/internal/command/common"
	"github.com/adrianpk/tyn/internal/model"
	"github.com/adrianpk/tyn/internal/svc"
	"github.com/spf13/cobra"
)

type (
	TasksCommand struct {
		common.BaseCommand
	}

	TasksListCommand struct {
		common.BaseCommand
		tagFilter    string
		placeFilter  string
		statusFilter string
	}

	TasksStatusCommand struct {
		common.BaseCommand
	}

	TasksStatusSetCommand struct {
		common.BaseCommand
	}

	TasksStatusNextCommand struct {
		common.BaseCommand
	}

	TasksStatusPrevCommand struct {
		common.BaseCommand
	}
)

func NewCommand(svc *svc.Svc) *cobra.Command {
	cmd := &TasksCommand{
		BaseCommand: common.BaseCommand{
			Svc:         svc,
			CommandName: "tasks",
		},
	}

	cobraCmd := &cobra.Command{
		Use:     "tasks",
		Aliases: []string{"t"},
		Short:   "Manage tasks",
		Long:    "Manage tasks (list, filter, change status) in the system",
		RunE: func(cobra *cobra.Command, args []string) error {
			flags := common.ExtractFlagsFromCommand(cobra)
			listCmd := newListCommand(svc)
			return cmd.Execute(cobra.Context(), args, flags,
				func(ctx context.Context, args []string, flags map[string]interface{}) error {
					return listCmd.RunE(listCmd, args)
				},
				func(args []string, flags map[string]interface{}) error {
					return listCmd.RunE(listCmd, args)
				})
		},
	}

	cobraCmd.AddCommand(newListCommand(svc))
	cobraCmd.AddCommand(newStatusCommand(svc))

	cmd.CobraCmd = cobraCmd
	return cobraCmd
}

func newStatusCommand(svc *svc.Svc) *cobra.Command {
	cmd := &TasksStatusCommand{
		BaseCommand: common.BaseCommand{
			Svc:         svc,
			CommandName: "status",
		},
	}

	cobraCmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"st", "s"},
		Short:   "Change task status",
		Long:    "Change the status of a task using its ID",
	}

	cobraCmd.AddCommand(newStatusSetCommand(svc))
	cobraCmd.AddCommand(newStatusNextCommand(svc))
	cobraCmd.AddCommand(newStatusPrevCommand(svc))

	cmd.CobraCmd = cobraCmd
	return cobraCmd
}

func newStatusSetCommand(svc *svc.Svc) *cobra.Command {
	cmd := &TasksStatusSetCommand{
		BaseCommand: common.BaseCommand{
			Svc:         svc,
			CommandName: "status_set",
		},
	}

	cobraCmd := &cobra.Command{
		Use:   "set <id> <status>",
		Short: "Set specific status",
		Long:  "Set a specific status for a task using its ID",
		Args:  cobra.ExactArgs(2),
		RunE: func(cobra *cobra.Command, args []string) error {
			flags := common.ExtractFlagsFromCommand(cobra)
			return cmd.Execute(cobra.Context(), args, flags, cmd.ExecuteDirect, cmd.ExecuteViaIPC)
		},
	}

	cmd.CobraCmd = cobraCmd
	return cobraCmd
}

func (c *TasksStatusSetCommand) ExecuteDirect(ctx context.Context, args []string, flags map[string]interface{}) error {
	log.Printf("Executing status set command directly")
	id := args[0]
	targetStatus := args[1]
	return changeTaskStatus(c.Svc, id, targetStatus, "set")
}

func (c *TasksStatusSetCommand) ExecuteViaIPC(args []string, flags map[string]interface{}) error {
	log.Printf("Executing status set command via IPC")
	id := args[0]
	targetStatus := args[1]

	params := bkg.StatusParams{
		ID:        id,
		Status:    targetStatus,
		Operation: "set",
	}

	resp, err := bkg.SendCommand("status", params)
	if err != nil {
		return fmt.Errorf("error communicating with daemon: %w", err)
	}

	return handleStatusResponse(resp)
}

func newStatusNextCommand(svc *svc.Svc) *cobra.Command {
	cmd := &TasksStatusNextCommand{
		BaseCommand: common.BaseCommand{
			Svc:         svc,
			CommandName: "status_next",
		},
	}

	cobraCmd := &cobra.Command{
		Use:   "next <id>",
		Short: "Move to next status",
		Long:  "Advance a task to the next status in the cycle",
		Args:  cobra.ExactArgs(1),
		RunE: func(cobra *cobra.Command, args []string) error {
			flags := common.ExtractFlagsFromCommand(cobra)
			return cmd.Execute(cobra.Context(), args, flags, cmd.ExecuteDirect, cmd.ExecuteViaIPC)
		},
	}

	cmd.CobraCmd = cobraCmd
	return cobraCmd
}

func (c *TasksStatusNextCommand) ExecuteDirect(ctx context.Context, args []string, flags map[string]interface{}) error {
	log.Printf("Executing status next command directly")
	id := args[0]
	return changeTaskStatus(c.Svc, id, "", "next")
}

func (c *TasksStatusNextCommand) ExecuteViaIPC(args []string, flags map[string]interface{}) error {
	log.Printf("Executing status next command via IPC")
	id := args[0]

	params := bkg.StatusParams{
		ID:        id,
		Status:    "",
		Operation: "next",
	}

	resp, err := bkg.SendCommand("status", params)
	if err != nil {
		return fmt.Errorf("error communicating with daemon: %w", err)
	}

	return handleStatusResponse(resp)
}

func newStatusPrevCommand(svc *svc.Svc) *cobra.Command {
	cmd := &TasksStatusPrevCommand{
		BaseCommand: common.BaseCommand{
			Svc:         svc,
			CommandName: "status_prev",
		},
	}

	cobraCmd := &cobra.Command{
		Use:   "prev <id>",
		Short: "Move to previous status",
		Long:  "Move a task to the previous status in the cycle",
		Args:  cobra.ExactArgs(1),
		RunE: func(cobra *cobra.Command, args []string) error {
			flags := common.ExtractFlagsFromCommand(cobra)
			return cmd.Execute(cobra.Context(), args, flags, cmd.ExecuteDirect, cmd.ExecuteViaIPC)
		},
	}

	cmd.CobraCmd = cobraCmd
	return cobraCmd
}

func (c *TasksStatusPrevCommand) ExecuteDirect(ctx context.Context, args []string, flags map[string]interface{}) error {
	log.Printf("Executing status prev command directly")
	id := args[0]
	return changeTaskStatus(c.Svc, id, "", "prev")
}

func (c *TasksStatusPrevCommand) ExecuteViaIPC(args []string, flags map[string]interface{}) error {
	log.Printf("Executing status prev command via IPC")
	id := args[0]

	params := bkg.StatusParams{
		ID:        id,
		Status:    "",
		Operation: "prev",
	}

	resp, err := bkg.SendCommand("status", params)
	if err != nil {
		return fmt.Errorf("error communicating with daemon: %w", err)
	}

	return handleStatusResponse(resp)
}

func handleStatusResponse(resp *bkg.Response) error {
	if !resp.Success {
		return fmt.Errorf("daemon returned error: %s", resp.Error)
	}

	var result struct {
		OriginalStatus string `json:"original_status"`
		NewStatus      string `json:"new_status"`
	}

	err := json.Unmarshal(resp.Data, &result)
	if err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Printf("Task status updated: '%s' → '%s'\n", result.OriginalStatus, result.NewStatus)
	displayStatusCycle(result.OriginalStatus, result.NewStatus)

	return nil
}

func newListCommand(svc *svc.Svc) *cobra.Command {
	cmd := &TasksListCommand{
		BaseCommand: common.BaseCommand{
			Svc:         svc,
			CommandName: "tasks_list",
		},
	}

	cobraCmd := &cobra.Command{
		Use:     "list [status] [#tag] [@place]",
		Aliases: []string{"ls", "l"},
		Short:   "List tasks with optional filtering",
		Long:    "List tasks with optional filtering by status, tags, and places",
		RunE: func(cobra *cobra.Command, args []string) error {
			for _, arg := range args {
				if strings.HasPrefix(arg, ":") {
					cmd.statusFilter = strings.TrimPrefix(arg, ":")
				} else if strings.HasPrefix(arg, "#") {
					cmd.tagFilter = strings.TrimPrefix(arg, "#")
				} else if strings.HasPrefix(arg, "@") {
					cmd.placeFilter = strings.TrimPrefix(arg, "@")
				}
			}

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

func (c *TasksListCommand) ExecuteDirect(ctx context.Context, args []string, flags map[string]interface{}) error {
	log.Printf("Executing tasks list command directly")
	nodes, err := c.Svc.Repo.List(ctx)
	if err != nil {
		return err
	}

	tasks := filterTasks(nodes, c.statusFilter, c.tagFilter, c.placeFilter)
	printTasks(tasks)
	return nil
}

func (c *TasksListCommand) ExecuteViaIPC(args []string, flags map[string]interface{}) error {
	log.Printf("Executing tasks list command via IPC")
	params := bkg.ListParams{
		Type:   "task",
		Status: c.statusFilter,
	}

	if c.tagFilter != "" {
		params.Tags = strings.Split(c.tagFilter, ",")
	}

	if c.placeFilter != "" {
		params.Places = strings.Split(c.placeFilter, ",")
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

	tasks := filterTasks(nodes, c.statusFilter, c.tagFilter, c.placeFilter)
	printTasks(tasks)
	return nil
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

	return fmt.Errorf("service not available")
}

// Las demás funciones auxiliares permanecen sin cambios
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
