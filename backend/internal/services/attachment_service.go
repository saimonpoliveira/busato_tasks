package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/dto"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/models"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/repositories"
	"gorm.io/gorm"
)

var (
	ErrAttachmentNotFound    = errors.New("attachment not found")
	ErrFileTooLarge          = errors.New("file exceeds maximum upload size")
	ErrInvalidFileType       = errors.New("invalid file type")
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
	"application/pdf": true,
	"text/plain":      true,
	"application/zip": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       true,
}

type AttachmentService interface {
	Upload(entityType models.EntityType, entityID uuid.UUID, file *multipart.FileHeader, userID uuid.UUID) (*dto.AttachmentResponse, error)
	GetByID(id uuid.UUID) (*dto.AttachmentResponse, error)
	Delete(id uuid.UUID) error
	List(filter dto.AttachmentFilterQuery) (dto.PaginatedResponse[dto.AttachmentResponse], error)
	GetFilePath(id uuid.UUID) (string, string, error)
}

type attachmentService struct {
	attachmentRepo repositories.AttachmentRepository
	ticketRepo     repositories.TicketRepository
	taskRepo       repositories.TaskRepository
	uploadDir      string
	maxSizeBytes   int64
}

func NewAttachmentService(
	attachmentRepo repositories.AttachmentRepository,
	ticketRepo repositories.TicketRepository,
	taskRepo repositories.TaskRepository,
	uploadDir string,
	maxSizeMB int64,
) AttachmentService {
	return &attachmentService{
		attachmentRepo: attachmentRepo,
		ticketRepo:     ticketRepo,
		taskRepo:       taskRepo,
		uploadDir:      uploadDir,
		maxSizeBytes:   maxSizeMB * 1024 * 1024,
	}
}

func (s *attachmentService) validateEntity(entityType models.EntityType, entityID uuid.UUID) error {
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

func (s *attachmentService) Upload(entityType models.EntityType, entityID uuid.UUID, file *multipart.FileHeader, userID uuid.UUID) (*dto.AttachmentResponse, error) {
	if err := s.validateEntity(entityType, entityID); err != nil {
		return nil, err
	}

	if file.Size > s.maxSizeBytes {
		return nil, ErrFileTooLarge
	}

	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	if !allowedMimeTypes[mimeType] {
		return nil, ErrInvalidFileType
	}

	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	filePath := filepath.Join(s.uploadDir, filename)

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	attachment := &models.Attachment{
		EntityType:   entityType,
		EntityID:     entityID,
		Filename:     filename,
		OriginalName: file.Filename,
		FilePath:     filePath,
		Size:         file.Size,
		MimeType:     mimeType,
		UploadedByID: userID,
	}

	if err := s.attachmentRepo.Create(attachment); err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save attachment record: %w", err)
	}

	created, err := s.attachmentRepo.FindByID(attachment.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created attachment: %w", err)
	}

	resp := dto.ToAttachmentResponse(created)
	return &resp, nil
}

func (s *attachmentService) GetByID(id uuid.UUID) (*dto.AttachmentResponse, error) {
	attachment, err := s.attachmentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAttachmentNotFound
		}
		return nil, fmt.Errorf("failed to find attachment: %w", err)
	}

	resp := dto.ToAttachmentResponse(attachment)
	return &resp, nil
}

func (s *attachmentService) Delete(id uuid.UUID) error {
	attachment, err := s.attachmentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAttachmentNotFound
		}
		return fmt.Errorf("failed to find attachment: %w", err)
	}

	if err := s.attachmentRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete attachment record: %w", err)
	}

	if err := os.Remove(attachment.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (s *attachmentService) List(filter dto.AttachmentFilterQuery) (dto.PaginatedResponse[dto.AttachmentResponse], error) {
	attachments, total, err := s.attachmentRepo.FindAll(filter)
	if err != nil {
		return dto.PaginatedResponse[dto.AttachmentResponse]{}, fmt.Errorf("failed to list attachments: %w", err)
	}

	return dto.NewPaginatedResponse(dto.ToAttachmentResponses(attachments), total, filter.Page, filter.PageSize), nil
}

func (s *attachmentService) GetFilePath(id uuid.UUID) (string, string, error) {
	attachment, err := s.attachmentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", ErrAttachmentNotFound
		}
		return "", "", fmt.Errorf("failed to find attachment: %w", err)
	}

	if !strings.HasPrefix(attachment.FilePath, s.uploadDir) {
		return "", "", fmt.Errorf("invalid file path")
	}

	return attachment.FilePath, attachment.OriginalName, nil
}
