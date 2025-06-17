package main

import (
	"log"
	"os"

	"github.com/adrianpk/tyn/internal/capture"
	"github.com/adrianpk/tyn/internal/list"
	"github.com/adrianpk/tyn/internal/repo/sqlite"
	"github.com/adrianpk/tyn/internal/svc"
	"github.com/spf13/cobra"
)

func main() {
	r, err := sqlite.NewTynRepo()
	if err != nil {
		log.Fatalf("repo db error: %v", err)
	}

	s := svc.New(r)

	rootCmd := &cobra.Command{Use: "tyn"}
	rootCmd.AddCommand(capture.NewCommand(s))
	rootCmd.AddCommand(list.NewCommand(s))

	if len(os.Args) > 1 {
		rootCmd.SetArgs(os.Args[1:])
	}

	if err := rootCmd.Execute(); err != nil {
		log.Printf("error: %v", err)
		return
	}
}
