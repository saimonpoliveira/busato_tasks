package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/dto"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/models"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/repositories"
	"gorm.io/gorm"
)

var ErrTicketNotFound = errors.New("ticket not found")

type TicketService interface {
	Create(req dto.CreateTicketRequest, reporterID uuid.UUID) (*dto.TicketResponse, error)
	GetByID(id uuid.UUID) (*dto.TicketResponse, error)
	Update(id uuid.UUID, req dto.UpdateTicketRequest) (*dto.TicketResponse, error)
	Delete(id uuid.UUID) error
	List(filter dto.TicketFilterQuery) (dto.PaginatedResponse[dto.TicketResponse], error)
}

type ticketService struct {
	ticketRepo  repositories.TicketRepository
	projectRepo repositories.ProjectRepository
}

func NewTicketService(ticketRepo repositories.TicketRepository, projectRepo repositories.ProjectRepository) TicketService {
	return &ticketService{
		ticketRepo:  ticketRepo,
		projectRepo: projectRepo,
	}
}

func (s *ticketService) Create(req dto.CreateTicketRequest, reporterID uuid.UUID) (*dto.TicketResponse, error) {
	_, err := s.projectRepo.FindByID(req.ProjectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("failed to verify project: %w", err)
	}

	status := req.Status
	if status == "" {
		status = models.TicketStatusOpen
	}
	priority := req.Priority
	if priority == "" {
		priority = models.TicketPriorityMedium
	}

	ticket := &models.Ticket{
		ProjectID:   req.ProjectID,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Priority:    priority,
		AssigneeID:  req.AssigneeID,
		ReporterID:  reporterID,
	}

	if err := s.ticketRepo.Create(ticket); err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	created, err := s.ticketRepo.FindByID(ticket.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created ticket: %w", err)
	}

	resp := dto.ToTicketResponse(created)
	return &resp, nil
}

func (s *ticketService) GetByID(id uuid.UUID) (*dto.TicketResponse, error) {
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to find ticket: %w", err)
	}

	resp := dto.ToTicketResponse(ticket)
	return &resp, nil
}

func (s *ticketService) Update(id uuid.UUID, req dto.UpdateTicketRequest) (*dto.TicketResponse, error) {
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to find ticket: %w", err)
	}

	if req.Title != nil {
		ticket.Title = *req.Title
	}
	if req.Description != nil {
		ticket.Description = *req.Description
	}
	if req.Status != nil {
		ticket.Status = *req.Status
	}
	if req.Priority != nil {
		ticket.Priority = *req.Priority
	}
	if req.AssigneeID != nil {
		ticket.AssigneeID = req.AssigneeID
	}

	if err := s.ticketRepo.Update(ticket); err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	updated, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated ticket: %w", err)
	}

	resp := dto.ToTicketResponse(updated)
	return &resp, nil
}

func (s *ticketService) Delete(id uuid.UUID) error {
	_, err := s.ticketRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTicketNotFound
		}
		return fmt.Errorf("failed to find ticket: %w", err)
	}

	return s.ticketRepo.Delete(id)
}

func (s *ticketService) List(filter dto.TicketFilterQuery) (dto.PaginatedResponse[dto.TicketResponse], error) {
	tickets, total, err := s.ticketRepo.FindAll(filter)
	if err != nil {
		return dto.PaginatedResponse[dto.TicketResponse]{}, fmt.Errorf("failed to list tickets: %w", err)
	}

	return dto.NewPaginatedResponse(dto.ToTicketResponses(tickets), total, filter.Page, filter.PageSize), nil
}
