package capture

import (
	"log"
	"strings"

	"github.com/adrianpk/tyn/internal/svc"
	"github.com/spf13/cobra"
)

func NewCommand(svc *svc.Svc) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "capture",
		Aliases: []string{"c"},
		Short:   "capture a new node",
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := strings.Join(args, " ")
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
		},
	}
	return cmd
}
