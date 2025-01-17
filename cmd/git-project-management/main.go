package main

import (
	"context"
	"fmt"
	"git-project-management/config"
	"git-project-management/internal/database"
	"git-project-management/internal/route"
	"git-project-management/internal/types"
	"log"
	"net/http"

	"gitea.com/logicamp/lc"
	"github.com/danielgtaylor/huma/v2"
	humaFiber "github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(config.GetConfig().JWT_PRIVATE_KEY)

func GetApi() {

}

func main() {
	config, _ := lc.GetConfig[config.Config](&config.Config{})
	fiberApp := fiber.New()
	fiberApp.Use(cors.New())

	// Or extend your config for customization
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	humaConfig := huma.DefaultConfig("Git Project Management", "1.0.0")
	humaConfig.Servers = []*huma.Server{{URL: config.BASE_URL}}
	humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"auth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}
	api := humaFiber.New(fiberApp, humaConfig)

	// database init ---------
	db := database.GetDB()
	defer db.Close()
	if err := db.Ping(context.Background()); err != nil {
		panic(err)
	}
	// ------------------------
	route.SetupUser(api)
	route.SetupCommit(api)

	api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		authHeader := ctx.Header("Authorization")

		if authHeader == "" {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Authorization header missing")
			return
		}
		tokenString := authHeader[len("Bearer "):]

		claims := &types.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid token")
			return
		}

		ctx = huma.WithValue(ctx, "user_id", claims.UserID)
		next(ctx)
	})
	route.SetupProject(api)
	route.SetupApiKey(api)
	route.SetupTask(api)
	route.SetupActivity(api)

	log.Fatal(fiberApp.Listen(fmt.Sprintf(":%s", config.PORT)))
}
