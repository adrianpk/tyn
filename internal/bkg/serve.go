package bkg

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adrianpk/tyn/internal/journal"
	"github.com/adrianpk/tyn/internal/repo/sqlite"
	"github.com/adrianpk/tyn/internal/svc"
)

const DefaultPollInterval = 30 * time.Second

// Service is the main daemon service
type Service struct {
	svc              *svc.Svc
	journalGenerator *journal.Generator
	lastJournalGen   time.Time
}

// ServeLoop is the main daemon loop that processes pending nodes
func ServeLoop(isDaemon bool) {
	if isDaemon {
		logFile, err := getLogFilePath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting log file path: %v\n", err)
			return
		}

		logFd, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
			return
		}
		defer logFd.Close()

		log.SetOutput(logFd)
	}

	log.Println("Starting Tyn background service...")

	repo, err := sqlite.NewTynRepo()
	if err != nil {
		log.Fatalf("Error initializing repository: %v", err)
	}

	service := &Service{
		svc:              svc.New(repo),
		journalGenerator: journal.New(repo),
	}

	err = HandleConnections(service.handleMessage)
	if err != nil {
		log.Fatalf("Error starting IPC handler: %v", err)
	}

	log.Println("IPC server started successfully")

	for {
		err = service.processPendingNodes()
		if err != nil {
			log.Printf("Error processing pending nodes: %v\n", err)
		}

		if time.Since(service.lastJournalGen) >= journal.JournalInterval {
			log.Println("Generating daily journal...")
			err = service.journalGenerator.GenerateDaily()
			if err != nil {
				log.Printf("Error generating journal: %v\n", err)
			} else {
				log.Println("Journal generated successfully")
				service.lastJournalGen = time.Now()
			}
		}

		time.Sleep(DefaultPollInterval)
	}
}

func (s *Service) processPendingNodes() error {
	log.Println("Checking for pending nodes...")

	// WIP: implement the actual processing logic:
	// * Query the database for pending nodes
	// * Process them

	return nil
}

func (s *Service) handleMessage(msg Message) Response {
	log.Printf("Received message: %s\n", msg.Command)

	switch msg.Command {
	case "capture":
		return s.handleCapture(msg.Params)
	case "list":
		return s.handleList(msg.Params)
	default:
		return Response{
			Success: false,
			Error:   fmt.Sprintf("unknown command: %s", msg.Command),
		}
	}
}
