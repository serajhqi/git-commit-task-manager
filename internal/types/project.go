package types

import (
	"time"
)

type ProjectEntity struct {
	tableName   struct{}  `pg:"tbl_project,alias:tbl_project"`
	ID          int64     `pg:"id,pk"`                    // Unique identifier
	Name        string    `pg:"name,notnull"`             // Project name
	Description string    `pg:"description"`              // Optional project description
	CreatedBy   int64     `pg:"created_by"`               // User ID who created the project
	CreatedAt   time.Time `pg:"created_at,default:now()"` // Timestamp when the project was created
}

type ProjectDTO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// ---
type CreateProjectRequest struct {
	Body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
}

// ---
type CreateTaskRequest struct {
	ProjectId int64 `path:"id"`
	Body      struct {
		Title       string       `json:"title"`
		ParentID    int64        `json:"parent_id,omitempty"`
		AssigneeID  int64        `json:"assignee_id,omitempty"`
		Description string       `json:"description,omitempty"`
		Status      TaskStatus   `json:"status,omitempty" enum:"open,closed"`
		Priority    TaskPriority `json:"priority,omitempty" enum:"high,medium,low"`
		DueDate     time.Time    `json:"due_date,omitempty"`
	}
}

type GetProjectTasksRequest struct {
	Authorization string `header:"Authorization"`
	ProjectID     int64  `path:"project_id"`
	Offset        int    `query:"offset,minimum=0"`
	Limit         int    `query:"limit,minimum=1,maximum=1000"`
}
