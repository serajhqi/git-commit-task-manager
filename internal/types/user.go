package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserEntity struct {
	tableName        struct{}  `pg:"user"`
	ID               int64     `pg:"id,pk"`                          // Unique identifier
	Name             string    `pg:"name,notnull"`                   // User's full name
	Email            string    `pg:"email,notnull,unique"`           // User's email address
	Password         string    `pg:"password,notnull"`               // User's hashed password
	Verified         bool      `pg:"verified,notnull,default:false"` // User's hashed password
	VerificationCode string    `pg:"verification_code"`
	CreatedAt        time.Time `pg:"created_at,default:now()"` // Timestamp when the user was created
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// ------
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// ------
type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
}
