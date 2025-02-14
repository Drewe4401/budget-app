# Budget API

A simple budgeting API built with Go, PostgreSQL, and Docker. This project provides endpoints for managing users, budgets, charges, and shares with JWT-based authentication. It also includes a static file server for a frontend (if available).

---

## Run on Docker

This project is containerized using Docker. Follow the steps below to build and run the application inside a Docker container.

# Frontend

# Backend

The backend is implemented in Go and offers a RESTful API with the following features:

## Key Features

### Authentication & Authorization
- **JWT-based authentication.**
- **Password hashing using bcrypt.**
- **Admin-only endpoints for managing users.**

### Data Models
- **User:** Contains `id`, `username`, `password` (bcrypt-hashed), and `permissions`.
- **Budget:** Represents a budget with details like `name`, `amount`, `description`, `period`, and `user_id`.
- **Charge:** Represents a charge with details including `name`, `amount`, `charge_type`, `periodical`, `user_id`, and `created_at`.
- **Share:** Handles sharing between users with `user_id`, `user_share_id`, and `access` level.

### Database Initialization
- Creates tables (`users`, `budgets`, `charges`, `shares`) if they don't exist.
- Creates a default admin user (username: `admin`, password: `admin`) if no admin is found.

## API Endpoints

### User Endpoints (Admin Only)
- **POST** `/api/users`  
  Create a new user.
- **PUT** `/api/users/{id}`  
  Update an existing user.
- **DELETE** `/api/users/{id}`  
  Delete a user.

### Authentication
- **POST** `/api/login`  
  Log in a user and return a JWT token.

### Budget Endpoints
- **GET** `/api/budgets`  
  Retrieve budgets belonging to the authenticated user.
- **POST** `/api/budgets`  
  Create a new budget for the authenticated user.
- **PUT** `/api/budgets/{id}`  
  Update an existing budget (only if it belongs to the authenticated user).
- **DELETE** `/api/budgets/{id}`  
  Delete a budget (only if it belongs to the authenticated user).

### Charge Endpoints
- **GET** `/api/charges`  
  Retrieve charges for the authenticated user.
- **POST** `/api/charges`  
  Create a new charge for the authenticated user.
- **PUT** `/api/charges/{id}`  
  Update an existing charge (only if it belongs to the authenticated user).
- **DELETE** `/api/charges/{id}`  
  Delete a charge (only if it belongs to the authenticated user).

### Share Endpoints
- **GET** `/api/shares`  
  Retrieve shares where the authenticated user is either the owner or recipient.
- **POST** `/api/shares`  
  Create a share with another user by providing the username and access level.
- **DELETE** `/api/shares/{id}`  
  Delete a share if the authenticated user is permitted to do so.

## How It Works

### Initialization
The application initializes by connecting to a PostgreSQL database and setting up the necessary tables. If no admin user exists, a default admin user is created.

### Routing
Routes are managed using the Gorilla Mux router. Each route is associated with its corresponding handler function.

### Security
- Passwords are hashed using bcrypt.
- JWT tokens are generated for authenticated sessions.
- Middleware checks ensure that sensitive operations (like user management) are restricted to admins.

### Static File Serving
The application serves static files from the `./public` directory, which allows integration with a frontend.

For further details, please refer to the inline comments within the source code.
