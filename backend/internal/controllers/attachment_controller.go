package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/dto"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/middlewares"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/models"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/services"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/utils"
)

type AttachmentController struct {
	attachmentService services.AttachmentService
}

func NewAttachmentController(attachmentService services.AttachmentService) *AttachmentController {
	return &AttachmentController{attachmentService: attachmentService}
}

func (ctrl *AttachmentController) Upload(c *gin.Context) {
	entityType := models.EntityType(c.PostForm("entity_type"))
	entityIDStr := c.PostForm("entity_id")

	if entityType != models.EntityTypeTicket && entityType != models.EntityTypeTask {
		utils.Error(c, http.StatusBadRequest, "invalid entity_type")
		return
	}

	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid entity_id")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "file is required")
		return
	}

	userID := c.MustGet(middlewares.ContextUserIDKey).(uuid.UUID)

	resp, err := ctrl.attachmentService.Upload(entityType, entityID, file, userID)
	if err != nil {
		if errors.Is(err, services.ErrTicketNotFound) || errors.Is(err, services.ErrTaskNotFound) {
			utils.NotFound(c, "entity")
			return
		}
		if errors.Is(err, services.ErrFileTooLarge) || errors.Is(err, services.ErrInvalidFileType) {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.InternalError(c, "failed to upload attachment")
		return
	}

	utils.Created(c, resp)
}

func (ctrl *AttachmentController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid attachment id")
		return
	}

	resp, err := ctrl.attachmentService.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrAttachmentNotFound) {
			utils.NotFound(c, "attachment")
			return
		}
		utils.InternalError(c, "failed to get attachment")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}

func (ctrl *AttachmentController) Download(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid attachment id")
		return
	}

	filePath, originalName, err := ctrl.attachmentService.GetFilePath(id)
	if err != nil {
		if errors.Is(err, services.ErrAttachmentNotFound) {
			utils.NotFound(c, "attachment")
			return
		}
		utils.InternalError(c, "failed to get attachment file")
		return
	}

	c.FileAttachment(filePath, originalName)
}

func (ctrl *AttachmentController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid attachment id")
		return
	}

	if err := ctrl.attachmentService.Delete(id); err != nil {
		if errors.Is(err, services.ErrAttachmentNotFound) {
			utils.NotFound(c, "attachment")
			return
		}
		utils.InternalError(c, "failed to delete attachment")
		return
	}

	utils.NoContent(c)
}

func (ctrl *AttachmentController) List(c *gin.Context) {
	var filter dto.AttachmentFilterQuery
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	resp, err := ctrl.attachmentService.List(filter)
	if err != nil {
		utils.InternalError(c, "failed to list attachments")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}
