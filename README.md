# 8-level-desent

REST API built for the [Desent Solutions API Quest](https://www.desent.io/coding-test-backend) backend challenge. Implements all 8 levels: ping, echo, CRUD, auth, pagination, search, error handling, and the speed-run boss fight.

## Tech Stack

- **Go 1.24** with **Gin** HTTP framework
- **JWT** authentication (`golang-jwt/jwt/v5`)
- **In-memory** storage (no database required)
- Repository pattern for clean separation of concerns

## Project Structure

```
├── api/
│   └── index.go                     # Vercel serverless handler
├── app/
│   ├── main.go                      # entry point
│   ├── config/                      # env-based configuration
│   ├── errors/                      # custom AppError type
│   ├── models/                      # request/response structs
│   ├── repositories/
│   │   ├── book_repository.go       # interface
│   │   └── inmemory/                # in-memory implementation
│   ├── services/                    # business logic
│   ├── controllers/
│   │   ├── http/
│   │   │   ├── response.go          # response helpers
│   │   │   └── v1/                  # versioned HTTP handlers + router
│   │   └── middlewares/             # JWT auth middleware
│   └── tests/                       # integration tests
├── deployment/
│   ├── local/
│   │   └── Dockerfile               # multi-stage Docker build
│   └── dev/
│       └── vercel.json              # Vercel route config (symlink)
└── vercel.json                      # Vercel config (source of truth)
```

## Getting Started

### Prerequisites

- Go 1.24+
- (Optional) Docker

### Run Locally

```bash
make run
```

The server starts on port `8080` by default. Override with the `PORT` env var.

### Build Binary

```bash
make build
./bin/server
```

### Run Tests

```bash
make test
```

### Docker

```bash
make docker-run
```

## Environment Variables

| Variable     | Default          | Description                |
|--------------|------------------|----------------------------|
| `PORT`       | `8080`           | Server listen port         |
| `JWT_SECRET` | random on start  | Secret for signing JWTs    |
| `GIN_MODE`   | `release`        | Gin mode (debug/release)   |

## API Endpoints

### Public

| Method | Path          | Description          |
|--------|---------------|----------------------|
| GET    | `/ping`       | Health check         |
| POST   | `/echo`       | Echo JSON body back  |
| POST   | `/auth/token` | Get a JWT token      |

### Protected (requires `Authorization: Bearer <token>`)

| Method | Path           | Description                     |
|--------|----------------|---------------------------------|
| POST   | `/books`       | Create a book                   |
| GET    | `/books`       | List books (search + paginate)  |
| GET    | `/books/:id`   | Get a book by ID                |
| PUT    | `/books/:id`   | Update a book                   |
| DELETE | `/books/:id`   | Delete a book                   |

### Query Parameters for `GET /books`

| Param    | Example               | Description              |
|----------|-----------------------|--------------------------|
| `author` | `?author=Alan`        | Filter by author name    |
| `page`   | `?page=1&limit=2`     | Page number (1-indexed)  |
| `limit`  | `?page=1&limit=2`     | Results per page         |

## Deployment

### Docker

```bash
make docker-run
```

Or manually:

```bash
docker build -f deployment/local/Dockerfile -t my-api .
docker run -p 8080:8080 -e JWT_SECRET=my-secret my-api
```

### Vercel

The `api/index.go` handler wraps the entire Gin app as a single serverless function. All routes are rewritten to it via `vercel.json`.

```bash
# Set JWT_SECRET so it persists across cold starts
vercel env add JWT_SECRET

# Deploy
make deploy
```

The app also works on Render, Railway, Fly.io, or any platform that supports Docker or Go binaries -- it reads `PORT` from the environment.
