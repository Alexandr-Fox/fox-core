package middleware

import (
	"github.com/Alexandr-Fox/fox-core/consts"
	"github.com/Alexandr-Fox/fox-core/database"
	"github.com/gofiber/fiber/v2"
)

func Auth(ctx *fiber.Ctx) error {
	if token := ctx.Query("token", consts.DefaultEmpty); token == consts.DefaultEmpty {
		return ctx.Status(fiber.StatusUnauthorized).JSON(database.ResponseError{Error: "invalid token"})
	}

	//ctx.Locals("user", models.NewUser())
	return ctx.Next()
}
