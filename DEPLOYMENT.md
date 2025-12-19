# Deployment Guide

## Google Cloud Run Deployment

This application is configured to deploy to Google Cloud Run using Application Default Credentials (ADC), which means no service account file is needed in the container.

### Prerequisites

1. Google Cloud Project with billing enabled
2. Firebase project with Firestore enabled
3. Cloud Run API enabled
4. Cloud Build API enabled (for automated builds)
5. gcloud CLI installed and authenticated

### Setup Steps

#### 1. Enable Required APIs

```bash
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable firestore.googleapis.com
```

#### 2. Set Up Service Account for Cloud Run

The Cloud Run service needs permissions to access Firestore. The default Compute Engine service account will be used with ADC.

```bash
# Get the service account email
PROJECT_ID=$(gcloud config get-value project)
SERVICE_ACCOUNT="${PROJECT_NUMBER}-compute@developer.gserviceaccount.com"

# Grant Firestore permissions
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:${SERVICE_ACCOUNT}" \
  --role="roles/datastore.user"
```

#### 3. Build and Deploy

**Option A: Using Cloud Build (Recommended)**

```bash
# Set substitution variables
gcloud builds submit --config=cloudbuild.yaml \
  --substitutions=_FIREBASE_PROJECT_ID=your-project-id,_ADMIN_PASSWORD=your-password
```

**Option B: Manual Build and Deploy**

```bash
# Build the image
gcloud builds submit --tag gcr.io/$PROJECT_ID/appdirect-workshop

# Deploy to Cloud Run
gcloud run deploy appdirect-workshop \
  --image gcr.io/$PROJECT_ID/appdirect-workshop \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars FIREBASE_PROJECT_ID=your-project-id,PORT=8080,FIRESTORE_SUBCOLLECTION_ID=workshop_attendees,ADMIN_PASSWORD=your-password,FIREBASE_SERVICE_ACCOUNT_PATH=ADC,FIRESTORE_DATABASE_ID=(default)
```

#### 4. Set Environment Variables

After deployment, you can update environment variables:

```bash
gcloud run services update appdirect-workshop \
  --update-env-vars FIREBASE_PROJECT_ID=your-project-id,ADMIN_PASSWORD=your-password
```

### Environment Variables

- `FIREBASE_PROJECT_ID`: Your Firebase project ID (required)
- `FIREBASE_SERVICE_ACCOUNT_PATH`: Set to `ADC` for Cloud Run (defaults to ADC if not set)
- `PORT`: Server port (default: 8080)
- `ADMIN_PASSWORD`: Password for admin login (required)
- `FIRESTORE_SUBCOLLECTION_ID`: Firestore collection ID (default: workshop_attendees)
- `FIRESTORE_DATABASE_ID`: Database ID if using named database (default: (default))

### Local Testing with Docker

```bash
# Build the image
docker build -t appdirect-workshop .

# Run locally (still needs service account file for local)
docker run -p 8080:8080 \
  -e FIREBASE_PROJECT_ID=your-project-id \
  -e FIREBASE_SERVICE_ACCOUNT_PATH=./serviceAccountKey.json \
  -e ADMIN_PASSWORD=your-password \
  -v $(pwd)/serviceAccountKey.json:/root/serviceAccountKey.json:ro \
  appdirect-workshop
```

### Health Check

The application includes a health check endpoint:

```bash
curl https://your-service-url.run.app/health
```

### Troubleshooting

1. **Permission Denied Errors**: Ensure the service account has `roles/datastore.user` role
2. **Connection Errors**: Verify `FIREBASE_PROJECT_ID` is correct
3. **Build Failures**: Check that all dependencies are in `go.mod` and `package.json`

### CI/CD Integration

The `cloudbuild.yaml` file can be used with Cloud Build triggers for automatic deployments on git push.

