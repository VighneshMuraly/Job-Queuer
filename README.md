# Job Queuer

A Go-based job queue and scheduler with REST API, cron job processing, and database-backed persistence. Supports job prioritization, retries, and up/down migrations.

## Features
- REST API for scheduling jobs and checking status
- Job prioritization (high, medium, low) and quota-based processing
- Automatic retries and job termination after 3 failures
- Cron job processor (runs every 30 seconds)
- PostgreSQL database with migration support (up/down)
- Dockerized for easy deployment
- Cobra CLI with multiple commands

## Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.21+

### Setup
1. **Clone the repository**
   ```sh
   git clone <repo-url>
   cd Job-Queuer
   ```

2. **Configure environment**
   - Edit `app.env` with your database credentials (used by Docker Compose and Go app):
     ```env
     DB_HOST=db
     DB_PORT=5432
     DB_USER=jobuser
     DB_PASSWORD=jobpass
     DB_NAME=jobdb
     ```

3. **Run with Docker Compose**
   ```sh
   docker-compose up
   ```
   - This starts both the Go app and a Postgres database.
   - The app listens on port `8080` by default.

4. **Database Migrations**
   - Migrations are in `app/migrations/` as `.up.sql` and `.down.sql` files.
   - Migrations run automatically on startup. If a migration fails, rollback is attempted.

## API Endpoints

### Health Check
- `GET /health`
  - Returns `{ "status": "ok" }` if the service is running.

### Schedule a Job
- `POST /schedule`
  - Request Body (JSON):
    ```json
    {
      "type": "email_notification",
      "payload": { "to": "user@example.com", "subject": "Hello" },
      "priority": "high" //"medium" / "low"
    }
    ```
  - Response: Job details with ID and status.

### Check Job Status
- `GET /status?id=<job_id>`
  - Returns job status and details for the given job ID.

## Job Processing Logic
- Jobs are processed by a cron every 30 seconds.
- Quotas: high (3), medium (2), low (1) jobs per cycle.
- Jobs are prioritized and sorted by last update time.
- Failed jobs are retried (priority increases for low/medium), terminated after 3 failures.

## CLI Usage
- The app uses Cobra for CLI commands.
- Example commands:
  - `./jobqueuer` (default server)
  - `./jobqueuer worker` (run worker logic)
  - Add more commands in `cmd/` as needed.

## Development
- Source code is organized in `app/` and `cmd/` folders.
- Handlers, services, repositories, and router are modular.
- Migrations are managed with [golang-migrate](https://github.com/golang-migrate/migrate).

## Testing
- You can use Postman or curl to interact with the API endpoints.
- Example:
  ```sh
  curl -X POST http://localhost:8080/schedule \
    -H "Content-Type: application/json" \
    -d '{ "type": "email", "payload": { "to": "user@example.com" }, "priority": "high" }'
  ```


