package response

import "github.com/gofiber/fiber/v2"

func Success(c *fiber.Ctx, message string, data interface{}) error {
	if data != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": message,
			"data":    data,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": message,
	})
}

func Error(c *fiber.Ctx, message string, data interface{}) error {
	if data != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"data":    data,
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}

func Custom(c *fiber.Ctx, status int, success bool, message string, data interface{}) error {
	if data != nil {
		return c.Status(status).JSON(fiber.Map{
			"success": success,
			"message": message,
			"data":    data,
		})
	}

	return c.Status(status).JSON(fiber.Map{
		"success": success,
		"message": message,
	})
}
