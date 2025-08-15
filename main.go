package main

import (
	"fmt"
	"github.com/Alexandr-Fox/fox-core/internal/boot"
	"github.com/Alexandr-Fox/fox-core/internal/database"
	"github.com/Alexandr-Fox/fox-core/internal/models"
	"github.com/Alexandr-Fox/fox-core/router"
	"github.com/gofiber/fiber/v2"
	"strings"
)

//go:generate go run generator/generator_methods.go
//go:generate go run generator/generator_migrate.go
//go:generate go run generator/generator_docs_route.go
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

	app.Options("/*", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderAccessControlAllowOrigin, "*")
		c.Set(fiber.HeaderAccessControlAllowMethods, "POST, GET, OPTIONS")
		c.Set(fiber.HeaderAccessControlAllowHeaders, strings.Join(
			[]string{
				fiber.HeaderContentType,
				fiber.HeaderAccept,
				fiber.HeaderAuthorization,
				fiber.HeaderETag,
				fiber.HeaderExpect,
				fiber.HeaderCacheControl,
				fiber.HeaderUpgrade,
				fiber.HeaderXRequestedWith,
				fiber.HeaderXRequestID,
				fiber.HeaderAcceptLanguage,
				fiber.HeaderAcceptEncoding,
				fiber.HeaderAcceptCharset,
				fiber.HeaderEarlyData,
				fiber.HeaderLastModified,
				fiber.HeaderLastEventID,
				fiber.HeaderIfMatch,
				fiber.HeaderIfNoneMatch,
				fiber.HeaderCookie,
				fiber.HeaderSetCookie,
			}, ", "))
		return c.SendStatus(fiber.StatusNoContent)
	})
	app.All("/*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusMethodNotAllowed).JSON(database.ResponseError{Error: fmt.Sprintf("Path or method not found")})
	})
	err = app.Listen(fmt.Sprintf("0.0.0.0:%d", config.App.Port))
	if err != nil {
		panic(err)
	}
}
