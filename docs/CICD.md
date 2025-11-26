# CI/CD Pipeline Documentation

## Overview

The opgl-data service uses GitHub Actions for continuous integration and deployment. The pipeline automatically runs when code is pushed or when pull requests are created.

## Pipeline Triggers

| Event | Branch | What Runs |
|-------|--------|-----------|
| Push | `main` | Tests + Docker Build & Push |
| Pull Request | `main` | Tests only |

## Job 1: Test & Build

This job runs on every push and pull request to `main`.

### Steps

1. **Checkout code** - Clones the repository
2. **Set up Go 1.21** - Installs Go with dependency caching enabled
3. **Download dependencies** - Runs `go mod download`
4. **Verify dependencies** - Runs `go mod verify` to ensure integrity
5. **Run go vet** - Static analysis to catch common errors
6. **Run tests with coverage** - Executes all tests with race detection
7. **Check coverage threshold** - Enforces 90% minimum coverage. If coverage is below 90%, the pipeline fails.
8. **Build application** - Compiles the binary to verify it builds
9. **Upload coverage report** - Stores coverage.out as an artifact

## Job 2: Docker Build & Push

This job only runs on pushes to `main` (not on PRs) and requires Job 1 to pass first.

### Steps

1. **Set up Docker Buildx** - Enables advanced Docker build features
2. **Login to GHCR** - Authenticates with GitHub Container Registry
3. **Build & Push** - Creates and publishes the Docker image

### Image Tags

When pushed to main, the image receives these tags:
- `ghcr.io/opglol/opgl-data:main` - Branch name
- `ghcr.io/opglol/opgl-data:<sha>` - Commit SHA
- `ghcr.io/opglol/opgl-data:latest` - Latest tag

## Coverage Requirements

The pipeline enforces a **90% minimum test coverage** threshold.

### What Gets Measured

| Package | Description |
|---------|-------------|
| `internal/api` | HTTP handlers and router |
| `internal/config` | Configuration loading |
| `internal/middleware` | Logging middleware |
| `internal/services` | Riot API service |

### What's Excluded

- `main.go` - Entry point (excluded from coverage calculation)
- `internal/models` - Data structures only (no logic to test)

### If Coverage Drops Below 90%

The pipeline will fail with this message:
```
Coverage XX.X% is below threshold 90%
```

You must add tests to increase coverage before the code can be merged.

## For Engineers: Workflow

### Before Pushing Code

1. Run tests locally:
   ```bash
   go test -v -race -coverprofile=coverage.out ./...
   ```

2. Check coverage:
   ```bash
   go tool cover -func=coverage.out | grep total:
   ```

3. Ensure coverage is at least 90%

### Creating a Pull Request

1. Push your branch
2. Create PR targeting `main`
3. Wait for CI checks to pass
4. Coverage report is available in the Actions artifacts

### Merging to Main

1. PR must have passing CI checks
2. After merge, Docker image is automatically built and pushed
3. Image is available at `ghcr.io/opglol/opgl-data:latest`

## Troubleshooting

### Test Failures

Check the "Run tests with coverage" step in the GitHub Actions log for details.

### Coverage Below Threshold

1. Download the coverage report artifact
2. Run locally to see uncovered lines:
   ```bash
   go tool cover -html=coverage.out
   ```
3. Add tests for uncovered code paths

### Docker Build Failures

- Ensure Dockerfile builds locally: `docker build -t opgl-data .`
- Check for missing dependencies in go.mod

## Environment Variables

The CI pipeline uses GitHub-provided secrets. No additional secrets configuration required.

| Variable | Description |
|----------|-------------|
| `GITHUB_TOKEN` | Auto-provided, used for GHCR authentication |
