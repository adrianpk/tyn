package journal

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/adrianpk/tyn/internal/model"
)

const (
	JournalInterval = 1 * time.Minute
	JournalBasePath = "~/Documents/tyn/journal"
)

type Generator struct {
	repo JournalRepo
}

type JournalRepo interface {
	GetNodesByDay(day time.Time) ([]model.Node, error)
	GetAllTasks(ctx context.Context) ([]model.Node, error)
	GetNotesAndLinksByDay(day time.Time) ([]model.Node, error)
}

func New(repo JournalRepo) *Generator {
	return &Generator{
		repo: repo,
	}
}

func (g *Generator) GenerateDaily() error {
	today := time.Now()
	log.Println("Journal: Starting daily journal generation...")

	// Get all tasks (regardless of date)
	ctx := context.Background()
	allTasks, err := g.repo.GetAllTasks(ctx)
	if err != nil {
		return fmt.Errorf("error fetching all tasks: %w", err)
	}
	log.Printf("Journal: Found %d total tasks", len(allTasks))

	// Get today's notes and links
	notesAndLinks, err := g.repo.GetNotesAndLinksByDay(today)
	if err != nil {
		return fmt.Errorf("error fetching today's notes and links: %w", err)
	}
	log.Printf("Journal: Found %d notes and links for today", len(notesAndLinks))

	notes := []model.Node{}
	links := []model.Node{}

	for _, node := range notesAndLinks {
		switch node.Type {
		case model.Type.Note:
			notes = append(notes, node)
		case model.Type.Link:
			links = append(links, node)
		}
	}
	log.Printf("Journal: Filtered %d notes and %d links for today", len(notes), len(links))

	content := genMarkdownContent(today, allTasks, notes, links)

	journalPath, err := saveJournal(today, content)
	if err != nil {
		return fmt.Errorf("error saving journal: %w", err)
	}
	log.Printf("Journal: Successfully saved journal to %s", journalPath)

	return nil
}

func genMarkdownContent(day time.Time, tasks, notes, links []model.Node) string {
	header := fmt.Sprintf("# %s\n\n", day.Format("060102"))

	tasksSection := "## Tasks\n\n"

	log.Println("Journal: StatusCycle values:")
	for i, status := range model.StatusCycle {
		log.Printf("Journal: StatusCycle[%d] = %q", i, status)
	}

	tasksByStatus := make(map[string][]model.Node)
	for _, task := range tasks {
		status := task.Status
		if status == "" {
			status = model.Status.Todo
		}
		tasksByStatus[status] = append(tasksByStatus[status], task)
	}

	log.Println("Journal: tasksByStatus map contents:")
	for status, tasks := range tasksByStatus {
		log.Printf("Journal: Status %q has %d tasks", status, len(tasks))
	}

	if len(tasks) > 0 {
		for _, statusValue := range model.StatusCycle {
			statusTasks, exists := tasksByStatus[statusValue]
			if !exists || len(statusTasks) == 0 {
				log.Printf("Journal: Skipping section for status %q (exists=%v, len=%d)",
					statusValue, exists, len(statusTasks))
				continue
			}

			log.Printf("Journal: Adding section for status %q with %d tasks",
				statusValue, len(statusTasks))
			tasksSection += fmt.Sprintf("### %s\n\n", model.Status.Label(statusValue))

			for _, task := range statusTasks {
				checkMark := " "
				if task.Status == model.Status.Done {
					checkMark = "x"
				}

				taskLine := fmt.Sprintf("- [%s] %s", checkMark, task.Content)
				tasksSection += taskLine + "\n"
			}
			tasksSection += "\n"
		}
	} else {
		tasksSection += "No tasks found.\n\n"
	}

	notesSection := "## Notes\n\n"
	if len(notes) > 0 {
		for _, note := range notes {
			notesSection += fmt.Sprintf("- %s\n", note.Content)
		}
	} else {
		notesSection += "No notes recorded today.\n"
	}
	notesSection += "\n"

	linksSection := "## Links\n\n"
	if len(links) > 0 {
		for _, link := range links {
			linksSection += fmt.Sprintf("- [%s](%s)\n", link.Content, link.Link)
		}
	} else {
		linksSection += "No links recorded today.\n"
	}

	return header + tasksSection + notesSection + linksSection
}

func saveJournal(day time.Time, content string) (string, error) {
	expandedPath := JournalBasePath
	if len(expandedPath) > 0 && expandedPath[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("error getting home directory: %w", err)
		}
		expandedPath = filepath.Join(home, expandedPath[1:])
	}

	year := day.Format("2006")
	month := day.Format("01")
	dirPath := filepath.Join(expandedPath, year, month)

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return "", fmt.Errorf("error creating journal directories: %w", err)
	}

	fileName := day.Format("20060102") + ".md"
	filePath := filepath.Join(dirPath, fileName)

	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("error writing journal file: %w", err)
	}

	return filePath, nil
}
