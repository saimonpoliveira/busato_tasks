package repositories

import (
	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/dto"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/models"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/utils"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	FindAll(filter dto.UserFilterQuery) ([]models.User, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

var userSortFields = map[string]string{
	"name":       "name",
	"email":      "email",
	"role":       "role",
	"created_at": "created_at",
}

func (r *userRepository) FindAll(filter dto.UserFilterQuery) ([]models.User, int64, error) {
	filter.Normalize()

	query := r.db.Model(&models.User{})

	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}
	if filter.Active != nil {
		query = query.Where("active = ?", *filter.Active)
	}

	query = utils.ApplySearch(query, filter.Search, "name", "email")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = utils.ApplySorting(query, filter.SortBy, filter.SortOrder, userSortFields)
	query = utils.ApplyPagination(query, filter.Page, filter.PageSize)

	var users []models.User
	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
