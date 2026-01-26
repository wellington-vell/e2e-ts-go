package routers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"server/internal"
	"strconv"
	"time"

	"github.com/wellington-vell/gorpc"
)

type Todo struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTodoInput struct {
	Text      string `json:"text" validate:"required,max=250"`
	Completed bool   `json:"completed"`
}

type UpdateTodoInput struct {
	ID        int    `json:"id" validate:"required,min=1"`
	Text      string `json:"text" validate:"required,max=250"`
	Completed bool   `json:"completed"`
}

type DeleteTodoInput struct {
	ID int `json:"id" validate:"required,min=1"`
}

type GetTodoByIdInput struct {
	ID int `json:"id" validate:"required,min=1"`
}

var TodoRouter = gorpc.Router{
	"getAll": gorpc.OS().
		Input(nil).
		Output([]Todo{}).
		Tag("todos").
		Meta(gorpc.Meta{
			Summary:     "Get all todos",
			Description: "Retrieve a list of all todo items",
		}).
		Route(gorpc.Route{
			Method: "GET",
			Path:   "/todos",
		}).
		Handler(func(ctx *gorpc.Context, input interface{}) (interface{}, error) {
			rows, err := internal.DB.Query("SELECT id, text, completed, created_at, updated_at FROM todos ORDER BY created_at DESC")
			if err != nil {
				return nil, gorpc.NewHTTPError(500, "Failed to fetch todos: "+err.Error())
			}
			defer rows.Close()

			var todos []Todo
			for rows.Next() {
				var todo Todo
				if err := rows.Scan(&todo.ID, &todo.Text, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
					return nil, gorpc.NewHTTPError(500, "Failed to scan todo: "+err.Error())
				}
				todos = append(todos, todo)
			}

			if err := rows.Err(); err != nil {
				return nil, gorpc.NewHTTPError(500, "Error iterating todos: "+err.Error())
			}

			return todos, nil
		}).
		Build(),

	"getById": gorpc.OS().
		Input(GetTodoByIdInput{}).
		Output(Todo{}).
		Tag("todos").
		Meta(gorpc.Meta{
			Summary:     "Get a todo by ID",
			Description: "Retrieve a todo item by its ID",
		}).
		Route(gorpc.Route{
			Method: "GET",
			Path:   "/todos/:id",
		}).
		Handler(func(ctx *gorpc.Context, input interface{}) (interface{}, error) {
			idStr, ok := ctx.Params["id"]
			if !ok {
				return nil, gorpc.NewHTTPError(400, "Missing todo ID parameter")
			}

			id, err := strconv.Atoi(idStr)
			if err != nil {
				return nil, gorpc.NewHTTPError(400, "Invalid todo ID: "+idStr)
			}
			var todo Todo
			err = internal.DB.QueryRow(
				"SELECT id, text, completed, created_at, updated_at FROM todos WHERE id = $1",
				id,
			).Scan(&todo.ID, &todo.Text, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, gorpc.NewHTTPError(404, "Todo not found")
				}
				return nil, gorpc.NewHTTPError(500, "Failed to fetch todo: "+err.Error())
			}

			return todo, nil
		}).
		Build(),

	"create": gorpc.OS().
		Input(CreateTodoInput{}).
		Output(Todo{}).
		Tag("todos").
		Meta(gorpc.Meta{
			Summary:     "Create a new todo",
			Description: "Create a new todo item",
		}).
		Route(gorpc.Route{
			Method: "POST",
			Path:   "/todos",
		}).
		Handler(func(ctx *gorpc.Context, input interface{}) (interface{}, error) {
			var req CreateTodoInput
			if input != nil {
				inputBytes, err := json.Marshal(input)
				if err != nil {
					return nil, gorpc.NewHTTPError(400, "Invalid input format: "+err.Error())
				}
				if err := json.Unmarshal(inputBytes, &req); err != nil {
					return nil, gorpc.NewHTTPError(400, "Invalid input structure: "+err.Error())
				}
			}

			var todo Todo
			err := internal.DB.QueryRow(
				"INSERT INTO todos (text, completed) VALUES ($1, $2) RETURNING id, text, completed, created_at, updated_at",
				req.Text, req.Completed,
			).Scan(&todo.ID, &todo.Text, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

			if err != nil {
				return nil, gorpc.NewHTTPError(500, "Failed to create todo: "+err.Error())
			}

			return todo, nil
		}).
		Build(),

	"update": gorpc.OS().
		Input(UpdateTodoInput{}).
		Output(Todo{}).
		Tag("todos").
		Meta(gorpc.Meta{
			Summary:     "Update a todo",
			Description: "Update an existing todo item by its ID",
		}).
		Route(gorpc.Route{
			Method: "PUT",
			Path:   "/todos/:id",
		}).
		Handler(func(ctx *gorpc.Context, input interface{}) (interface{}, error) {
			idStr, ok := ctx.Params["id"]
			if !ok {
				return nil, gorpc.NewHTTPError(400, "Missing todo ID parameter")
			}
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return nil, gorpc.NewHTTPError(400, "Invalid todo ID: "+idStr)
			}

			var req UpdateTodoInput
			if input != nil {
				inputBytes, err := json.Marshal(input)
				if err != nil {
					return nil, gorpc.NewHTTPError(400, "Invalid input format: "+err.Error())
				}
				if err := json.Unmarshal(inputBytes, &req); err != nil {
					return nil, gorpc.NewHTTPError(400, "Invalid input structure: "+err.Error())
				}
			}

			var todo Todo
			err = internal.DB.QueryRow(
				"UPDATE todos SET text = $1, completed = $2, updated_at = NOW() WHERE id = $3 RETURNING id, text, completed, created_at, updated_at",
				req.Text, req.Completed, id,
			).Scan(&todo.ID, &todo.Text, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, gorpc.NewHTTPError(404, "Todo not found")
				}
				return nil, gorpc.NewHTTPError(500, "Failed to update todo: "+err.Error())
			}

			return todo, nil
		}).
		Build(),

	"delete": gorpc.OS().
		Input(DeleteTodoInput{}).
		Output(nil).
		Tag("todos").
		Meta(gorpc.Meta{
			Summary:     "Delete a todo",
			Description: "Delete a todo item by its ID",
		}).
		Route(gorpc.Route{
			Method: "DELETE",
			Path:   "/todos/:id",
		}).
		Handler(func(ctx *gorpc.Context, input interface{}) (interface{}, error) {
			idStr, ok := ctx.Params["id"]
			if !ok {
				return nil, gorpc.NewHTTPError(400, "Missing todo ID parameter")
			}
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return nil, gorpc.NewHTTPError(400, "Invalid todo ID: "+idStr)
			}

			result, err := internal.DB.Exec("DELETE FROM todos WHERE id = $1", id)
			if err != nil {
				return nil, gorpc.NewHTTPError(500, "Failed to delete todo: "+err.Error())
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return nil, gorpc.NewHTTPError(500, "Failed to check deletion: "+err.Error())
			}

			if rowsAffected == 0 {
				return nil, gorpc.NewHTTPError(404, "Todo not found")
			}

			return map[string]string{"message": "Todo deleted successfully"}, nil
		}).
		Build(),
}
