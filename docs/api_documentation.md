
# Task Management API Documentation

This document provides details for the Task Management REST API.

**Base URL:** `http://localhost:5000`

---

## Endpoints

### 1. Create a New Task

- **Endpoint:** `POST /tasks`
- **Description:** Adds a new task to the system.
- **Request Body (JSON):**
  ```json
  {
      "title": "string (required)",
      "description": "string",
      "due_date": "datetime (RFC3339 format, e.g., 2025-12-31T15:00:00Z)",
      "status": "string (e.g., 'Pending', 'In Progress', 'Completed')"
  }
  ```
- **Success Response:**
  - **Code:** `201 Created`
  - **Content:** The newly created task object, including its unique ID.
- **Error Response:**
  - **Code:** `400 Bad Request` if the payload is invalid or the title is missing.

### 2. Get All Tasks

- **Endpoint:** `GET /tasks`
- **Description:** Retrieves a list of all tasks in the system.
- **Success Response:**
  - **Code:** `200 OK`
  - **Content:** An array of task objects.
  ```json
  [
      {
          "id": "...",
          "title": "...",
          ...
      }
  ]
  ```

### 3. Get a Specific Task

- **Endpoint:** `GET /tasks/:id`
- **Description:** Retrieves the details of a single task by its ID.
- **URL Parameters:**
  - `id` (string, required): The unique identifier of the task.
- **Success Response:**
  - **Code:** `200 OK`
  - **Content:** A single task object.
- **Error Response:**
  - **Code:** `404 Not Found` if a task with the specified ID does not exist.

### 4. Update a Task

- **Endpoint:** `PUT /tasks/:id`
- **Description:** Updates the details of an existing task. You can provide one or more fields to update.
- **URL Parameters:**
  - `id` (string, required): The unique identifier of the task to update.
- **Request Body (JSON):**
  ```json
  {
      "title": "string",
      "description": "string",
      "due_date": "datetime",
      "status": "string"
  }
  ```
- **Success Response:**
  - **Code:** `200 OK`
  - **Content:** The fully updated task object.
- **Error Response:**
  - **Code:** `400 Bad Request` if the payload is invalid.
  - **Code:** `404 Not Found` if a task with the specified ID does not exist.

### 5. Delete a Task

- **Endpoint:** `DELETE /tasks/:id`
- **Description:** Deletes a task from the system.
- **URL Parameters:**
  - `id` (string, required): The unique identifier of the task to delete.
- **Success Response:**
  - **Code:** `204 No Content`
- **Error Response:**
  - **Code:** `404 Not Found` if a task with the specified ID does not exist.

