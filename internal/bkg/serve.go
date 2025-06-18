package bkg

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adrianpk/tyn/internal/journal"
	"github.com/adrianpk/tyn/internal/notify"
	"github.com/adrianpk/tyn/internal/repo/sqlite"
	"github.com/adrianpk/tyn/internal/svc"
)

const DefaultPollInterval = 30 * time.Second

type Service struct {
	svc                   *svc.Svc
	journalGenerator      *journal.Generator
	lastJournalGen        time.Time
	lastNotificationCheck time.Time
	notifiedTaskIDs       map[string]bool
}

// ServeLoop is the main daemon loop that processes pending nodes
func ServeLoop(isDaemon bool) {
	if isDaemon {
		logFile, err := logFilePath()
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
		notifiedTaskIDs:  make(map[string]bool),
	}

	err = HandleConnections(service.handleMessage)
	if err != nil {
		log.Fatalf("Error starting IPC handler: %v", err)
	}

	log.Println("IPC server started successfully")

	log.Println("Initial journal generation on startup...")
	err = service.journalGenerator.GenerateDaily()
	if err != nil {
		log.Printf("Error generating initial journal: %v\n", err)
	} else {
		log.Println("Initial journal generated successfully")
		service.lastJournalGen = time.Now()
	}

	for {
		err = service.processPendingNodes()
		if err != nil {
			log.Printf("Error processing pending nodes: %v\n", err)
		}

		err = service.checkOverdueTasks()
		if err != nil {
			log.Printf("Error checking overdue tasks: %v\n", err)
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

// checkOverdueTasks checks for tasks with past due dates and sends notifications
func (s *Service) checkOverdueTasks() error {
	// NOTE: Check for notifications once per hour
	if !s.lastNotificationCheck.IsZero() && time.Since(s.lastNotificationCheck) < 1*time.Hour {
		log.Println("Skipping overdue task check, last check was less than 1 hour ago")
		return nil
	}

	log.Println("Checking for overdue tasks...")

	ctx := context.Background()
	overdueTasks, err := s.svc.GetOverdueTasks(ctx)
	if err != nil {
		return fmt.Errorf("error getting overdue tasks: %w", err)
	}

	if len(overdueTasks) == 0 {
		log.Println("No overdue tasks found")
		return nil
	}

	log.Printf("Found %d overdue tasks", len(overdueTasks))

	for _, task := range overdueTasks {
		// NOTE: Skip if we've already notified about this task in this daemon session
		if s.notifiedTaskIDs[task.ID] {
			log.Printf("Task %s was already notified in this session, skipping", task.ID)
			continue
		}

		isNew, err := s.svc.NotifyOverdueTask(ctx, task.ID)
		if err != nil {
			log.Printf("Error recording notification for task %s: %v", task.ID, err)
			continue
		}

		s.notifiedTaskIDs[task.ID] = true

		dueDateStr := "unknown"
		if task.DueDate != nil {
			localDueDate := task.DueDate.In(time.Local)
			dueDateStr = localDueDate.Format("2006-01-02 15:04")
		}

		message := fmt.Sprintf("Due date: %s - %s", dueDateStr, task.Content)
		if isNew {
			err = notify.NotifyDueDate(task.Content, message)
		} else {
			err = notify.NotifyDueDateReminder(task.Content, message)
		}

		if err != nil {
			log.Printf("Error sending notification for task %s: %v", task.ID, err)
		} else {
			log.Printf("Successfully sent notification for task %s", task.ID)
		}
	}

	s.lastNotificationCheck = time.Now()

	return nil
}
