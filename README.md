# Task Management REST API

![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)

This project is an implementation of a task from the A2SV backend learning with Golang (as specified in [task](./Task.md)). It is a simple and robust RESTful API for managing tasks, built with Go and the high-performance Gin web framework, demonstrating a clean, layered architecture for building maintainable and scalable web services.

The API supports full CRUD (Create, Read, Update, Delete) functionality for tasks and uses an in-memory data store, requiring no external database setup.

## Features

-   **Create, Read, Update, Delete (CRUD)** operations for tasks.
-   RESTful endpoints for easy integration with any client.
-   In-memory data storage (no external database required).
-   Structured logging and clear error responses.
-   Dependency injection for decoupled and testable components.

## Project Structure

The project follows a layered architecture to separate concerns, making the codebase clean and easy to navigate:

-   `main.go`: The application's entry point, responsible for wiring all components together.
-   `/models`: Defines the core data structures (e.g., `Task`).
-   `/data`: The service layer, containing all business logic and data manipulation.
-   `/controllers`: The presentation layer, handling HTTP requests and responses.
-   `/router`: Defines all API routes and connects them to controller handlers.
-   `/docs`: Contains project documentation.

## Getting Started

Follow these instructions to get the project running on your local machine.

### Prerequisites

-   [Go](https://go.dev/doc/install) (version 1.24 or later)

### Installation & Running

1.  **Clone the repository:**
    *(If you were using Git, you would clone it. For now, just navigate to your project folder.)*
    ```sh
    cd /path/to/task_manager
    ```

2.  **Install dependencies:**
    The Go module system will automatically handle dependencies when you build or run the project. You can also install them explicitly if needed.
    ```sh
    go mod tidy
    ```

3.  **Run the application:**
    This command will compile and run the `main.go` file, starting the web server on `http://localhost:5000`.
    ```sh
    go run main.go
    ```

The API is now running and ready to accept requests!

## API Documentation

For detailed information on all available endpoints, including request payloads, response formats, status codes, and usage examples, please refer to the official API documentation.

➡️ **[View API Documentation](./docs/api_documentation.md)**

## Testing the API

You can test the endpoints using any API client like [Postman](https://www.postman.com/) or a command-line tool like `curl`.

**Example: Get all tasks using `curl`**
```sh
curl -X GET http://localhost:8080/tasks
```

**Example: Create a new task using `curl`**
```sh
curl -X POST http://localhost:8080/tasks \
-H "Content-Type: application/json" \
-d '{
    "title": "My New Task",
    "description": "This is a test task.",
    "due_date": "2025-12-31T15:00:00Z"
}'
```

Refer to the [API Documentation](./docs/api_documentation.md) for more examples and details on how to use each endpoint.