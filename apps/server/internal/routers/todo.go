package routers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

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
func HandleGetTodos(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, text, status, label, priority, estimated_hours, actual_hours, progress, cost, due_date, completed_at, created_at, updated_at FROM todos ORDER BY created_at DESC`
	rows, err := db.Query.QueryContext(r.Context(), query)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Text, &t.Status, &t.Label, &t.Priority, &t.EstimatedHours, &t.ActualHours, &t.Progress, &t.Cost, &t.DueDate, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		todos = append(todos, t)
	}
	if err := rows.Err(); err != nil {
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
func HandleGetTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	query := `SELECT id, text, status, label, priority, estimated_hours, actual_hours, progress, cost, due_date, completed_at, created_at, updated_at FROM todos WHERE id = $1`
	var t models.Todo
	err = db.Query.QueryRowContext(r.Context(), query, id).Scan(&t.ID, &t.Text, &t.Status, &t.Label, &t.Priority, &t.EstimatedHours, &t.ActualHours, &t.Progress, &t.Cost, &t.DueDate, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt)
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
		log.Printf("JSON encode error in HandleGetTodo: %v", err)
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
func HandleCreateTodo(w http.ResponseWriter, r *http.Request) {
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

	query := `INSERT INTO todos (text, status, label, priority, estimated_hours, actual_hours, progress, cost, due_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, text, status, label, priority, estimated_hours, actual_hours, progress, cost, due_date, completed_at, created_at, updated_at`
	var t models.Todo
	err := db.Query.QueryRowContext(r.Context(), query, req.Text, status, req.Label, priority, estimatedHours, actualHours, progress, cost, req.DueDate).Scan(&t.ID, &t.Text, &t.Status, &t.Label, &t.Priority, &t.EstimatedHours, &t.ActualHours, &t.Progress, &t.Cost, &t.DueDate, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, db.ErrUniqueViolation) {
			http.Error(w, "Duplicate entry", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&t); err != nil {
		log.Printf("JSON encode error in HandleCreateTodo: %v", err)
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
func HandleUpdateTodo(w http.ResponseWriter, r *http.Request) {
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
	query := `
		UPDATE todos SET
			text = COALESCE($1, text),
			status = COALESCE($2, status),
			label = COALESCE($3, label),
			priority = COALESCE($4, priority),
			estimated_hours = COALESCE($5, estimated_hours),
			actual_hours = COALESCE($6, actual_hours),
			progress = COALESCE($7, progress),
			cost = COALESCE($8, cost),
			due_date = COALESCE($9, due_date),
			completed_at = COALESCE($10, completed_at),
			updated_at = NOW()
		WHERE id = $11
		RETURNING id, text, status, label, priority, estimated_hours, actual_hours, progress, cost, due_date, completed_at, created_at, updated_at
	`
	var t models.Todo
	err = db.Query.QueryRowContext(r.Context(), query,
		req.Text,
		req.Status,
		req.Label,
		req.Priority,
		req.EstimatedHours,
		req.ActualHours,
		req.Progress,
		req.Cost,
		req.DueDate,
		req.CompletedAt,
		id,
	).Scan(&t.ID, &t.Text, &t.Status, &t.Label, &t.Priority, &t.EstimatedHours, &t.ActualHours, &t.Progress, &t.Cost, &t.DueDate, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt)

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
		log.Printf("JSON encode error in HandleUpdateTodo: %v", err)
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
func HandleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	query := `DELETE FROM todos WHERE id = $1`
	result, err := db.Query.ExecContext(r.Context(), query, id)
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
