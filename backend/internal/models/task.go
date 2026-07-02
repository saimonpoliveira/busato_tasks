package models

import "github.com/google/uuid"

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

type Task struct {
	BaseModel
	TicketID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"ticket_id"`
	Ticket      Ticket     `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	Title       string     `gorm:"not null;size:500" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Status      TaskStatus `gorm:"type:varchar(20);not null;default:'todo'" json:"status"`
	AssigneeID  *uuid.UUID `gorm:"type:uuid;index" json:"assignee_id,omitempty"`
	Assignee    *User      `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
	Order       int        `gorm:"default:0" json:"order"`
}
