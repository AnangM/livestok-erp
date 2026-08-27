# Livestock ERP (Backend API)

Livestock ERP is a Go-based backend service designed for livestock farm management. It provides livestock pedigree and lineage tracking (sire/dam relationships), health/status tracking, harvest and meat yield records, and multi-tenant authentication powered by Supabase.

---

## 🏗️ Architecture & Core Folder Structure

The project follows standard Go Clean / Layered Architecture principles to maintain separation of concerns and testability:

```text
livestok-erp/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point & dependency wiring
├── internal/
│   ├── domain/                  # Core domain models and entities (e.g., Animal)
│   ├── repository/              # Data persistence layer (PostgreSQL implementations)
│   ├── service/                 # Business logic layer and interface definitions
│   └── transport/
│       └── http/
│           ├── handler/         # HTTP request handlers (Auth, Animals)
│           └── middleware/      # HTTP middlewares (e.g., Supabase JWT Auth)
├── migrations/                  # Database migration files (up/down SQL scripts)
├── supabase/                    # Supabase local configuration & CLI settings
├── go.mod                       # Go module dependencies
└── README.md                    # Project documentation
```

### Layer Breakdown

- **`internal/domain`**: Defines business domain entities (`Animal`, etc.) free from external dependencies.
- **`internal/repository`**: Handles database interactions and queries (`PostgresAnimalRepository`).
- **`internal/service`**: Defines business contracts via interfaces (`AuthServiceInterface`, `AnimalServiceInterface`) and implements business rules.
- **`internal/transport/http`**: Exposes RESTful HTTP endpoints using `chi` router, JSON serializations, and middleware for Supabase authentication.
- **`cmd/api`**: Wires all dependencies (database connection, repositories, services, handlers) and starts the HTTP server.

---

## ⚙️ Prerequisites

- **Go**: Version `1.22+` (or latest)
- **PostgreSQL**: PostgreSQL database instance (or Supabase local/hosted instance)
- **Supabase CLI** *(optional, for local Supabase emulation)*

---

## 🔐 Environment Variables

The application requires the following environment variables to run:

| Variable | Description | Example |
| :--- | :--- | :--- |
| `DB_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:54322/postgres?sslmode=disable` |
| `SUPABASE_URL` | Supabase Project URL | `https://your-project-id.supabase.co` (or `http://localhost:54321`) |
| `SUPABASE_ANON_KEY` | Supabase public/anonymous API key | `eyJhbGciOi...` |
| `SUPABASE_JWT_SECRET` | Supabase JWT Secret used to verify auth tokens | `your-supabase-jwt-secret` |

---

## 🚀 Getting Started & How to Run

### 1. Clone the repository

```bash
git clone https://github.com/AnangM/livestok-erp.git
cd livestok-erp
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Run Database Migrations

Apply the initial schema from the `migrations/` directory to your PostgreSQL database:

```bash
# Using migrate CLI or psql
psql "$DB_URL" -f migrations/000001_init_schema.up.sql
```

### 4. Set Environment Variables

You can export them in your current terminal session:

```bash
export DB_URL="postgres://postgres:postgres@localhost:54322/postgres?sslmode=disable"
export SUPABASE_URL="http://127.0.0.1:54321"
export SUPABASE_ANON_KEY="your-anon-key"
export SUPABASE_JWT_SECRET="your-jwt-secret"
```

### 5. Run the API Server

```bash
go run cmd/api/main.go
```

The server will start listening on port `:8080`.

---

## 📡 API Endpoints

### Public Routes

- **`GET /health`**: Health check probe.
  - **Response**: `200 OK` (`"OK"`)

- **`POST /api/v1/login`**: Authenticate user via Supabase Auth.
  - **Request Body**:
    ```json
    {
      "email": "user@example.com",
      "password": "password123"
    }
    ```
  - **Response**:
    ```json
    {
      "access_token": "jwt_token_here"
    }
    ```

### Protected Routes (Requires `Authorization: Bearer <token>`)

- **`POST /api/v1/animals`**: Register a new animal under the authenticated farm.
  - **Headers**: `Authorization: Bearer <token>`
  - **Request Body**:
    ```json
    {
      "tag_number": "RAB-001",
      "breed": "New Zealand White",
      "gender": "Buck",
      "sire_id": null,
      "dam_id": null,
      "birth_date": "2026-01-15T00:00:00Z",
      "status": "Active"
    }
    ```
  - **Response**: `201 Created`

---

## 🧪 Build & Test

```bash
# Compile and build the binary
go build -v ./...

# Run test suite
go test -v ./...
```
