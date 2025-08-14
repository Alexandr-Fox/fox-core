package main

import (
	"fmt"
	"github.com/Alexandr-Fox/fox-core/boot"
	"github.com/Alexandr-Fox/fox-core/models"
	"github.com/Alexandr-Fox/fox-core/router"
	"github.com/gofiber/fiber/v2"
)

//go:generate go run generator/generator_methods.go
//go:generate go run generator/generator_migrate.go
//go:generate go run generator/generator_routers_rest.go

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	app := fiber.New()

	models.Migrate()
	config, err := boot.LoadConfig("config.yaml")
	if err != nil {
		panic(config)
	}

	router.RegisterRESTRoutes(app)

	err = app.Listen(fmt.Sprintf("%s:%d", config.App.Host, config.App.Port))
	if err != nil {
		panic(err)
	}
}
