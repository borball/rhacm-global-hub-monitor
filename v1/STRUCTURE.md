# Project Structure - v1

## Root Directory
```
v1/
├── README.md                  Main project documentation
├── QUICKSTART.md             5-minute getting started guide
├── VERSION.md                Version-specific information
├── STRUCTURE.md              This file
├── LICENSE                   Apache 2.0 license
└── .gitignore               Git ignore rules
```

## Source Code
```
backend/                      Go backend application
├── cmd/server/              Main application entry
├── pkg/
│   ├── api/                 API routes
│   ├── auth/                Authentication
│   ├── client/              Kubernetes clients
│   ├── handlers/            HTTP handlers
│   └── models/              Data models
└── internal/
    ├── config/              Configuration
    └── middleware/          HTTP middleware

frontend/                     React TypeScript application
├── src/
│   ├── components/          React components
│   ├── pages/               Page components
│   ├── services/            API services
│   ├── types/               TypeScript types
│   └── utils/               Utilities
└── package.json             Dependencies

frontend-static/              Static HTML version (deployed)
├── index.html               Main page
├── app.js                   Application logic
├── styles.css               Styling
└── httpd-with-proxy.conf    Proxy configuration
```

## Deployment
```
deployment/
├── k8s/                     Kubernetes manifests
│   ├── namespace.yaml
│   ├── serviceaccount.yaml
│   ├── clusterrole.yaml
│   ├── backend-*.yaml
│   ├── frontend-*.yaml
│   └── route.yaml
└── docker/                  Docker configurations
    ├── Dockerfile.backend
    ├── Dockerfile.frontend
    └── nginx.conf

operator/                     OpenShift Operator
├── config/
│   ├── crd/                 Custom Resource Definitions
│   ├── rbac/                RBAC manifests
│   └── samples/             Sample configurations
└── PROJECT                  Operator metadata
```

## Documentation
```
docs/
├── README.md                Documentation index
├── API.md                   API reference
├── DEPLOYMENT.md            Deployment guide
├── DEVELOPMENT.md           Developer guide
├── ARCHITECTURE.md          System architecture
├── BUILD_AND_DEPLOY.md      Build instructions
│
├── guides/                  User guides
│   ├── API_EXAMPLES.md      Real API examples
│   ├── ACCESS_INSTRUCTIONS.md
│   └── DOCKER_REGISTRY_CHANGES.md
│
├── project-status/          Project documentation
│   ├── PROJECT_SUMMARY.md
│   ├── SUCCESS_REPORT.md
│   ├── DEPLOYMENT_COMPLETE.md
│   └── FINAL_SUMMARY.md
│
└── test-results/            Test documentation
    ├── TESTING_RESULTS.md
    ├── COMPLETE_API_TEST_RESULTS.md
    ├── FINAL_TEST_RESULTS.md
    ├── DEPLOYMENT_TEST_RESULTS.md
    └── SPOKE_CLUSTER_TEST_RESULTS.md
```

## Build Artifacts
```
backend/bin/                 Compiled binaries
frontend/dist/               Built frontend
*.tar.gz                     Archives (gitignored)
```

---

**Clean, organized structure for easy navigation and maintenance!**
