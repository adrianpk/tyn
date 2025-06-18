package journal

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrianpk/tyn/internal/model"
)

const (
	// JournalInterval is the interval between journal updates in minutes (hardcoded for now)
	JournalInterval = 1 * time.Minute

	// JournalBasePath is the base directory for journal files (hardcoded for now)
	JournalBasePath = "~/Documents/tyn/journal"
)

// Generator is the journal generator service
type Generator struct {
	repo JournalRepo
}

// JournalRepo defines the repository interface needed by the journal generator
type JournalRepo interface {
	// GetNodesByDay returns all nodes created on a specific day
	GetNodesByDay(day time.Time) ([]model.Node, error)
}

// New creates a new journal generator
func New(repo JournalRepo) *Generator {
	return &Generator{
		repo: repo,
	}
}

// GenerateDaily generates the daily journal for the current day
func (g *Generator) GenerateDaily() error {
	// Use current day for journal generation
	today := time.Now()

	// Get nodes for today
	nodes, err := g.repo.GetNodesByDay(today)
	if err != nil {
		return fmt.Errorf("error fetching today's nodes: %w", err)
	}

	// Organize nodes by type
	tasks := []model.Node{}
	notes := []model.Node{}
	links := []model.Node{}

	for _, node := range nodes {
		switch node.Type {
		case model.Type.Task:
			tasks = append(tasks, node)
		case model.Type.Note:
			notes = append(notes, node)
		case model.Type.Link:
			links = append(links, node)
		}
	}

	content := genMarkdownContent(today, tasks, notes, links)

	err = saveJournal(today, content)
	if err != nil {
		return fmt.Errorf("error saving journal: %w", err)
	}

	return nil
}

func genMarkdownContent(day time.Time, tasks, notes, links []model.Node) string {
	header := fmt.Sprintf("# %s\n\n", day.Format("060102"))

	tasksSection := "## Tasks\n\n"
	if len(tasks) > 0 {
		for _, task := range tasks {
			status := " "
			if task.Status == "done" {
				status = "x"
			}

			taskLine := fmt.Sprintf("- [%s] %s", status, task.Content)

			if task.Status != "" {
				taskLine += fmt.Sprintf(" `:%s`", task.Status)
			}

			tasksSection += taskLine + "\n"
		}
	} else {
		tasksSection += "No tasks recorded today.\n"
	}
	tasksSection += "\n"

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

func saveJournal(day time.Time, content string) error {
	expandedPath := JournalBasePath
	if len(expandedPath) > 0 && expandedPath[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %w", err)
		}
		expandedPath = filepath.Join(home, expandedPath[1:])
	}

	year := day.Format("2006")
	month := day.Format("01")
	dirPath := filepath.Join(expandedPath, year, month)

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("error creating journal directories: %w", err)
	}

	fileName := day.Format("20060102") + ".md"
	filePath := filepath.Join(dirPath, fileName)

	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error writing journal file: %w", err)
	}

	return nil
}
