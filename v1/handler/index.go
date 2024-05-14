package handler

import (
	"github.com/gofiber/fiber/v2"
)

type handleFunc func(c *fiber.Ctx) error

// Define struct
type other struct {
	CreateContact handleFunc
}

type auth struct {
	Login, Refresh handleFunc
}

type user struct {
	GetUser, CreateUser, UpdateUser, DeleteUser, ResetUser handleFunc
}

// Use
var Other = other{
	CreateContact,
}

var Auth = auth{
	Login,
	Refresh,
}

var User = user{
	GetUser,
	CreateUser,
	UpdateUser,
	DeleteUser,
	ResetUser,
}
