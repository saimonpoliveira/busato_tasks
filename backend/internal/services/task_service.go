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

var ErrTaskNotFound = errors.New("task not found")

type TaskService interface {
	Create(req dto.CreateTaskRequest) (*dto.TaskResponse, error)
	GetByID(id uuid.UUID) (*dto.TaskResponse, error)
	Update(id uuid.UUID, req dto.UpdateTaskRequest) (*dto.TaskResponse, error)
	Delete(id uuid.UUID) error
	List(filter dto.TaskFilterQuery) (dto.PaginatedResponse[dto.TaskResponse], error)
}

type taskService struct {
	taskRepo   repositories.TaskRepository
	ticketRepo repositories.TicketRepository
}

func NewTaskService(taskRepo repositories.TaskRepository, ticketRepo repositories.TicketRepository) TaskService {
	return &taskService{
		taskRepo:   taskRepo,
		ticketRepo: ticketRepo,
	}
}

func (s *taskService) Create(req dto.CreateTaskRequest) (*dto.TaskResponse, error) {
	_, err := s.ticketRepo.FindByID(req.TicketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, fmt.Errorf("failed to verify ticket: %w", err)
	}

	status := req.Status
	if status == "" {
		status = models.TaskStatusTodo
	}

	order := 0
	if req.Order != nil {
		order = *req.Order
	}

	task := &models.Task{
		TicketID:    req.TicketID,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		AssigneeID:  req.AssigneeID,
		Order:       order,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	created, err := s.taskRepo.FindByID(task.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created task: %w", err)
	}

	resp := dto.ToTaskResponse(created)
	return &resp, nil
}

func (s *taskService) GetByID(id uuid.UUID) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to find task: %w", err)
	}

	resp := dto.ToTaskResponse(task)
	return &resp, nil
}

func (s *taskService) Update(id uuid.UUID, req dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to find task: %w", err)
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.AssigneeID != nil {
		task.AssigneeID = req.AssigneeID
	}
	if req.Order != nil {
		task.Order = *req.Order
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	updated, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated task: %w", err)
	}

	resp := dto.ToTaskResponse(updated)
	return &resp, nil
}

func (s *taskService) Delete(id uuid.UUID) error {
	_, err := s.taskRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTaskNotFound
		}
		return fmt.Errorf("failed to find task: %w", err)
	}

	return s.taskRepo.Delete(id)
}

func (s *taskService) List(filter dto.TaskFilterQuery) (dto.PaginatedResponse[dto.TaskResponse], error) {
	tasks, total, err := s.taskRepo.FindAll(filter)
	if err != nil {
		return dto.PaginatedResponse[dto.TaskResponse]{}, fmt.Errorf("failed to list tasks: %w", err)
	}

	return dto.NewPaginatedResponse(dto.ToTaskResponses(tasks), total, filter.Page, filter.PageSize), nil
}
