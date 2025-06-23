package svc

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/adrianpk/tyn/internal/model"
)

var (
	tagPattern    = regexp.MustCompile(`#(\w+)`)
	placePattern  = regexp.MustCompile(`@(\w+)`)
	statusPattern = regexp.MustCompile(`:([a-zA-Z0-9-]+)`)
	datePattern   = regexp.MustCompile(`\^([\d\-T:_]+)`)
	urlPattern    = regexp.MustCompile(`https?://[^\s]+`)
	draftPattern  = regexp.MustCompile(`\+([a-zA-Z0-9-_]+)`)
)

func Parse(input string) (model.Node, error) {
	log.Printf("Parsing input: %s", input)

	node := model.Node{
		Date: time.Now(),
	}
	node.GenID()

	draftMatch := draftPattern.FindStringSubmatch(input)
	if len(draftMatch) > 1 {
		draft := draftMatch[1]
		node.Draft = draft
		node.Type = model.Type.Draft
		input = draftPattern.ReplaceAllString(input, "")
	}

	// Process tags - available for all node types including drafts
	tagMatches := tagPattern.FindAllStringSubmatch(input, -1)
	for _, match := range tagMatches {
		if len(match) > 1 {
			node.Tags = append(node.Tags, match[1])
		}
	}
	input = tagPattern.ReplaceAllString(input, "")

	// Process places - also available for all node types
	placeMatches := placePattern.FindAllStringSubmatch(input, -1)
	for _, match := range placeMatches {
		if len(match) > 1 {
			node.Places = append(node.Places, match[1])
		}
	}
	input = placePattern.ReplaceAllString(input, "")

	// Process URLs - for drafts, we'll still extract and store links, but not change the type
	urls := urlPattern.FindAllString(input, -1)
	if len(urls) > 0 {
		node.Link = urls[0]
		if strings.TrimSpace(strings.ReplaceAll(input, node.Link, "")) == "" && node.Type == "" {
			node.Type = model.Type.Link
		}
	}
	input = urlPattern.ReplaceAllString(input, "")

	// Process status - for drafts, we can still have a status but won't change the type
	statusMatch := statusPattern.FindStringSubmatch(input)
	if len(statusMatch) > 1 {
		node.Status = statusMatch[1]
		if node.Type == "" {
			node.Type = model.Type.Task
		}
	}
	input = statusPattern.ReplaceAllString(input, "")

	// Process due date
	dateMatch := datePattern.FindStringSubmatch(input)
	if len(dateMatch) > 1 {
		dateStr := strings.TrimSpace(dateMatch[1])
		log.Printf("Captured date string: %s", dateStr)

		var dueDate time.Time
		var err error

		formats := []string{
			"2006-01-02-15-04-05",
			"2006-01-02T15:04:05",
			"2006-01-02_15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}

		parsed := false
		for _, format := range formats {
			dueDate, err = time.ParseInLocation(format, dateStr, time.Local)
			if err == nil {
				parsed = true
				break
			}
		}

		if !parsed {
			log.Printf("Failed to parse date with any format: %s", dateStr)
			return model.Node{}, fmt.Errorf("invalid due date format: unable to parse %s", dateStr)
		}

		log.Printf("Parsed due date (in local timezone): %v", dueDate)
		node.DueDate = &dueDate
		log.Printf("Set node.DueDate = %v", *node.DueDate)
	}
	input = datePattern.ReplaceAllString(input, "")

	node.Content = strings.TrimSpace(input)

	// Set the default type if not already set (draft type will be preserved)
	if node.Type == "" {
		if node.Status != "" {
			node.Type = model.Type.Task
		} else if node.Link != "" {
			node.Type = model.Type.Link
		} else {
			node.Type = model.Type.Note
		}
	}

	return node, nil
}
