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

type TicketController struct {
	ticketService services.TicketService
}

func NewTicketController(ticketService services.TicketService) *TicketController {
	return &TicketController{ticketService: ticketService}
}

func (ctrl *TicketController) Create(c *gin.Context) {
	var req dto.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if details := validators.ValidateStruct(req); details != nil {
		utils.ValidationError(c, details)
		return
	}

	reporterID := c.MustGet(middlewares.ContextUserIDKey).(uuid.UUID)

	resp, err := ctrl.ticketService.Create(req, reporterID)
	if err != nil {
		if errors.Is(err, services.ErrProjectNotFound) {
			utils.NotFound(c, "project")
			return
		}
		utils.InternalError(c, "failed to create ticket")
		return
	}

	utils.Created(c, resp)
}

func (ctrl *TicketController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid ticket id")
		return
	}

	resp, err := ctrl.ticketService.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrTicketNotFound) {
			utils.NotFound(c, "ticket")
			return
		}
		utils.InternalError(c, "failed to get ticket")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}

func (ctrl *TicketController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid ticket id")
		return
	}

	var req dto.UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if details := validators.ValidateStruct(req); details != nil {
		utils.ValidationError(c, details)
		return
	}

	resp, err := ctrl.ticketService.Update(id, req)
	if err != nil {
		if errors.Is(err, services.ErrTicketNotFound) {
			utils.NotFound(c, "ticket")
			return
		}
		utils.InternalError(c, "failed to update ticket")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}

func (ctrl *TicketController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid ticket id")
		return
	}

	if err := ctrl.ticketService.Delete(id); err != nil {
		if errors.Is(err, services.ErrTicketNotFound) {
			utils.NotFound(c, "ticket")
			return
		}
		utils.InternalError(c, "failed to delete ticket")
		return
	}

	utils.NoContent(c)
}

func (ctrl *TicketController) List(c *gin.Context) {
	var filter dto.TicketFilterQuery
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	resp, err := ctrl.ticketService.List(filter)
	if err != nil {
		utils.InternalError(c, "failed to list tickets")
		return
	}

	utils.Success(c, http.StatusOK, resp)
}
