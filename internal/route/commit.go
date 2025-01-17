package route

import (
	"git-project-management/internal/controller"
	"git-project-management/internal/repository"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

func SetupCommit(api huma.API) {

	apiKeyMiddleware := func(ctx huma.Context, next func(huma.Context)) {
		apiKey := ctx.Header("Authorization")
		if apiKey == "" {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "No API key provided")
			return
		} else {
			apiKey, err := repository.GetApiKey(apiKey)
			if err != nil {
				huma.WriteErr(api, ctx, http.StatusUnauthorized, "API key not found")
				return
			}

			if apiKey.ExpiresAt.Before(time.Now()) {
				huma.WriteErr(api, ctx, http.StatusUnauthorized, "Renew your API key")
				return
			}

			ctx = huma.WithValue(ctx, "user_id", apiKey.UserID)
		}

		next(ctx)
	}

	ctrl := controller.NewCommitController()

	huma.Register(api, huma.Operation{
		OperationID: "get-one-activity",
		Method:      http.MethodPost,
		Path:        "/commit",
		Summary:     "commit",
		Description: "git commit",
		Tags:        []string{"Commit"},
		Middlewares: huma.Middlewares{apiKeyMiddleware},
	}, ctrl.Commit)
}
