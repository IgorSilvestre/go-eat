# AI Agent Guide - go-eat

This guide provides context, architecture, and development rules for AI agents working on the `go-eat` project.

## 🚀 Project Overview
`go-eat` is a REST API for restaurant management, focusing on robustness and scalability.
- **Language**: Go 1.24+
- **Framework**: Fiber v2
- **Database**: MongoDB (with automated migrations)
- **Architecture**: Hexagonal (Ports and Adapters)
- **Documentation**: Swagger/Scalar

## 🛠️ Technical Stack
- **HTTP Framework**: `github.com/gofiber/fiber/v2`
- **Database**: `go.mongodb.org/mongo-driver` (MongoDB)
- **Migrations**: `github.com/golang-migrate/migrate/v4`
- **Storage**: S3 (via Minio-go client)
- **Testing**: `github.com/stretchr/testify`
- **API Docs**: `github.com/swaggo/swag`, `github.com/yokeTH/gofiber-scalar`

## 🏗️ Architecture: Hexagonal (Clean Architecture)
The project strictly follows the Hexagonal Architecture pattern. Logic is separated into **Core** (Business Logic) and **Adapters** (External Infrastructure).

### 1. Internal Core (`internal/core`)
This is the heart of the application, independent of any external framework or database.
- **`domain/`**: Entities and models (e.g., `User`, `Product`, `Order`).
- **`ports/`**: Interfaces defining how the core interacts with the world.
    - **Repository Ports**: Interfaces for database access (e.g., `UserRepository`).
    - **Service Ports**: Interfaces for business logic operations (e.g., `UserService`).
    - **Storage Ports**: Interfaces for external storage (e.g., `StorageService`).
- **`services/`**: Implementations of the Service Ports. Contains business logic.

### 2. Internal Adapters (`internal/adapters`)
Implementation details that connect the core to external systems.
- **`handlers/http/`**: Fiber handlers that process HTTP requests and call core services.
- **`repositories/mongodb/`**: MongoDB implementations of the repository ports.
- **`storage/`**: External storage implementations (e.g., S3).

### Dependency Rule
- **Inner layers (Core) must NOT depend on outer layers (Adapters).**
- Business logic in `services/` should only depend on `domain/` and `ports/`.
- Concrete implementations like `mongodb/` or `fiber/` handlers depend on `ports/` and `domain/`.

## 📌 Development Guidelines (Modern Go 1.24)
Follow modern Go idioms and the project's established patterns:

### Go Idioms
- **Error Checking**: Use `errors.Is(err, target)` instead of `err == target`.
- **Types**: Use `any` instead of `interface{}`.
- **Slices & Maps**: Use `slices` and `maps` standard library packages (e.g., `slices.Contains`, `maps.Clone`).
- **Loops**: Use `for i := range n` for simple range-based loops.
- **JSON Tags**: Use `omitzero` instead of `omitempty` for modern Go 1.24 JSON handling when appropriate.
- **Context**: Use `t.Context()` in tests when a context is needed.

### Code Style
- Keep `internal/core/domain` free of infrastructure-specific tags (like MongoDB or JSON) unless strictly necessary for the whole application.
- Use dependency injection: services should receive repositories via interfaces (ports).
- Errors should be handled or wrapped, never ignored.
- **IMPORTANT**: Always update API documentation (Swagger/Scalar) and `api.http` when an endpoint or its contract changes.

## 🛠️ Common Tasks

### Database Migrations
Migrations are stored in `migrations/` and use `golang-migrate`.
- `up.json`: Applied when migrating up.
- `down.json`: Applied when rolling back.

### Running & Testing
- **Run**: `go run cmd/api/main.go`
- **Test**: `go test ./...`
- **Manual Testing**: `api.http` contains sample requests for all endpoints.
- **API Documentation**: Accessible at `http://localhost:7000/docs/` when running.

### Configuration
The app uses environment variables for configuration.
- `DB_MONGO_URL`: MongoDB connection URI.
- `MINIO_ENDPOINT`: Endpoint for S3 storage.
- `MINIO_ROOT_USER`/`MINIO_ROOT_PASSWORD`: Credentials for S3.
- `MINIO_BUCKET_IMAGES_PUBLIC`: Bucket name for images.

## 📂 Project Structure Map
```text
.
├── cmd/api/             # Entry point (Main)
├── internal/
│   ├── core/
│   │   ├── domain/      # Entities
│   │   ├── ports/       # Interfaces
│   │   └── services/    # Business Logic
│   └── adapters/
│       ├── handlers/    # HTTP Controllers
│       ├── repositories/# DB implementations
│       └── storage/     # S3 implementations
├── migrations/          # DB Migrations
├── nixpacks.toml        # Nixpacks deployment config
└── docs/                # API Specs
```
