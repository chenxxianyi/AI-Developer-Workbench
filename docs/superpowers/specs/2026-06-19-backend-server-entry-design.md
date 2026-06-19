# Backend Server Entry Design

## Goal

Restore the missing Go backend executable entry point so the existing API can be built and started with:

```powershell
cd backend
go run ./cmd/server
```

The entry point will use the existing real OpenAI-compatible service. Mock AI behavior is explicitly out of scope.

## Root Cause

The repository documentation, Dockerfile, and development plans all reference `backend/cmd/server/main.go`, but the initial repository commit never included the `cmd` directory. The internal packages compile independently, but there is no executable package to assemble and start them.

## Architecture

Create a single entry point at `backend/cmd/server/main.go`. It will remain an application composition root: it will initialize existing dependencies and register routes, but it will not contain business logic.

Startup order:

1. Load and validate configuration from `.env` and environment variables.
2. Configure structured logging and Gin mode.
3. Create the configured upload and temporary directories.
4. Connect to MySQL and run the existing optional GORM migrations.
5. Construct repositories.
6. Construct shared services, including the real `OpenAICompatibleService`.
7. Construct all five tool services and HTTP handlers.
8. Register middleware and all existing routes under `/api`.
9. Start a configured `http.Server`.
10. On `SIGINT` or `SIGTERM`, gracefully stop HTTP traffic and close the database connection.

## Dependency Assembly

The entry point will instantiate the existing components in dependency order:

- Repositories: reports, generated files, and report assets.
- Shared services: report, file, ZIP, export, and OpenAI-compatible AI services.
- Tool services: Agent Config, DB Schema, UI Review, Project Doctor, and API Doc.
- Handlers: health, system, dashboard, tools, reports, exports, and tool execution.

No new service abstraction or `internal/app` package will be introduced because one composition file is sufficient for the current scope.

## HTTP Behavior

The Gin engine will use the existing middleware in this order:

1. Request ID
2. Recovery
3. Request logger
4. CORS

All existing handler registration functions will attach to an `/api` route group. The HTTP server will set header, read, write, and idle timeouts. The write timeout will exceed the configured AI request timeout so valid AI calls are not terminated prematurely.

## Configuration and Failure Handling

- The existing `config.LoadConfig(".env")` behavior will be used.
- `AI_API_KEY` remains required by the current configuration validation.
- `AI_MOCK_MODE` will not be added or read.
- Missing required configuration, directory creation failure, database connection failure, or migration failure will stop startup with a clear error.
- Secrets, passwords, API keys, and complete DSNs will not be logged.
- Server shutdown will use a bounded timeout; forced termination will be logged if graceful shutdown cannot finish in time.

## Testing

The implementation will be developed test-first:

1. Add a structural test that requires the server composition function to build a Gin engine with the expected routes and middleware behavior without opening a listening socket.
2. Verify the test fails because the entry point/application composition does not exist.
3. Add the minimal server entry implementation.
4. Run the focused test, all backend tests, and `go build ./cmd/server`.

Database-dependent startup will not be faked with mock application data. Tests will isolate route assembly from the live listener and use lightweight dependency seams only where required to verify startup composition.

## Scope

Included:

- Missing executable entry point
- Existing dependency assembly
- Existing route registration
- Real AI provider wiring
- HTTP timeouts and graceful shutdown
- Build and startup-structure verification

Excluded:

- Mock AI service or mock report data
- Frontend changes
- New endpoints
- Business logic changes
- Database schema redesign
- Unrelated refactoring
