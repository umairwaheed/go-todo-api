package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/umair/go-todo-api/models"
)

// TodoHandler handles HTTP requests for todo operations
type TodoHandler struct {
	todoModel *models.TodoModel
}

// NewTodoHandler creates a new TodoHandler instance
func NewTodoHandler(todoModel *models.TodoModel) *TodoHandler {
	return &TodoHandler{
		todoModel: todoModel,
	}
}

// GetTodos handles GET /todos - retrieves all todos
func (h *TodoHandler) GetTodos(c *gin.Context) {
	todos, err := h.todoModel.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve todos"})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// GetTodo handles GET /todos/:id - retrieves a specific todo
func (h *TodoHandler) GetTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	todo, err := h.todoModel.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// CreateTodo handles POST /todos - creates a new todo
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req models.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	todo, err := h.todoModel.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// UpdateTodo handles PUT /todos/:id - updates an existing todo
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	var req models.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	todo, err := h.todoModel.Update(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// DeleteTodo handles DELETE /todos/:id - deletes a todo
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	// Check if todo exists before deleting
	_, err = h.todoModel.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	err = h.todoModel.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}

// CompleteTodo handles PATCH /todos/:id/complete - marks a todo as completed
func (h *TodoHandler) CompleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	todo, err := h.todoModel.ToggleComplete(id, true)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// UncompleteTodo handles PATCH /todos/:id/uncomplete - marks a todo as incomplete
func (h *TodoHandler) UncompleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	todo, err := h.todoModel.ToggleComplete(id, false)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}
