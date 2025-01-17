package controller

import (
	"context"
	"git-project-management/config"
	"git-project-management/internal/controller/utils"
	"git-project-management/internal/repository"
	"git-project-management/internal/types"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"

	"gitea.com/logicamp/lc"
)

type UserController struct{}

func NewUserController() UserController {
	return UserController{}
}

var jwtKey = []byte(config.GetConfig().JWT_PRIVATE_KEY)

func (uc UserController) Login(ctx context.Context, req *lc.ReqWithBody[types.LoginRequest]) (*lc.RespBody[types.LoginResponse], error) {

	user, err := repository.GetUser(req.Body.Email, utils.ToSha256(req.Body.Password))
	if err != nil {
		return nil, huma.Error404NotFound("incorrect email or password")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &types.Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, lc.SendInternalErrorResponse(err, "[user] failed to create token string")
	}

	return &lc.RespBody[types.LoginResponse]{
		Body: types.LoginResponse{
			Token: tokenString,
		},
	}, nil
}

func (uc UserController) SignUp(ctx context.Context, req *lc.ReqWithBody[types.SignUpRequest]) (*lc.RespBody[types.SignUpResponse], error) {
	user, err := repository.GetUserByEmail(req.Body.Email)
	if user != nil && user.ID > 0 {
		return nil, huma.Error403Forbidden("email already exist")
	}

	code, _ := gonanoid.New(6)
	newUser := types.UserEntity{
		Name:             req.Body.Name,
		Email:            req.Body.Email,
		Password:         utils.ToSha256(req.Body.Password),
		VerificationCode: code,
	}

	user, err = repository.Create(ctx, newUser)
	if err != nil {
		return nil, repository.HandleError(err)
	}

	return &lc.RespBody[types.SignUpResponse]{
		Body: types.SignUpResponse{},
	}, nil
}
