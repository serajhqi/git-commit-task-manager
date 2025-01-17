package route

import (
	"git-project-management/internal/controller"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func SetupUser(api huma.API) {
	ctrl := controller.NewUserController()

	huma.Register(api, huma.Operation{
		OperationID: "user-login",
		Method:      http.MethodPost,
		Path:        "/user/login",
		Summary:     "login",
		Description: "",
		Tags:        []string{"User"},
	}, ctrl.Login)

	huma.Register(api, huma.Operation{
		OperationID: "user-signup",
		Method:      http.MethodPost,
		Path:        "/user/singup",
		Summary:     "sign up",
		Description: "",
		Tags:        []string{"User"},
	}, ctrl.SignUp)
}
