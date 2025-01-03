package controller

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func ToActivityDTO(model ActivityEntity) ActivityDTO {
	return ActivityDTO{
		ID:          model.ID,
		TaskId:      model.TaskID,
		Title:       model.Title,
		Description: model.Description,
		Duration:    model.Duration,
		CreatedBy:   model.CreatedBy,
		CreateAt:    model.CreatedAt,
	}
}

func ParseFullCommitMessage(fullCommitMessage string) map[string]string {
	// Create a map to store the extracted values
	values := make(map[string]string)

	// Split the message into key-value pairs using " | " as the delimiter
	parts := strings.Split(fullCommitMessage, " ||| ")

	// Iterate over each part and split it into key and value
	for _, part := range parts {
		keyValue := strings.SplitN(part, ": ", 2) // Split by ": " into key and value
		if len(keyValue) == 2 {
			key := strings.TrimSpace(keyValue[0])
			value := strings.TrimSpace(keyValue[1])
			values[key] = value
		}
	}

	return values
}

type ParsedData struct {
	ActivityTitle       string
	ActivityDescription string
	TaskID              string
	TaskStatus          string
	ProjectID           string
	NewTaskTitle        string
	NewTaskDescription  string
	NewTaskStatus       string
	TimeSpentMinutes    int
}

func ParseCommitMessage(commitMessage string) (ParsedData, error) {
	var parsed ParsedData
	lines := strings.Split(commitMessage, "\n")

	// Check if the first line starts with a tag
	tagRegex := regexp.MustCompile(`^\[(task|proj|spent)-[a-zA-Z0-9]+\]`)
	if len(lines) > 0 && tagRegex.MatchString(lines[0]) {
		return parsed, errors.New("error: commit message cannot start with a tag")
	}

	// Capture the activity title (first line)
	if len(lines) > 0 {
		parsed.ActivityTitle = strings.TrimSpace(lines[0])
	}

	// Updated Activity Description Parsing
	var descriptionLines []string
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)

		// Skip lines with tags
		if strings.HasPrefix(line, "[task-") || strings.HasPrefix(line, "[proj-") || strings.HasPrefix(line, "[spent]") {
			continue
		}

		descriptionLines = append(descriptionLines, line)
	}

	// Combine lines to form the Activity Description
	parsed.ActivityDescription = strings.Join(descriptionLines, "\n")
	parsed.ActivityDescription = strings.TrimSpace(parsed.ActivityDescription)

	spentRegex := regexp.MustCompile(`\[spent\](?:\s*((?:\d+h\s*)?(?:\d+m)?)|\s*~)`)
	taskRegex := regexp.MustCompile(`\[task-([a-zA-Z0-9]+)\]\s*(\w+)?`)
	projectRegex := regexp.MustCompile(`\[proj-([a-zA-Z0-9]+)\]\s*([^|]+)?(?:\s*\|\s*([^|]+))?(?:\s*\|\s*(\w+))?`)

	validTaskStatuses := map[string]bool{"open": true, "done": true, "in-progress": true}

	// Track occurrences
	seenTags := map[string]bool{"[spent]": false, "[task]": false, "[proj]": false}

	// Process the lines
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)

		// Parse [spent]
		if spentMatches := spentRegex.FindStringSubmatch(line); spentMatches != nil {
			if seenTags["[spent]"] {
				return parsed, errors.New("error: multiple [spent] tags found")
			}
			seenTags["[spent]"] = true

			if strings.Contains(line, "~") && strings.TrimSpace(spentMatches[1]) != "" {
				return parsed, errors.New("error: [spent] cannot contain both a time format and ~")
			}

			if strings.Contains(line, "~") {
				parsed.TimeSpentMinutes = -1
			} else if spentMatches[1] != "" {
				parsed.TimeSpentMinutes = convertTimeToMinutes(spentMatches[1])
			} else {
				return parsed, errors.New("error: [spent] must be followed by either a time format or ~")
			}
		}

		// Parse [task-XX]
		if taskMatches := taskRegex.FindStringSubmatch(line); taskMatches != nil {
			if seenTags["[task]"] {
				return parsed, errors.New("error: multiple [task] tags found")
			}
			if seenTags["[proj]"] {
				return parsed, errors.New("error: [task] and [proj] tags cannot coexist in the same commit message")
			}
			seenTags["[task]"] = true

			parsed.TaskID = taskMatches[1] // Alphanumeric Task ID

			if len(taskMatches) > 2 && taskMatches[2] != "" {
				if !validTaskStatuses[taskMatches[2]] {
					return parsed, errors.New("error: invalid task status for [task]")
				}
				parsed.TaskStatus = taskMatches[2]
			}
		}

		// Parse [proj-XX]
		if projectMatches := projectRegex.FindStringSubmatch(line); projectMatches != nil {
			if seenTags["[proj]"] {
				return parsed, errors.New("error: multiple [proj] tags found")
			}
			if seenTags["[task]"] {
				return parsed, errors.New("error: [task] and [proj] tags cannot coexist in the same commit message")
			}
			seenTags["[proj]"] = true

			_, err := strconv.Atoi(projectMatches[1])
			if err != nil {
				return parsed, errors.New("error: malformed [proj-id]")
			}
			parsed.ProjectID = projectMatches[1]

			if len(projectMatches) > 2 && projectMatches[2] != "" {
				parsed.NewTaskTitle = strings.TrimSpace(projectMatches[2])
			}

			if len(projectMatches) > 3 && projectMatches[3] != "" {
				parsed.NewTaskDescription = strings.TrimSpace(projectMatches[3])
			}

			if len(projectMatches) > 4 && projectMatches[4] != "" {
				parsed.NewTaskStatus = strings.TrimSpace(projectMatches[4])
			}
		}
	}

	// Final check: ensure either [task] or [proj] exists
	if !seenTags["[task]"] && !seenTags["[proj]"] {
		return parsed, errors.New("error: commit message must contain either [task] or [proj] tag")
	}

	return parsed, nil
}

func convertTimeToMinutes(timeStr string) int {
	hoursRegex := regexp.MustCompile(`(\d+)h`)
	minutesRegex := regexp.MustCompile(`(\d+)m`)

	hours := 0
	minutes := 0

	if hoursMatch := hoursRegex.FindStringSubmatch(timeStr); hoursMatch != nil {
		h, _ := strconv.Atoi(hoursMatch[1])
		hours = h
	}

	if minutesMatch := minutesRegex.FindStringSubmatch(timeStr); minutesMatch != nil {
		m, _ := strconv.Atoi(minutesMatch[1])
		minutes = m
	}

	return hours*60 + minutes
}

func parseInt(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}
