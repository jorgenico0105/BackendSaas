# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SaaS Médico - Backend API para una plataforma SaaS médica multi-clínica. Un usuario (profesional de salud) puede pertenecer a una o más clínicas; los pacientes pertenecen a la clínica. Soporta múltiples especialidades médicas.

**Tech Stack:** Go 1.24 · Gin · GORM · MySQL 8 · JWT (`golang-jwt/jwt/v5`) · Jasper Reports (binario externo)

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

# Seed de catálogos de nutrición (idempotente — alimentos, grupos, tipos de comida)
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

### Orden de inicio (cmd/api/main.go)

```
config.LoadConfig() → database.Connect() → auth.Setup() → register routes → router.Run()
```

`auth.Setup()` debe llamarse antes de `auth.GetAuthMiddleware()`. `database.RunMigrations()` está **comentado por defecto** — las migraciones son opt-in (ver sección Migraciones).

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
- Tokens JWT con `UserID uint` (no UUID). Claims en contexto Gin: `userID`, `email`, `rolID`, `rolName`, `permisos`.
- Los handlers recuperan el user con `c.GetUint("userID")`.
- Constantes de roles en `internal/modules/auth/models/rol.go` (`models.RolSuperAdmin`, `models.RolAdmin`, etc.).
- Refresh token rotation: el token actual se revoca en cada `/refresh`, logout y cambio de contraseña.
- Login es en dos pasos: primer POST devuelve clínicas disponibles, segundo POST con `clinica_id` completa el login.

**Pacientes (app móvil):**
- Sistema JWT separado con claims distintos: `user_id` → `paciente_id`, `rol_name` → `"paciente"`, más `clinica_id` y `aplicacion_id`.
- Contraseña por defecto al crear paciente: `Usuario123`.
- Ver `NUTRICION_PACIENTE_CONTEXT.md` para endpoints y contexto de la API de pacientes.

### Agregar nuevos módulos

1. Crear `internal/modules/nuevo_modulo/` con la estructura estándar.
2. En `routes.go`: `func RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware)`.
3. Agregar los modelos a `internal/database/migrate.go` en el orden correcto de FKs.
4. Registrar en `cmd/api/main.go`.

### Middleware

```go
router.Use(authMiddleware.RequireAuth())
router.Use(authMiddleware.RequireRoles("super_admin", "admin"))
router.Use(authMiddleware.RequirePermissions("psicologia.ver"))
router.Use(authMiddleware.RequireAnyPermission("psicologia.ver", "nutricion.ver"))
```

CORS se aplica globalmente en `internal/middleware/cors.go`. Por ahora todos los módulos usan `RequireAuth()` + `RequireRoles(RolSuperAdmin, RolAdmin)`.

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

### Migraciones

**GORM AutoMigrate:** `internal/database/migrate.go` define `RunMigrations()` con grupos comentados por módulo. **La llamada en `main.go` está comentada** (`//database.RunMigrations()`): las migraciones son opt-in, no automáticas. Para migrar, descomenta el grupo deseado en `migrate.go` y activa temporalmente la llamada en `main.go`.

**SQL manuales:** `resources/migrations/` contiene migraciones SQL que deben correrse directamente. Ejemplo: `002_nutricion_menu_detalle.sql` renombra tablas de nutrición (`nutricion_dieta_detalle` → `nutricion_menu_detalle`, `nutricion_dieta_alimentos` → `nutricion_menu_alimentos`).

Orden de dependencias FK para GORM AutoMigrate:
1. Auth (Rol, User, UsuarioRol, RefreshToken)
2. Admin (Clinica, EstiloClinica, Transaccion, RolTransaccion → luego Sucursal, Consultorio, Profesion, UsuarioClinica, UsuarioConsultorio, PlanSaas, EstadoSuscripcion, Suscripcion, BloqueoAcceso)
3. Pacientes (Paciente, PacienteUsuario, Aplicacion, PacienteAplicacion)
4. Agenda (TipoCita, EstadoCita, Cita, Sesion, HorarioMedico, BloqueoAgenda)
5. Cobros (EstadoCobro, MedioPago, TipoEgreso, CobroSesion, Pago, Egreso)
6. Historia — catálogos y formularios (Formulario, TipoFormulario, AlergiaCatalogo, …)
7. Historia — registros paciente (HistoriaClinica, PacienteAlergia, PacienteDiagnostico, …)
8. Tests psicológicos (TestRegla, Test, TestRespuesta, SesionTest, …)
9. Nutrición — catálogos (NutricionTipoComida, NutricionAlimento, NutricionDietaCatalogo, NutricionEjercicioCatalogo, NutricionLogroCatalogo)
10. Nutrición — dieta/menú (NutricionDietaPaciente, NutricionMenu, NutricionMenuDetalle, NutricionMenuAlimento) — ojo: tablas renombradas por `resources/migrations/002_nutricion_menu_detalle.sql`
11. Nutrición — registros y seguimiento (NutricionR24H, NutricionRegistroComida, NutricionProgresoPaciente, NutricionLogroPaciente, NutricionPacienteXP, …)

### Módulos pendientes de implementar

Pendientes según `SaasMedico_esquema_contexto.md`:
- **tratamiento/** — planes de tratamiento, items, cobros de plan
- **odontologia/** — odontogramas, piezas, caras, eventos (expandir stub)
- **psicologia/** — tests psicológicos, reglas de puntaje (expandir stub; la lógica de tests psicológicos ya está en `tests/`)
- **documentos/** — consentimientos, prescripciones
- **tareas/** — tareas paciente, progreso diario, observaciones de sesión
- **notificaciones/** — cola de notificaciones WhatsApp/SMS/email
- **recursos/** — recursos psicoeducativos, plantillas de intervención

### Jasper Reports

`internal/shared/reports/JasperService` invoca un binario externo **jasper-starter**. Requiere:
- Template `.jasper` compilado en `resources/jasper_templates/`
- JDBC MySQL driver JAR
- Ruta al binario jasper-starter

Reportes generados en `storage/reports/`.

### Configuración

Variables de entorno (`.env`):
```
DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
SERVER_PORT, ENVIRONMENT
JWT_SECRET, JWT_EXPIRATION_HOURS, JWT_REFRESH_DAYS
OPEN_AI_API_KEY   # opcional, para funciones de IA
```

### Esquema completo de BD

El esquema completo de la base de datos está en `SaasMedico_esquema_contexto.md`. Consultar al implementar nuevos módulos.
