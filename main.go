package main

import (
	"bytes"
	"log"
	"os"

	"github.com/adrianpk/tyn/internal/capture"
)

func main() {
	cmd := capture.NewCommand()

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
