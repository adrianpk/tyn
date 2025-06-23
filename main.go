package main

import (
	"log"
	"os"

	"github.com/adrianpk/tyn/internal/command/root"
	"github.com/adrianpk/tyn/internal/config"
	"github.com/adrianpk/tyn/internal/repo/sqlite"
	"github.com/adrianpk/tyn/internal/svc"
)

func main() {
	cfg := config.Load()

	var s *svc.Svc

	hasArgs := len(os.Args) > 1

	if hasArgs && os.Args[1] == "serve" {
		s = setup(cfg)
	}

	rootCmd := root.NewCommand(s, cfg)

	if hasArgs {
		rootCmd.SetArgs(os.Args[1:])
	}

	err := rootCmd.Execute()
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
}

func setup(config *config.Config) *svc.Svc {
	repo, err := sqlite.NewTynRepo(config)
	if err != nil {
		log.Fatalf("repo db error: %v", err)
	}

	return svc.New(repo, config)
}
