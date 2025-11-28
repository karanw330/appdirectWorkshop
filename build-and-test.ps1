# PowerShell script to build and test Docker image
# Make sure Docker Desktop is running before executing this script

Write-Host "=== AppDirect Workshop - Docker Build & Test ===" -ForegroundColor Cyan
Write-Host ""

# Check if Docker is running
Write-Host "Checking Docker..." -ForegroundColor Yellow
try {
    $dockerVersion = docker --version
    Write-Host "✓ Docker found: $dockerVersion" -ForegroundColor Green
} catch {
    Write-Host "✗ Docker is not installed or not in PATH" -ForegroundColor Red
    Write-Host "Please install Docker Desktop and ensure it's running" -ForegroundColor Yellow
    exit 1
}

# Check if Docker daemon is running
try {
    docker ps | Out-Null
    Write-Host "✓ Docker daemon is running" -ForegroundColor Green
} catch {
    Write-Host "✗ Docker daemon is not running" -ForegroundColor Red
    Write-Host "Please start Docker Desktop" -ForegroundColor Yellow
    exit 1
}

Write-Host ""
Write-Host "=== Step 1: Building Docker Image ===" -ForegroundColor Cyan
docker build -t appdirect-workshop:latest .

if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ Docker build failed" -ForegroundColor Red
    exit 1
}

Write-Host "✓ Docker image built successfully" -ForegroundColor Green
Write-Host ""

# Check if .env file exists to get environment variables
$envFile = ".env"
$firebaseProjectId = "india-tech-meetup-2025"
$adminPassword = "testpass"
$serviceAccountPath = ""

if (Test-Path $envFile) {
    Write-Host "Reading environment variables from .env file..." -ForegroundColor Yellow
    $envContent = Get-Content $envFile
    foreach ($line in $envContent) {
        if ($line -match "^FIREBASE_PROJECT_ID=(.+)") {
            $firebaseProjectId = $matches[1]
        }
        if ($line -match "^ADMIN_PASSWORD=(.+)") {
            $adminPassword = $matches[1]
        }
        if ($line -match "^FIREBASE_SERVICE_ACCOUNT_PATH=(.+)") {
            $serviceAccountPath = $matches[1] -replace '"', ''
        }
    }
}

Write-Host ""
Write-Host "=== Step 2: Running Docker Container ===" -ForegroundColor Cyan
Write-Host "Using environment variables:" -ForegroundColor Yellow
Write-Host "  FIREBASE_PROJECT_ID: $firebaseProjectId" -ForegroundColor Gray
Write-Host "  ADMIN_PASSWORD: [HIDDEN]" -ForegroundColor Gray
Write-Host ""

# Stop and remove existing container if it exists
docker stop appdirect-workshop-test 2>$null
docker rm appdirect-workshop-test 2>$null

# Build docker run command
$dockerRunCmd = "docker run -d --name appdirect-workshop-test -p 8080:8080"
$dockerRunCmd += " -e FIREBASE_PROJECT_ID=$firebaseProjectId"
$dockerRunCmd += " -e ADMIN_PASSWORD=$adminPassword"
$dockerRunCmd += " -e PORT=8080"
$dockerRunCmd += " -e FIRESTORE_SUBCOLLECTION_ID=workshop_attendees"

# If service account path is provided and file exists, mount it
if ($serviceAccountPath -and (Test-Path $serviceAccountPath)) {
    Write-Host "Mounting service account file: $serviceAccountPath" -ForegroundColor Yellow
    $dockerRunCmd += " -e FIREBASE_SERVICE_ACCOUNT_PATH=/root/serviceAccountKey.json"
    $dockerRunCmd += " -v ${serviceAccountPath}:/root/serviceAccountKey.json:ro"
} else {
    Write-Host "No service account file found. Using ADC (will fail without proper setup)" -ForegroundColor Yellow
    $dockerRunCmd += " -e FIREBASE_SERVICE_ACCOUNT_PATH=ADC"
}

$dockerRunCmd += " appdirect-workshop:latest"

Write-Host "Starting container..." -ForegroundColor Yellow
Invoke-Expression $dockerRunCmd

if ($LASTEXITCODE -ne 0) {
    Write-Host "✗ Failed to start container" -ForegroundColor Red
    exit 1
}

Write-Host "✓ Container started" -ForegroundColor Green
Write-Host ""

# Wait for container to be ready
Write-Host "=== Step 3: Waiting for service to be ready ===" -ForegroundColor Cyan
Start-Sleep -Seconds 5

$maxRetries = 12
$retryCount = 0
$isReady = $false

while ($retryCount -lt $maxRetries -and -not $isReady) {
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -UseBasicParsing -TimeoutSec 2 -ErrorAction Stop
        if ($response.StatusCode -eq 200) {
            $isReady = $true
            Write-Host "✓ Service is ready!" -ForegroundColor Green
        }
    } catch {
        $retryCount++
        Write-Host "Waiting for service... ($retryCount/$maxRetries)" -ForegroundColor Yellow
        Start-Sleep -Seconds 2
    }
}

if (-not $isReady) {
    Write-Host "✗ Service did not become ready" -ForegroundColor Red
    Write-Host "Container logs:" -ForegroundColor Yellow
    docker logs appdirect-workshop-test
    exit 1
}

Write-Host ""
Write-Host "=== Step 4: Testing Endpoints ===" -ForegroundColor Cyan

# Test health endpoint
Write-Host "Testing /health endpoint..." -ForegroundColor Yellow
try {
    $healthResponse = Invoke-WebRequest -Uri "http://localhost:8080/health" -UseBasicParsing
    Write-Host "✓ Health check: $($healthResponse.StatusCode)" -ForegroundColor Green
    Write-Host "  Response: $($healthResponse.Content)" -ForegroundColor Gray
} catch {
    Write-Host "✗ Health check failed: $_" -ForegroundColor Red
}

# Test API info endpoint
Write-Host "Testing /api endpoint..." -ForegroundColor Yellow
try {
    $apiResponse = Invoke-WebRequest -Uri "http://localhost:8080/api" -UseBasicParsing
    Write-Host "✓ API info: $($apiResponse.StatusCode)" -ForegroundColor Green
    $apiContent = $apiResponse.Content | ConvertFrom-Json
    Write-Host "  API Version: $($apiContent.version)" -ForegroundColor Gray
} catch {
    Write-Host "✗ API endpoint failed: $_" -ForegroundColor Red
}

# Test attendees endpoint
Write-Host "Testing /api/attendees endpoint..." -ForegroundColor Yellow
try {
    $attendeesResponse = Invoke-WebRequest -Uri "http://localhost:8080/api/attendees" -UseBasicParsing
    Write-Host "✓ Attendees endpoint: $($attendeesResponse.StatusCode)" -ForegroundColor Green
} catch {
    Write-Host "✗ Attendees endpoint failed: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "=== Summary ===" -ForegroundColor Cyan
Write-Host "Container Name: appdirect-workshop-test" -ForegroundColor Gray
Write-Host "Container URL: http://localhost:8080" -ForegroundColor Gray
Write-Host ""
Write-Host "To view logs: docker logs appdirect-workshop-test" -ForegroundColor Yellow
Write-Host "To stop container: docker stop appdirect-workshop-test" -ForegroundColor Yellow
Write-Host "To remove container: docker rm appdirect-workshop-test" -ForegroundColor Yellow
Write-Host ""
Write-Host "✓ Docker build and test completed!" -ForegroundColor Green

