package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/umair/go-todo-api/database"
	"github.com/umair/go-todo-api/handlers"
	"github.com/umair/go-todo-api/models"
)

// TestTodoModel tests the TodoModel database operations
func TestTodoModel(t *testing.T) {
	// Use a test database
	dbPath := "test_todo.db"
	defer os.Remove(dbPath)

	db, err := database.InitDB(dbPath)
	assert.NoError(t, err)
	defer database.CloseDB(db)

	todoModel := models.NewTodoModel(db)

	t.Run("Create Todo", func(t *testing.T) {
		req := models.CreateTodoRequest{
			Title:       "Test Todo",
			Description: "Test Description",
		}

		todo, err := todoModel.Create(req)
		assert.NoError(t, err)
		assert.NotNil(t, todo)
		assert.Equal(t, req.Title, todo.Title)
		assert.Equal(t, req.Description, todo.Description)
		assert.False(t, todo.Completed)
		assert.NotZero(t, todo.ID)
		assert.NotZero(t, todo.CreatedAt)
		assert.NotZero(t, todo.UpdatedAt)
	})

	t.Run("Get Todo By ID", func(t *testing.T) {
		// Create a todo first
		req := models.CreateTodoRequest{
			Title:       "Get Test Todo",
			Description: "Get Test Description",
		}
		createdTodo, err := todoModel.Create(req)
		assert.NoError(t, err)

		// Get the todo by ID
		todo, err := todoModel.GetByID(createdTodo.ID)
		assert.NoError(t, err)
		assert.NotNil(t, todo)
		assert.Equal(t, createdTodo.ID, todo.ID)
		assert.Equal(t, createdTodo.Title, todo.Title)
	})

	t.Run("Get Non-existent Todo", func(t *testing.T) {
		todo, err := todoModel.GetByID(999)
		assert.Error(t, err)
		assert.Nil(t, todo)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Get All Todos", func(t *testing.T) {
		// Create multiple todos
		req1 := models.CreateTodoRequest{Title: "Todo 1", Description: "Desc 1"}
		req2 := models.CreateTodoRequest{Title: "Todo 2", Description: "Desc 2"}

		_, err := todoModel.Create(req1)
		assert.NoError(t, err)
		_, err = todoModel.Create(req2)
		assert.NoError(t, err)

		todos, err := todoModel.GetAll()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(todos), 2)
	})

	t.Run("Update Todo", func(t *testing.T) {
		// Create a todo first
		req := models.CreateTodoRequest{
			Title:       "Original Title",
			Description: "Original Description",
		}
		createdTodo, err := todoModel.Create(req)
		assert.NoError(t, err)

		// Update the todo
		updateReq := models.UpdateTodoRequest{
			Title:       "Updated Title",
			Description: "Updated Description",
		}

		updatedTodo, err := todoModel.Update(createdTodo.ID, updateReq)
		assert.NoError(t, err)
		assert.NotNil(t, updatedTodo)
		assert.Equal(t, updateReq.Title, updatedTodo.Title)
		assert.Equal(t, updateReq.Description, updatedTodo.Description)
		assert.True(t, updatedTodo.UpdatedAt.After(createdTodo.UpdatedAt))
	})

	t.Run("Delete Todo", func(t *testing.T) {
		// Create a todo first
		req := models.CreateTodoRequest{
			Title:       "Delete Test Todo",
			Description: "Delete Test Description",
		}
		createdTodo, err := todoModel.Create(req)
		assert.NoError(t, err)

		// Delete the todo
		err = todoModel.Delete(createdTodo.ID)
		assert.NoError(t, err)

		// Verify it's deleted
		todo, err := todoModel.GetByID(createdTodo.ID)
		assert.Error(t, err)
		assert.Nil(t, todo)
	})

	t.Run("Toggle Complete", func(t *testing.T) {
		// Create a todo first
		req := models.CreateTodoRequest{
			Title:       "Complete Test Todo",
			Description: "Complete Test Description",
		}
		createdTodo, err := todoModel.Create(req)
		assert.NoError(t, err)
		assert.False(t, createdTodo.Completed)

		// Mark as complete
		completedTodo, err := todoModel.ToggleComplete(createdTodo.ID, true)
		assert.NoError(t, err)
		assert.True(t, completedTodo.Completed)

		// Mark as incomplete
		incompletedTodo, err := todoModel.ToggleComplete(createdTodo.ID, false)
		assert.NoError(t, err)
		assert.False(t, incompletedTodo.Completed)
	})
}

// TestTodoHandlers tests the HTTP handlers
func TestTodoHandlers(t *testing.T) {
	// Use a test database
	dbPath := "test_handlers.db"
	defer os.Remove(dbPath)

	db, err := database.InitDB(dbPath)
	assert.NoError(t, err)
	defer database.CloseDB(db)

	todoModel := models.NewTodoModel(db)
	todoHandler := handlers.NewTodoHandler(todoModel)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("Create Todo Handler", func(t *testing.T) {
		router := gin.New()
		router.POST("/todos", todoHandler.CreateTodo)

		reqBody := models.CreateTodoRequest{
			Title:       "Handler Test Todo",
			Description: "Handler Test Description",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Todo
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, reqBody.Title, response.Title)
		assert.Equal(t, reqBody.Description, response.Description)
		assert.False(t, response.Completed)
	})

	t.Run("Get Todos Handler", func(t *testing.T) {
		router := gin.New()
		router.GET("/todos", todoHandler.GetTodos)

		// Create a todo first
		reqBody := models.CreateTodoRequest{
			Title:       "Get Handler Test",
			Description: "Get Handler Test Desc",
		}
		_, err := todoModel.Create(reqBody)
		assert.NoError(t, err)

		req, _ := http.NewRequest("GET", "/todos", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.Todo
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(response), 1)
	})

	t.Run("Get Todo By ID Handler", func(t *testing.T) {
		router := gin.New()
		router.GET("/todos/:id", todoHandler.GetTodo)

		// Create a todo first
		reqBody := models.CreateTodoRequest{
			Title:       "Get By ID Test",
			Description: "Get By ID Test Desc",
		}
		createdTodo, err := todoModel.Create(reqBody)
		assert.NoError(t, err)

		// Use strconv.Itoa to convert ID to string for URL
		idStr := strconv.Itoa(createdTodo.ID)
		req, _ := http.NewRequest("GET", "/todos/"+idStr, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Todo
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, createdTodo.ID, response.ID)
	})

	t.Run("Update Todo Handler", func(t *testing.T) {
		router := gin.New()
		router.PUT("/todos/:id", todoHandler.UpdateTodo)

		// Create a todo first
		reqBody := models.CreateTodoRequest{
			Title:       "Update Handler Test",
			Description: "Update Handler Test Desc",
		}
		createdTodo, err := todoModel.Create(reqBody)
		assert.NoError(t, err)

		// Update the todo
		updateReq := models.UpdateTodoRequest{
			Title:       "Updated Handler Title",
			Description: "Updated Handler Description",
		}
		jsonBody, _ := json.Marshal(updateReq)

		// Use strconv.Itoa to convert ID to string for URL
		idStr := strconv.Itoa(createdTodo.ID)
		req, _ := http.NewRequest("PUT", "/todos/"+idStr, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Todo
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, updateReq.Title, response.Title)
		assert.Equal(t, updateReq.Description, response.Description)
	})

	t.Run("Delete Todo Handler", func(t *testing.T) {
		router := gin.New()
		router.DELETE("/todos/:id", todoHandler.DeleteTodo)

		// Create a todo first
		reqBody := models.CreateTodoRequest{
			Title:       "Delete Handler Test",
			Description: "Delete Handler Test Desc",
		}
		createdTodo, err := todoModel.Create(reqBody)
		assert.NoError(t, err)

		// Use strconv.Itoa to convert ID to string for URL
		idStr := strconv.Itoa(createdTodo.ID)
		req, _ := http.NewRequest("DELETE", "/todos/"+idStr, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify it's deleted
		_, err = todoModel.GetByID(createdTodo.ID)
		assert.Error(t, err)
	})

	t.Run("Complete Todo Handler", func(t *testing.T) {
		router := gin.New()
		router.PATCH("/todos/:id/complete", todoHandler.CompleteTodo)

		// Create a todo first
		reqBody := models.CreateTodoRequest{
			Title:       "Complete Handler Test",
			Description: "Complete Handler Test Desc",
		}
		createdTodo, err := todoModel.Create(reqBody)
		assert.NoError(t, err)
		assert.False(t, createdTodo.Completed)

		// Use strconv.Itoa to convert ID to string for URL
		idStr := strconv.Itoa(createdTodo.ID)
		req, _ := http.NewRequest("PATCH", "/todos/"+idStr+"/complete", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Todo
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Completed)
	})

	t.Run("Uncomplete Todo Handler", func(t *testing.T) {
		router := gin.New()
		router.PATCH("/todos/:id/uncomplete", todoHandler.UncompleteTodo)

		// Create a completed todo first
		reqBody := models.CreateTodoRequest{
			Title:       "Uncomplete Handler Test",
			Description: "Uncomplete Handler Test Desc",
		}
		createdTodo, err := todoModel.Create(reqBody)
		assert.NoError(t, err)

		// Mark as complete first
		_, err = todoModel.ToggleComplete(createdTodo.ID, true)
		assert.NoError(t, err)

		// Use strconv.Itoa to convert ID to string for URL
		idStr := strconv.Itoa(createdTodo.ID)
		req, _ := http.NewRequest("PATCH", "/todos/"+idStr+"/uncomplete", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Todo
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Completed)
	})
}

// TestValidation tests input validation
func TestValidation(t *testing.T) {
	// Use a test database
	dbPath := "test_validation.db"
	defer os.Remove(dbPath)

	db, err := database.InitDB(dbPath)
	assert.NoError(t, err)
	defer database.CloseDB(db)

	todoModel := models.NewTodoModel(db)
	todoHandler := handlers.NewTodoHandler(todoModel)

	gin.SetMode(gin.TestMode)

	t.Run("Create Todo with Empty Title", func(t *testing.T) {
		router := gin.New()
		router.POST("/todos", todoHandler.CreateTodo)

		reqBody := models.CreateTodoRequest{
			Title:       "", // Empty title should fail
			Description: "Test Description",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Todo ID", func(t *testing.T) {
		router := gin.New()
		router.GET("/todos/:id", todoHandler.GetTodo)

		req, _ := http.NewRequest("GET", "/todos/invalid", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// Benchmark tests for performance
func BenchmarkCreateTodo(b *testing.B) {
	dbPath := "benchmark.db"
	defer os.Remove(dbPath)

	db, err := database.InitDB(dbPath)
	if err != nil {
		b.Fatal(err)
	}
	defer database.CloseDB(db)

	todoModel := models.NewTodoModel(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := models.CreateTodoRequest{
			Title:       "Benchmark Todo",
			Description: "Benchmark Description",
		}
		_, err := todoModel.Create(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetAllTodos(b *testing.B) {
	dbPath := "benchmark_get.db"
	defer os.Remove(dbPath)

	db, err := database.InitDB(dbPath)
	if err != nil {
		b.Fatal(err)
	}
	defer database.CloseDB(db)

	todoModel := models.NewTodoModel(db)

	// Create some todos first
	for i := 0; i < 100; i++ {
		req := models.CreateTodoRequest{
			Title:       "Benchmark Todo",
			Description: "Benchmark Description",
		}
		_, err := todoModel.Create(req)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := todoModel.GetAll()
		if err != nil {
			b.Fatal(err)
		}
	}
} 