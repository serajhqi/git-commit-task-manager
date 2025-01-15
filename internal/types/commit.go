package types

const (
	DEFAULT_TASK_STATUS   TaskStatus   = TASK_STATUS_IN_PROGRESS
	DEFAULT_TASK_PRIORITY TaskPriority = TASK_PRIORITY_MEDIUM
	DEFAULT_TASK_WEIGHT   uint         = 3
)

type CommitMessage struct {
	Title       string `json:"title"`       // Commit title (first line or entire message)
	Description string `json:"description"` // Commit description (remaining lines)
	Hash        string `json:"hash"`
	Branch      string `json:"branch"`
}

// CommitValidationResult represents the detailed result of the validation
type CommitValidationResult struct {
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

// ------------------------------------
type GetAllRequest struct {
	Offset int `query:"offset"`
	Limit  int `query:"limit"`
}

type CommitRequest struct {
	Authorization string `header:"Authorization"`
	Body          struct {
		ProjectID     int64  `json:"project_id"`
		CommitMessage string `json:"commit_message"`
		CommitHash    string `json:"commit_hash"`
		Branch        string `json:"branch"`
	}
}

type ImportHistory struct {
	Authorization string `header:"Authorization"`
	Body          struct {
		Commits []CommitMessage
	}
}

type CommitResponse struct{}
