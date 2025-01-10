package types

const (
	DEFAULT_TASK_STATUS   TaskStatus   = TASK_STATUS_IN_PROGRESS
	DEFAULT_TASK_PRIORITY TaskPriority = TASK_PRIORITY_MEDIUM
)

type GetAllRequest struct {
	Offset int `query:"offset"`
	Limit  int `query:"limit"`
}

type TokenizedRequest[T any] struct {
	Token string `header:"token"`
	Body  T
}
