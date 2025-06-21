package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

	TasksTextCommand struct {
		common.BaseCommand
	}

	TasksTagCommand struct {
		common.BaseCommand
	}

	TasksPlaceCommand struct {
		common.BaseCommand
	}

	TasksDateCommand struct {
		common.BaseCommand
	}

	TasksUpdateCommand struct {
		common.BaseCommand
		tags   []string
		places []string
		due    string
		text   string
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
	cobraCmd.AddCommand(newUpdateCommand(svc))
	cobraCmd.AddCommand(newTagCommand(svc))
	cobraCmd.AddCommand(newPlaceCommand(svc))
	cobraCmd.AddCommand(newDateCommand(svc))

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
	nodes, err := c.Svc.Repo.List(context.TODO())
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
		nodes, err := svc.Repo.List(context.TODO())
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

		err = svc.Repo.Update(context.TODO(), targetNode)
		if err != nil {
			return err
		}

		fmt.Printf("Task status updated: '%s' → '%s'\n", originalStatus, newStatus)
		displayStatusCycle(originalStatus, newStatus)

		return nil
	}

	return fmt.Errorf("service not available")
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

func newTextCommand(svc *svc.Svc) *cobra.Command {
	cmd := &TasksTextCommand{
		BaseCommand: common.BaseCommand{
			Svc:         svc,
			CommandName: "text",
		},
	}

	cobraCmd := &cobra.Command{
		Use:   "text <id> <new_text>",
		Short: "Edit task text",
		Args:  cobra.ExactArgs(2),
		RunE: func(cobra *cobra.Command, args []string) error {
			flags := common.ExtractFlagsFromCommand(cobra)
			return cmd.Execute(cobra.Context(), args, flags, cmd.ExecuteDirect, cmd.ExecuteViaIPC)
		},
	}

	cmd.CobraCmd = cobraCmd
	return cobraCmd
}

func (c *TasksTextCommand) ExecuteDirect(ctx context.Context, args []string, flags map[string]interface{}) error {
	log.Printf("Executing text command directly")
	id := args[0]
	newText := args[1]

	task, err := c.Svc.Repo.GetTaskByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error fetching task: %w", err)
	}

	task.Content = newText
	if err := c.Svc.Repo.Update(ctx, task); err != nil {
		return fmt.Errorf("error updating task: %w", err)
	}

	log.Printf("Task text updated successfully: %s", task.ID)
	return nil
}

func (c *TasksTextCommand) ExecuteViaIPC(args []string, flags map[string]interface{}) error {
	log.Printf("Executing text command via IPC")
	id := args[0]
	newText := args[1]

	params := bkg.TextParams{
		ID:   id,
		Text: newText,
	}

	resp, err := bkg.SendCommand("text", params)
	if err != nil {
		return fmt.Errorf("error communicating with daemon: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("daemon returned error: %s", resp.Error)
	}

	log.Printf("Task text updated successfully via daemon: %s", id)
	return nil
}

func newTagCommand(svc *svc.Svc) *cobra.Command {
	cmd := &TasksTagCommand{
		BaseCommand: common.BaseCommand{
			Svc:         svc,
			CommandName: "tag",
		},
	}

	cobraCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage task tags",
	}

	cobraCmd.AddCommand(
		newTagAddCommand(svc),
		newTagRemoveCommand(svc),
		newTagClearCommand(svc),
	)

	cmd.CobraCmd = cobraCmd
	return cobraCmd
}

func newTagAddCommand(svc *svc.Svc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <id> <tag>",
		Short: "Add a tag to a task",
		Args:  cobra.ExactArgs(2),
		RunE: func(cobra *cobra.Command, args []string) error {
			id := args[0]
			tag := args[1]

			// Check environment to determine execution mode
			if os.Getenv("TYN_DEV") == "1" {
				// Execute via IPC in dev mode
				params := bkg.TagCmdParams{
					ID:        id,
					Tags:      []string{tag},
					Operation: "add",
				}

				resp, err := bkg.SendCommand("tag", params)
				if err != nil {
					return fmt.Errorf("error communicating with daemon: %w", err)
				}

				if !resp.Success {
					return fmt.Errorf("daemon returned error: %s", resp.Error)
				}

				var result struct {
					Message string `json:"message"`
				}

				err = json.Unmarshal(resp.Data, &result)
				if err != nil {
					return fmt.Errorf("error parsing response: %w", err)
				}

				fmt.Println(result.Message)
				return nil
			}

			// Execute directly in non-dev mode
			task, err := svc.Repo.GetTaskByID(cobra.Context(), id)
			if err != nil {
				return err
			}
			task.Tags = append(task.Tags, tag)
			err = svc.Repo.Update(cobra.Context(), task)
			if err != nil {
				return err
			}
			fmt.Printf("Added tag '%s' to task %s\n", tag, id)
			return nil
		},
	}

	return cmd
}

func newTagRemoveCommand(svc *svc.Svc) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <id> <tag>",
		Short: "Remove a tag from a task",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			tag := args[1]

			// Check environment to determine execution mode
			if os.Getenv("TYN_DEV") == "1" {
				// Execute via IPC in dev mode
				params := bkg.TagCmdParams{
					ID:        id,
					Tags:      []string{tag},
					Operation: "remove",
				}

				resp, err := bkg.SendCommand("tag", params)
				if err != nil {
					return fmt.Errorf("error communicating with daemon: %w", err)
				}

				if !resp.Success {
					return fmt.Errorf("daemon returned error: %s", resp.Error)
				}

				var result struct {
					Message string `json:"message"`
				}

				err = json.Unmarshal(resp.Data, &result)
				if err != nil {
					return fmt.Errorf("error parsing response: %w", err)
				}

				fmt.Println(result.Message)
				return nil
			}

			// Execute directly in non-dev mode
			task, err := svc.Repo.GetTaskByID(cmd.Context(), id)
			if err != nil {
				return err
			}
			filteredTags := []string{}
			for _, t := range task.Tags {
				if t != tag {
					filteredTags = append(filteredTags, t)
				}
			}
			task.Tags = filteredTags
			err = svc.Repo.Update(cmd.Context(), task)
			if err != nil {
				return err
			}
			fmt.Printf("Removed tag '%s' from task %s\n", tag, id)
			return nil
		},
	}
}

func newTagClearCommand(svc *svc.Svc) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <id>",
		Short: "Clear all tags from a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			// Check environment to determine execution mode
			if os.Getenv("TYN_DEV") == "1" {
				// Execute via IPC in dev mode
				params := bkg.TagCmdParams{
					ID:        id,
					Tags:      []string{},
					Operation: "clear",
				}

				resp, err := bkg.SendCommand("tag", params)
				if err != nil {
					return fmt.Errorf("error communicating with daemon: %w", err)
				}

				if !resp.Success {
					return fmt.Errorf("daemon returned error: %s", resp.Error)
				}

				var result struct {
					Message string `json:"message"`
				}

				err = json.Unmarshal(resp.Data, &result)
				if err != nil {
					return fmt.Errorf("error parsing response: %w", err)
				}

				fmt.Println(result.Message)
				return nil
			}

			// Execute directly in non-dev mode
			task, err := svc.Repo.GetTaskByID(cmd.Context(), id)
			if err != nil {
				return err
			}
			task.Tags = []string{}
			err = svc.Repo.Update(cmd.Context(), task)
			if err != nil {
				return err
			}
			fmt.Printf("Cleared all tags from task %s\n", id)
			return nil
		},
	}
}

func newPlaceCommand(svc *svc.Svc) *cobra.Command {
	cmd := &TasksPlaceCommand{
		BaseCommand: common.BaseCommand{
			Svc:         svc,
			CommandName: "place",
		},
	}

	cobraCmd := &cobra.Command{
		Use:   "place",
		Short: "Manage task places",
	}

	cobraCmd.AddCommand(
		newPlaceAddCommand(svc),
		newPlaceRemoveCommand(svc),
		newPlaceClearCommand(svc),
	)

	cmd.CobraCmd = cobraCmd
	return cobraCmd
}

func newPlaceAddCommand(svc *svc.Svc) *cobra.Command {
	return &cobra.Command{
		Use:   "add <id> <place>",
		Short: "Add a place to a task",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			place := args[1]

			// Check environment to determine execution mode
			if os.Getenv("TYN_DEV") == "1" {
				// Execute via IPC in dev mode
				params := bkg.PlaceCmdParams{
					ID:        id,
					Places:    []string{place},
					Operation: "add",
				}

				resp, err := bkg.SendCommand("place", params)
				if err != nil {
					return fmt.Errorf("error communicating with daemon: %w", err)
				}

				if !resp.Success {
					return fmt.Errorf("daemon returned error: %s", resp.Error)
				}

				var result struct {
					Message string `json:"message"`
				}

				err = json.Unmarshal(resp.Data, &result)
				if err != nil {
					return fmt.Errorf("error parsing response: %w", err)
				}

				fmt.Println(result.Message)
				return nil
			}

			// Execute directly in non-dev mode
			task, err := svc.Repo.GetTaskByID(cmd.Context(), id)
			if err != nil {
				return err
			}
			task.Places = append(task.Places, place)
			err = svc.Repo.Update(cmd.Context(), task)
			if err != nil {
				return err
			}
			fmt.Printf("Added place '%s' to task %s\n", place, id)
			return nil
		},
	}
}

func newPlaceRemoveCommand(svc *svc.Svc) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <id> <place>",
		Short: "Remove a place from a task",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			place := args[1]

			// Check environment to determine execution mode
			if os.Getenv("TYN_DEV") == "1" {
				// Execute via IPC in dev mode
				params := bkg.PlaceCmdParams{
					ID:        id,
					Places:    []string{place},
					Operation: "remove",
				}

				resp, err := bkg.SendCommand("place", params)
				if err != nil {
					return fmt.Errorf("error communicating with daemon: %w", err)
				}

				if !resp.Success {
					return fmt.Errorf("daemon returned error: %s", resp.Error)
				}

				var result struct {
					Message string `json:"message"`
				}

				err = json.Unmarshal(resp.Data, &result)
				if err != nil {
					return fmt.Errorf("error parsing response: %w", err)
				}

				fmt.Println(result.Message)
				return nil
			}

			// Execute directly in non-dev mode
			task, err := svc.Repo.GetTaskByID(cmd.Context(), id)
			if err != nil {
				return err
			}
			filteredPlaces := []string{}
			for _, p := range task.Places {
				if p != place {
					filteredPlaces = append(filteredPlaces, p)
				}
			}
			task.Places = filteredPlaces
			err = svc.Repo.Update(cmd.Context(), task)
			if err != nil {
				return err
			}
			fmt.Printf("Removed place '%s' from task %s\n", place, id)
			return nil
		},
	}
}

func newPlaceClearCommand(svc *svc.Svc) *cobra.Command {
	return &cobra.Command{
		Use:   "clear <id>",
		Short: "Clear all places from a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			// Check environment to determine execution mode
			if os.Getenv("TYN_DEV") == "1" {
				// Execute via IPC in dev mode
				params := bkg.PlaceCmdParams{
					ID:        id,
					Places:    []string{},
					Operation: "clear",
				}

				resp, err := bkg.SendCommand("place", params)
				if err != nil {
					return fmt.Errorf("error communicating with daemon: %w", err)
				}

				if !resp.Success {
					return fmt.Errorf("daemon returned error: %s", resp.Error)
				}

				var result struct {
					Message string `json:"message"`
				}

				err = json.Unmarshal(resp.Data, &result)
				if err != nil {
					return fmt.Errorf("error parsing response: %w", err)
				}

				fmt.Println(result.Message)
				return nil
			}

			// Execute directly in non-dev mode
			task, err := svc.Repo.GetTaskByID(cmd.Context(), id)
			if err != nil {
				return err
			}
			task.Places = []string{}
			err = svc.Repo.Update(cmd.Context(), task)
			if err != nil {
				return err
			}
			fmt.Printf("Cleared all places from task %s\n", id)
			return nil
		},
	}
}

func newDateCommand(svc *svc.Svc) *cobra.Command {
	cmd := &TasksDateCommand{
		BaseCommand: common.BaseCommand{
			Svc:         svc,
			CommandName: "date",
		},
	}

	cobraCmd := &cobra.Command{
		Use:   "date",
		Short: "Manage task due dates",
	}

	cobraCmd.AddCommand(
		newDateSetCommand(svc),
		newDateRemoveCommand(svc),
	)

	cmd.CobraCmd = cobraCmd
	return cobraCmd
}

func newDateSetCommand(svc *svc.Svc) *cobra.Command {
	return &cobra.Command{
		Use:   "set <id> <date>",
		Short: "Set a due date for a task",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			dateStr := args[1]

			// Check environment to determine execution mode
			if os.Getenv("TYN_DEV") == "1" {
				// Execute via IPC in dev mode
				params := bkg.DateCmdParams{
					ID:        id,
					Date:      dateStr,
					Operation: "set",
				}

				resp, err := bkg.SendCommand("date", params)
				if err != nil {
					return fmt.Errorf("error communicating with daemon: %w", err)
				}

				if !resp.Success {
					return fmt.Errorf("daemon returned error: %s", resp.Error)
				}

				var result struct {
					Message      string `json:"message"`
					OriginalDate string `json:"original_date"`
					NewDate      string `json:"new_date"`
				}

				err = json.Unmarshal(resp.Data, &result)
				if err != nil {
					return fmt.Errorf("error parsing response: %w", err)
				}

				fmt.Println(result.Message)
				return nil
			}

			// Execute directly in non-dev mode
			task, err := svc.Repo.GetTaskByID(cmd.Context(), id)
			if err != nil {
				return err
			}
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return fmt.Errorf("invalid date format: %v", err)
			}
			task.DueDate = &date
			err = svc.Repo.Update(cmd.Context(), task)
			if err != nil {
				return err
			}
			fmt.Printf("Set due date to %s for task %s\n", dateStr, id)
			return nil
		},
	}
}

func newDateRemoveCommand(svc *svc.Svc) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <id>",
		Short: "Remove the due date from a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			// Check environment to determine execution mode
			if os.Getenv("TYN_DEV") == "1" {
				// Execute via IPC in dev mode
				params := bkg.DateCmdParams{
					ID:        id,
					Operation: "remove",
				}

				resp, err := bkg.SendCommand("date", params)
				if err != nil {
					return fmt.Errorf("error communicating with daemon: %w", err)
				}

				if !resp.Success {
					return fmt.Errorf("daemon returned error: %s", resp.Error)
				}

				var result struct {
					Message      string `json:"message"`
					OriginalDate string `json:"original_date"`
				}

				err = json.Unmarshal(resp.Data, &result)
				if err != nil {
					return fmt.Errorf("error parsing response: %w", err)
				}

				fmt.Println(result.Message)
				return nil
			}

			// Execute directly in non-dev mode
			task, err := svc.Repo.GetTaskByID(cmd.Context(), id)
			if err != nil {
				return err
			}
			task.DueDate = nil
			err = svc.Repo.Update(cmd.Context(), task)
			if err != nil {
				return err
			}
			fmt.Printf("Removed due date from task %s\n", id)
			return nil
		},
	}
}

func newUpdateCommand(svc *svc.Svc) *cobra.Command {
	cmd := &TasksUpdateCommand{
		BaseCommand: common.BaseCommand{
			Svc:         svc,
			CommandName: "update",
		},
	}

	cobraCmd := &cobra.Command{
		Use:   "update <id> [flags]",
		Short: "Update task attributes",
		Long:  "Update task attributes such as tags, places, or other metadata",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cobra *cobra.Command, args []string) error {
			flags := common.ExtractFlagsFromCommand(cobra)
			return cmd.Execute(cobra.Context(), args, flags, cmd.ExecuteDirect, cmd.ExecuteViaIPC)
		},
	}

	cobraCmd.Flags().StringSliceVarP(&cmd.tags, "tags", "t", nil, "update tags")
	cobraCmd.Flags().StringSliceVarP(&cmd.places, "places", "p", nil, "update places")
	cobraCmd.Flags().StringVarP(&cmd.due, "due", "d", "", "set due date (format: YYYY-MM-DD)")
	cobraCmd.Flags().StringVar(&cmd.text, "text", "", "update task text content")

	cmd.CobraCmd = cobraCmd
	return cobraCmd
}

func (c *TasksUpdateCommand) ExecuteDirect(ctx context.Context, args []string, flags map[string]interface{}) error {
	log.Printf("Executing update command directly")
	id := args[0]
	task, err := c.Svc.Repo.GetTaskByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error fetching task: %w", err)
	}

	if c.tags != nil {
		task.Tags = c.tags
	}

	if c.places != nil {
		task.Places = c.places
	}

	if c.due != "" {
		dueDate, err := time.Parse("2006-01-02", c.due)
		if err != nil {
			return fmt.Errorf("invalid due date format: %w", err)
		}
		task.DueDate = &dueDate
	}

	if c.text != "" {
		task.Content = c.text
	}

	if err := c.Svc.Repo.Update(ctx, task); err != nil {
		return fmt.Errorf("error updating task: %w", err)
	}

	fmt.Println("Task updated successfully.")
	return nil
}

func (c *TasksUpdateCommand) ExecuteViaIPC(args []string, flags map[string]interface{}) error {
	log.Printf("Executing update command via IPC")
	id := args[0]

	params := bkg.UpdateParams{
		ID:     id,
		Tags:   c.tags,
		Places: c.places,
		Due:    c.due,
		Text:   c.text,
	}

	resp, err := bkg.SendCommand("update", params)
	if err != nil {
		return fmt.Errorf("error communicating with daemon: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("daemon returned error: %s", resp.Error)
	}

	fmt.Println("Task updated successfully via IPC.")
	return nil
}

func NewTasksCommand(svc *svc.Svc) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tasks",
		Short:   "Manage tasks",
		Aliases: []string{"t"},
	}

	cmd.AddCommand(
		newListCommand(svc),
		newStatusCommand(svc),
		newUpdateCommand(svc),
		newTextCommand(svc),
		newTagCommand(svc),
		newPlaceCommand(svc),
		newDateCommand(svc),
	)

	return cmd
}
