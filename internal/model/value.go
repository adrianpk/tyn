package model

// Common date format constants
const (
	DateTimeFormat = "2006-01-02 15:04:05"
)

type TypeVal string
type StatusVal string

type typeVal struct {
	Note  string
	Task  string
	Link  string
	Draft string
}

var Type = typeVal{
	Note:  "note",
	Task:  "task",
	Link:  "link",
	Draft: "draft",
}

func (t typeVal) Values() []string {
	return []string{
		t.Note,
		t.Task,
		t.Link,
		t.Draft,
	}
}

func (t typeVal) Validate(v string) bool {
	for _, t := range t.Values() {
		if t == v {
			return true
		}
	}

	return false
}

func (t typeVal) Label(v string) string {
	switch v {
	case t.Note:
		return "Note"
	case t.Task:
		return "Task"
	case t.Link:
		return "Link"
	case t.Draft:
		return "Draft"
	default:
		return v
	}
}

// NOTE: For simplicity, these status values and their cycling sequence are currently fixed.
// In the future, they will be customizable, allowing users to edit both the status values
// and their cycling order according to their workflow needs.
type statusVal struct {
	Todo       string
	Ready      string
	InProgress string
	Blocked    string
	OnHold     string
	Review     string
	Done       string
	Canceled   string
	Waiting    string
}

var Status = statusVal{
	Todo:       "todo",
	Ready:      "ready",
	InProgress: "wip",
	Blocked:    "blocked",
	OnHold:     "on-hold",
	Review:     "review",
	Done:       "done",
	Canceled:   "canceled",
	Waiting:    "waiting",
}

func (s statusVal) Values() []string {
	return []string{
		s.Todo,
		s.Ready,
		s.InProgress,
		s.Blocked,
		s.OnHold,
		s.Review,
		s.Done,
		s.Canceled,
		s.Waiting,
	}
}

func (s statusVal) Validate(v string) bool {
	for _, s := range s.Values() {
		if s == v {
			return true
		}
	}

	return false
}

func (s statusVal) Label(v string) string {
	switch v {
	case s.Todo:
		return "Todo"
	case s.Ready:
		return "Ready"
	case s.InProgress:
		return "In Progress"
	case s.Blocked:
		return "Blocked"
	case s.OnHold:
		return "On Hold"
	case s.Review:
		return "Review"
	case s.Done:
		return "Done"
	case s.Canceled:
		return "Canceled"
	case s.Waiting:
		return "Waiting"
	default:
		return v
	}
}

var StatusCycle = []string{
	Status.Todo,
	Status.Ready,
	Status.InProgress,
	Status.Blocked,
	Status.OnHold,
	Status.Review,
	Status.Done,
	Status.Canceled,
	Status.Waiting,
}

func ValidStatus(status string) bool {
	for _, s := range StatusCycle {
		if s == status {
			return true
		}
	}
	return false
}

func NextStatus(currentStatus string) string {
	for i, status := range StatusCycle {
		if status == currentStatus {
			if i == len(StatusCycle)-1 {
				return StatusCycle[0]
			}
			return StatusCycle[i+1]
		}
	}
	return StatusCycle[0]
}

func PreviousStatus(currentStatus string) string {
	for i, status := range StatusCycle {
		if status == currentStatus {
			if i == 0 {
				return StatusCycle[len(StatusCycle)-1]
			}
			return StatusCycle[i-1]
		}
	}
	return StatusCycle[0]
}
