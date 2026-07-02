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

type TaskController struct {
	taskService services.TaskService
}

func NewTaskController(taskService services.TaskService) *TaskController {
	return &TaskController{taskService: taskService}
}

func (ctrl *TaskController) Create(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if details := validators.ValidateStruct(req); details != nil {
		utils.ValidationError(c, details)
		return
	}

	resp, err := ctrl.taskService.Create(req)
	if err != nil {
		if errors.Is(err, services.ErrTicketNotFound) {
			utils.NotFound(c, "ticket")
			return
		}
		utils.InternalError(c, "failed to create task")
		return
	}

	utils.Created(c, resp)
}

func (ctrl *TaskController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	resp, err := ctrl.taskService.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrTaskNotFound) {
			utils.NotFound(c, "task")
			return
		}
		utils.InternalError(c, "failed to get task")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}

func (ctrl *TaskController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if details := validators.ValidateStruct(req); details != nil {
		utils.ValidationError(c, details)
		return
	}

	resp, err := ctrl.taskService.Update(id, req)
	if err != nil {
		if errors.Is(err, services.ErrTaskNotFound) {
			utils.NotFound(c, "task")
			return
		}
		utils.InternalError(c, "failed to update task")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}

func (ctrl *TaskController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	if err := ctrl.taskService.Delete(id); err != nil {
		if errors.Is(err, services.ErrTaskNotFound) {
			utils.NotFound(c, "task")
			return
		}
		utils.InternalError(c, "failed to delete task")
		return
	}

	utils.NoContent(c)
}

func (ctrl *TaskController) List(c *gin.Context) {
	var filter dto.TaskFilterQuery
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	resp, err := ctrl.taskService.List(filter)
	if err != nil {
		utils.InternalError(c, "failed to list tasks")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}
