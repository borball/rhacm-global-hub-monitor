# Development Guide

This guide helps developers set up their development environment and contribute to the RHACM Global Hub Monitor project.

## Table of Contents

1. [Development Environment Setup](#development-environment-setup)
2. [Project Structure](#project-structure)
3. [Backend Development](#backend-development)
4. [Frontend Development](#frontend-development)
5. [Testing](#testing)
6. [Building](#building)
7. [Contributing](#contributing)

## Development Environment Setup

### Prerequisites

- **Go**: Version 1.22 or later
- **Node.js**: Version 20 or later
- **npm**: Version 10 or later
- **Docker**: For building container images
- **kubectl/oc**: For Kubernetes/OpenShift access
- **Git**: For version control
- **IDE**: VS Code, GoLand, or your preferred IDE

### Clone the Repository

```bash
git clone https://github.com/your-org/rhacm-global-hub-monitor.git
cd rhacm-global-hub-monitor
```

### Access to a Cluster

You need access to an OpenShift cluster with RHACM installed for full development and testing.

For local development without a cluster, you can:
- Mock the Kubernetes client
- Disable authentication
- Use sample data

## Project Structure

```
rhacm-global-hub-monitor/
├── backend/                    # Go backend application
│   ├── cmd/
│   │   └── server/            # Main application entry point
│   ├── pkg/
│   │   ├── api/               # API route definitions
│   │   ├── auth/              # Authentication logic
│   │   ├── client/            # Kubernetes/RHACM clients
│   │   ├── handlers/          # HTTP handlers
│   │   └── models/            # Data models
│   ├── internal/
│   │   ├── config/            # Configuration
│   │   └── middleware/        # HTTP middleware
│   ├── go.mod                 # Go module definition
│   ├── go.sum                 # Go dependencies
│   └── Makefile              # Build automation
├── frontend/                   # React frontend application
│   ├── src/
│   │   ├── components/        # React components
│   │   ├── pages/             # Page components
│   │   ├── services/          # API client
│   │   ├── types/             # TypeScript types
│   │   └── utils/             # Utility functions
│   ├── package.json           # npm dependencies
│   ├── tsconfig.json          # TypeScript config
│   └── vite.config.ts         # Vite config
├── operator/                   # Kubernetes operator
│   └── config/                # Operator manifests
├── deployment/                 # Deployment configurations
│   ├── k8s/                   # Kubernetes manifests
│   └── docker/                # Dockerfiles
├── docs/                       # Documentation
└── README.md                  # Main documentation
```

## Backend Development

### Setup

1. **Install Go dependencies:**

```bash
cd backend
go mod download
```

2. **Set up environment variables:**

Create a `.env` file in the `backend` directory:

```bash
PORT=8080
KUBECONFIG=/path/to/your/kubeconfig
ENABLE_AUTH=false
LOG_LEVEL=debug
```

3. **Run the backend:**

```bash
go run cmd/server/main.go
```

The API will be available at `http://localhost:8080/api`

### Code Organization

#### Adding a New API Endpoint

1. **Define the model** in `pkg/models/types.go`
2. **Create handler** in `pkg/handlers/`
3. **Register route** in `pkg/api/router.go`

Example:

```go
// pkg/handlers/example.go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/rhacm-global-hub-monitor/backend/pkg/models"
)

type ExampleHandler struct {
    // dependencies
}

func NewExampleHandler() *ExampleHandler {
    return &ExampleHandler{}
}

func (h *ExampleHandler) GetExample(c *gin.Context) {
    c.JSON(http.StatusOK, models.APIResponse{
        Success: true,
        Data:    "example data",
    })
}

// pkg/api/router.go
func SetupRouter(...) *gin.Engine {
    // ...
    exampleHandler := handlers.NewExampleHandler()
    v1.GET("/example", exampleHandler.GetExample)
    // ...
}
```

### Development Best Practices

- **Error Handling**: Always return structured errors
- **Logging**: Use structured logging with appropriate levels
- **Testing**: Write unit tests for all handlers
- **Documentation**: Add godoc comments to exported functions

### Debugging

#### Using Delve

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug cmd/server/main.go
```

#### VS Code Debug Configuration

Create `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Backend",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/backend/cmd/server",
      "env": {
        "PORT": "8080",
        "ENABLE_AUTH": "false"
      }
    }
  ]
}
```

## Frontend Development

### Setup

1. **Install dependencies:**

```bash
cd frontend
npm install
```

2. **Run development server:**

```bash
npm run dev
```

The UI will be available at `http://localhost:3000`

### Development Server

The Vite dev server includes:
- Hot Module Replacement (HMR)
- API proxy to backend
- Fast refresh

### Code Organization

#### Adding a New Page

1. **Create page component** in `src/pages/`
2. **Add route** in `src/App.tsx`
3. **Create supporting components** in `src/components/`

Example:

```tsx
// src/pages/NewPage.tsx
import React from 'react'
import { PageSection, Title } from '@patternfly/react-core'

const NewPage: React.FC = () => {
  return (
    <>
      <PageSection>
        <Title headingLevel="h1">New Page</Title>
      </PageSection>
      <PageSection>
        {/* Content */}
      </PageSection>
    </>
  )
}

export default NewPage

// src/App.tsx
import NewPage from './pages/NewPage'

function App() {
  return (
    <Routes>
      {/* ... */}
      <Route path="/new-page" element={<NewPage />} />
    </Routes>
  )
}
```

#### Adding API Integration

1. **Define types** in `src/types/index.ts`
2. **Add API function** in `src/services/api.ts`
3. **Use React Query** in components

Example:

```tsx
// src/types/index.ts
export interface NewResource {
  id: string
  name: string
}

// src/services/api.ts
export const newResourceAPI = {
  list: async (): Promise<NewResource[]> => {
    const response = await api.get<APIResponse<NewResource[]>>('/new-resources')
    return response.data.data || []
  },
}

// Component usage
import { useQuery } from 'react-query'
import { newResourceAPI } from '@/services/api'

function MyComponent() {
  const { data, isLoading, error } = useQuery('newResources', newResourceAPI.list)
  // ...
}
```

### Development Best Practices

- **TypeScript**: Use strict type checking
- **Components**: Keep components small and focused
- **Hooks**: Follow React hooks rules
- **Testing**: Write tests for critical functionality
- **Accessibility**: Follow WCAG guidelines

## Testing

### Backend Tests

Run all tests:

```bash
cd backend
make test
```

Run with coverage:

```bash
make coverage
```

Run specific tests:

```bash
go test ./pkg/handlers -v
```

#### Writing Tests

```go
// pkg/handlers/example_test.go
package handlers

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestExampleHandler(t *testing.T) {
    handler := NewExampleHandler()
    // Test implementation
    assert.NotNil(t, handler)
}
```

### Frontend Tests

Run all tests:

```bash
cd frontend
npm test
```

Run with coverage:

```bash
npm run test:coverage
```

Run in UI mode:

```bash
npm run test:ui
```

#### Writing Tests

```tsx
// src/components/Example.test.tsx
import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import Example from './Example'

describe('Example Component', () => {
  it('renders correctly', () => {
    render(<Example />)
    expect(screen.getByText('Example')).toBeInTheDocument()
  })
})
```

## Building

### Backend Build

```bash
cd backend
make build
```

Output: `backend/bin/server`

### Frontend Build

```bash
cd frontend
npm run build
```

Output: `frontend/dist/`

### Docker Images

Build backend image:

```bash
docker build -f deployment/docker/Dockerfile.backend -t rhacm-monitor-backend:dev .
```

Build frontend image:

```bash
docker build -f deployment/docker/Dockerfile.frontend -t rhacm-monitor-frontend:dev .
```

### Local Testing with Docker

```bash
# Run backend
docker run -p 8080:8080 \
  -e ENABLE_AUTH=false \
  rhacm-monitor-backend:dev

# Run frontend
docker run -p 3000:80 \
  rhacm-monitor-frontend:dev
```

## Contributing

### Git Workflow

1. **Create a feature branch:**

```bash
git checkout -b feature/my-new-feature
```

2. **Make your changes**

3. **Run tests:**

```bash
cd backend && make test
cd frontend && npm test
```

4. **Commit your changes:**

```bash
git add .
git commit -m "Add my new feature"
```

5. **Push and create PR:**

```bash
git push origin feature/my-new-feature
```

### Code Style

#### Backend (Go)

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Run `gofmt` before committing
- Use meaningful variable names
- Add comments for exported functions

```bash
make fmt
make lint
```

#### Frontend (TypeScript/React)

- Follow [Airbnb Style Guide](https://github.com/airbnb/javascript/tree/master/react)
- Use Prettier for formatting
- Use ESLint for linting

```bash
npm run format
npm run lint
```

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat: add new feature`
- `fix: fix bug in handler`
- `docs: update documentation`
- `test: add tests for component`
- `refactor: refactor code`

### Pull Request Process

1. Ensure all tests pass
2. Update documentation if needed
3. Add tests for new functionality
4. Request review from maintainers
5. Address review feedback
6. Squash commits if requested

## Useful Commands

### Backend

```bash
# Run tests
make test

# Run with race detection
go test -race ./...

# Generate coverage report
make coverage

# Lint code
make lint

# Format code
make fmt

# Build binary
make build

# Clean build artifacts
make clean
```

### Frontend

```bash
# Install dependencies
npm install

# Run dev server
npm run dev

# Build for production
npm run build

# Run tests
npm test

# Run tests with UI
npm run test:ui

# Lint code
npm run lint

# Format code
npm run format

# Type check
npm run type-check
```

## Troubleshooting

### Backend Issues

**Import errors:**
```bash
cd backend
go mod tidy
```

**Build errors:**
```bash
cd backend
go clean -cache
go build ./...
```

### Frontend Issues

**Dependency issues:**
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

**Build issues:**
```bash
cd frontend
rm -rf dist
npm run build
```

## Resources

- [Go Documentation](https://golang.org/doc/)
- [React Documentation](https://react.dev/)
- [PatternFly Documentation](https://www.patternfly.org/)
- [RHACM Documentation](https://access.redhat.com/documentation/en-us/red_hat_advanced_cluster_management_for_kubernetes/)
- [Kubernetes Client Go](https://github.com/kubernetes/client-go)

