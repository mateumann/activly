package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mateumann/activly/core/services"
)

type FiberHandler struct {
	userService *services.UserService
}

func NewFiberHandler(userService *services.UserService) *FiberHandler {
	return &FiberHandler{
		userService: userService,
	}
}

func (h *FiberHandler) CreateUser(c *fiber.Ctx) error {
	response := "creating users\n"

	type User struct {
		Name     string                 `json:"name" xml:"name" form:"name" bson:"name" asn1:"name"`
		Settings map[string]interface{} `json:"settings" xml:"settings" form:"settings" bson:"settings" asn1:"settings"`
	}

	u := new(User)
	if err := c.BodyParser(u); err != nil {
		return fmt.Errorf("cannot parse request body: %w", err)
	}

	err := h.userService.CreateUser(u.Name, u.Settings)
	if err != nil {
		return fmt.Errorf("cannot create a user: %w", err)
	}

	response += "user created\n"

	err = c.SendString(response)
	if err != nil {
		return fmt.Errorf("cannot send string with fiber: %w", err)
	}

	return nil
}

func (h *FiberHandler) ListUsers(c *fiber.Ctx) error {
	response := "listing users\n"

	users, err := h.userService.ListUsers()
	if err != nil {
		return fmt.Errorf("cannot list users: %w", err)
	}

	for _, u := range users {
		response += fmt.Sprintf("  · %s\n", u)
	}

	err = c.SendString(response)
	if err != nil {
		return fmt.Errorf("cannot send string with fiber: %w", err)
	}

	return nil
}
