package route

import (
	"git-project-management/internal/controller"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func SetupApiKey(api huma.API) {
	ctrl := controller.NewApiKeyController()

	huma.Register(api, huma.Operation{
		OperationID: "get-api-key",
		Method:      http.MethodGet,
		Path:        "/api-key",
		Summary:     "get api key",
		Description: "get api key",
		Tags:        []string{"API Key"},
	}, ctrl.GetNewApiKey)
}
