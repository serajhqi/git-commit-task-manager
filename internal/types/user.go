package types

import "time"

type UserEntity struct {
	tableName struct{}  `pg:"user"`
	ID        int64     `pg:"id,pk"`                    // Unique identifier
	Name      string    `pg:"name,notnull"`             // User's full name
	Email     string    `pg:"email,notnull,unique"`     // User's email address
	Password  string    `pg:"password,notnull"`         // User's hashed password
	CreatedAt time.Time `pg:"created_at,default:now()"` // Timestamp when the user was created
}
