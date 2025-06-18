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
	statusPattern = regexp.MustCompile(`:(\w+)`)
	datePattern   = regexp.MustCompile(`\^([\d\-T:_]+)`)
	urlPattern    = regexp.MustCompile(`https?://[^\s]+`)
)

func Parse(input string) (model.Node, error) {
	log.Printf("Parsing input: %s", input)

	node := model.Node{
		Date: time.Now(),
	}
	node.GenID()

	urls := urlPattern.FindAllString(input, -1)
	if len(urls) > 0 {
		node.Link = urls[0]
		if strings.TrimSpace(strings.ReplaceAll(input, node.Link, "")) == "" && node.Type == "" {
			node.Type = model.Type.Link
		}
	}
	input = urlPattern.ReplaceAllString(input, "")

	tagMatches := tagPattern.FindAllStringSubmatch(input, -1)
	for _, match := range tagMatches {
		if len(match) > 1 {
			node.Tags = append(node.Tags, match[1])
		}
	}
	input = tagPattern.ReplaceAllString(input, "")

	placeMatches := placePattern.FindAllStringSubmatch(input, -1)
	for _, match := range placeMatches {
		if len(match) > 1 {
			node.Places = append(node.Places, match[1])
		}
	}
	input = placePattern.ReplaceAllString(input, "")

	statusMatch := statusPattern.FindStringSubmatch(input)
	if len(statusMatch) > 1 {
		node.Status = statusMatch[1]
		if node.Type == "" {
			node.Type = model.Type.Task
		}
	}
	input = statusPattern.ReplaceAllString(input, "")

	dateMatch := datePattern.FindStringSubmatch(input)
	if len(dateMatch) > 1 {
		dateStr := strings.TrimSpace(dateMatch[1])
		log.Printf("Captured date string: %s", dateStr)

		var dueDate time.Time
		var err error

		formats := []string{
			"2006-01-02-15-04-05", // Format with all hyphens (from makefile)
			"2006-01-02T15:04:05", // ISO8601
			"2006-01-02_15:04:05", // Underscore format
			"2006-01-02 15:04:05", // Space format
			"2006-01-02",          // Date only
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
