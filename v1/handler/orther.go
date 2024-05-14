package handler

import (
	"backend/v1/model"
	"backend/v1/repository"
	"backend/v1/response"
	"net/smtp"

	"github.com/gofiber/fiber/v2"
)

func sendMail(to []string, msg []byte) error {
	auth := smtp.PlainAuth("", "email", "pass", "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "email", to, msg)

	return err
}

func CreateContact(c *fiber.Ctx) error {
	body := new(model.Contact)
	if err := c.BodyParser(&body); err != nil {
		return response.Error(c, "Trường dữ liệu không hợp lệ!", nil)
	}

	if err := repository.CreateContact(body); err != nil {
		return response.Custom(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}

	msg := []byte(
		"Subject: Có một yêu cầu liên hệ mới từ Brand name \r\n" +
			"\r\n" +
			"Họ và tên: " + body.Fullname + "\r\n" +
			"Email: " + body.Email + "\r\n" +
			"Số điện thoại: " + body.Phone + "\r\n" +
			"Nội dung: " + body.Content)
	if err := sendMail([]string{"email"}, msg); err != nil {
		return response.Custom(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}

	return response.Success(c, "OK", body)

}
