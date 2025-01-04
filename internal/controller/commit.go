package controller

import (
	"context"
	"encoding/base64"
	"fmt"
	"git-project-management/internal/types"
	"regexp"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

// ValidationResult represents the detailed result of the validation
type ValidationResult struct {
	IsValid    bool   `json:"is_valid"`     // Whether the commit message is valid
	Message    string `json:"message"`      // A message describing the validation result
	IsNoCommit bool   `json:"is_no_commit"` // Whether the commit should be ignored (nc)
	IsNewTask  bool   `json:"is_new_task"`  // Whether the commit is for a new task (t)
	TaskID     string `json:"task_id"`      // The task ID (if applicable)
	HasTimelog bool   `json:"has_timelog"`  // Whether the commit includes a timelog (l)
	Timelog    string `json:"timelog"`      // The timelog (if applicable)
	TaskTitle  string `json:"task_title"`   // The task title (if applicable)
	ProjectID  string `json:"project_id"`   // The project ID (if provided)
}

type CommitController struct{}

func NewCommitController() CommitController {
	return CommitController{}
}

func Commit(ctx context.Context, req *types.CommitRequest) (*struct {
	Message string `json:"message"`
}, error) {
	response := "OK"

	decodedData, err := base64.StdEncoding.DecodeString(req.Body.CommitMessage)
	if err != nil {
		fmt.Println("Error decoding string:", err)
		return nil, huma.Error400BadRequest("")
	}

	// Convert the decoded byte slice to a string
	decodedString := string(decodedData)

	lines := strings.Split(decodedString, "\n")
	lastLine := ""
	if len(lines) > 0 {
		lastLine = strings.TrimSpace(lines[len(lines)-1]) // Trim spaces for safety
	}

	fmt.Printf("Last line of commit message: '%s'\n", lastLine)

	_ = []string{
		lastLine,
		"[t] task title",
		"[t,l] task title | 2h 30m",
		"[t,l] task title | invalid timelog",
		"[t-8]",
		"[t-8,l] | 1h",
		"[t-8,l] 3h 45m",
		"[nc]",
		"[nc,l] | 4h",
		"[nc,l] 5h 30m",
		"[nc,t] task title",
		"[nc,t,l] task title | 6h",
		"[nc,t,l] task title | invalid timelog",
		"[nc,t-8]",
		"[nc,t-8,l] | 7h",
		"[nc,t-8,l] 8h 45m",
		"[t]",
		"[t-8] task title",
		"[nc,t]",
		"[nc,t-8] task title",
		"[t,t-8] task title", // Invalid: t and t-<task-id> coexist
		"[nc,nc] task title", // Invalid: multiple nc
		"[t,l,l] task title", // Invalid: multiple l
		"[t,l] | 2.5h",       // Valid: supports decimal hours
		"[t,l] | 150m",       // Valid: supports minutes only
	}

	// Validate the last line of the commit message
	validationResult := isValidCommitLine(lastLine)
	if validationResult.ProjectID == "" {
		validationResult.ProjectID = req.Body.ProjectID
	}
	if validationResult.IsValid {
		fmt.Println("✅ Valid")
	} else {
		fmt.Println(validationResult)
		fmt.Printf("Invalid commit message: %s %s\n", validationResult.Message, validationResult.ProjectID)
		fmt.Println("❌ Invalid")
	}

	runAction(validationResult, 1)

	return &struct {
		Message string "json:\"message\""
	}{Message: response}, nil
}

func isValidCommitLine(line string) ValidationResult {
	// Regex to match the general structure
	pattern := `^\[([^]]+)\](.*)$`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(line)
	if len(matches) != 3 {
		return ValidationResult{
			IsValid: false,
			Message: "Commit message must start with '[' and end with ']'.",
		}
	}

	params := matches[1]                  // Parameters inside []
	text := strings.TrimSpace(matches[2]) // Text after []

	// Initialize the validation result
	result := ValidationResult{
		IsValid:    true,
		Message:    "Commit message is valid.",
		IsNoCommit: strings.Contains(params, "nc"),
		HasTimelog: strings.Contains(params, "l"),
	}

	// Extract project ID if provided (e.g., [p-123,t] task title)
	projectIDRegex := regexp.MustCompile(`p-(\d+)`)
	projectIDMatch := projectIDRegex.FindStringSubmatch(params)
	if len(projectIDMatch) > 1 {
		result.ProjectID = projectIDMatch[1]
	}

	// Check if t or t-<task-id> is present
	hasT := strings.Contains(params, "t,") || strings.HasSuffix(params, "t")
	hasTWithID := strings.Contains(params, "t-")

	// 1. Check for mutual exclusivity of t and t-<task-id>
	if hasT && hasTWithID {
		return ValidationResult{
			IsValid: false,
			Message: "'t' and 't-<task-id>' cannot coexist.",
		}
	}

	// 2. Check if project ID and task ID coexist
	if result.ProjectID != "" && hasTWithID {
		return ValidationResult{
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
			result.TaskID = taskIDMatch[1]
		}
	} else if hasT {
		result.IsNewTask = true
	}

	// 4. Task title rules
	if hasT && !hasTWithID {
		// t is present, so a task title is required
		if text == "" {
			return ValidationResult{
				IsValid: false,
				Message: "Task title is required for 't'.",
			}
		}
		result.TaskTitle = text
	} else if hasTWithID {
		// t-<task-id> is present, so a task title is not allowed
		if text != "" {
			return ValidationResult{
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

func isValidTimelog(timelog string) bool {
	// Regex to match timelog format (e.g., 23h 34m, 2.5h, 150m)
	pattern := `^(\d+h \d+m|\d+(\.\d+)?h|\d+m)$`
	matched, _ := regexp.MatchString(pattern, timelog)
	return matched
}

func runAction(validationResult ValidationResult, userId int64) error {
	if !validationResult.IsValid {
		return huma.Error400BadRequest(validationResult.Message)
	}

	// check project id
	// is it his~

	if validationResult.IsNewTask {
		// create a task with title0
	} else {
		// does this task belong to this project and person
	}

	if validationResult.HasTimelog {
		if validationResult.Timelog == "" {

		} else {

		}
	}

	if validationResult.IsNoCommit {

	} else {

	}

	return nil
}
