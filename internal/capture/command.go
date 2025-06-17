package capture

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

func NewCommand(repo Repo) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "capture",
		Aliases: []string{"c"},
		Short:   "capture a new node",
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := strings.Join(args, " ")
			node, err := NewNode(input)
			if err != nil {
				return err
			}

			// WIP: The persistence mechanism will remain, but will be implemented in a more polished way in future versions.
			err = repo.Create(cmd.Context(), node)
			if err != nil {
				return err
			}

			log.Printf("%+v", node)
			return nil
		},
	}
	return cmd
}
