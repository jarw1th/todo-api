# ToDo API

A RESTful API for managing personal todo tasks with JWT-based authentication and Swagger documentation.

---

## Features

- User registration and login with JWT access & refresh tokens.
- CRUD operations for todo items (Create, Read, Update, Delete).
- Each user can access only their own todos.
- Filtering, sorting, pagination support.
- Automatic todo history tracking.
- Swagger UI for API documentation and testing.
- Centralized error handling with JSON responses.

---

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **Router:** Gorilla Mux
- **Authentication:** JWT (access + refresh tokens)
- **Swagger UI:** Swaggo
- **Hosting:** Render.com or any Go-compatible server

---

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- `swag` CLI for generating Swagger docs

### Installation

1. Clone the repository:

```bash
git clone https://github.com/your-username/todo-api.git
cd todo-api
```

2.	Set environment variables:

```bash
export JWT_SECRET_KEY="your_secret_key"
export DB_HOST="127.0.0.1"
export DB_PORT="5432"
export DB_USER="postgres"
export DB_PASSWORD="your_db_password"
export DB_NAME="todo"
```

3.	Install dependencies:

```bash
go mod tidy
```

4.	Generate Swagger documentation:

```bash
swag init -g main.go
```

5.	Run the server:

```bash
go run main.go
```

- Server will start on http://localhost:8080.

Swagger UI will be available at:

```bash
http://localhost:8080/swagger/index.html
```

---

## API Endpoints

### Authentication

| Method | Endpoint    | Description           |
|--------|------------|---------------------|
| POST   | `/register` | Register a new user  |
| POST   | `/login`    | Login and get access & refresh tokens |

### Todos (Requires Authorization)

| Method | Endpoint       | Description             |
|--------|---------------|------------------------|
| GET    | `/todos`       | List all todos for the authenticated user |
| POST   | `/todos`       | Create a new todo       |
| PUT    | `/todos/{id}`  | Replace a todo completely |
| PATCH  | `/todos/{id}`  | Update a todo partially |
| DELETE | `/todos/{id}`  | Delete a todo           |

**All requests must include an `Authorization: Bearer <access_token>` header.**

---

## Database Schema

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    title TEXT NOT NULL,
    description TEXT,
    done BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE todo_history (
    id SERIAL PRIMARY KEY,
    todo_id INT REFERENCES todos(id),
    user_id INT REFERENCES users(id),
    old_value TEXT,
    new_value TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

---

## Usage Examples

### Register a User

```bash
curl -X POST http://localhost:8080/register \
-H "Content-Type: application/json" \
-d '{"username": "Alice", "password": "securepassword"}'
```

### Login

```bash
curl -X POST http://localhost:8080/login \
-H "Content-Type: application/json" \
-d '{"username": "Alice", "password": "securepassword"}'
```

Response:

```json
{
  "user_id": 1,
  "access_token": "<access_token>",
  "refresh_token": "<refresh_token>"
}
```

### Create a Todo

```bash
curl -X POST http://localhost:8080/todos \
-H "Authorization: Bearer <access_token>" \
-H "Content-Type: application/json" \
-d '{"title": "Buy milk", "description": "Get milk from the store"}'
```

### List Todos

```bash
curl -X GET http://localhost:8080/todos \
-H "Authorization: Bearer <access_token>"
```

### Update a Todo

- **PUT** (replace completely):

```bash
curl -X PUT http://localhost:8080/todos/1 \
-H "Authorization: Bearer <access_token>" \
-H "Content-Type: application/json" \
-d '{"title": "Buy eggs", "description": "From supermarket", "done": false}'
```

- **PATCH** (update partially):

```bash
curl -X PATCH http://localhost:8080/todos/1 \
-H "Authorization: Bearer <access_token>" \
-H "Content-Type: application/json" \
-d '{"done": true}'
```

### Delete a Todo

```bash
curl -X DELETE http://localhost:8080/todos/1 \
-H "Authorization: Bearer <access_token>"
```

---

## Notes

- Access tokens expire after 24 hours. Use refresh tokens to generate new access tokens with the `/refresh` endpoint.
- Only the owner of a todo can modify or delete it.
- API responses are always in JSON format.
- Swagger UI provides interactive documentation at `/swagger/index.html`.
- Centralized error handling ensures consistent JSON error responses.
- Todo history is automatically tracked in the `todo_history` table.

---

## Contributing

1. Fork the repository.  
2. Create a feature branch (`git checkout -b feature-name`).  
3. Commit your changes (`git commit -m 'Add feature'`).  
4. Push to the branch (`git push origin feature-name`).  
5. Open a Pull Request.

---

## License

Apache License © Ruslan Parastaev