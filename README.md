# go-eat

`go-eat` is a robust and scalable REST API for restaurant management, built with Go 1.24 and MongoDB. The project follows a Clean Architecture (Hexagonal Architecture) to ensure code maintainability, testability, and separation of concerns.

## 🚀 Features

- **User Management**: Create, list, retrieve, update, and delete users.
- **Ingredient Management**: Manage restaurant ingredients.
- **Product Management**: Create and manage products (dishes) with ingredients.
- **Order Management**: Process and track customer orders.
- **Database Migrations**: Automated MongoDB schema migrations using `golang-migrate`.
- **API Documentation**: Interactive documentation using Scalar and Swagger.
- **Health Check**: Simple endpoint to verify server status.

## 🛠️ Tech Stack

- **Go 1.24**: Core programming language.
- **Fiber v2**: High-performance HTTP web framework.
- **MongoDB**: NoSQL database for flexible data storage.
- **Golang Migrate**: Tool for managing database migrations.
- **Swagger/Scalar**: For interactive API documentation.

## 🏗️ Architecture

The project follows the **Hexagonal Architecture (Ports and Adapters)** pattern:

- **`cmd/`**: Entry point of the application.
- **`internal/core/domain/`**: Domain models (entities).
- **`internal/core/ports/`**: Interface definitions for services and repositories.
- **`internal/core/services/`**: Business logic implementation.
- **`internal/adapters/handlers/http/`**: HTTP controllers (Fiber handlers).
- **`internal/adapters/repositories/mongodb/`**: MongoDB implementation for repositories.
- **`migrations/`**: MongoDB migration files.

## ⚙️ Getting Started

### Prerequisites

- Go 1.24 or higher
- MongoDB instance (Local or Atlas)
- Git

### Environment Variables

The application requires the following environment variables to be set in a `.env` file or in your environment:

```env
DB_MONGO_URL=mongodb://your-username:your-password@your-host:your-port/restaurant?authSource=admin
DB_MONGO_NAME=restaurant
```

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/dunas/go-eat.git
   cd go-eat
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

### Running the Application

To start the API server:

```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:7000`.

## 📖 API Documentation

The project includes built-in API documentation using Scalar. Once the server is running, you can access it at:

- **Interactive API Docs (Scalar)**: [http://localhost:7000/docs/](http://localhost:7000/docs/)
- **Health Check**: [http://localhost:7000/health](http://localhost:7000/health)

## 📂 Project Structure

```text
.
├── cmd/
│   └── api/                # Application entry point
├── docs/                   # Swagger/OpenAPI documentation
├── internal/
│   ├── adapters/           # Implementation details (DB, HTTP)
│   │   ├── handlers/       # HTTP Request handlers
│   │   └── repositories/   # Database access implementation
│   └── core/               # Business Logic
│       ├── domain/         # Entities and models
│       ├── ports/          # Interfaces (Repository/Service)
│       └── services/       # Service implementation
├── migrations/             # MongoDB migration files
└── go.mod                  # Go module definition
```
