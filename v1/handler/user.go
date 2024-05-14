package handler

import (
	"backend/v1/model"
	"backend/v1/repository"
	"backend/v1/response"
	"backend/v1/utils"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}

	isAll := c.Query("isAll") == "true"
	id := c.Query("id")
	username := c.Query("username")

	if (id == "" && username == "") || isAll {
		if isAdmin := utils.IsAdmin(c); !isAdmin {
			return response.Custom(c, fiber.StatusForbidden, false, "Không có quyền!", nil)
		}
	}

	users, total, err := repository.GetUser(id, username, page, limit, isAll)
	if err != nil {
		return response.Custom(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}

	result := struct {
		Total int                `json:"total"`
		Data  []model.UserExport `json:"data"`
	}{
		total,
		users,
	}

	return response.Success(c, "OK", result)
}

func CreateUser(c *fiber.Ctx) error {
	body := new(model.User)
	if err := c.BodyParser(&body); err != nil {
		return response.Error(c, "Trường dữ liệu không hợp lệ!", nil)
	}

	body.Username = strings.ToLower(body.Username)
	if body.Fullname == "" {
		body.Fullname = body.Username
	}
	if isAdmin := utils.IsAdmin(c); !isAdmin {
		isBanned := true
		body.IsBanned = &isBanned
	}

	if err := utils.Validator.Username(body.Username); err != "" {
		return response.Error(c, err, nil)
	}
	if err := utils.Validator.Password(body.Password); err != "" {
		return response.Error(c, err, nil)
	}

	if body.Password != "" {
		hash, err := utils.HashPassword(body.Password)
		if err != nil {
			return response.Error(c, "Trường dữ liệu không hợp lệ!", nil)
		}
		body.Password = hash
	}

	isUserExists, err := repository.UserExists("", body.Username)
	if err != nil {
		return response.Custom(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}
	if isUserExists {
		return response.Error(c, "Tài khoản đã tồn tại!", nil)
	}

	isEmptyUsers, err := repository.EmptyUsers()
	if err != nil {
		return response.Custom(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}
	if isEmptyUsers {
		body.Role = "Admin"
		isBanned := false
		body.IsBanned = &isBanned
	}

	if err := repository.CreateUser(body); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry ") {
			return response.Custom(c, fiber.StatusBadRequest, false, "Tài khoản đã tồn tại!", nil)
		}
		return response.Custom(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}
	return response.Success(c, "OK", body)
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	isAdmin := utils.IsAdmin(c)

	if isMyRequest := utils.IsMyRequest(c, id); !isMyRequest && !isAdmin {
		return response.Custom(c, fiber.StatusForbidden, false, "Không có quyền!", nil)
	}

	body := new(model.User)
	if err := c.BodyParser(&body); err != nil {
		return response.Error(c, "Trường dữ liệu không hợp lệ!", nil)
	}

	if body.Role != "" && body.Role != "Admin" {
		body.Role = "Member"
	}
	if !isAdmin && body.Role == "Admin" {
		return response.Custom(c, fiber.StatusForbidden, false, "Không có quyền!", nil)
	}

	isUserExists, err := repository.UserExists(id, "")
	if err != nil {
		return response.Custom(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}
	if !isUserExists {
		return response.Error(c, "Tài khoản không tồn tại!", nil)
	}

	if body.Password != "" {
		hash, err := utils.HashPassword(body.Password)
		if err != nil {
			return response.Error(c, "Trường dữ liệu không hợp lệ!", nil)
		}
		body.Password = hash
	}

	err = repository.UpdateUser(id, body)
	if err != nil {
		return response.Custom(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}

	return response.Success(c, "Cập nhật thành công", nil)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if isAdmin := utils.IsAdmin(c); !isAdmin {
		fmt.Println(isAdmin)
		if isMyRequest := utils.IsMyRequest(c, id); !isMyRequest {
			return response.Custom(c, fiber.StatusForbidden, false, "Không có quyền!", nil)
		}
		return response.Custom(c, fiber.StatusForbidden, false, "Không có quyền!", nil)
	}

	isUserExists, err := repository.UserExists(id, "")
	if err != nil {
		return response.Custom(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}
	if !isUserExists {
		return response.Custom(c, fiber.StatusNotFound, false, "Tài khoản không tồn tại!", nil)
	}

	if err := repository.DeleteUser(id); err != nil {
		return response.Custom(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}

	return response.Success(c, "Xóa tài khoản thành công", nil)
}

func ResetUser(c *fiber.Ctx) error {
	if isAdmin := utils.IsAdmin(c); !isAdmin {
		return response.Custom(c, fiber.StatusForbidden, false, "Không có quyền!", nil)
	}

	body := struct {
		ListIds []string
	}{}
	if err := c.BodyParser(&body); err != nil {
		return response.Error(c, "Trường dữ liệu không hợp lệ!", nil)
	}

	if err := repository.ResetUser(body.ListIds); err != nil {
		return response.Custom(c, fiber.StatusInternalServerError, false, err.Error(), nil)
	}

	return response.Success(c, "Khôi phục tài khoản thành công", nil)
}
