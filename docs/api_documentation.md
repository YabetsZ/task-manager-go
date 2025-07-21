# Task Management API Documentation

This document provides details for the Task Management REST API.

**Base URL:** `http://localhost:5000`

---

## Authentication Endpoints

### 1. Register a New User

-   **Endpoint:** `POST /register`
-   **Description:** Creates a new user account.
-   **Request Body (JSON):**

    ```json
    {
        "username": "string (required)",
        "password": "string (required)"
    }
    ```

-   **Success Response:**
    -   **Code:** `201 Created`
-   **Error Responses:**
    -   **Code:** `400 Bad Request` if the payload is invalid or username already exists
    -   **Code:** `500 Internal Server Error` for unexpected errors.

### 2. Login

-   **Endpoint:** `POST /login`
-   **Description:** Authenticates a user and returns a JWT token.
-   **Request Body (JSON):**

    ```json
    {
        "username": "string (required)",
        "password": "string (required)"
    }
    ```

-   **Success Response:**
    -   **Code:** `200 OK`
    -   **Content:**

        ```json
        {
            "token": "string (JWT token)"
        }
        ```

-   **Error Responses:**
    -   **Code:** `400 Bad Request` if the payload is invalid.
    -   **Code:** `401 Unauthorized` if the credentials are invalid.
    -   **Code:** `500 Internal Server Error` for unexpected errors.

### 3. Promote a User

-   **Endpoint:** `POST /api/promote/:id`
-   **Description:** Promotes a user to an admin role. This endpoint requires admin privileges.
-   **URL Parameters:**
    -   `id` (string, required): The unique identifier of the user to promote.
-   **Success Response:**
    -   **Code:** `200 OK`
-   **Error Responses:**
    -   **Code:** `400 Bad Request` if the user ID format is invalid.
    -   **Code:** `401 Unauthorized` if the user is not an admin.
    -   **Code:** `404 Not Found` if a user with the specified ID does not exist.
    -   **Code:** `500 Internal Server Error` for unexpected errors.

## Task Management Endpoints

### 1. Create a New Task

-   **Endpoint:** `POST /api/tasks`
-   **Description:** Adds a new task to the system. This endpoint requires admin privileges.
-   **Request Body (JSON):**

    ```json
    {
        "title": "string (required)",
        "description": "string",
        "due_date": "datetime (RFC3339 format, e.g., 2025-12-31T15:00:00Z)",
        "status": "string (e.g., 'Pending', 'In Progress', 'Completed')"
    }
    ```

-   **Success Response:**
    -   **Code:** `201 Created`
    -   **Content:** The newly created task object, including its unique ID.
-   **Error Responses:**
    -   **Code:** `400 Bad Request` if the payload is invalid or the title is missing.
    -   **Code:** `401 Unauthorized` if the user is not an admin.

### 2. Get All Tasks

-   **Endpoint:** `GET /api/tasks`
-   **Description:** Retrieves a list of all tasks in the system. This endpoint is accessible to all authenticated users.
-   **Success Response:**
    -   **Code:** `200 OK`
    -   **Content:** An array of task objects.

-   **Error Responses:**
    -   **Code:** `400 Bad Request` for a bad request.
    -   **Code:** `401 Unauthorized` if the user is not authenticated.

### 3. Get a Specific Task

-   **Endpoint:** `GET /api/tasks/:id`
-   **Description:** Retrieves the details of a single task by its ID. This endpoint is accessible to all authenticated users.
-   **URL Parameters:**
    -   `id` (string, required): The unique identifier of the task.
-   **Success Response:**
    -   **Code:** `200 OK`
    -   **Content:** A single task object.
-   **Error Responses:**
    -   **Code:** `401 Unauthorized` if the user is not authenticated.
    -   **Code:** `404 Not Found` if a task with the specified ID does not exist.

### 4. Update a Task

-   **Endpoint:** `PUT /api/tasks/:id`
-   **Description:** Updates the details of an existing task. This endpoint requires admin privileges.
-   **URL Parameters:**
    -   `id` (string, required): The unique identifier of the task to update.
-   **Request Body (JSON):**

    ```json
    {
        "title": "string",
        "description": "string",
        "due_date": "datetime",
        "status": "string"
    }
    ```

-   **Success Response:**
    -   **Code:** `200 OK`
    -   **Content:** The fully updated task object.
-   **Error Responses:**
    -   **Code:** `400 Bad Request` if the payload is invalid.
    -   **Code:** `401 Unauthorized` if the user is not an admin.
    -   **Code:** `404 Not Found` if a task with the specified ID does not exist.

### 5. Delete a Task

-   **Endpoint:** `DELETE /api/tasks/:id`
-   **Description:** Deletes a task from the system. This endpoint requires admin privileges.
-   **URL Parameters:**
    -   `id` (string, required): The unique identifier of the task to delete.
-   **Success Response:**
    -   **Code:** `204 No Content`
-   **Error Responses:**
    -   **Code:** `401 Unauthorized` if the user is not an admin.
    -   **Code:** `404 Not Found` if a task with the specified ID does not exist.

