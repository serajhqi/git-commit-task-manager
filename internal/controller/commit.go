package controller

import (
	"context"
	"encoding/base64"
	"fmt"
	"git-project-management/internal/repository"
	"git-project-management/internal/types"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gitea.com/logicamp/lc"
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

// CommitMessage represents the parsed commit message
type CommitMessage struct {
	Title       string `json:"title"`       // Commit title (first line or entire message)
	Description string `json:"description"` // Commit description (remaining lines)
}

// ValidationResult represents the detailed result of the validation
type ValidationResult struct {
	IsValid    bool   `json:"is_valid"`     // Whether the commit message is valid
	Message    string `json:"message"`      // A message describing the validation result
	IsNoCommit bool   `json:"is_no_commit"` // Whether the commit should be ignored (nc)
	IsNewTask  bool   `json:"is_new_task"`  // Whether the commit is for a new task (t)
	TaskID     int64  `json:"task_id"`      // The task ID (if applicable)
	HasTimelog bool   `json:"has_timelog"`  // Whether the commit includes a timelog (l)
	Timelog    string `json:"timelog"`      // The timelog (if applicable)
	TaskTitle  string `json:"task_title"`   // The task title (if applicable)
	ProjectID  int64  `json:"project_id"`   // The project ID (if provided)
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
	commit, customPattern := parseCommitMessage(decodedString)
	fmt.Printf("Last line of commit message: '%s'\n", customPattern)
	lc.Logger.Debug("commit message", zap.Any("commit", commit))

	// Validate the last line of the commit message
	validationResult := isValidCommitLine(customPattern)
	if validationResult.ProjectID == 0 {
		validationResult.ProjectID = req.Body.ProjectID
	}
	if validationResult.IsValid {
		fmt.Println("✅ Valid")
	} else {
		fmt.Println(validationResult)
		fmt.Printf("Invalid commit message: %s %s\n", validationResult.Message, validationResult.ProjectID)
		fmt.Println("❌ Invalid")
	}

	runAction(validationResult, commit, 1)

	return &struct {
		Message string "json:\"message\""
	}{Message: response}, nil
}

func parseCommitMessage(message string) (CommitMessage, string) {
	// Split the commit message into lines
	lines := strings.Split(message, "\n")

	// Initialize the commit message
	var commit CommitMessage

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
		projectID, err := parseInt64(projectIDMatch[1])
		if err != nil {
			return ValidationResult{
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
		return ValidationResult{
			IsValid: false,
			Message: "'t' and 't-<task-id>' cannot coexist.",
		}
	}

	// 2. Check if project ID and task ID coexist
	if result.ProjectID != 0 && hasTWithID {
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
			taskID, err := parseInt64(taskIDMatch[1])
			if err != nil {
				return ValidationResult{
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

func convertTimelogToMinutes(timelog string) (int64, error) {
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
	hours := int64(0)
	if hoursStr != "" {
		var err error
		hours, err = strconv.ParseInt(hoursStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid hours in timelog: %s", hoursStr)
		}
	}

	// Convert minutes to minutes
	minutes := int64(0)
	if minutesStr != "" {
		var err error
		minutes, err = strconv.ParseInt(minutesStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid minutes in timelog: %s", minutesStr)
		}
	}

	// Calculate total minutes
	totalMinutes := hours*60 + minutes
	return totalMinutes, nil
}

func runAction(validationResult ValidationResult, commit CommitMessage, userId int64) error {
	if !validationResult.IsValid {
		return huma.Error400BadRequest(validationResult.Message)
	}

	// check project
	_, err := repository.GetUserProject(validationResult.ProjectID, userId)
	if err != nil {
		lc.Logger.Error("[commit] get project error", zap.Error(err))
		return repository.HandleError(err)
	}

	// check existing task
	var task *types.TaskEntity
	if !validationResult.IsNewTask {
		task, err = repository.GetUserTask(validationResult.TaskID, userId)
		if err != nil {
			lc.Logger.Error("[commit] get task error", zap.Error(err))
			return repository.HandleError(err)
		}

		_, err = repository.Update(context.Background(), task.ID, task)
		if err != nil {
			lc.Logger.Error("[commit] update task error", zap.Error(err))
			return repository.HandleError(err)
		}
	} else {

		task = &types.TaskEntity{
			Title:      validationResult.TaskTitle,
			Status:     types.DEFAULT_TASK_STATUS,
			Priority:   types.DEFAUL_TASK_PRIORITY,
			AssigneeID: userId,
			ProjectID:  validationResult.ProjectID,
			CreatedBy:  userId,
			UpdatedAt:  time.Now(),
		}
		task, err = repository.Create(context.Background(), *task)
		if err != nil {
			lc.Logger.Error("[commit] failed to create new task", zap.Error(err))
			return repository.HandleError(err)
		}
	}

	var timelog int
	if validationResult.HasTimelog {
		if validationResult.Timelog == "" {

		} else {

		}
	}

	types.ActivityEntity{
		TaskID:      task.ID,
		CommitHash:  ,
		Branch:      "",
		Title:       "",
		Description: "",
		Duration:    new(int),
		CreatedBy:   0,
		CreatedAt:   time.Time{},
	}
	return nil
}
