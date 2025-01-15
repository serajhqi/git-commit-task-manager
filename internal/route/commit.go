package route

import (
	"git-project-management/internal/controller"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func SetupCommit(api huma.API) {
	ctrl := controller.NewCommitController()

	huma.Register(api, huma.Operation{
		OperationID: "get-one-activity",
		Method:      http.MethodPost,
		Path:        "/commit",
		Summary:     "commit",
		Description: "git commit",
		Tags:        []string{"Commit"},
	}, ctrl.Commit)
}
