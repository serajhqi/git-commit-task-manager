package types

import "time"

type ApiKeyEntity struct {
	tableName struct{}  `pg:"api_key"`
	ID        int64     `pg:"id,pk"`
	UserID    int64     `pg:"user_id"`
	Key       string    `pg:"key,unique"`
	ExpiresAt time.Time `pg:"expires_at"`
	CreatedAt time.Time `pg:"created_at,default:now()"`
}

type GetNewApiKeyRequest struct{}

type GetNewApiKeyResponse struct {
	Body struct {
		API_KEY string `json:"api_key"`
	}
}

type RemoveApiKeyRequest struct {
	Body struct {
		API_KEY string `json:"api_key"`
	}
}

type RemoveApiKeyResponse struct{}
