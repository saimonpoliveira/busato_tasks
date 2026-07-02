package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/dto"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/services"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/utils"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/validators"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if details := validators.ValidateStruct(req); details != nil {
		utils.ValidationError(c, details)
		return
	}

	resp, err := ctrl.authService.Login(req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) || errors.Is(err, services.ErrUserInactive) {
			utils.Unauthorized(c, err.Error())
			return
		}
		utils.InternalError(c, "failed to login")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if details := validators.ValidateStruct(req); details != nil {
		utils.ValidationError(c, details)
		return
	}

	resp, err := ctrl.authService.Register(req)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			utils.Error(c, http.StatusConflict, err.Error())
			return
		}
		utils.InternalError(c, "failed to register")
		return
	}

	utils.Created(c, resp)
}
