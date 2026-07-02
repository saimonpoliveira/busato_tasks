package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/dto"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/middlewares"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/services"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/utils"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/validators"
)

type ProjectController struct {
	projectService services.ProjectService
}

func NewProjectController(projectService services.ProjectService) *ProjectController {
	return &ProjectController{projectService: projectService}
}

func (ctrl *ProjectController) Create(c *gin.Context) {
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if details := validators.ValidateStruct(req); details != nil {
		utils.ValidationError(c, details)
		return
	}

	ownerID := c.MustGet(middlewares.ContextUserIDKey).(uuid.UUID)

	resp, err := ctrl.projectService.Create(req, ownerID)
	if err != nil {
		utils.InternalError(c, "failed to create project")
		return
	}

	utils.Created(c, resp)
}

func (ctrl *ProjectController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	resp, err := ctrl.projectService.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrProjectNotFound) {
			utils.NotFound(c, "project")
			return
		}
		utils.InternalError(c, "failed to get project")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}

func (ctrl *ProjectController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if details := validators.ValidateStruct(req); details != nil {
		utils.ValidationError(c, details)
		return
	}

	resp, err := ctrl.projectService.Update(id, req)
	if err != nil {
		if errors.Is(err, services.ErrProjectNotFound) {
			utils.NotFound(c, "project")
			return
		}
		utils.InternalError(c, "failed to update project")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}

func (ctrl *ProjectController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	if err := ctrl.projectService.Delete(id); err != nil {
		if errors.Is(err, services.ErrProjectNotFound) {
			utils.NotFound(c, "project")
			return
		}
		utils.InternalError(c, "failed to delete project")
		return
	}

	utils.NoContent(c)
}

func (ctrl *ProjectController) List(c *gin.Context) {
	var filter dto.ProjectFilterQuery
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	resp, err := ctrl.projectService.List(filter)
	if err != nil {
		utils.InternalError(c, "failed to list projects")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}
