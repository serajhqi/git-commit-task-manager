package main

import (
	"context"
	"fmt"
	"git-project-management/config"
	"git-project-management/internal/database"
	"git-project-management/internal/route"
	"log"

	"gitea.com/logicamp/lc"
	"github.com/danielgtaylor/huma/v2"
	humaFiber "github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

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

	route.SetupProject(api)
	route.SetupTask(api)
	route.SetupActivity(api)
	route.SetupCommit(api)

	log.Fatal(fiberApp.Listen(fmt.Sprintf(":%s", config.PORT)))
}
