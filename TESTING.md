# Testing Guide

## Running Tests

### Backend Tests

**Unit Tests:**
```bash
make test-backend
# or
go test -v ./internal/handlers/... -coverprofile=coverage.out
```

**Integration Tests (requires Firestore):**
```bash
make test-integration
# or
go test -v -tags=integration ./internal/handlers/...
```

**With Coverage:**
```bash
make test-coverage
# Opens coverage.html in browser
```

### Frontend Tests

```bash
make test-frontend
# or
npm run test
```

**With UI:**
```bash
npm run test:ui
```

**With Coverage:**
```bash
npm run test:coverage
```

### All Tests

```bash
make test
```

## Test Structure

### Backend Tests

- **Unit Tests** (`handlers_test.go`): Test handlers with mocked dependencies
- **Integration Tests** (`handlers_integration_test.go`): Test with real Firestore (requires `-tags=integration`)

### Frontend Tests

- **Component Tests**: Located in `src/components/__tests__/`
- Uses Vitest and React Testing Library
- Mocks API calls using `vi.mock()`

## Writing New Tests

### Backend Test Example

```go
func TestMyHandler(t *testing.T) {
    handler := &Handlers{adminPassword: "test"}
    
    req := httptest.NewRequest("GET", "/api/endpoint", nil)
    w := httptest.NewRecorder()
    
    handler.MyHandler(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

### Frontend Test Example

```jsx
import { render, screen } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import MyComponent from '../MyComponent'

describe('MyComponent', () => {
  it('renders correctly', () => {
    render(<MyComponent />)
    expect(screen.getByText('Hello')).toBeInTheDocument()
  })
})
```

## Test Coverage Goals

- Backend: >80% coverage
- Frontend: >70% coverage for critical components

