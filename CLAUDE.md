# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SaaS Médico — Backend API para una plataforma SaaS médica multi-clínica. Un usuario (profesional de salud) puede pertenecer a una o más clínicas; los pacientes pertenecen a la clínica. Soporta múltiples especialidades médicas.

**Tech Stack:** Go 1.26 · Gin · GORM · MySQL 8 · JWT (`golang-jwt/jwt/v5`) · Redis · OpenAI (`openai-go/v3`) · Jasper Reports (binario externo) · maroto/v2 (PDF)

## Common Commands

```bash
go run cmd/api/main.go          # Start API (port from .env SERVER_PORT)
air                             # Hot reload
go build -o bin/api cmd/api/main.go
go mod tidy
go test ./...
go test ./internal/modules/auth/...   # Single module tests

# Seed data (idempotent)
go run cmd/seed/main.go
mysql -u$DB_USER -p$DB_PASSWORD $DB_NAME < resources/seed_agenda.sql
mysql -u$DB_USER -p$DB_PASSWORD $DB_NAME < resources/seed_menu.sql
```

Health check: `GET /ping` → `{"message":"pong","status":"healthy"}`

## Architecture

### Module layout

```
internal/modules/
├── auth/          # JWT auth, usuarios, roles, permisos, refresh tokens
├── admin/         # Clínicas, sucursales, consultorios, profesiones, planes SaaS, suscripciones
├── pacientes/     # Pacientes, pre-pacientes, aplicaciones
├── agenda/        # Citas, sesiones, horarios médico, bloqueos
├── cobros/        # Cobros por sesión, pagos, egresos
├── historia/      # Historia clínica, formularios dinámicos, alergias, diagnósticos
├── tests/         # Tests psicológicos, reglas de puntaje, sesiones de test
├── nutricion/     # Planes de dieta, menús, alimentos, ejercicios, progreso, XP/logros
├── psicologia/    # [stub]
└── odontologia/   # [stub]
```

Each module follows this structure:
```
module_name/
├── handlers/      # One file: <module>_handler.go
├── services/      # One file: <module>_service.go
├── repositories/  # One file: <module>_repository.go
├── models/        # GORM models + DTOs (separate files: <model>.go + dto.go)
└── routes.go      # DI: repo→service→handler, route registration
```

Each layer uses **one file per module** (not one file per entity). DI happens **inside `routes.go`** via `database.GetDB()`. No global service singletons.

### Startup order (`cmd/api/main.go`)

```
config.LoadConfig() → database.Connect() → redis.NewClient() → auth.Setup() → register routes → router.Run()
```

- `auth.Setup()` must be called before `auth.GetAuthMiddleware()`.
- `database.RunMigrations()` is **commented out** by default — uncomment only when needed, then re-comment.
- Redis is **hardcoded** in `main.go` (`127.0.0.1:6379`, password `nico1234.`); `REDIS_ADDR` from `.env` is not used.
- Server binds on `127.0.0.1:<port>` (localhost only, behind Nginx); trusted proxies set to `127.0.0.1` and `::1`.

### RegisterRoutes signatures

```go
// Most modules
RegisterRoutes(api *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware)

// Nutrición (extended — needs Redis + OpenAI)
RegisterRoutes(api, authMiddleware, rdb *redis.Client, openiaService *openia.OpenIaService)
```

## Model Conventions

- **IDs:** `uint` with `autoIncrement` (no UUID).
- **Soft delete:** `State string` `char(1)` — `'A'` active, `'I'` inactive. No `gorm.DeletedAt`. Use constants `models.StateActivo` / `models.StateInactivo` (`internal/shared/models/base.go`).
- **Base struct:** embed `Base` from `internal/shared/models/base.go` (ID + State + CreadoEn + ActualizadoEn) or define fields directly.
- **Timestamps:** `CreadoEn`/`ActualizadoEn` with `gorm:"autoCreateTime/autoUpdateTime"`. Auth models use `CreatedAt`/`UpdatedAt` (standard GORM).
- **Queries:** always filter `WHERE state = 'A'` in list queries. Soft delete → `UPDATE SET state = 'I'`.
- **TableName():** add whenever GORM's plural doesn't match the schema (e.g. `Rol` → `roles`).
- **Two status fields:** some models have both `State char(1)` (record active/inactive) **and** `Estado string` (business state like `ACTIVA/COMPLETADA/CANCELADA`). Don't conflate them.

## Authentication

**Staff (profesionales/admin):**
- JWT claims in Gin context: `userID`, `clinicaID`, `email`, `rolID`, `rolName`, `permisos`.
- Retrieve with `c.GetUint("userID")` and `c.GetUint("clinicaID")`.
- Role constants in `internal/modules/auth/models/rol.go`.
- Two-step login: first POST returns available clinics; second POST with `clinica_id` completes auth.
- Refresh token rotation on every `/refresh`, logout, and password change.

**Pacientes (mobile app):**
- Separate JWT system with claims: `user_id` → `paciente_id`, `rol_name` → `"paciente"`, plus `clinica_id` and `aplicacion_id`.
- Default password on creation: `Usuario123`.
- Patient-facing nutrition routes live under `/api/v1/paciente/nutricion/` in the nutricion module.

## Middleware

```go
router.Use(authMiddleware.RequireAuth())
router.Use(authMiddleware.RequireRoles("super_admin", "admin"))
router.Use(authMiddleware.RequirePermissions("psicologia.ver"))
router.Use(authMiddleware.RequireAnyPermission("psicologia.ver", "nutricion.ver"))
```

CORS applied globally in `internal/middleware/cors.go`.

## HTTP Responses

Always use `internal/shared/responses` — never raw `c.JSON`:

```go
responses.Success(c, message, data)
responses.Created(c, message, data)
responses.Paginated(c, data, page, pageSize, totalItems)
responses.BadRequest(c, message)
responses.NotFound(c, message)
responses.Unauthorized(c, message)
responses.Forbidden(c, message)
responses.InternalError(c, message)
```

## Shared Utilities

**Pagination** (`internal/shared/utils`):
```go
page, pageSize := utils.GetPaginationParams(c)   // reads ?page= and ?page_size= (max 100)
offset := utils.GetOffset(page, pageSize)
```

**File uploads** (`internal/shared/uploads`):
```go
result, err := uploads.SaveFile(c, fileHeader, "subdir", uploads.AllowedImageTypes)
// Saves to storage/uploads/<subdir>/. Limit: 10 MB.
uploads.DeleteFile(result.FilePath)
```
Static files served at `/storage` → `./storage`.

**PDF — two systems:**

1. `internal/shared/pdfbuilder/` — code-generated PDFs with maroto/v2. Implement `UseCase` interface in `pdfbuilder/usecases/` per document type.
```go
pdfS := pdfbuilder.NewPdfBuilder(watermarkPath)
m, err := pdfS.GeneratePdfBuilder()
uc := usecases.NewMenuPdfUseCase(dieta, menu, m, logoPath, outputPath)
uc.CreatePdf()
```

2. `internal/shared/reports/` — Jasper template reports via `jasper-starter` binary. Templates in `resources/jasper_templates/`, output in `storage/reports/`.
```go
jasperSvc := reports.NewJasperService(jasperPath, jdbcPath, host, port, db, user, pass)
outPath, err := jasperSvc.GenerateReport(reports.ReportParams{
    TemplateName: "nombre_template",
    OutputName:   "archivo_salida",
    Format:       reports.FormatPDF,
    Parameters:   map[string]interface{}{"clinica_id": 1},
})
reports.DeleteReport(outPath)
```

**OpenAI / Chat** (`internal/shared/openia/`): chat history per patient stored in Redis, key `conv:<paciente_id>` (use `openia.BuildConversationKey(pacienteID)`). Service instantiated once in `main.go`.

**Scheduler** (`internal/shared/scheduler/`): `StartCron(job JobFunc)` runs daily at midnight `America/Guayaquil`. Called from `nutricion.RegisterRoutes` to deactivate old menus.

## Migrations

`internal/database/migrate.go` has `RunMigrations()` with groups commented per module. To migrate a new group, uncomment its block, run, then re-comment.

Manual SQL migrations live in `resources/migrations/` and must be run directly against the DB.

FK dependency order for AutoMigrate:
1. Auth → 2. Admin → 3. Pacientes → 4. Agenda → 5. Cobros → 6. Historia (catálogos) → 7. Historia (registros) → 8. Tests → 9. Nutrición (catálogos) → 10. Nutrición (dieta/menú) → 11. Nutrición (registros/seguimiento)

## Configuration

`.env` variables:
```
DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
SERVER_PORT, ENVIRONMENT
JWT_SECRET, JWT_EXPIRATION_HOURS, JWT_REFRESH_DAYS
OPEN_AI_API_KEY
```

## Adding a New Module

1. Create `internal/modules/nuevo_modulo/` with the standard structure.
2. In `routes.go`: `func RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware)`.
3. Add models to `internal/database/migrate.go` in FK order.
4. Register in `cmd/api/main.go`.

## Pending Modules

- **tratamiento/** — treatment plans, items, plan billing
- **odontologia/** — odontograms, teeth, faces, events (expand stub)
- **psicologia/** — expand stub (test logic already in `tests/`)
- **documentos/** — consents, prescriptions
- **tareas/** — patient tasks, daily progress, session observations
- **notificaciones/** — WhatsApp/SMS/email notification queue
- **recursos/** — psychoeducational resources, intervention templates
