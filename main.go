package main

import (
	"log"
	"os"

	"github.com/adrianpk/tyn/internal/repo/sqlite"
	"github.com/adrianpk/tyn/internal/root"
	"github.com/adrianpk/tyn/internal/svc"
)

func main() {
	var s *svc.Svc

	hasArgs := len(os.Args) > 1

	if hasArgs && os.Args[1] == "serve" {
		s = setup()
	}

	rootCmd := root.NewCommand(s)

	if hasArgs {
		rootCmd.SetArgs(os.Args[1:])
	}

	err := rootCmd.Execute()
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
}

func setup() *svc.Svc {
	repo, err := sqlite.NewTynRepo()
	if err != nil {
		log.Fatalf("repo db error: %v", err)
	}
	return svc.New(repo)
}
