# pecatu-be

Backend API for the pecatu project, built using Golang and following Clean Architecture principles.

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- PostgreSQL database

## Getting Started

1. **Clone the repository** (if you haven't already).
2. **Environment Variables**:
   Copy `.env.example` to `.env` and fill in your actual database credentials.
   ```bash
   cp .env.example .env
   ```

## Running the Application Locally

You can run the Go application directly on your host machine:

```bash
# Download dependencies
go mod tidy

# Run the API
go run ./cmd/api/main.go
```
The server will start on port `6201` (or whatever `PORT` you defined in your `.env`).

## Running with Docker Compose

To build and run the application using Docker:

```bash
# Start the container in detached mode
docker-compose up --build -d
```
The API will be available at `http://localhost:6201`.

## Available Endpoints

- **Health Check**
  ```bash
  curl -i http://localhost:6201/health
  ```
  Returns `{"status": "ok"}`.

- **Get Donate** (Dummy Data)
  ```bash
  curl -i http://localhost:6201/donate/123
  ```
  Returns the mocked entity `{"id": "123", "amount": 10000, "status": "success"}`.
