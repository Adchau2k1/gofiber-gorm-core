package handler

import (
	"backend/v1/config"
	"backend/v1/database"
	"backend/v1/model"
	"backend/v1/response"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var secretKey = config.GetConfigByKey("SECRET_KEY")

func Login(c *fiber.Ctx) error {
	var body struct {
		Username string
		Password string
	}

	db := database.DB
	err := c.BodyParser(&body)
	if err != nil {
		return response.Error(c, "Trường dữ liệu không hợp lệ!", nil)
	}

	userInfo := model.User{}
	result := db.First(&userInfo, "username = ?", body.Username)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return response.Error(c, "Tài khoản mật khẩu không chính xác!", nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(body.Password))
	if err != nil {
		return response.Error(c, "Tài khoản mật khẩu không chính xác!", nil)
	}

	if *userInfo.IsBanned {
		return response.Custom(c, fiber.StatusForbidden, false, "Tài khoản của bạn đang bị khóa!", nil)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"id":       userInfo.ID,
		"username": userInfo.Username,
		"role":     userInfo.Role,
		"isBanned": userInfo.IsBanned,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	rClaims := jwt.MapClaims{
		"id":       userInfo.ID,
		"username": userInfo.Username,
		"role":     userInfo.Role,
		"isBanned": userInfo.IsBanned,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rClaims)
	t, err := token.SignedString([]byte(secretKey))
	rt, rErr := rToken.SignedString([]byte(secretKey))
	if err != nil || rErr != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return response.Success(c, "OK", fiber.Map{
		"id":           userInfo.ID,
		"username":     userInfo.Username,
		"fullname":     userInfo.Fullname,
		"role":         userInfo.Role,
		"image":        userInfo.Image,
		"isBanned":     userInfo.IsBanned,
		"accessToken":  t,
		"refreshToken": rt,
	})
}

func Refresh(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}
	var err = c.BodyParser(&body)
	if err != nil {
		return response.Error(c, "Trường dữ liệu không hợp lệ!", nil)
	}

	// Auth refresh token
	rfToken, err := jwt.Parse(body.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !rfToken.Valid {
		return response.Custom(c, fiber.StatusUnauthorized, false, "Invalid refresh token!", nil)
	}

	claims, ok := rfToken.Claims.(jwt.MapClaims)
	if !ok {
		if err != nil || !rfToken.Valid {
			return response.Custom(c, fiber.StatusUnauthorized, false, "Invalid refresh token!", nil)
		}
	}

	// Create the Claims
	newClaims := jwt.MapClaims{
		"id":       claims["id"],
		"username": claims["username"],
		"role":     claims["role"],
		"isBanned": claims["isBanned"],
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Minute * 30).Unix(),
	}

	// Create token
	clToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	accessToken, err := clToken.SignedString([]byte(secretKey))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	data := struct {
		AccessToken string `json:"accessToken"`
	}{accessToken}
	return response.Success(c, "OK", data)
}
