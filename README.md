# Job Manager

A background job queue service written in Go. Jobs are created via REST API and processed asynchronously by a worker pool.

Built as a pet project to practice Go concurrency, layered architecture, and REST API design.

## Features

- Create jobs and track their lifecycle: `pending -> running -> done / failed`
- Concurrent job processing via a worker pool (3 workers by default)
- Filter jobs by status
- Workers finish current jobs before stopping
- In-memory storage with `sync.RWMutex` for thread safety

## Getting Started

**Requirements:** Go 1.22 or later

```bash
git clone https://github.com/fourseasonsm/job_manager
cd job_manager
go mod tidy
go run cmd/server/main.go
```

Server starts on `http://localhost:8080`.
### Create a job

```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{}'
```

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "pending",
  "created_at": "2026-04-22T10:00:00Z",
  "started_at": "0001-01-01T00:00:00Z",
  "finished_at": "0001-01-01T00:00:00Z"
}
```

### Get all jobs

```bash
curl http://localhost:8080/jobs
```

### Filter by status

```bash
curl http://localhost:8080/jobs?status=pending
curl http://localhost:8080/jobs?status=running
curl http://localhost:8080/jobs?status=done
curl http://localhost:8080/jobs?status=failed
```

### Get job by ID

```bash
curl http://localhost:8080/jobs/{id}
```

## Job Lifecycle

```
POST /jobs
    │
    ▼
 pending    job created, waiting in queue
    │
    ▼
 running    picked up by a worker
    │
    ├──▶ done    (85% of the time)
    │
    └──▶ failed  (15% of the time, simulated)
```

Workers simulate work with a random delay of 1–4 seconds. 15% of jobs fail randomly to demonstrate error handling.

## Planned Improvements

- [ ] `DELETE /jobs/{id}` — cancel a pending job
- [ ] PostgreSQL storage as an alternative to in-memory
- [ ] Job types and payload
- [ ] Retry logic for failed jobs
- [ ] `GET /jobs/{id}/logs` — job execution history