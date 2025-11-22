.PHONY: help install-frontend install-backend dev-frontend dev-backend build-frontend build-backend run-backend test test-backend test-frontend test-integration clean docker-build docker-up docker-down

help:
	@echo "Available commands:"
	@echo "  make install-frontend  - Install frontend dependencies"
	@echo "  make install-backend   - Install backend dependencies"
	@echo "  make dev-frontend      - Run frontend development server"
	@echo "  make dev-backend       - Run backend development server"
	@echo "  make build-frontend   - Build frontend for production"
	@echo "  make build-backend    - Build backend for production"
	@echo "  make run-backend      - Run backend server"
	@echo "  make test             - Run all tests"
	@echo "  make test-backend     - Run backend unit tests"
	@echo "  make test-frontend    - Run frontend tests"
	@echo "  make test-integration - Run integration tests (requires Firestore)"
	@echo "  make docker-build    - Build Docker image"
	@echo "  make docker-up        - Start services with docker-compose"
	@echo "  make docker-down      - Stop services"
	@echo "  make clean            - Clean build artifacts"

install-frontend:
	cd . && npm install

install-backend:
	go mod download

dev-frontend:
	npm run dev

dev-backend:
	go run cmd/server/main.go

build-frontend:
	npm run build

build-backend:
	CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

run-backend:
	./server

test: test-backend test-frontend

test-backend:
	go test -v ./internal/handlers/... -coverprofile=coverage-backend.out

test-frontend:
	npm run test

test-integration:
	go test -v -tags=integration ./internal/handlers/... -coverprofile=coverage-integration.out

test-coverage:
	go test -v ./internal/handlers/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

docker-build:
	docker build -t appdirect-workshop:latest .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

clean:
	rm -rf dist node_modules server coverage.out coverage.html coverage-backend.out coverage-integration.out
