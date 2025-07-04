package root

import (
	"github.com/adrianpk/tyn/internal/bkg"
	"github.com/adrianpk/tyn/internal/command/capture"
	"github.com/adrianpk/tyn/internal/command/list"
	"github.com/adrianpk/tyn/internal/command/tasks"
	"github.com/adrianpk/tyn/internal/config"
	"github.com/adrianpk/tyn/internal/svc"
	"github.com/spf13/cobra"
)

func NewCommand(s *svc.Svc, cfg *config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "tn",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == "serve" {
				return nil
			}
			return bkg.EnsureDaemon()
		},
	}

	rootCmd.AddCommand(capture.NewCommand(s))
	rootCmd.AddCommand(list.NewCommand(s))
	rootCmd.AddCommand(tasks.NewCommand(s))
	rootCmd.AddCommand(newServeCommand(cfg))

	return rootCmd
}

func newServeCommand(config *config.Config) *cobra.Command {
	serveCmd := &cobra.Command{
		Use:    "serve",
		Short:  "Start the background service",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			isDaemon, _ := cmd.Flags().GetBool("daemon")
			bkg.ServeLoop(isDaemon, config)
		},
	}
	serveCmd.Flags().Bool("daemon", false, "Run as a daemon process")

	return serveCmd
}
