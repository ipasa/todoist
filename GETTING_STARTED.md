# Getting Started with Todoist

This guide will help you set up and run the Todoist application on your local machine.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.21+** - [Download](https://golang.org/dl/)
- **Node.js 18+** - [Download](https://nodejs.org/)
- **Docker & Docker Compose** - [Download](https://docs.docker.com/get-docker/)
- **Git** - [Download](https://git-scm.com/downloads)

## Project Structure

```
todoist/
├── backend/
│   ├── auth-service/        # Authentication service (Port 8001)
│   ├── api-gateway/         # API Gateway (Port 8000)
│   ├── pkg/                 # Shared packages
│   └── [other services]     # To be implemented in Phase 2
├── frontend/                # React application (Port 3000)
└── infrastructure/
    └── docker/              # Docker Compose configuration
```

## Quick Start

### Step 1: Clone the Repository

```bash
git clone <your-repo-url>
cd todoist
```

### Step 2: Start Infrastructure Services

Start PostgreSQL, RabbitMQ, and Redis using Docker Compose:

```bash
cd infrastructure/docker
docker-compose up -d postgres rabbitmq redis mailhog
```

Wait for the services to be healthy (about 30 seconds):

```bash
docker-compose ps
```

### Step 3: Initialize Databases

The database will be automatically initialized when PostgreSQL starts. You can verify with:

```bash
docker exec -it todoist-postgres psql -U todoist -c "\l"
```

You should see: `auth_db`, `task_db`, `project_db`, `notification_db`

### Step 4: Run Database Migrations

The Auth Service will automatically run migrations on startup. Alternatively, you can run them manually:

```bash
cd ../../backend/auth-service
# Install golang-migrate tool first if needed
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path infrastructure/persistence/postgres/migrations -database "postgres://todoist:todoist_dev@localhost:5432/auth_db?sslmode=disable" up
```

### Step 5: Start Backend Services

#### Option A: Using Docker Compose (Recommended)

```bash
cd infrastructure/docker
docker-compose up auth-service api-gateway
```

#### Option B: Running Locally (for Development)

**Terminal 1 - Auth Service:**
```bash
cd backend/auth-service
go mod download
go run cmd/main.go
```

**Terminal 2 - API Gateway:**
```bash
cd backend/api-gateway
go mod download
go run cmd/main.go
```

### Step 6: Start Frontend

**Terminal 3 - Frontend:**
```bash
cd frontend
npm install
npm run dev
```

## Access the Application

Once all services are running:

- **Frontend:** http://localhost:3000
- **API Gateway:** http://localhost:8000
- **Auth Service:** http://localhost:8001
- **RabbitMQ Management:** http://localhost:15672 (guest/guest)
- **MailHog (Email Testing):** http://localhost:8025

## Testing the Application

### 1. Register a New User

1. Open http://localhost:3000
2. Click "create a new account"
3. Fill in:
   - Full Name: "John Doe"
   - Email: "john@example.com"
   - Password: "password123"
4. Click "Create account"

You should be automatically logged in and redirected to the dashboard.

### 2. Login

1. Logout from the dashboard
2. Click "Sign in"
3. Enter your credentials
4. Click "Sign in"

### 3. Test API Directly

You can also test the API directly using curl:

**Register:**
```bash
curl -X POST http://localhost:8000/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8000/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

## Troubleshooting

### Port Already in Use

If you get port conflict errors:

```bash
# Find process using the port
lsof -i :8000  # or :8001, :3000, etc.

# Kill the process
kill -9 <PID>
```

### Database Connection Error

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check database logs
docker logs todoist-postgres

# Restart PostgreSQL
docker-compose restart postgres
```

### RabbitMQ Connection Error

```bash
# Check if RabbitMQ is running
docker ps | grep rabbitmq

# Check RabbitMQ logs
docker logs todoist-rabbitmq

# Restart RabbitMQ
docker-compose restart rabbitmq
```

### Frontend Not Loading

```bash
# Clear node_modules and reinstall
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run dev
```

## Development Workflow

### Backend Development

1. Make changes to the code
2. The service will automatically restart (if using `go run`)
3. Or rebuild the Docker image: `docker-compose build auth-service`
4. Test your changes

### Frontend Development

Vite provides hot module replacement (HMR), so changes are reflected immediately without page reload.

### Database Migrations

To create a new migration:

```bash
cd backend/auth-service
migrate create -ext sql -dir infrastructure/persistence/postgres/migrations -seq add_new_field
```

This creates two files:
- `XXX_add_new_field.up.sql` - Apply migration
- `XXX_add_new_field.down.sql` - Rollback migration

## Environment Variables

### Backend Services

All backend services use environment variables for configuration. See `.env.example` in `infrastructure/docker/`

Key variables:
- `PORT` - Service port
- `DATABASE_URL` - PostgreSQL connection string
- `RABBITMQ_URL` - RabbitMQ connection string
- `JWT_SECRET` - Secret for JWT tokens
- `JWT_EXPIRY` - Access token expiry (default: 15m)
- `REFRESH_TOKEN_EXPIRY` - Refresh token expiry (default: 168h)

### Frontend

Create a `.env` file in the `frontend` directory:

```bash
cp frontend/.env.example frontend/.env
```

Edit if needed:
```
VITE_API_URL=http://localhost:8000/v1
VITE_WS_URL=ws://localhost:8000/ws
```

## Stopping the Application

### Stop All Services

```bash
cd infrastructure/docker
docker-compose down
```

### Stop and Remove Volumes (Database Data)

```bash
docker-compose down -v
```

### Stop Individual Services

```bash
# Stop backend services
Ctrl+C in each terminal

# Or using Docker
docker-compose stop auth-service api-gateway
```

## Next Steps

Now that you have the authentication system running, you can:

1. **Test the authentication flow** - Register users, login, logout
2. **Explore the code** - Review the DDD architecture in the auth service
3. **Check the database** - View user data in PostgreSQL
4. **Monitor events** - Check RabbitMQ management UI for published events
5. **Review logs** - Check application logs for debugging

## Phase 2: Task Management

The next phase will implement:
- Task Service (CRUD operations)
- Project Service (project management)
- Real-time updates with WebSocket Gateway
- Full task management UI

Stay tuned!

## Need Help?

- Check the main [README.md](./README.md) for architecture details
- Review the code in `backend/auth-service` for DDD patterns
- Check Docker logs: `docker-compose logs -f <service-name>`
- Open an issue on GitHub

Happy coding!
