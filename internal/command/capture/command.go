package capture

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/adrianpk/tyn/internal/bkg"
	"github.com/adrianpk/tyn/internal/command/common"
	"github.com/adrianpk/tyn/internal/model"
	"github.com/adrianpk/tyn/internal/svc"
	"github.com/spf13/cobra"
)

type CaptureCommand struct {
	common.BaseCommand
}

func NewCommand(svc *svc.Svc) *cobra.Command {
	cmd := &CaptureCommand{
		BaseCommand: common.BaseCommand{
			Svc:         svc,
			CommandName: "capture",
		},
	}

	cobraCmd := &cobra.Command{
		Use:     "capture",
		Aliases: []string{"cap", "c"},
		Short:   "capture a new node",
		Args:    cobra.ArbitraryArgs,
		RunE: func(cobra *cobra.Command, args []string) error {
			flags := common.ExtractFlagsFromCommand(cobra)
			return cmd.Execute(cobra.Context(), args, flags, cmd.ExecuteDirect, cmd.ExecuteViaIPC)
		},
	}

	cmd.CobraCmd = cobraCmd
	return cobraCmd
}

func (c *CaptureCommand) ExecuteDirect(ctx context.Context, args []string, flags map[string]interface{}) error {
	input := strings.Join(args, " ")

	node, err := c.Svc.Parser(input)
	if err != nil {
		return err
	}

	node.GenID()

	log.Printf("Captured node: %s", node.ID)

	err = c.Svc.Repo.Create(ctx, node)
	if err != nil {
		return err
	}

	log.Printf("%+v", node)
	return nil
}

func (c *CaptureCommand) ExecuteViaIPC(args []string, flags map[string]interface{}) error {
	input := strings.Join(args, " ")

	params := bkg.CaptureParams{
		Text: input,
	}

	resp, err := bkg.SendCommand("capture", params)
	if err != nil {
		return fmt.Errorf("error communicating with daemon: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("daemon returned error: %s", resp.Error)
	}

	var node model.Node
	err = json.Unmarshal(resp.Data, &node)
	if err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	log.Printf("Captured node via daemon: %s", node.ID)

	if node.DueDate != nil {
		localDueDate := node.DueDate.In(time.Local)
		node.DueDate = &localDueDate
	}

	log.Printf("%+v", node)
	return nil
}
