# Go Todo API

A simple and efficient Todo API built with Go, Gin framework, and SQLite database.

## Features

- ✅ Create, Read, Update, Delete (CRUD) operations for todos
- ✅ SQLite database for data persistence
- ✅ RESTful API design
- ✅ Comprehensive test coverage
- ✅ JSON request/response format
- ✅ Input validation
- ✅ Error handling

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/todos` | Get all todos |
| GET | `/todos/:id` | Get a specific todo by ID |
| POST | `/todos` | Create a new todo |
| PUT | `/todos/:id` | Update an existing todo |
| DELETE | `/todos/:id` | Delete a todo |
| PATCH | `/todos/:id/complete` | Mark a todo as completed |
| PATCH | `/todos/:id/uncomplete` | Mark a todo as incomplete |

## Todo Model

```json
{
  "id": 1,
  "title": "Buy groceries",
  "description": "Get milk, bread, and eggs",
  "completed": false,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

## Prerequisites

- Go 1.21 or higher
- SQLite3 (usually comes with Go)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd go-todo-api
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Database

The application uses SQLite as the database. The database file (`todo.db`) will be created automatically when you first run the application.

### Project Structure

```
go-todo-api/
├── main.go              # Application entry point
├── models/
│   └── todo.go          # Todo model and database operations
├── handlers/
│   └── todo.go          # HTTP request handlers
├── database/
│   └── sqlite.go        # Database connection and initialization
├── tests/
│   └── todo_test.go     # Test files
├── go.mod               # Go module file
├── go.sum               # Go module checksums
├── README.md            # This file
└── .gitignore           # Git ignore file
```

## API Examples

### Create a Todo

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Go",
    "description": "Study Go programming language"
  }'
```

### Get All Todos

```bash
curl http://localhost:8080/api/v1/todos
```

### Get a Specific Todo

```bash
curl http://localhost:8080/api/v1/todos/1
```

### Update a Todo

```bash
curl -X PUT http://localhost:8080/api/v1/todos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Go Programming",
    "description": "Study Go programming language and build projects"
  }'
```

### Mark Todo as Complete

```bash
curl -X PATCH http://localhost:8080/api/v1/todos/1/complete
```

### Delete a Todo

```bash
curl -X DELETE http://localhost:8080/api/v1/todos/1
```

## Error Responses

The API returns appropriate HTTP status codes and error messages:

- `400 Bad Request` - Invalid input data
- `404 Not Found` - Todo not found
- `500 Internal Server Error` - Server error

Example error response:
```json
{
  "error": "Todo not found"
}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE). 