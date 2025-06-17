package main

import (
	"bytes"
	"log"
	"os"

	"github.com/adrianpk/tyn/internal/capture"
	"github.com/adrianpk/tyn/internal/repo/sqlite"
)

func main() {
	dsn := os.Getenv("TYN_SQLITE_DSN")
	if dsn == "" {
		dsn = "tyn.db"
	}

	r, err := sqlite.NewTynRepo(dsn)
	if err != nil {
		log.Fatalf("repo db error: %v", err)
	}

	cmd := capture.NewCommand(r)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	if len(os.Args) > 1 {
		cmd.SetArgs(os.Args[1:])
	}

	if err := cmd.Execute(); err != nil {
		log.Printf("error: %v", err)
		return
	}

	log.Printf("raw output: %s", buf.String())
}
