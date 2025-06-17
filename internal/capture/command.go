package capture

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "capture",
		Aliases: []string{"c"},
		Short:   "capture a new node",
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			input := strings.Join(args, " ")
			node, err := Parse(input)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%+v", node)
			return nil
		},
	}
	return cmd
}
