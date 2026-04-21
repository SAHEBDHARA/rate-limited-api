# Rate-Limited API Service

A production-considerate API service built with Go, featuring atomic rate limiting and persistent statistics using Redis.

## 🚀 Features

- **Hexagonal Architecture**: Clean separation of business logic from infrastructure.
- **Atomic Rate Limiting**: Fixed-window rate limiting (5 req/min) using Redis `INCR` to ensure accuracy under high concurrency.
- **Persistent Stats**: Per-user request counts stored in Redis.
- **Dockerized**: Ready-to-go environment with `docker-compose`.
- **Dockerized**: Ready-to-go environment with `docker-compose`.

## 🛠 Tech Stack

- **Language**: Go 1.24+
- **Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
- **Database**: [Redis](https://redis.io/)
- **Infrastructure**: Docker & Docker Compose

---

## 🏗 Architecture

The project follows a **Simplified Hexagonal Architecture** (Ports and Adapters):

- **`/internal/core/domain`**: Pure data models (Request, Stats).
- **`/internal/core/ports`**: Interfaces defining the contracts for the core.
- **`/internal/core/services`**: Pure business logic (the "Brain").
- **`/internal/adapters`**: Specific technology implementations (Redis, HTTP/Gin).

This structure ensures that the business logic is independent of the database or web framework, making it highly testable and maintainable.

---

## 🚦 Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)

### Running the Application

1. **Clone the repository**:
   ```bash
   git clone <repo-url>
   cd rate-limited-api
   ```

2. **Start the stack**:
   ```bash
   docker-compose up --build
   ```
   *This starts the API on port `8080` and a Redis instance on port `6379`.*

---

## 📖 API Endpoints

- **`POST /api/request`**: Submit a payload.
  - **Auth**: Pass `user_id` in JSON body, or `X-User-ID` header.
  - **Limit**: Max 5 requests per minute.
- **`GET /api/stats?user_id=...`**: Retrieve total successful requests for a user.

---

## ⚙️ Configuration

The application can be configured via environment variables (found in `docker-compose.yml` or `.env`):

| Variable | Description | Default |
|----------|-------------|---------|
| `API_PORT` | Port for the API server | `8080` |
| `REDIS_ADDR` | Address of the Redis server | `127.0.0.1:6379` |
| `RATE_LIMIT` | Max requests per minute | `5` |

---

## 🧪 Testing Concurrency

The rate limiter uses Redis's atomic `INCR` operation. You can test it by firing parallel requests:

```bash
# Example using curl in parallel (bash)
for i in {1..10}; do curl -X POST http://localhost:8080/api/request -d '{"user_id":"demo","payload":"test"}' & done
```
You will notice that exactly 5 requests succeed, and the rest return `429 Too Many Requests`.
