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

type CommentController struct {
	commentService services.CommentService
}

func NewCommentController(commentService services.CommentService) *CommentController {
	return &CommentController{commentService: commentService}
}

func (ctrl *CommentController) Create(c *gin.Context) {
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if details := validators.ValidateStruct(req); details != nil {
		utils.ValidationError(c, details)
		return
	}

	userID := c.MustGet(middlewares.ContextUserIDKey).(uuid.UUID)

	resp, err := ctrl.commentService.Create(req, userID)
	if err != nil {
		if errors.Is(err, services.ErrTicketNotFound) || errors.Is(err, services.ErrTaskNotFound) {
			utils.NotFound(c, "entity")
			return
		}
		utils.InternalError(c, "failed to create comment")
		return
	}

	utils.Created(c, resp)
}

func (ctrl *CommentController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid comment id")
		return
	}

	resp, err := ctrl.commentService.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrCommentNotFound) {
			utils.NotFound(c, "comment")
			return
		}
		utils.InternalError(c, "failed to get comment")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}

func (ctrl *CommentController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid comment id")
		return
	}

	var req dto.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if details := validators.ValidateStruct(req); details != nil {
		utils.ValidationError(c, details)
		return
	}

	userID := c.MustGet(middlewares.ContextUserIDKey).(uuid.UUID)

	resp, err := ctrl.commentService.Update(id, req, userID)
	if err != nil {
		if errors.Is(err, services.ErrCommentNotFound) {
			utils.NotFound(c, "comment")
			return
		}
		if errors.Is(err, services.ErrCommentUnauthorized) {
			utils.Forbidden(c, err.Error())
			return
		}
		utils.InternalError(c, "failed to update comment")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}

func (ctrl *CommentController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid comment id")
		return
	}

	userID := c.MustGet(middlewares.ContextUserIDKey).(uuid.UUID)

	if err := ctrl.commentService.Delete(id, userID); err != nil {
		if errors.Is(err, services.ErrCommentNotFound) {
			utils.NotFound(c, "comment")
			return
		}
		if errors.Is(err, services.ErrCommentUnauthorized) {
			utils.Forbidden(c, err.Error())
			return
		}
		utils.InternalError(c, "failed to delete comment")
		return
	}

	utils.NoContent(c)
}

func (ctrl *CommentController) List(c *gin.Context) {
	var filter dto.CommentFilterQuery
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	resp, err := ctrl.commentService.List(filter)
	if err != nil {
		utils.InternalError(c, "failed to list comments")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}
