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
