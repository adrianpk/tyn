package model

type TypeVal string

type typeVal struct {
	Note string
	Task string
	Link string
}

var Type = typeVal{
	Note: "note",
	Task: "task",
	Link: "link",
}

func (t typeVal) Values() []string {
	return []string{
		t.Note,
		t.Task,
		t.Link,
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
	default:
		return v
	}
}
