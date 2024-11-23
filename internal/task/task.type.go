package task

import (
	"git-project-management/internal/activity"
	"time"
)

type TaskEntity struct {
	tableName   struct{}     `pg:"task"`
	ID          int64        `pg:"id,pk"`                    // Unique identifier
	ParentID    int64        `pg:"parent_id"`                // Unique identifier
	Title       string       `pg:"title,notnull"`            // Task title
	Description string       `pg:"description"`              // Task description
	Status      TaskStatus   `pg:"status,notnull"`           // Task status (e.g., "Pending", "Completed")
	Priority    TaskPriority `pg:"priority,notnull"`         // Task priority (e.g., "Low", "Medium", "High")
	AssigneeID  int64        `pg:"assignee_id"`              // ID of the user assigned to this task
	ProjectID   int64        `pg:"project_id,notnull"`       // ID of the project this task belongs to
	DueDate     time.Time    `pg:"due_date"`                 // Due date for the task
	CreatedBy   int64        `pg:"created_by"`               // User ID who created the task
	CreatedAt   time.Time    `pg:"created_at,default:now()"` // Timestamp when the task was created
}

// Get All ---
type TaskStatus string

const (
	TODO        TaskStatus = "todo"
	IN_PROGRESS TaskStatus = "in_progress"
	DONE        TaskStatus = "done"
	CANCELED    TaskStatus = "cancelled"
)

type TaskPriority string

const (
	HIGH   TaskPriority = "high"
	MEDIUM TaskPriority = "medium"
	LOW    TaskPriority = "low"
)

type TaskDTO struct {
	ID          int64        `json:"id"`
	ParentID    int64        `json:"parent_id,omitempty"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	AssigneeID  int64        `json:"assignee_id"`
	ProjectID   int64        `json:"project_id"`
	DueDate     time.Time    `json:"due_date,omitempty"`
	CreatedBy   int64        `json:"created_by"`
	CreatedAt   time.Time    `json:"created_at"`
}

// ---
type GetOneRequest struct {
	Id int64 `path:"id"`
}

type GetOneResponse struct {
	Body TaskDTO
}

// ---

type GetAllActivitiesRequest struct {
	Limit  int   `query:"limit"`
	Offset int   `query:"offset"`
	TaskID int64 `path:"id"`
}
type GetAllActivitiesResponse struct {
	Body []activity.ActivityDTO
}

// ---

type SetTaskStatusRequest struct {
	TaskID int64 `path:"id"`
	Body   struct {
		Status TaskStatus `json:"status" enum:"open,closed"`
	}
}
type SetTaskStatusResponse struct {
	Body TaskDTO
}
