package schemas

import (
	"time"

	"github.com/google/uuid"
)

type (
	TodoStatus   string
	TodoLabel    *string
	TodoPriority string
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

type Todo struct {
	ID             uuid.UUID    `json:"id"          description:"Unique identifier for the todo"`
	Text           string       `json:"text"        validate:"min=1" description:"The text content of the todo"`
	Status         TodoStatus   `json:"status"      description:"Current status of the todo"`
	Label          TodoLabel    `json:"label,omitempty" description:"Label category for the todo"`
	Priority       TodoPriority `json:"priority"   description:"Priority level of the todo"`
	EstimatedHours float64      `json:"estimatedHours" description:"Estimated hours to complete"`
	ActualHours    float64      `json:"actualHours"   description:"Actual hours spent"`
	Progress       int          `json:"progress"    description:"Progress percentage (0-100)"`
	Cost           float64      `json:"cost"        description:"Cost associated with the todo"`
	DueDate        *time.Time   `json:"dueDate,omitempty" description:"Due date for the todo"`
	CompletedAt    *time.Time   `json:"completedAt,omitempty" description:"Timestamp when completed"`
	CreatedAt      time.Time    `json:"createdAt"   description:"Timestamp when created"`
	UpdatedAt      time.Time    `json:"updatedAt"   description:"Timestamp when last updated"`
}

type CreateTodoRequest struct {
	Text           string        `json:"text"            validate:"required,min=1" description:"The text content for the new todo"`
	Status         *TodoStatus   `json:"status,omitempty" description:"Initial status"`
	Label          TodoLabel     `json:"label,omitempty" description:"Label category"`
	Priority       *TodoPriority `json:"priority,omitempty" description:"Priority level"`
	EstimatedHours *float64      `json:"estimatedHours,omitempty" description:"Estimated hours"`
	ActualHours    *float64      `json:"actualHours,omitempty" description:"Actual hours"`
	Progress       *int          `json:"progress,omitempty" description:"Progress percentage"`
	Cost           *float64      `json:"cost,omitempty" description:"Cost"`
	DueDate        *time.Time    `json:"dueDate,omitempty" description:"Due date"`
}

type UpdateTodoRequest struct {
	Text           *string       `json:"text,omitempty" description:"The updated text content"`
	Status         *TodoStatus   `json:"status,omitempty" description:"Updated status"`
	Label          TodoLabel     `json:"label,omitempty" description:"Label category"`
	Priority       *TodoPriority `json:"priority,omitempty" description:"Priority level"`
	EstimatedHours *float64      `json:"estimatedHours,omitempty" description:"Estimated hours"`
	ActualHours    *float64      `json:"actualHours,omitempty" description:"Actual hours"`
	Progress       *int          `json:"progress,omitempty" description:"Progress percentage"`
	Cost           *float64      `json:"cost,omitempty" description:"Cost"`
	DueDate        *time.Time    `json:"dueDate,omitempty" description:"Due date"`
	CompletedAt    *time.Time    `json:"completedAt,omitempty" description:"Completion timestamp"`
}
