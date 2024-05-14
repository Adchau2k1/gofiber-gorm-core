package utils

import (
	"backend/v1/config"
	"strconv"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func IsMyRequest(c *fiber.Ctx, id string) bool {
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return false
	}

	claims := token.Claims.(jwt.MapClaims)
	claimsId, ok := claims["id"].(float64)
	if !ok {
		return false
	}

	newId, strErr := strconv.ParseFloat(id, 64)
	if strErr != nil {
		return false
	}

	return claimsId == newId
}

func IsAdmin(c *fiber.Ctx) bool {
	secretKey := config.GetConfigByKey("SECRET_KEY")
	requestToken := strings.Replace(c.Get("Authorization"), "Bearer ", "", 1)

	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	role := claims["role"].(string)

	return role == "Admin"
}

func RemoveAccent(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

func TextToAlias(text string) string {
	alias := strings.ToLower(strings.ReplaceAll(RemoveAccent(text), " ", "-"))
	return alias
}
