package controller

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"git-project-management/internal/repository"
	"git-project-management/internal/types"
	"time"

	"gitea.com/logicamp/lc"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"
)

type ApiKeyController struct{}

func NewApiKeyController() ApiKeyController {
	return ApiKeyController{}
}

func (akc *ApiKeyController) GetNewApiKey(ctx context.Context, req *struct{}) (*types.GetNewApiKeyResponse, error) {
	key, _ := gonanoid.New(64)
	hash := md5.Sum([]byte(key))
	hashedKey := hex.EncodeToString(hash[:])
	_, err := repository.Create(ctx, types.ApiKeyEntity{
		UserID:    1,
		Key:       hashedKey,
		ExpiresAt: time.Now().AddDate(1, 0, 0),
	})

	if err != nil {
		lc.Logger.Error("[api_key] get api key", zap.Error(err))
		return nil, repository.HandleError(err)
	}
	return &types.GetNewApiKeyResponse{
		Body: struct {
			API_KEY string `json:"api_key"`
		}{
			API_KEY: key,
		},
	}, nil
}

func (akc *ApiKeyController) RemoveApiKey(ctx context.Context, req types.RemoveApiKeyRequest) (*types.RemoveApiKeyResponse, error) {

	err := repository.RemoveApiKey(req.Body.API_KEY)
	if err != nil {
		lc.Logger.Error("[api_key] get api key", zap.Error(err))
		return nil, repository.HandleError(err)
	}
	return &types.RemoveApiKeyResponse{}, nil
}
