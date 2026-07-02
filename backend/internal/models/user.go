package models

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleMember UserRole = "member"
)

type User struct {
	BaseModel
	Email        string   `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash string   `gorm:"not null" json:"-"`
	Name         string   `gorm:"not null;size:255" json:"name"`
	Role         UserRole `gorm:"type:varchar(20);not null;default:'member'" json:"role"`
	Active       bool     `gorm:"default:true" json:"active"`
}
