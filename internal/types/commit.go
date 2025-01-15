package types

const (
	DEFAULT_TASK_STATUS   TaskStatus   = TASK_STATUS_IN_PROGRESS
	DEFAULT_TASK_PRIORITY TaskPriority = TASK_PRIORITY_MEDIUM
)

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

type CommitResponse struct{}
