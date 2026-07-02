package models

import "github.com/google/uuid"

type TicketStatus string

const (
	TicketStatusOpen       TicketStatus = "open"
	TicketStatusInProgress TicketStatus = "in_progress"
	TicketStatusResolved   TicketStatus = "resolved"
	TicketStatusClosed     TicketStatus = "closed"
)

type TicketPriority string

const (
	TicketPriorityLow      TicketPriority = "low"
	TicketPriorityMedium   TicketPriority = "medium"
	TicketPriorityHigh     TicketPriority = "high"
	TicketPriorityCritical TicketPriority = "critical"
)

type Ticket struct {
	BaseModel
	ProjectID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	Project     Project        `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Title       string         `gorm:"not null;size:500" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Status      TicketStatus   `gorm:"type:varchar(20);not null;default:'open'" json:"status"`
	Priority    TicketPriority `gorm:"type:varchar(20);not null;default:'medium'" json:"priority"`
	AssigneeID  *uuid.UUID     `gorm:"type:uuid;index" json:"assignee_id,omitempty"`
	Assignee    *User          `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
	ReporterID  uuid.UUID      `gorm:"type:uuid;not null" json:"reporter_id"`
	Reporter    User           `gorm:"foreignKey:ReporterID" json:"reporter,omitempty"`
}
