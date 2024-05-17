# Learnify - A Platform for Educational Courses and Lessons

## Introduction

Learnify is a platform designed to facilitate the creation, management, and delivery of educational courses and lessons. It enables teachers to create courses and lessons, and students to enroll in and access these lessons.

## Team

- Arai Bayaly

## API Endpoints

### Health Check

- **GET /health** - Checks the health of the API.

### Authentication

- **POST /register** - Registers a new user.
- **POST /login** - Authenticates a user and returns a JWT token.

### Courses

- **POST /api/v1/courses** - Creates a new course.
- **GET /api/v1/courses/{id}** - Retrieves a course by ID.
- **PUT /api/v1/courses/{id}** - Updates a course by ID.
- **DELETE /api/v1/courses/{id}** - Deletes a course by ID.
- **GET /api/v1/courses** - Retrieves all courses with pagination, filtering, and sorting.

### Lessons

- **POST /api/v1/lessons** - Creates a new lesson.
- **GET /api/v1/lessons/{id}** - Retrieves a lesson by ID.
- **PUT /api/v1/lessons/{id}** - Updates a lesson by ID.
- **DELETE /api/v1/lessons/{id}** - Deletes a lesson by ID.
- **GET /api/v1/lessons** - Retrieves all lessons with pagination, filtering, and sorting.

## Database Structure

### Table: users

| Column        | Type         | Constraints                           |
| ------------- | ------------ | ------------------------------------- |
| id            | SERIAL       | PRIMARY KEY                           |
| first_name    | VARCHAR(255) | NOT NULL                              |
| last_name     | VARCHAR(255) | NOT NULL                              |
| email         | VARCHAR(255) | UNIQUE, NOT NULL                      |
| password_hash | TEXT         | NOT NULL                              |
| role          | VARCHAR(50)  | NOT NULL                              |
| created_at    | TIMESTAMP    | DEFAULT CURRENT_TIMESTAMP             |
| updated_at    | TIMESTAMP    | DEFAULT CURRENT_TIMESTAMP             |

### Table: courses

| Column      | Type         | Constraints                                        |
| ----------- | ------------ | -------------------------------------------------- |
| id          | SERIAL       | PRIMARY KEY                                        |
| title       | VARCHAR(255) | NOT NULL                                           |
| description | TEXT         |                                                    |
| teacher_id  | INT          | NOT NULL, FOREIGN KEY REFERENCES users(id)         |
| created_at  | TIMESTAMP    | DEFAULT CURRENT_TIMESTAMP                          |
| updated_at  | TIMESTAMP    | DEFAULT CURRENT_TIMESTAMP                          |

### Table: lessons

| Column      | Type         | Constraints                                                |
| ----------- | ------------ | ---------------------------------------------------------- |
| id          | SERIAL       | PRIMARY KEY                                                |
| title       | VARCHAR(255) | NOT NULL                                                   |
| content     | TEXT         |                                                            |
| course_id   | INT          | NOT NULL, FOREIGN KEY REFERENCES courses(id) ON DELETE CASCADE |
| created_at  | TIMESTAMP    | DEFAULT CURRENT_TIMESTAMP                                  |
| updated_at  | TIMESTAMP    | DEFAULT CURRENT_TIMESTAMP                                  |

## Relations Between Entities

### Course (n) ⇄ Lesson (n)

- **Course (1) ⇄ (n) Lessons**
  - Each course can have multiple lessons.
  - Each lesson belongs to a single course.
  - Represented by `course_id` in the `lessons` table.

### learnify
