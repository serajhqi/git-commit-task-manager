package controller

import (
	"context"
	"git-project-management/internal/types"
)

type UserController struct{}

func NewUserController() UserController {
	return UserController{}
}

func (uc *UserController) Login(ctx context.Context, req *types.LoginRequest) (*types.LoginResponse, error) {
	return nil, nil
}

func (uc *UserController) SignUp(ctx context.Context, req *types.SignUpRequest) (*types.SignUpResponse, error) {

	return nil, nil
}
