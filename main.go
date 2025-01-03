package main

import (
	"context"
	"fmt"
	"git-project-management/config"
	"log"

	"gitea.com/logicamp/lc"
	"github.com/danielgtaylor/huma/v2"
	humaFiber "github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/go-pg/pg/v10"
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
	db := pg.Connect(&pg.Options{
		Addr:     config.PG_HOST,
		User:     config.PG_USER,
		Password: config.PG_PASSWORD,
		Database: config.PG_DATABASE,
	})
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		panic(err)
	}
	// ------------------------

	app.Setup(&api, db)

	log.Fatal(fiberApp.Listen(fmt.Sprintf(":%s", config.PORT)))
}
