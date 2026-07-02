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

var ErrProjectNotFound = errors.New("project not found")

type ProjectService interface {
	Create(req dto.CreateProjectRequest, ownerID uuid.UUID) (*dto.ProjectResponse, error)
	GetByID(id uuid.UUID) (*dto.ProjectResponse, error)
	Update(id uuid.UUID, req dto.UpdateProjectRequest) (*dto.ProjectResponse, error)
	Delete(id uuid.UUID) error
	List(filter dto.ProjectFilterQuery) (dto.PaginatedResponse[dto.ProjectResponse], error)
}

type projectService struct {
	projectRepo repositories.ProjectRepository
}

func NewProjectService(projectRepo repositories.ProjectRepository) ProjectService {
	return &projectService{projectRepo: projectRepo}
}

func (s *projectService) Create(req dto.CreateProjectRequest, ownerID uuid.UUID) (*dto.ProjectResponse, error) {
	status := req.Status
	if status == "" {
		status = models.ProjectStatusActive
	}

	project := &models.Project{
		Name:        req.Name,
		Description: req.Description,
		Status:      status,
		OwnerID:     ownerID,
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	created, err := s.projectRepo.FindByID(project.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created project: %w", err)
	}

	resp := dto.ToProjectResponse(created)
	return &resp, nil
}

func (s *projectService) GetByID(id uuid.UUID) (*dto.ProjectResponse, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("failed to find project: %w", err)
	}

	resp := dto.ToProjectResponse(project)
	return &resp, nil
}

func (s *projectService) Update(id uuid.UUID, req dto.UpdateProjectRequest) (*dto.ProjectResponse, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("failed to find project: %w", err)
	}

	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.Description != nil {
		project.Description = *req.Description
	}
	if req.Status != nil {
		project.Status = *req.Status
	}

	if err := s.projectRepo.Update(project); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	updated, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated project: %w", err)
	}

	resp := dto.ToProjectResponse(updated)
	return &resp, nil
}

func (s *projectService) Delete(id uuid.UUID) error {
	_, err := s.projectRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return fmt.Errorf("failed to find project: %w", err)
	}

	return s.projectRepo.Delete(id)
}

func (s *projectService) List(filter dto.ProjectFilterQuery) (dto.PaginatedResponse[dto.ProjectResponse], error) {
	projects, total, err := s.projectRepo.FindAll(filter)
	if err != nil {
		return dto.PaginatedResponse[dto.ProjectResponse]{}, fmt.Errorf("failed to list projects: %w", err)
	}

	return dto.NewPaginatedResponse(dto.ToProjectResponses(projects), total, filter.Page, filter.PageSize), nil
}
