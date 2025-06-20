package capture

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/adrianpk/tyn/internal/bkg"
	"github.com/adrianpk/tyn/internal/model"
	"github.com/adrianpk/tyn/internal/svc"
	"github.com/spf13/cobra"
)

func NewCommand(svc *svc.Svc) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "capture",
		Aliases: []string{"cap", "c"},
		Short:   "capture a new node",
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := strings.Join(args, " ")

			// Direct svc use if available
			if svc != nil {
				node, err := svc.Parser(input)
				if err != nil {
					return err
				}

				node.GenID()

				log.Printf("Captured node: %s", node.ID)

				err = svc.Repo.Create(cmd.Context(), node)
				if err != nil {
					return err
				}

				log.Printf("%+v", node)
				return nil
			}

			// IPC otherwise
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

			// Convert DueDate to local timezone before logging
			if node.DueDate != nil {
				localDueDate := node.DueDate.In(time.Local)
				node.DueDate = &localDueDate
			}

			log.Printf("%+v", node)
			return nil
		},
	}
	return cmd
}
