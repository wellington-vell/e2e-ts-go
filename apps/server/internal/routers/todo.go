package routers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"server/internal/db"
	"server/internal/models"

	"github.com/google/uuid"
)

// @Summary Get all todos
// @Description Retrieves all todo items
// @Tags Todos
// @Produce json
// @Success 200 {array} models.Todo
// @Router /api/v1/todos [get]
func GetTodos(w http.ResponseWriter, r *http.Request) {
	var todos []models.Todo
	err := db.DB.NewSelect().
		Model(&todos).
		Order("created_at DESC").
		Scan(r.Context())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if todos == nil {
		todos = []models.Todo{}
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		log.Printf("JSON encode error in HandleGetTodos: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// @Summary Get a todo by ID
// @Description Retrieves a single todo item by its ID
// @Tags Todos
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Todo not found"
// @Router /api/v1/todos/{id} [get]
func GetTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var t models.Todo
	err = db.DB.NewSelect().
		Model(&t).
		Where("id = ?", id).
		Scan(r.Context())
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&t); err != nil {
		log.Printf("JSON encode error in GetTodo: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// @Summary Create a new todo
// @Description Creates a new todo item with the provided text
// @Tags Todos
// @Accept json
// @Produce json
// @Param request body models.CreateTodoRequest true "Todo creation request"
// @Success 201 {object} models.Todo
// @Failure 400 {string} string "Invalid request body or missing text"
// @Router /api/v1/todos [post]
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	status := models.TodoStatusBacklog
	if req.Status != nil {
		status = *req.Status
	}
	priority := models.TodoPriorityMedium
	if req.Priority != nil {
		priority = *req.Priority
	}
	estimatedHours := 0.0
	if req.EstimatedHours != nil {
		estimatedHours = *req.EstimatedHours
	}
	actualHours := 0.0
	if req.ActualHours != nil {
		actualHours = *req.ActualHours
	}
	progress := 0
	if req.Progress != nil {
		progress = *req.Progress
	}
	cost := 0.0
	if req.Cost != nil {
		cost = *req.Cost
	}

	t := models.Todo{
		Text:           req.Text,
		Status:         status,
		Label:          req.Label,
		Priority:       priority,
		EstimatedHours: estimatedHours,
		ActualHours:    actualHours,
		Progress:       progress,
		Cost:           cost,
		DueDate:        req.DueDate,
	}

	_, err := db.DB.NewInsert().
		Model(&t).
		Returning("*").
		Exec(r.Context())
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			http.Error(w, "Duplicate entry", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&t); err != nil {
		log.Printf("JSON encode error in CreateTodo: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// @Summary Update a todo
// @Description Updates an existing todo item by its ID
// @Tags Todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param request body models.UpdateTodoRequest true "Todo update request"
// @Success 200 {object} models.Todo
// @Failure 400 {string} string "Invalid ID or request body"
// @Failure 404 {string} string "Todo not found"
// @Router /api/v1/todos/{id} [put]
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var t models.Todo
	query := db.DB.NewUpdate().
		Model(&t).
		Where("id = ?", id).
		Returning("*")

	if req.Text != nil {
		query.Set("text = ?", *req.Text)
	}
	if req.Status != nil {
		query.Set("status = ?", *req.Status)
	}
	if req.Label != nil {
		query.Set("label = ?", req.Label)
	}
	if req.Priority != nil {
		query.Set("priority = ?", *req.Priority)
	}
	if req.EstimatedHours != nil {
		query.Set("estimated_hours = ?", *req.EstimatedHours)
	}
	if req.ActualHours != nil {
		query.Set("actual_hours = ?", *req.ActualHours)
	}
	if req.Progress != nil {
		query.Set("progress = ?", *req.Progress)
	}
	if req.Cost != nil {
		query.Set("cost = ?", *req.Cost)
	}
	if req.DueDate != nil {
		query.Set("due_date = ?", req.DueDate)
	}
	if req.CompletedAt != nil {
		query.Set("completed_at = ?", req.CompletedAt)
	}

	query.Set("updated_at = NOW()")

	err = query.Scan(r.Context())
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&t); err != nil {
		log.Printf("JSON encode error in UpdateTodo: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// @Summary Delete a todo
// @Description Deletes a todo item by its ID
// @Tags Todos
// @Param id path string true "Todo ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Todo not found"
// @Router /api/v1/todos/{id} [delete]
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	result, err := db.DB.NewDelete().
		Model((*models.Todo)(nil)).
		Where("id = ?", id).
		Exec(r.Context())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
