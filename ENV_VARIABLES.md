# Environment Variables Reference

## Backend Environment Variables

### Required Variables

1. **FIREBASE_PROJECT_ID**
   - **Description**: Your Firebase project ID
   - **Example**: `india-tech-meetup-2025`
   - **Used in**: `cmd/server/main.go`
   - **Required**: Yes

2. **ADMIN_PASSWORD**
   - **Description**: Password for admin login
   - **Example**: `your-secure-password-here`
   - **Used in**: `internal/handlers/handlers.go`
   - **Required**: Yes (defaults to `admin123` if not set - development only)

### Optional Variables

3. **PORT**
   - **Description**: Server port number
   - **Default**: `8080`
   - **Example**: `8080`
   - **Used in**: `cmd/server/main.go`

4. **FIREBASE_SERVICE_ACCOUNT_PATH**
   - **Description**: Path to Firebase service account JSON file (local) or `ADC` for Cloud Run
   - **Default**: `./serviceAccountKey.json` (local) or `ADC` (Cloud Run)
   - **Example (local)**: `./serviceAccountKey.json` or `C:\Users\Karan\Downloads\service-account.json`
   - **Example (Cloud Run)**: `ADC`
   - **Used in**: `cmd/server/main.go`, `internal/firestore/client.go`
   - **Note**: Set to `ADC` when deploying to Cloud Run to use Application Default Credentials

5. **FIRESTORE_SUBCOLLECTION_ID**
   - **Description**: Firestore subcollection identifier
   - **Default**: `workshop_attendees`
   - **Example**: `workshop_attendees`
   - **Used in**: `cmd/server/main.go`

6. **K_SERVICE** (Auto-detected in Cloud Run)
   - **Description**: Automatically set by Google Cloud Run
   - **Used in**: `cmd/server/main.go` to detect Cloud Run environment

## Frontend Environment Variables

1. **VITE_API_URL**
   - **Description**: Backend API URL
   - **Default**: `/api` (relative URL, uses Vite proxy in dev)
   - **Example (local)**: `http://localhost:8080/api`
   - **Example (production)**: `/api` (relative)
   - **Used in**: `src/services/api.js`
   - **Note**: In production, frontend is served by backend, so relative URL works

## Environment Variable Summary

### Local Development (.env file)
```env
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_SERVICE_ACCOUNT_PATH=./serviceAccountKey.json
FIRESTORE_SUBCOLLECTION_ID=workshop_attendees
PORT=8080
ADMIN_PASSWORD=your-secure-password-here
VITE_API_URL=http://localhost:8080/api
```

### Docker/Container
```env
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_SERVICE_ACCOUNT_PATH=./serviceAccountKey.json  # or path in container
FIRESTORE_SUBCOLLECTION_ID=workshop_attendees
PORT=8080
ADMIN_PASSWORD=your-secure-password-here
```

### Google Cloud Run
```env
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_SERVICE_ACCOUNT_PATH=ADC  # Uses Application Default Credentials
FIRESTORE_SUBCOLLECTION_ID=workshop_attendees
PORT=8080
ADMIN_PASSWORD=your-secure-password-here
```

## Variable Usage by Component

| Variable | Backend | Frontend | Required | Default |
|----------|---------|----------|----------|---------|
| FIREBASE_PROJECT_ID | ✅ | ❌ | Yes | - |
| ADMIN_PASSWORD | ✅ | ❌ | Yes* | `admin123` |
| PORT | ✅ | ❌ | No | `8080` |
| FIREBASE_SERVICE_ACCOUNT_PATH | ✅ | ❌ | No | `./serviceAccountKey.json` or `ADC` |
| FIRESTORE_SUBCOLLECTION_ID | ✅ | ❌ | No | `workshop_attendees` |
| VITE_API_URL | ❌ | ✅ | No | `/api` |
| K_SERVICE | ✅ | ❌ | Auto | - |

*Required in production, has default for development
