package svc

import (
	"errors"
	"strings"
	"time"

	"github.com/adrianpk/tyn/internal/model"
)

const (
	tagDelim    = "#"
	placeDelim  = "@"
	statusDelim = ":"
	dateDelim   = "^"
	dateFmt     = "2006-01-02"
	httpScheme  = "http://"
	httpsScheme = "https://"
)

func Parse(input string) (model.Node, error) {
	item := model.Node{Date: time.Now()}
	tokens := strings.Fields(input)
	var content []string

	for _, tok := range tokens {
		switch {
		case strings.HasPrefix(tok, tagDelim) && len(tok) > 1:
			item.Tags = append(item.Tags, strings.ToLower(tok[1:]))

		case strings.HasPrefix(tok, placeDelim) && len(tok) > 1:
			item.Places = append(item.Places, tok[1:])

		case strings.HasPrefix(tok, statusDelim) && len(tok) > 1 && item.Status == "":
			item.Status = tok[1:]

		case strings.HasPrefix(tok, dateDelim) && len(tok) > 1 && item.OverrideDate == nil:
			t, err := time.Parse(dateFmt, tok[1:])
			if err != nil {
				return item, errors.New("invalid date format: " + tok[1:])
			}
			item.OverrideDate = &t

		case (strings.HasPrefix(tok, httpScheme) || strings.HasPrefix(tok, httpsScheme)) && item.Link == "":
			item.Link = tok

		default:
			content = append(content, tok)
		}
	}

	item.Content = strings.Join(content, " ")

	switch {
	case item.Status != "":
		item.Type = model.Type.Task
	case item.Link != "":
		item.Type = model.Type.Link
	default:
		item.Type = model.Type.Note
	}

	return item, nil
}
