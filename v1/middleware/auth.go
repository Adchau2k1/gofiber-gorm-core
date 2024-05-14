package middleware

import (
	"backend/v1/response"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func handleError(c *fiber.Ctx, err error) error {
	return response.Custom(c, fiber.StatusUnauthorized, false, "Unauthorized", nil)
}

// Middleware JWT function
func Auth(secretKey string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(secretKey),
		},
		ErrorHandler: handleError,
	})
}
