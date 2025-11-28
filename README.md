# Todoist-like Todo Application

A production-ready Todo application built with microservices architecture, featuring real-time updates and modern UI.

## Tech Stack

### Backend
- **Language:** Go 1.21+
- **Architecture:** Microservices with DDD
- **Message Queue:** RabbitMQ
- **Database:** PostgreSQL
- **Cache:** Redis
- **Auth:** JWT + OAuth2 (Google, GitHub)

### Frontend
- **Framework:** React 18 + TypeScript
- **Build Tool:** Vite
- **Styling:** Tailwind CSS + Headless UI
- **State Management:** Zustand + TanStack Query
- **Real-time:** WebSocket

## Architecture

### Microservices
1. **Auth Service** (Port 8001) - Authentication, OAuth2, JWT
2. **Task Service** (Port 8002) - Task CRUD, labels, comments
3. **Project Service** (Port 8003) - Project management, sharing
4. **Notification Service** (Port 8004) - Notifications, reminders
5. **WebSocket Gateway** (Port 8005) - Real-time updates
6. **API Gateway** (Port 8000) - Request routing, middleware

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Node.js 18 or higher
- Docker & Docker Compose

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd todoist
```

2. Start infrastructure services:
```bash
cd infrastructure/docker
docker-compose up -d
```

3. Run database migrations:
```bash
make migrate
```

4. Start backend services:
```bash
make run-backend
```

5. Start frontend:
```bash
cd frontend
npm install
npm run dev
```

6. Access the application:
- Frontend: http://localhost:3000
- API Gateway: http://localhost:8000
- RabbitMQ Management: http://localhost:15672 (guest/guest)

## Project Structure

```
todoist/
├── backend/
│   ├── auth-service/       # Authentication service
│   ├── task-service/       # Task management service
│   ├── project-service/    # Project management service
│   ├── notification-service/ # Notification service
│   ├── websocket-gateway/  # WebSocket gateway
│   ├── api-gateway/        # API gateway
│   └── pkg/                # Shared packages
├── frontend/               # React frontend
├── infrastructure/         # Docker, K8s, scripts
└── docs/                   # Documentation
```

## Development

### Running Tests
```bash
make test
```

### Building Services
```bash
make build
```

### Linting
```bash
make lint
```

## API Documentation

API documentation is available at:
- Swagger UI: http://localhost:8000/swagger
- OpenAPI spec: `/docs/api/openapi.yaml`

## License

MIT
