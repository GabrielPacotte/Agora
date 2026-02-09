package userhttp

import "github.com/GabrielPacotte/Agora/internal/services"

type Handler struct {
	userService services.UserService
}

func NewHandler(userService services.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}
