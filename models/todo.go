package models

import (
	"database/sql"
	"time"
)

// Todo represents a todo item
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTodoRequest represents the request body for creating a todo
type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// UpdateTodoRequest represents the request body for updating a todo
type UpdateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// TodoModel handles database operations for todos
type TodoModel struct {
	DB *sql.DB
}

// NewTodoModel creates a new TodoModel instance
func NewTodoModel(db *sql.DB) *TodoModel {
	return &TodoModel{DB: db}
}

// Create inserts a new todo into the database
func (m *TodoModel) Create(req CreateTodoRequest) (*Todo, error) {
	query := `
		INSERT INTO todos (title, description, completed, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := m.DB.Exec(query, req.Title, req.Description, false, now, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Todo{
		ID:          int(id),
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// GetByID retrieves a todo by its ID
func (m *TodoModel) GetByID(id int) (*Todo, error) {
	query := `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos WHERE id = ?
	`

	todo := &Todo{}
	err := m.DB.QueryRow(query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return todo, nil
}

// GetAll retrieves all todos from the database
func (m *TodoModel) GetAll() ([]*Todo, error) {
	query := `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos ORDER BY created_at DESC
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*Todo
	for rows.Next() {
		todo := &Todo{}
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// Update modifies an existing todo
func (m *TodoModel) Update(id int, req UpdateTodoRequest) (*Todo, error) {
	query := `
		UPDATE todos 
		SET title = ?, description = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	_, err := m.DB.Exec(query, req.Title, req.Description, now, id)
	if err != nil {
		return nil, err
	}

	// Return the updated todo
	return m.GetByID(id)
}

// Delete removes a todo from the database
func (m *TodoModel) Delete(id int) error {
	query := `DELETE FROM todos WHERE id = ?`
	_, err := m.DB.Exec(query, id)
	return err
}

// ToggleComplete toggles the completed status of a todo
func (m *TodoModel) ToggleComplete(id int, completed bool) (*Todo, error) {
	query := `
		UPDATE todos 
		SET completed = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	_, err := m.DB.Exec(query, completed, now, id)
	if err != nil {
		return nil, err
	}

	// Return the updated todo
	return m.GetByID(id)
}
