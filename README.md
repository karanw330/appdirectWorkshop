# AppDirect India AI Workshop - Event Registration System

A production-ready React SPA with Golang backend for event registration and management.

## Features

- **Hero Section** with call-to-action buttons
- **Sessions & Speakers** grid display
- **Registration Form** with live attendee count
- **Location** with embedded Google Maps
- **Admin Dashboard** with password protection
  - Attendee management
  - Speaker management
  - Session management
  - Analytics (Pie chart by designation)

## Tech Stack

- **Frontend**: React 18, Vite, Tailwind CSS, Framer Motion, Recharts
- **Backend**: Golang, Gorilla Mux, Firebase Admin SDK
- **Database**: Google Firestore
- **Security**: Environment variables, password-protected admin

## Setup Instructions

### Prerequisites

- Node.js 18+
- Go 1.21+
- Firebase project with Firestore enabled
- Firebase service account key (JSON file)

### Backend Setup

1. Copy `.env.example` to `.env` and fill in your values:
   ```bash
   cp .env.example .env
   ```

2. Place your Firebase service account key as `serviceAccountKey.json` in the root directory

3. Install Go dependencies:
   ```bash
   go mod download
   ```

4. Run the backend:
   ```bash
   go run cmd/server/main.go
   ```

### Frontend Setup

1. Install dependencies:
   ```bash
   npm install
   ```

2. Run the development server:
   ```bash
   npm run dev
   ```

3. Build for production:
   ```bash
   npm run build
   ```

## Environment Variables

- `FIREBASE_PROJECT_ID`: Your Firebase project ID
- `FIREBASE_SERVICE_ACCOUNT_PATH`: Path to service account JSON
- `FIRESTORE_SUBCOLLECTION_ID`: Firestore subcollection identifier
- `PORT`: Backend server port (default: 8080)
- `ADMIN_PASSWORD`: Password for admin login
- `VITE_API_URL`: Backend API URL

## Docker Setup

1. Build and run with Docker Compose:
   ```bash
   docker-compose up --build
   ```

2. Or build the Docker image manually:
   ```bash
   docker build -t appdirect-workshop .
   docker run -p 8080:8080 --env-file .env appdirect-workshop
   ```

## Project Structure

```
appdirectWorkshop/
├── cmd/server/          # Golang backend server
├── internal/
│   ├── handlers/        # HTTP handlers
│   └── firestore/       # Firestore client
├── src/                 # React frontend
│   ├── components/      # React components
│   ├── pages/           # Page components
│   ├── services/        # API services
│   └── context/         # React context
├── dist/                # Frontend build output
└── static/              # Static files for production
```

## API Endpoints

### Attendees
- `GET /api/attendees` - Get all attendees
- `POST /api/attendees` - Register new attendee
- `GET /api/attendees/count` - Get attendee count

### Speakers
- `GET /api/speakers` - Get all speakers
- `POST /api/speakers` - Create speaker
- `PUT /api/speakers/{id}` - Update speaker
- `DELETE /api/speakers/{id}` - Delete speaker

### Sessions
- `GET /api/sessions` - Get all sessions
- `POST /api/sessions` - Create session
- `PUT /api/sessions/{id}` - Update session
- `DELETE /api/sessions/{id}` - Delete session

### Admin
- `POST /api/admin/login` - Admin login

## Firestore Collections

The application uses the following Firestore collections:
- `attendees` - Registered attendees
- `speakers` - Speaker profiles
- `sessions` - Workshop sessions

## Development

### Using Makefile

```bash
# Install dependencies
make install-frontend
make install-backend

# Run development servers
make dev-frontend    # Terminal 1
make dev-backend     # Terminal 2

# Build for production
make build-frontend
make build-backend
```

### Manual Development

1. Start backend (Terminal 1):
   ```bash
   go run cmd/server/main.go
   ```

2. Start frontend (Terminal 2):
   ```bash
   npm run dev
   ```

3. Access the application:
   - Frontend: http://localhost:3000 (or port shown by Vite)
   - Backend API: http://localhost:8080/api
   - Admin: http://localhost:3000/admin/login

## Security Notes

- Never commit `.env` or `serviceAccountKey.json` to version control
- Use strong passwords for admin access
- Keep Firebase credentials secure
- The `.gitignore` file is configured to exclude sensitive files

## Production Deployment

1. Set all environment variables
2. Build frontend: `npm run build`
3. Build backend: `go build -o server ./cmd/server`
4. Ensure `serviceAccountKey.json` is in the same directory
5. Run: `./server`

The server will serve the frontend from the `./static` directory if it exists.

