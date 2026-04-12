# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SaaS Médico - Backend API para una plataforma SaaS médica multi-clínica. Un usuario (profesional de salud) puede pertenecer a una o más clínicas; los pacientes pertenecen a la clínica. Soporta múltiples especialidades médicas.

**Tech Stack:** Go 1.24 · Gin · GORM · MySQL 8 · JWT (`golang-jwt/jwt/v5`) · Redis · OpenAI (`gpt-5.4` vía `openai-go/v3`) · Jasper Reports (binario externo)

## Common Commands

```bash
# Ejecutar la aplicación
go run cmd/api/main.go

# Hot reload (requiere air instalado)
air

# Build
go build -o bin/api cmd/api/main.go

# Dependencias
go mod tidy

# Tests
go test ./...
go test ./internal/modules/auth/...   # módulo específico

# Seed de catálogos de nutrición (idempotente)
go run cmd/seed/main.go

# Seed de agenda y menú del sistema (correr una vez contra la BD)
mysql -u$DB_USER -p$DB_PASSWORD $DB_NAME < resources/seed_agenda.sql
mysql -u$DB_USER -p$DB_PASSWORD $DB_NAME < resources/seed_menu.sql
```

Health check: `GET /ping` → `{"message":"pong","status":"healthy"}`

## Architecture

### Módulos actuales

```
internal/modules/
├── auth/          # JWT auth, usuarios, roles, permisos, refresh tokens
├── admin/         # Clínicas, sucursales, consultorios, profesiones, roles de transacción, planes SaaS, suscripciones
├── pacientes/     # Pacientes, pre-pacientes, aplicaciones
├── agenda/        # Citas, sesiones, horarios médico, bloqueos de agenda
├── cobros/        # Cobros por sesión, pagos, egresos
├── historia/      # Historia clínica, formularios dinámicos, alergias, antecedentes, diagnósticos
├── tests/         # Tests psicológicos, reglas de puntaje, sesiones de test
├── nutricion/     # Planes de dieta, menús, alimentos, ejercicios, progreso, gamificación (XP/logros)
├── psicologia/    # [stub] listo para expandir
└── odontologia/   # [stub] listo para expandir
```

### Estructura de cada módulo

```
module_name/
├── handlers/      # HTTP handlers — parsea request, llama service, retorna response
├── services/      # Lógica de negocio
├── repositories/  # Acceso a datos GORM
├── models/        # Modelos GORM + DTOs (en archivos separados)
└── routes.go      # Instancia repo→service→handler y registra rutas
```

La inyección de dependencias ocurre **dentro de `routes.go`**: cada `RegisterRoutes` instancia su propio repo, service y handler usando `database.GetDB()`.

### Orden de inicio (`cmd/api/main.go`)

```
config.LoadConfig() → database.Connect() → [RunMigrations() — comentado] → redis.NewClient() → auth.Setup() → register routes → router.Run()
```

`auth.Setup()` debe llamarse antes de `auth.GetAuthMiddleware()`. `database.RunMigrations()` está **comentado por defecto** en `main.go`; descoméntalo solo cuando necesites correr migraciones, luego vuélvelo a comentar.

**Redis:** la dirección está hardcodeada en `main.go` (`162.243.161.156:6379`), no lee `REDIS_ADDR` del `.env`.

**OpenAI:** inicializado en `internal/shared/openia/openia.go` usando `openai-go/v3`. El servicio se pasa a `nutricion.RegisterRoutes`. El modelo usado es `gpt-5.4`.

**Scheduler:** `internal/shared/scheduler/scheduler.go` — `StartCron(job JobFunc)` ejecuta el job diariamente a medianoche (zona `America/Guayaquil`). Se llama desde `nutricion.RegisterRoutes` para desactivar menús viejos.

### Convenciones de modelos

- **IDs:** `uint` con `autoIncrement` en todos los modelos (no UUID)
- **Soft delete:** campo `State string` con `char(1)` — `'A'` activo, `'I'` inactivo. No se usa `gorm.DeletedAt`. Constantes `models.StateActivo` / `models.StateInactivo` en `internal/shared/models/base.go`.
- **Base struct:** `internal/shared/models/base.go` expone `Base` (ID + State + CreadoEn + ActualizadoEn). Los nuevos modelos pueden embeber `Base` o definir los campos directamente.
- **Timestamps:** `CreadoEn`/`ActualizadoEn` con `gorm:"autoCreateTime/autoUpdateTime"`. Los modelos de auth usan `CreatedAt`/`UpdatedAt` (convención GORM estándar).
- **Queries:** siempre filtrar `WHERE state = 'A'` en listas. Soft delete → `UPDATE SET state = 'I'`.
- **TableName():** añadir cuando el plural GORM no coincide con el esquema (ej: `Rol` → `roles`, `HorarioMedico` → `horarios_medico`).
- **Dos campos de estado:** algunos modelos usan `State char(1)` (activo/inactivo de registro) **y** `Estado string` (estado de negocio como `ACTIVA/COMPLETADA/CANCELADA`). No confundirlos.

### Autenticación

**Staff (profesionales/admin):**
- Tokens JWT con `UserID uint`. Claims en contexto Gin: `userID`, `clinicaID`, `email`, `rolID`, `rolName`, `permisos`.
- Los handlers recuperan el user con `c.GetUint("userID")` y la clínica activa con `c.GetUint("clinicaID")`.
- Constantes de roles en `internal/modules/auth/models/rol.go` (`models.RolSuperAdmin`, `models.RolAdmin`, etc.).
- Refresh token rotation: el token actual se revoca en cada `/refresh`, logout y cambio de contraseña.
- Login en dos pasos: primer POST devuelve clínicas disponibles, segundo POST con `clinica_id` completa el login.

**Pacientes (app móvil):**
- Sistema JWT separado con claims distintos: `user_id` → `paciente_id`, `rol_name` → `"paciente"`, más `clinica_id` y `aplicacion_id`.
- Contraseña por defecto al crear paciente: `Usuario123`.
- Ver `NUTRICION_PACIENTE_CONTEXT.md` para endpoints y contexto de la API de pacientes.

### Middleware

```go
router.Use(authMiddleware.RequireAuth())
router.Use(authMiddleware.RequireRoles("super_admin", "admin"))
router.Use(authMiddleware.RequirePermissions("psicologia.ver"))
router.Use(authMiddleware.RequireAnyPermission("psicologia.ver", "nutricion.ver"))
```

CORS se aplica globalmente en `internal/middleware/cors.go`. El módulo `historia` usa roles: `RolAdmin, RolNutriologo, RolPsicologo, RolSuperAdmin`.

### Firma de RegisterRoutes

La mayoría de módulos: `RegisterRoutes(api *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware)`.

Nutrición tiene firma extendida: `RegisterRoutes(api, authMiddleware, rdb *redis.Client, openiaService *openia.OpenIaService)`.

### Respuestas HTTP

Usar siempre `internal/shared/responses`:

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

### Utilidades compartidas

**Paginación** (`internal/shared/utils`):
```go
page, pageSize := utils.GetPaginationParams(c)   // lee ?page= y ?page_size= (máx 100)
offset := utils.GetOffset(page, pageSize)
```

**Subida de archivos** (`internal/shared/uploads`):
```go
result, err := uploads.SaveFile(c, fileHeader, "subdir", uploads.AllowedImageTypes)
// Guarda en storage/uploads/<subdir>/. Límite: 10 MB.
// AllowedImageTypes: .jpg .jpeg .png .gif .webp
// AllowedDocTypes:   .pdf .doc .docx .xls .xlsx
uploads.DeleteFile(result.FilePath)
```

Archivos estáticos servidos en `/storage` → `./storage` (configurado en `main.go`).

**Generación de PDFs — dos sistemas:**

1. `internal/shared/pdfbuilder/` — PDFs por código con `maroto/v2` (Poppins font, logo watermark). Implementar la interfaz `UseCase` en `pdfbuilder/usecases/` para cada tipo de documento (ejemplo: `nutricion-use-case.go` genera el PDF del menú semanal).
```go
m, err := pdfSvc.GeneratePdfBuilder()   // obtiene core.Maroto configurado
uc := usecases.NewMenuPdfUseCase(dieta, menu, m, logoPath, outputPath)
uc.CreatePdf()
```

2. `internal/shared/reports/` — Reportes desde templates `.jasper` vía binario `jasper-starter`. Templates en `resources/jasper_templates/`, salida en `storage/reports/`.
```go
jasperSvc := reports.NewJasperService(jasperPath, jdbcPath, host, port, db, user, pass)
outPath, err := jasperSvc.GenerateReport(reports.ReportParams{
    TemplateName: "nombre_template",   // sin extensión
    OutputName:   "archivo_salida",
    Format:       reports.FormatPDF,
    Parameters:   map[string]interface{}{"clinica_id": 1},
})
reports.DeleteReport(outPath)   // limpiar después de servir
```

**OpenAI / Conversaciones IA** (`internal/shared/openia/`): el historial de chat por paciente se almacena en Redis con clave `conv:<paciente_id>` (helper `openia.BuildConversationKey(pacienteID)`). El servicio se instancia una sola vez en `main.go` y se pasa a `nutricion.RegisterRoutes`.

### Agregar nuevos módulos

1. Crear `internal/modules/nuevo_modulo/` con la estructura estándar.
2. En `routes.go`: `func RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware)`.
3. Agregar los modelos a `internal/database/migrate.go` en el orden correcto de FKs.
4. Registrar en `cmd/api/main.go`.

### Migraciones

`internal/database/migrate.go` define `RunMigrations()` con grupos comentados por módulo. Actualmente activos: `historiaModels.TipoFormulario{}` y `nutricionModels.NutricionMenuPlantilla{}`. Para migrar un nuevo grupo, descomenta el bloque correspondiente.

**SQL manuales:** `resources/migrations/` contiene migraciones SQL que deben correrse directamente. Ejemplo: `002_nutricion_menu_detalle.sql` renombra tablas de nutrición (`nutricion_dieta_detalle` → `nutricion_menu_detalle`).

Orden de dependencias FK para GORM AutoMigrate:
1. Auth → 2. Admin → 3. Pacientes → 4. Agenda → 5. Cobros → 6. Historia (catálogos/formularios) → 7. Historia (registros paciente) → 8. Tests psicológicos → 9. Nutrición (catálogos) → 10. Nutrición (dieta/menú) → 11. Nutrición (registros y seguimiento)

### Configuración

Variables de entorno (`.env`):
```
DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
SERVER_PORT, ENVIRONMENT
JWT_SECRET, JWT_EXPIRATION_HOURS, JWT_REFRESH_DAYS
OPEN_AI_API_KEY
REDIS_ADDR        # actualmente ignorado — ver hardcode en main.go
```

### Esquema completo de BD

El esquema completo está en `SaasMedico_esquema_contexto.md`. Consultar al implementar nuevos módulos.

### Módulos pendientes de implementar

- **tratamiento/** — planes de tratamiento, items, cobros de plan
- **odontologia/** — odontogramas, piezas, caras, eventos (expandir stub)
- **psicologia/** — expandir stub (lógica de tests ya está en `tests/`)
- **documentos/** — consentimientos, prescripciones
- **tareas/** — tareas paciente, progreso diario, observaciones de sesión
- **notificaciones/** — cola de notificaciones WhatsApp/SMS/email
- **recursos/** — recursos psicoeducativos, plantillas de intervención
