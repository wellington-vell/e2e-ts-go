package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type (
	TodoStatus   string
	TodoLabel    *string
	TodoPriority string
	Date         struct {
		time.Time
	}
)

const (
	TodoStatusBacklog    TodoStatus = "backlog"
	TodoStatusTodo       TodoStatus = "todo"
	TodoStatusInProgress TodoStatus = "in_progress"
	TodoStatusDone       TodoStatus = "done"
	TodoStatusCanceled   TodoStatus = "canceled"
)

const (
	TodoLabelBug     string = "bug"
	TodoLabelFeature string = "feature"
	TodoLabelDoc     string = "doc"
)

const (
	TodoPriorityLow    TodoPriority = "low"
	TodoPriorityMedium TodoPriority = "medium"
	TodoPriorityHigh   TodoPriority = "high"
)

func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format("2006-01-02"))
}

func (d *Date) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	return d.Time, nil
}

type Todo struct {
	ID             uuid.UUID    `json:"id"                        description:"Unique identifier for the todo"`
	Text           string       `json:"text"        validate:"required,min=1,nonwhitespace" description:"The text content of the todo"`
	Status         TodoStatus   `json:"status"      description:"Current status of the todo"`
	Label          TodoLabel    `json:"label,omitempty" description:"Label category for the todo"`
	Priority       TodoPriority `json:"priority"   description:"Priority level of the todo"`
	EstimatedHours float64      `json:"estimatedHours" description:"Estimated hours to complete"`
	ActualHours    float64      `json:"actualHours"   description:"Actual hours spent"`
	Progress       int          `json:"progress"    description:"Progress percentage (0-100)"`
	Cost           float64      `json:"cost"        description:"Cost associated with the todo"`
	DueDate        *Date        `json:"dueDate,omitempty" swaggertype:"string" format:"date" description:"Due date for the todo"`
	CompletedAt    *time.Time   `json:"completedAt,omitempty" swaggertype:"string" format:"date" description:"Timestamp when completed"`
	CreatedAt      time.Time    `json:"createdAt" swaggertype:"string" format:"date"  description:"Timestamp when created"`
	UpdatedAt      time.Time    `json:"updatedAt" swaggertype:"string" format:"date"   description:"Timestamp when last updated"`
}

type CreateTodoRequest struct {
	Text           string        `json:"text"            validate:"required,min=1,nonwhitespace" description:"The text content for the new todo"`
	Status         *TodoStatus   `json:"status,omitempty" validate:"omitempty,oneof=backlog todo in_progress done canceled" description:"Initial status"`
	Label          TodoLabel     `json:"label,omitempty" validate:"omitempty,oneof=bug feature doc" description:"Label category"`
	Priority       *TodoPriority `json:"priority,omitempty" validate:"omitempty,oneof=low medium high" description:"Priority level"`
	EstimatedHours *float64      `json:"estimatedHours,omitempty" validate:"omitempty,gte=0" description:"Estimated hours"`
	ActualHours    *float64      `json:"actualHours,omitempty" validate:"omitempty,gte=0" description:"Actual hours"`
	Progress       *int          `json:"progress,omitempty" validate:"omitempty,gte=0,lte=100" description:"Progress percentage"`
	Cost           *float64      `json:"cost,omitempty" validate:"omitempty,gte=0" description:"Cost"`
	DueDate        *Date         `json:"dueDate,omitempty" swaggertype:"string" format:"date" description:"Due date"`
}

type UpdateTodoRequest struct {
	Text           *string       `json:"text,omitempty" validate:"omitempty,min=1,nonwhitespace" description:"The updated text content"`
	Status         *TodoStatus   `json:"status,omitempty" validate:"omitempty,oneof=backlog todo in_progress done canceled" description:"Updated status"`
	Label          TodoLabel     `json:"label,omitempty" validate:"omitempty,oneof=bug feature doc" description:"Label category"`
	Priority       *TodoPriority `json:"priority,omitempty" validate:"omitempty,oneof=low medium high" description:"Priority level"`
	EstimatedHours *float64      `json:"estimatedHours,omitempty" validate:"omitempty,gte=0" description:"Estimated hours"`
	ActualHours    *float64      `json:"actualHours,omitempty" validate:"omitempty,gte=0" description:"Actual hours"`
	Progress       *int          `json:"progress,omitempty" validate:"omitempty,gte=0,lte=100" description:"Progress percentage"`
	Cost           *float64      `json:"cost,omitempty" validate:"omitempty,gte=0" description:"Cost"`
	DueDate        *Date         `json:"dueDate,omitempty" swaggertype:"string" format:"date" description:"Due date"`
	CompletedAt    *time.Time    `json:"completedAt,omitempty" description:"Completion timestamp"`
}
