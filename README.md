# E-cercise Backend

E-cercise is a full-stack fitness equipment platform. This backend is built with **Go (Fiber)** and powered by **PostgreSQL**, **Docker**, and a Python-based **Recommendation Service**.

## Features

- Equipment management (CRUD)
- Image upload via Cloudinary
- Equipment options, features, attributes
- User registration and goal tracking
- Muscle group targeting
- Personalized recommendation system
- JWT-based authentication
- Admin dashboard support
- RESTful APIs
- Dockerized environment

## Technologies

- **Go** (with Fiber, Gorm, UUID)
- **PostgreSQL**
- **Docker + docker-compose**
- **Cloudinary** (for image upload)
- **Python (FastAPI)** recommender
- **GitHub Actions** for CI/CD

## Getting Started

### Prerequisites

- Go 1.19+
- Docker & Docker Compose
- PostgreSQL (or use Docker)

### Environment Setup

Copy the example environment file and configure your own:

```bash
cp .env.example .env
```

### Run with Docker

```bash
docker-compose up --build
```

### Run Locally

```bash
go mod tidy
go run main.go
```

## Folder Structure

- `main.go` – Entry point
- `/config` – Config loaders
- `/controller` – Route handlers
- `/model` – GORM models
- `/repository` – Data access
- `/service` – Business logic
- `/data` – Request/response types
- `/scripts` – SQL seeders
- `/templates` – Email or HTML templates
- `.github/workflows` – GitHub CI/CD

## API Documentation

Coming soon or available via Swagger/Postman (if provided).

## Recommendation Service

Make sure your Python recommendation service is running and accessible at the configured URL (default: `http://localhost:8000`).

## Deployment

The project supports deployment to Google Cloud Run or any Docker-compatible platform.

## License

MIT – feel free to use and modify.
