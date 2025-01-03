package types

import "time"

type NotificationEntity struct {
	tableName  struct{}  `pg:"notification"`
	ID         int64     `pg:"id,pk"`                    // Unique identifier
	UserID     int64     `pg:"user_id,notnull"`          // ID of the user receiving the notification
	EntityType string    `pg:"entity_type,notnull"`      // Type of entity triggering the notification (e.g., "Task", "Project")
	EntityID   int64     `pg:"entity_id,notnull"`        // ID of the entity triggering the notification
	Message    string    `pg:"message,notnull"`          // Notification message
	Type       string    `pg:"type,notnull"`             // Type of notification (e.g., "Info", "Warning", "Critical")
	IsRead     bool      `pg:"is_read"`                  // Whether the notification has been read
	CreatedBy  int64     `pg:"created_by"`               // User ID who created the notification
	CreatedAt  time.Time `pg:"created_at,default:now()"` // Timestamp when the notification was created
}
