package utils

import (
	"fmt"
	"git-project-management/internal/types"
	"regexp"
	"strconv"
	"strings"
)

func ParseCommitMessage(message string) (types.CommitMessage, string) {
	// Split the commit message into lines
	lines := strings.Split(message, "\n")

	// Initialize the commit message
	var commit types.CommitMessage

	// Extract the title (first line)
	if len(lines) > 0 {
		commit.Title = strings.TrimSpace(lines[0])
	}

	// Extract the description (remaining lines, excluding the custom pattern)
	if len(lines) > 1 {
		descriptionLines := lines[1 : len(lines)-1] // Exclude the last line (custom pattern)
		commit.Description = strings.TrimSpace(strings.Join(descriptionLines, "\n"))
	}

	// Extract the custom pattern (last line)
	customPattern := ""
	if len(lines) > 0 {
		customPattern = strings.TrimSpace(lines[len(lines)-1])
	}

	return commit, customPattern
}

func IsValidCommitLine(line string) types.CommitValidationResult {
	// Regex to match the general structure
	pattern := `^\[([^]]+)\](.*)$`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(line)
	if len(matches) != 3 {
		return types.CommitValidationResult{
			IsValid: false,
			Message: "Commit message must start with '[' and end with ']'.",
		}
	}

	params := matches[1]                  // Parameters inside []
	text := strings.TrimSpace(matches[2]) // Text after []

	// Initialize the validation result
	result := types.CommitValidationResult{
		IsValid:    true,
		Message:    "Commit message is valid.",
		IsNoCommit: strings.Contains(params, "nc"),
		HasTimelog: strings.Contains(params, "l"),
	}

	// Extract project ID if provided (e.g., [p-123,t] task title)
	projectIDRegex := regexp.MustCompile(`p-(\d+)`)
	projectIDMatch := projectIDRegex.FindStringSubmatch(params)
	if len(projectIDMatch) > 1 {
		projectID, err := parseInt64(projectIDMatch[1])
		if err != nil {
			return types.CommitValidationResult{
				IsValid: false,
				Message: fmt.Sprintf("Invalid project ID: %v", err),
			}
		}
		result.ProjectID = projectID
	}

	// Check if t or t-<task-id> is present
	hasT := strings.Contains(params, "t,") || strings.HasSuffix(params, "t")
	hasTWithID := strings.Contains(params, "t-")

	// 1. Check for mutual exclusivity of t and t-<task-id>
	if hasT && hasTWithID {
		return types.CommitValidationResult{
			IsValid: false,
			Message: "'t' and 't-<task-id>' cannot coexist.",
		}
	}

	// 2. Check if project ID and task ID coexist
	if result.ProjectID != 0 && hasTWithID {
		return types.CommitValidationResult{
			IsValid: false,
			Message: "'p-<project-id>' and 't-<task-id>' cannot coexist.",
		}
	}

	// 3. Extract task type and ID
	if hasTWithID {
		result.IsNewTask = false
		// Extract task ID
		taskIDRegex := regexp.MustCompile(`t-(\d+)`)
		taskIDMatch := taskIDRegex.FindStringSubmatch(params)
		if len(taskIDMatch) > 1 {
			taskID, err := parseInt64(taskIDMatch[1])
			if err != nil {
				return types.CommitValidationResult{
					IsValid: false,
					Message: fmt.Sprintf("Invalid task ID: %v", err),
				}
			}
			result.TaskID = taskID
		}
	} else if hasT {
		result.IsNewTask = true
	}

	// 4. Task title rules
	if hasT && !hasTWithID {
		// t is present, so a task title is required
		if text == "" {
			return types.CommitValidationResult{
				IsValid: false,
				Message: "Task title is required for 't'.",
			}
		}
		result.TaskTitle = text
	} else if hasTWithID {
		// t-<task-id> is present, so a task title is not allowed
		if text != "" {
			return types.CommitValidationResult{
				IsValid: false,
				Message: "Task title is not allowed for 't-<task-id>'.",
			}
		}
	}

	// 5. Timelog rules
	if result.HasTimelog {
		// Split text by "|" to separate task title and timelog
		parts := strings.Split(text, "|")
		if len(parts) > 1 {
			timelog := strings.TrimSpace(parts[1])
			if isValidTimelog(timelog) {
				result.Timelog = timelog
			} else {
				result.Message = fmt.Sprintf("'%s' is not a valid timelog. Treating it as a task title.", timelog)
			}
		}
	}

	return result
}

// parseInt64 converts a string to int64
func parseInt64(s string) (int64, error) {
	var result int64
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", s)
	}
	return result, nil
}

func isValidTimelog(timelog string) bool {
	// Regex to match timelog format (e.g., 23h 34m, 2.5h, 150m)
	pattern := `^(\d+h \d+m|\d+(\.\d+)?h|\d+m)$`
	matched, _ := regexp.MatchString(pattern, timelog)
	return matched
}

func ConvertTimelogToMinutes(timelog string) (int, error) {
	// Regex to match hours and minutes in the timelog
	pattern := `^(?:(\d+)h)?\s*(?:(\d+)m)?$`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(timelog)

	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid timelog format: %s", timelog)
	}

	// Extract hours and minutes
	hoursStr := matches[1]
	minutesStr := matches[2]

	// Convert hours to minutes
	hours := int(0)
	if hoursStr != "" {
		var err error
		hours, err = strconv.Atoi(hoursStr)
		if err != nil {
			return 0, fmt.Errorf("invalid hours in timelog: %s", hoursStr)
		}
	}

	// Convert minutes to minutes
	minutes := int(0)
	if minutesStr != "" {
		var err error
		minutes, err = strconv.Atoi(minutesStr)
		if err != nil {
			return 0, fmt.Errorf("invalid minutes in timelog: %s", minutesStr)
		}
	}

	// Calculate total minutes
	totalMinutes := hours*60 + minutes
	return totalMinutes, nil
}
