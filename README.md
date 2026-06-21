# Student Management System API

A production-ready Student Management System built with Go, Gin, PostgreSQL, Redis, Docker, and JWT Authentication.

## Features

### Authentication & Authorization

* User Registration
* User Login
* JWT Access Tokens
* JWT Refresh Tokens
* Logout Support
* Protected Routes
* Authentication Middleware

### Student Management

* Create Student
* Get Student By ID
* List Students
* Update Student
* Partial Update (PATCH)
* Delete Student

### Course Management

* Create Course
* Get Course By ID
* List Courses
* Update Course
* Delete Course

### Enrollment System

* Student Course Enrollment
* Transaction Support
* Seat Availability Validation
* Duplicate Enrollment Prevention

### Performance & Reliability

* Redis Caching
* Redis Rate Limiting
* Graceful Shutdown
* Database Connection Pooling
* Dockerized Deployment

### Developer Experience

* Swagger Documentation
* Integration Testing
* Docker Compose
* Clean Architecture
* Repository Pattern

## Tech Stack

### Backend

* Go 1.26
* Gin

### Database

* PostgreSQL

### Cache

* Redis

### Authentication

* JWT

### Documentation

* Swagger / OpenAPI

### Testing

* Go Testing Package

### Containerization

* Docker
* Docker Compose

### Deployment

* Render

## ER Diagram
<img width="1002" height="675" alt="er-diagram-student-management" src="https://github.com/user-attachments/assets/26d9b8d0-b681-4906-b5c2-a271c210692b" />

## Architecture Diagram
<img width="4707" height="8192" alt="image" src="https://github.com/user-attachments/assets/151324c8-06e3-4cd6-afc4-67a1de6334c8" />
<img width="4707" height="8192" alt="image" src="https://github.com/user-attachments/assets/28e8f93d-8937-4254-a94e-b4b925a46274" />

## Project Structure

```text
cmd/
└── server/

internal/
├── apperrors/
├── auth/
├── cache/
├── config/
├── database/
├── dto/
├── handlers/
├── middleware/
├── models/
├── repositories/
├── routes/
├── services/
└── tests/

migrations/
docs/
```

## API Documentation

Swagger Documentation:

```text[
https://YOUR_RENDER_URL/swagger/index.html](https://golang-student-management.onrender.com/swagger/index.html
```

## Environment Variables

```env
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=

REDIS_ADDR=

JWT_ACCESS_SECRET=
JWT_REFRESH_SECRET=
```

## Running Locally

### Clone Repository

```bash
git clone https://github.com/Abhi78k/golang-student-management.git
cd project-1
```

### Start Services

```bash
docker compose up --build
```

Application:

```text
http://localhost:8080
```

Swagger:

```text
http://localhost:8080/swagger/index.html
```

## Running Tests

```bash
go test ./internal/tests -v
```

## Architecture Highlights

* Clean Architecture
* Repository Pattern
* Dependency Injection
* Redis Caching Layer
* JWT Authentication
* PostgreSQL Transactions
* Graceful Shutdown
* Integration Testing
