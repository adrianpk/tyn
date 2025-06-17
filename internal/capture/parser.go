package capture

import (
	"errors"
	"strings"
	"time"
)

func Parse(input string) (Node, error) {
	item := Node{Date: time.Now()}
	tokens := strings.Fields(input)
	var content []string

	for _, tok := range tokens {
		switch {
		case strings.HasPrefix(tok, "#") && len(tok) > 1:
			item.Tags = append(item.Tags, strings.ToLower(tok[1:]))
		case strings.HasPrefix(tok, "@") && len(tok) > 1:
			item.Places = append(item.Places, tok[1:])
		case strings.HasPrefix(tok, ":") && len(tok) > 1 && item.Status == "":
			item.Status = tok[1:]
		case strings.HasPrefix(tok, "^") && len(tok) > 1 && item.OverrideDate == nil:
			t, err := time.Parse("2006-01-02", tok[1:])
			if err != nil {
				return item, errors.New("invalid date format: " + tok[1:])
			}
			item.OverrideDate = &t
		case (strings.HasPrefix(tok, "http://") || strings.HasPrefix(tok, "https://")) && item.Link == "":
			item.Link = tok
		default:
			content = append(content, tok)
		}
	}

	item.Content = strings.Join(content, " ")

	switch {
	case item.Status != "":
		item.Type = Type.Task
	case item.Link != "":
		item.Type = Type.Link
	default:
		item.Type = Type.Note
	}

	if !Type.Validate(item.Type) {
		return item, errors.New("invalid item type: " + item.Type)
	}

	return item, nil
}
