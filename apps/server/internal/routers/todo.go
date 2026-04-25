package routers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"server/internal/db"
	"server/internal/schemas"
)

// @Summary Get all todos
// @Description Retrieves all todo items
// @Tags Todos
// @Produce json
// @Success 200 {array} schemas.Todo
// @Router /api/v1/todos [get]
func HandleGetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := GetAllTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if todos == nil {
		todos = []schemas.Todo{}
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		return
	}
}

func GetAllTodos() ([]schemas.Todo, error) {
	query := `SELECT id, text, status, label, priority, estimated_hours, actual_hours, progress, cost, due_date, completed_at, created_at, updated_at FROM todos ORDER BY created_at DESC`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []schemas.Todo
	for rows.Next() {
		var t schemas.Todo
		if err := rows.Scan(&t.ID, &t.Text, &t.Status, &t.Label, &t.Priority, &t.EstimatedHours, &t.ActualHours, &t.Progress, &t.Cost, &t.DueDate, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

// @Summary Get a todo by ID
// @Description Retrieves a single todo item by its ID
// @Tags Todos
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} schemas.Todo
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
	todo, err := GetTodoByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if todo == nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		return
	}
}

func GetTodoByID(id uuid.UUID) (*schemas.Todo, error) {
	query := `SELECT id, text, status, label, priority, estimated_hours, actual_hours, progress, cost, due_date, completed_at, created_at, updated_at FROM todos WHERE id = $1`
	var t schemas.Todo
	err := db.DB.QueryRow(query, id).Scan(&t.ID, &t.Text, &t.Status, &t.Label, &t.Priority, &t.EstimatedHours, &t.ActualHours, &t.Progress, &t.Cost, &t.DueDate, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// @Summary Create a new todo
// @Description Creates a new todo item with the provided text
// @Tags Todos
// @Accept json
// @Produce json
// @Param request body schemas.CreateTodoRequest true "Todo creation request"
// @Success 201 {object} schemas.Todo
// @Failure 400 {string} string "Invalid request body or missing text"
// @Router /api/v1/todos [post]
func HandleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Text == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}
	todo, err := CreateTodo(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		return
	}
}

func CreateTodo(req *schemas.CreateTodoRequest) (*schemas.Todo, error) {
	status := schemas.TodoStatusBacklog
	if req.Status != nil {
		status = *req.Status
	}
	priority := schemas.TodoPriorityMedium
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
	var t schemas.Todo
	err := db.DB.QueryRow(query, req.Text, status, req.Label, priority, estimatedHours, actualHours, progress, cost, req.DueDate).Scan(&t.ID, &t.Text, &t.Status, &t.Label, &t.Priority, &t.EstimatedHours, &t.ActualHours, &t.Progress, &t.Cost, &t.DueDate, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// @Summary Update a todo
// @Description Updates an existing todo item by its ID
// @Tags Todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param request body schemas.UpdateTodoRequest true "Todo update request"
// @Success 200 {object} schemas.Todo
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
	var req schemas.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	todo, err := UpdateTodo(id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if todo == nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		return
	}
}

func UpdateTodo(id uuid.UUID, req *schemas.UpdateTodoRequest) (*schemas.Todo, error) {
	t, err := GetTodoByID(id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}

	if req.Text != nil {
		t.Text = *req.Text
	}
	if req.Status != nil {
		t.Status = *req.Status
	}
	if req.Label != nil {
		t.Label = req.Label
	}
	if req.Priority != nil {
		t.Priority = *req.Priority
	}
	if req.EstimatedHours != nil {
		t.EstimatedHours = *req.EstimatedHours
	}
	if req.ActualHours != nil {
		t.ActualHours = *req.ActualHours
	}
	if req.Progress != nil {
		t.Progress = *req.Progress
	}
	if req.Cost != nil {
		t.Cost = *req.Cost
	}
	if req.DueDate != nil {
		t.DueDate = req.DueDate
	}
	if req.CompletedAt != nil {
		t.CompletedAt = req.CompletedAt
	}
	t.UpdatedAt = time.Now()

	query := `UPDATE todos SET text = $1, status = $2, label = $3, priority = $4, estimated_hours = $5, actual_hours = $6, progress = $7, cost = $8, due_date = $9, completed_at = $10, updated_at = $11 WHERE id = $12`
	_, err = db.DB.Exec(query, t.Text, t.Status, t.Label, t.Priority, t.EstimatedHours, t.ActualHours, t.Progress, t.Cost, t.DueDate, t.CompletedAt, t.UpdatedAt, id)
	if err != nil {
		return nil, err
	}

	return t, nil
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
	deleted, err := DeleteTodo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteTodo(id uuid.UUID) (bool, error) {
	query := `DELETE FROM todos WHERE id = $1`
	result, err := db.DB.Exec(query, id)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}
