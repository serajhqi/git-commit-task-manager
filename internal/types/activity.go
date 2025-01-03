package types

import "time"

type ActivityEntity struct {
	tableName   struct{}  `pg:"activity"`
	ID          int64     `pg:"id,pk"`                    // Unique identifier
	TaskID      int64     `pg:"task_id,notnull"`          // ID of the associated task
	CommitHash  string    `pg:"commit_hash"`              // ID of the associated task
	Branch      string    `pg:"branch"`                   // ID of the associated task
	Title       string    `pg:"title,notnull"`            // Action performed (e.g., "Worked On sth")
	Description string    `pg:"description"`              // Optional description of the action
	Duration    *int      `pg:"duration"`                 // Duration in minutes (optional for non-timelog actions)
	CreatedBy   int64     `pg:"created_by,notnull"`       // CreatedBy of the user who performed the action
	CreatedAt   time.Time `pg:"created_at,default:now()"` // Timestamp when the action occurred
}

type IdBody struct {
	Id int64 `json:"id"`
}

// Create ---
type CreateActivityRequest struct {
	Body struct {
		CommitMessage string `json:"commit_message"`
	}
}
type CreateActivityResponse struct {
	Body IdBody
}

// Get ---
type ActivityDTO struct {
	ID          int64     `json:"id,pk"`
	TaskId      int64     `json:"task_id"`
	Title       string    `json:"title"`
	CommitHash  string    `json:"commit_hash"`
	Branch      string    `json:"branch"`
	Description string    `json:"description"`
	Duration    *int      `json:"duration"`
	CreatedBy   int64     `json:"created_by"`
	CreateAt    time.Time `json:"created_at"`
}

type GetAllRequest struct {
	Limit  int   `query:"limit"`
	Offset int   `query:"offset"`
	TaskId int64 `path:"task_id"`
}
type GetAllResponse struct {
	Body []ActivityDTO
}

// ---
type GetOneRequest struct {
	Id int64 `path:"id"`
}

type GetOneResponse struct {
	Body ActivityDTO
}

type GetAllActivitiesRequest struct {
	Limit  int   `query:"limit"`
	Offset int   `query:"offset"`
	TaskID int64 `query:"task_id"`
}
type GetAllActivitiesResponse struct {
	Body []ActivityDTO
}

// ---
