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

var (
	ErrCommentNotFound      = errors.New("comment not found")
	ErrCommentUnauthorized  = errors.New("not authorized to modify this comment")
)

type CommentService interface {
	Create(req dto.CreateCommentRequest, userID uuid.UUID) (*dto.CommentResponse, error)
	GetByID(id uuid.UUID) (*dto.CommentResponse, error)
	Update(id uuid.UUID, req dto.UpdateCommentRequest, userID uuid.UUID) (*dto.CommentResponse, error)
	Delete(id uuid.UUID, userID uuid.UUID) error
	List(filter dto.CommentFilterQuery) (dto.PaginatedResponse[dto.CommentResponse], error)
}

type commentService struct {
	commentRepo repositories.CommentRepository
	ticketRepo  repositories.TicketRepository
	taskRepo    repositories.TaskRepository
}

func NewCommentService(
	commentRepo repositories.CommentRepository,
	ticketRepo repositories.TicketRepository,
	taskRepo repositories.TaskRepository,
) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		ticketRepo:  ticketRepo,
		taskRepo:    taskRepo,
	}
}

func (s *commentService) validateEntity(entityType models.EntityType, entityID uuid.UUID) error {
	switch entityType {
	case models.EntityTypeTicket:
		_, err := s.ticketRepo.FindByID(entityID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTicketNotFound
			}
			return fmt.Errorf("failed to verify ticket: %w", err)
		}
	case models.EntityTypeTask:
		_, err := s.taskRepo.FindByID(entityID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTaskNotFound
			}
			return fmt.Errorf("failed to verify task: %w", err)
		}
	default:
		return fmt.Errorf("invalid entity type")
	}
	return nil
}

func (s *commentService) Create(req dto.CreateCommentRequest, userID uuid.UUID) (*dto.CommentResponse, error) {
	if err := s.validateEntity(req.EntityType, req.EntityID); err != nil {
		return nil, err
	}

	comment := &models.Comment{
		EntityType: req.EntityType,
		EntityID:   req.EntityID,
		UserID:     userID,
		Content:    req.Content,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	created, err := s.commentRepo.FindByID(comment.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created comment: %w", err)
	}

	resp := dto.ToCommentResponse(created)
	return &resp, nil
}

func (s *commentService) GetByID(id uuid.UUID) (*dto.CommentResponse, error) {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}
		return nil, fmt.Errorf("failed to find comment: %w", err)
	}

	resp := dto.ToCommentResponse(comment)
	return &resp, nil
}

func (s *commentService) Update(id uuid.UUID, req dto.UpdateCommentRequest, userID uuid.UUID) (*dto.CommentResponse, error) {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}
		return nil, fmt.Errorf("failed to find comment: %w", err)
	}

	if comment.UserID != userID {
		return nil, ErrCommentUnauthorized
	}

	comment.Content = req.Content

	if err := s.commentRepo.Update(comment); err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	updated, err := s.commentRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated comment: %w", err)
	}

	resp := dto.ToCommentResponse(updated)
	return &resp, nil
}

func (s *commentService) Delete(id uuid.UUID, userID uuid.UUID) error {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotFound
		}
		return fmt.Errorf("failed to find comment: %w", err)
	}

	if comment.UserID != userID {
		return ErrCommentUnauthorized
	}

	return s.commentRepo.Delete(id)
}

func (s *commentService) List(filter dto.CommentFilterQuery) (dto.PaginatedResponse[dto.CommentResponse], error) {
	comments, total, err := s.commentRepo.FindAll(filter)
	if err != nil {
		return dto.PaginatedResponse[dto.CommentResponse]{}, fmt.Errorf("failed to list comments: %w", err)
	}

	return dto.NewPaginatedResponse(dto.ToCommentResponses(comments), total, filter.Page, filter.PageSize), nil
}
