package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/dto"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/models"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/repositories"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/utils"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type UserService interface {
	Create(req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetByID(id uuid.UUID) (*dto.UserResponse, error)
	Update(id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(id uuid.UUID) error
	List(filter dto.UserFilterQuery) (dto.PaginatedResponse[dto.UserResponse], error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Create(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	role := req.Role
	if role == "" {
		role = models.RoleMember
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: hash,
		Name:         req.Name,
		Role:         role,
		Active:       true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	resp := dto.ToUserResponse(user)
	return &resp, nil
}

func (s *userService) GetByID(id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	resp := dto.ToUserResponse(user)
	return &resp, nil
}

func (s *userService) Update(id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if req.Email != nil {
		existing, err := s.userRepo.FindByEmail(*req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if existing != nil && existing.ID != id {
			return nil, ErrEmailAlreadyExists
		}
		user.Email = *req.Email
	}
	if req.Password != nil {
		hash, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = hash
	}
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.Active != nil {
		user.Active = *req.Active
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	resp := dto.ToUserResponse(user)
	return &resp, nil
}

func (s *userService) Delete(id uuid.UUID) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	return s.userRepo.Delete(id)
}

func (s *userService) List(filter dto.UserFilterQuery) (dto.PaginatedResponse[dto.UserResponse], error) {
	users, total, err := s.userRepo.FindAll(filter)
	if err != nil {
		return dto.PaginatedResponse[dto.UserResponse]{}, fmt.Errorf("failed to list users: %w", err)
	}

	return dto.NewPaginatedResponse(dto.ToUserResponses(users), total, filter.Page, filter.PageSize), nil
}
