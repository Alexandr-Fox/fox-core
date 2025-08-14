package middleware

import (
	"github.com/Alexandr-Fox/fox-core/consts"
	"github.com/Alexandr-Fox/fox-core/database"
	"github.com/gofiber/fiber/v2"
)

func RolesRequire(roles []string) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		for _, role := range roles {
			if role == consts.DefaultEmpty {
				return ctx.Status(fiber.StatusUnauthorized).JSON(database.ResponseError{Error: "invalid role"})

			}
		}

		return ctx.Next()
	}
}
