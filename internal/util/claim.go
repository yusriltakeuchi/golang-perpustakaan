package util

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// GetUserClaim retrieves the JWT claims from the Fiber context
func GetUserClaim(ctx *fiber.Ctx) map[string]interface{} {
	claims := ctx.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	return claims
}
