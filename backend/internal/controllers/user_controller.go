package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/dto"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/services"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/utils"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/validators"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (ctrl *UserController) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if details := validators.ValidateStruct(req); details != nil {
		utils.ValidationError(c, details)
		return
	}

	resp, err := ctrl.userService.Create(req)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			utils.Error(c, http.StatusConflict, err.Error())
			return
		}
		utils.InternalError(c, "failed to create user")
		return
	}

	utils.Created(c, resp)
}

func (ctrl *UserController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	resp, err := ctrl.userService.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			utils.NotFound(c, "user")
			return
		}
		utils.InternalError(c, "failed to get user")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}

func (ctrl *UserController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if details := validators.ValidateStruct(req); details != nil {
		utils.ValidationError(c, details)
		return
	}

	resp, err := ctrl.userService.Update(id, req)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			utils.NotFound(c, "user")
			return
		}
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			utils.Error(c, http.StatusConflict, err.Error())
			return
		}
		utils.InternalError(c, "failed to update user")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}

func (ctrl *UserController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := ctrl.userService.Delete(id); err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			utils.NotFound(c, "user")
			return
		}
		utils.InternalError(c, "failed to delete user")
		return
	}

	utils.NoContent(c)
}

func (ctrl *UserController) List(c *gin.Context) {
	var filter dto.UserFilterQuery
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	resp, err := ctrl.userService.List(filter)
	if err != nil {
		utils.InternalError(c, "failed to list users")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}

func (ctrl *UserController) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	resp, err := ctrl.userService.GetByID(userID.(uuid.UUID))
	if err != nil {
		utils.InternalError(c, "failed to get current user")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}
