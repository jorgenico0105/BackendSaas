# Contexto: App del Paciente — Nutrición

Este documento describe la autenticación del paciente y todos los endpoints
de nutrición disponibles para la **app del paciente** (mobile/web).

---

## Stack y convenciones

- **Backend:** Go 1.24 · Gin · GORM · MySQL 8
- **Base URL API:** `/api/v1`
- **Auth:** JWT Bearer. Adjuntar en cada request:
  `Authorization: Bearer <access_token>`
- **Fechas en request:** siempre string `"YYYY-MM-DD"`.
- **IDs:** siempre `uint` (entero positivo), nunca UUID.
- **Respuesta exitosa:**
  ```json
  { "success": true, "message": "...", "data": { ... } }
  ```
- **Respuesta paginada:**
  ```json
  { "success": true, "data": [...], "page": 1, "page_size": 10, "total": 50 }
  ```
- **Error:**
  ```json
  { "success": false, "message": "descripción del error" }
  ```
- **Paginación:** `?page=1&page_size=20` en todos los GET de listas.

---

## 1. Autenticación del Paciente

El paciente tiene un sistema de login **separado** del médico/admin.
Cada clínica tiene sus propias apps (`Aplicacion`), y el paciente
debe tener acceso asignado a la app específica antes de poder ingresar.

### 1.1 Login del paciente

```
POST /api/v1/pacientes/login
```
**Sin token requerido (pública).**

```json
{
  "username": "1234567890",    // número de documento o teléfono — requerido
  "password": "Usuario123",    // mín 6 caracteres — requerido
  "clinica_id": 1,             // ID de la clínica — requerido
  "aplicacion_id": 2           // ID de la app (ej: app nutrición) — requerido
}
```

**Respuesta exitosa:**
```json
{
  "success": true,
  "message": "Login exitoso",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "paciente_id": 11,
    "clinica_id": 1,
    "aplicacion_id": 2,
    "username": "Juan Pérez"
  }
}
```

**Errores posibles:**
- `401 Unauthorized` — credenciales inválidas (usuario o contraseña incorrectos)
- `403 Forbidden` — el paciente no tiene acceso a esa aplicación

### 1.2 Estructura del JWT del paciente

El token del paciente contiene los siguientes claims:

| Claim | Valor |
|---|---|
| `user_id` | `paciente_id` — ID del paciente (no del médico) |
| `rol_name` | `"paciente"` — siempre este valor |
| `clinica_id` | ID de la clínica a la que pertenece |
| `aplicacion_id` | ID de la app a la que tiene acceso |

**Importante:** El backend extrae del token:
- `c.GetUint("userID")` → `paciente_id`
- `c.GetUint("clinicaID")` → `clinica_id`
- `c.GetUint("aplicacionID")` → `aplicacion_id`

Para todos los endpoints de nutrición, el `pacienteId` en la URL
**debe ser el mismo** que el `user_id` del token.

### 1.3 Contraseña por defecto

Cuando un médico asigna una app a un paciente, se crea automáticamente
un usuario con:
- **username:** número de documento (o teléfono si no tiene documento)
- **password:** `Usuario123`


---

## 2. Modelos de la app (Aplicación)

La `Aplicacion` es la app móvil/web habilitada por la clínica.

```
id, clinica_id,
codigo (ej: "NUTRICION", "PSICOLOGIA"),
nombre, descripcion?,
state (A|I), creado_en
```

---

## 3. Plan de Dieta del Paciente

### 3.1 Listar dietas
```
GET /api/v1/nutricion/pacientes/:pacienteId/dietas
```

### 3.2 Obtener dieta por ID (con menús)
```
GET /api/v1/nutricion/pacientes/:pacienteId/dietas/:dietaId
```

**Modelo `NutricionDietaPaciente`:**
```
id, paciente_id, medico_id,
dieta_catalogo_id? (plantilla base),
nombre, descripcion?, objetivo?, resultado_esperado?,
fecha_inicio (date), duracion_dias (default 7), num_comidas (default 5),
fecha_fin? (date),
calorias_dia_objetivo?, proteinas_g_dia?,
carbohidratos_g_dia?, grasas_g_dia?, fibra_g_dia?,
estado: ACTIVA | COMPLETADA | CANCELADA | PAUSADA,
state (A|I), creado_en, actualizado_en
```

---

## 4. Menú Semanal

### 4.1 Listar menús de una dieta
```
GET /api/v1/nutricion/pacientes/:pacienteId/dietas/:dietaId/menus
```

### 4.2 Obtener menú
```
GET /api/v1/nutricion/pacientes/:pacienteId/menus/:menuId
```

### 4.3 Ver detalles del menú (días y comidas)
```
GET /api/v1/nutricion/pacientes/:pacienteId/menus/:menuId/detalles
```

Retorna la estructura completa:
`Menu → detalles[] (día+tipo_comida) → alimentos[]`

### 4.4 Ver alimentos de un detalle específico
```
GET /api/v1/nutricion/pacientes/:pacienteId/menu-detalles/:detalleId/alimentos
```

**Modelo `NutricionMenu`:**
```
id, dieta_paciente_id, semana_numero, fecha_inicio (date), fecha_fin (date),
nombre?, notas?,
estado: PENDIENTE | ACTIVO | COMPLETADO,
state (A|I), creado_en, actualizado_en,
detalles[] → NutricionMenuDetalle
```

**Modelo `NutricionMenuDetalle`:**
```
id, menu_id, tipo_comida_id, dia_numero (1=Lun…7=Dom),
nombre_comida?, instrucciones?,
calorias_total?, proteinas_g_total?, carbohidratos_g_total?, grasas_g_total?,
state (A|I), creado_en,
alimentos[] → NutricionMenuAlimento
```

**Modelo `NutricionMenuAlimento`:**
```
id, menu_detalle_id, alimento_id,
gramos_asignados, calorias_calc?,
proteinas_g_calc?, carbohidratos_g_calc?, grasas_g_calc?,
observacion?, state (A|I)
```

**Tipos de comida (catálogo fijo):**
```
DES = Desayuno       (07:00)
MMA = Media Mañana   (10:00)
ALM = Almuerzo       (13:00)
MTA = Media Tarde    (16:00)
MER = Merienda/Cena  (19:00)
```

---

## 5. Registros Diarios de Comida

El paciente registra lo que comió cada día. Un registro = una comida del día.

### 5.1 Listar registros
```
GET /api/v1/nutricion/pacientes/:pacienteId/registros-comida
GET /api/v1/nutricion/pacientes/:pacienteId/registros-comida?fecha=2026-03-24
```

### 5.2 Crear registro de comida
```
POST /api/v1/nutricion/pacientes/:pacienteId/registros-comida
```
```json
{
  "fecha": "2026-03-24",          // requerido
  "tipo_comida_id": 1,            // requerido — ID del tipo de comida
  "menu_detalle_id": 10,          // opcional — si cumplió un ítem del plan
  "fuera_de_plan": false,         // true si comió algo no planificado
  "descripcion_libre": "Sandwich de pollo",
  "calorias_consumidas": 450.0,
  "foto_comida": "ruta/foto.jpg",
  "notas": "..."
}
```

### 5.3 Agregar alimento a un registro
```
POST /api/v1/nutricion/pacientes/:pacienteId/registros-comida/:registroId/alimentos
```
```json
{
  "alimento_id": 5,              // opcional — FK al catálogo de alimentos
  "nombre_libre": "Arroz casero", // opcional — si no está en catálogo
  "gramos_consumidos": 200.0     // requerido, > 0
}
```
> Si se envía `alimento_id`, las calorías y macros se calculan automáticamente.

**Modelo `NutricionRegistroComida`:**
```
id, paciente_id, fecha (date), tipo_comida_id,
menu_detalle_id?, fuera_de_plan,
descripcion_libre?, calorias_consumidas?,
proteinas_g?, carbohidratos_g?, grasas_g?,
porcentaje_cumplido?, foto_comida?, notas?,
state (A|I), creado_en
```

**Modelo `NutricionRegistroAlimento`:**
```
id, registro_comida_id,
alimento_id?, nombre_libre?,
gramos_consumidos,
calorias_calc?, proteinas_g_calc?,
carbohidratos_g_calc?, grasas_g_calc?,
state (A|I), creado_en
```

---

## 6. Registros de Ejercicio

### 6.1 Listar registros
```
GET /api/v1/nutricion/pacientes/:pacienteId/registros-ejercicio
```

### 6.2 Crear registro de ejercicio
```
POST /api/v1/nutricion/pacientes/:pacienteId/registros-ejercicio
```
```json
{
  "fecha": "2026-03-24",           // requerido
  "ejercicio_paciente_id": 3,      // opcional — FK al ejercicio asignado por el médico
  "ejercicio_id": 7,               // opcional — FK al catálogo de ejercicios
  "nombre_libre": "Caminata",      // opcional — texto libre si no usa catálogo
  "duracion_min_real": 45,
  "series_real": 3,
  "repeticiones_real": 12,
  "peso_kg_real": 20.0,
  "calorias_quemadas": 300.0,
  "frecuencia_cardiaca_max": 145,
  "nivel_esfuerzo": 7,             // 1–10
  "notas": "..."
}
```

**Modelo `NutricionRegistroEjercicio`:**
```
id, paciente_id, fecha (date),
ejercicio_paciente_id?, ejercicio_id?, nombre_libre?,
duracion_min_real?, series_real?, repeticiones_real?,
peso_kg_real?, calorias_quemadas?,
frecuencia_cardiaca_max?, nivel_esfuerzo? (1-10),
notas?, state (A|I), creado_en
```

---

## 7. Ejercicios Asignados (vista del paciente)

El médico los asigna; el paciente los consulta para saber qué debe hacer.

### 7.1 Listar ejercicios asignados
```
GET /api/v1/nutricion/pacientes/:pacienteId/ejercicios
```

**Modelo `NutricionEjercicioPaciente`:**
```
id, paciente_id, medico_id, dieta_paciente_id?,
ejercicio_id → NutricionEjercicioCatalogo,
dia_numero? (1-7), dia_semana? (texto),
duracion_min?, series?, repeticiones?,
peso_kg?, descanso_seg?,
calorias_estimadas?, instrucciones?,
estado: PENDIENTE | COMPLETADO | SALTADO,
state (A|I)
```

---

## 8. Recordatorio 24 Horas (R24H)

Herramienta clínica donde se registra todo lo consumido en las últimas 24h.
Puede usarse desde la app del paciente o por el médico.

### 8.1 Listar R24H
```
GET /api/v1/nutricion/pacientes/:pacienteId/r24h
```

### 8.2 Crear R24H (encabezado)
```
POST /api/v1/nutricion/pacientes/:pacienteId/r24h
```
```json
{
  "fecha": "2026-03-24",
  "observaciones": "..."
}
```

### 8.3 Ver ítems de un R24H
```
GET /api/v1/nutricion/pacientes/:pacienteId/r24h/:r24hId/items
```

### 8.4 Agregar ítem al R24H
```
POST /api/v1/nutricion/pacientes/:pacienteId/r24h/:r24hId/items
```
```json
{
  "hora_aprox": "07:30",
  "tipo_comida": "Desayuno",      // requerido — texto libre
  "alimento": "Avena con leche", // requerido — texto libre
  "cantidad": "1 taza",
  "calorias_est": 250.0,
  "notas": "..."
}
```

---

## 9. Síntomas

### 9.1 Listar síntomas
```
GET /api/v1/nutricion/pacientes/:pacienteId/sintomas
```

### 9.2 Registrar síntoma
```
POST /api/v1/nutricion/pacientes/:pacienteId/sintomas
```
```json
{
  "fecha": "2026-03-24",                    // requerido
  "tipo": "GASTROINTESTINAL",               // opcional: GASTROINTESTINAL | ENERGETICO | DIGESTIVO | OTRO
  "descripcion": "Náuseas tras el almuerzo", // requerido
  "intensidad": 6,                           // opcional, 1–10
  "alimento_posible": "Leche entera"
}
```

---

## 10. Preferencias Alimentarias

### 10.1 Listar preferencias
```
GET /api/v1/nutricion/pacientes/:pacienteId/preferencias
```

### 10.2 Registrar preferencia
```
POST /api/v1/nutricion/pacientes/:pacienteId/preferencias
```
```json
{
  "alimento_id": 5,              // opcional — si está en el catálogo
  "nombre_libre": "Brócoli",     // opcional — si no está en catálogo
  "tipo": "DISGUSTO",            // requerido: GUSTO | DISGUSTO | INTOLERANCIA | ALERGIA
  "notas": "Produce gases"
}
```

### 10.3 Eliminar preferencia
```
DELETE /api/v1/nutricion/pacientes/:pacienteId/preferencias/:id
```

---

## 11. Progreso Físico

### 11.1 Listar registros de progreso
```
GET /api/v1/nutricion/pacientes/:pacienteId/progreso
```

### 11.2 Agregar registro de progreso
```
POST /api/v1/nutricion/pacientes/:pacienteId/progreso
```
```json
{
  "fecha": "2026-03-24",          // requerido
  "dieta_paciente_id": 2,         // opcional
  "peso_kg": 75.5,
  "altura_cm": 170.0,
  "grasa_corporal_pct": 22.5,
  "masa_muscular_kg": 35.0,
  "cintura_cm": 85.0,
  "cadera_cm": 98.0,
  "pecho_cm": 95.0,
  "brazo_cm": 32.0,
  "muslo_cm": 55.0,
  "hidratacion_litros": 2.5,
  "sueno_horas": 7.5,
  "energia_nivel": 8,             // 1–10
  "pct_cumplimiento_dieta": 85,
  "foto_progreso": "ruta/foto.jpg",
  "notas": "..."
}
```

---

## 12. Gamificación — XP y Logros

### 12.1 Obtener XP del paciente
```
GET /api/v1/nutricion/pacientes/:pacienteId/xp
```
```json
{
  "paciente_id": 11,
  "xp_total": 350,
  "nivel": 4,
  "racha_actual": 5,
  "racha_maxima": 12,
  "ultimo_registro": "2026-03-23"
}
```

### 12.2 Listar logros obtenidos
```
GET /api/v1/nutricion/pacientes/:pacienteId/logros
```
Retorna lista de logros obtenidos por el paciente con datos del catálogo.

**Catálogo de logros (`NutricionLogroCatalogo`):**
```
id, codigo, nombre, descripcion?, icono?, categoria?,
condicion_tipo?, condicion_valor?,
puntos_xp, state (A|I)
```

---

## 13. Catálogos (solo lectura)

```
GET /api/v1/nutricion/alimentos              — lista de alimentos
GET /api/v1/nutricion/alimentos/:id          — detalle de alimento
GET /api/v1/nutricion/ejercicios-catalogo    — lista de ejercicios disponibles
GET /api/v1/nutricion/logros-catalogo        — lista de todos los logros posibles
GET /api/v1/nutricion/dietas-catalogo        — plantillas de dieta
```

**Modelo `NutricionAlimento`:**
```
id, nombre, descripcion?, grupo_id?,
categoria?, gramos_porcion (default 100g),
calorias, proteinas_g, carbohidratos_g, grasas_g,
fibra_g?, azucares_g?, sodio_mg?,
desayuno (bool), almuerzo (bool), media_tarde_mana (bool), merienda (bool),
state (A|I)
```

---

## 14. Cálculo de fórmulas nutricionales (sin persistencia)

```
POST /api/v1/nutricion/formulas
```
```json
{
  "sexo": "M",           // requerido: M | F
  "edad_anos": 30,
  "altura_cm": 170.0,
  "peso_kg": 75.0,
  "cintura_cm": 85.0,
  "cadera_cm": 98.0,
  "factor_actividad": 1.55  // 1.2 sedentario · 1.375 leve · 1.55 moderado · 1.725 activo · 1.9 muy activo
}
```
Retorna IMC, clasificación, ICC, riesgo metabólico, TMB y GET. No guarda nada en BD.

---

## Resumen de endpoints de la app del paciente

| Método | Endpoint | Descripción |
|---|---|---|
| POST | `/pacientes/login` | Login (público) |
| GET | `/nutricion/pacientes/:id/dietas` | Ver mis dietas |
| GET | `/nutricion/pacientes/:id/dietas/:dietaId` | Ver dieta |
| GET | `/nutricion/pacientes/:id/dietas/:dietaId/menus` | Ver menús de dieta |
| GET | `/nutricion/pacientes/:id/menus/:menuId` | Ver menú |
| GET | `/nutricion/pacientes/:id/menus/:menuId/detalles` | Ver plan semana completo |
| GET | `/nutricion/pacientes/:id/menu-detalles/:detalleId/alimentos` | Ver alimentos de comida |
| GET/POST | `/nutricion/pacientes/:id/registros-comida` | Mis registros de comida |
| POST | `/nutricion/pacientes/:id/registros-comida/:rId/alimentos` | Agregar alimento |
| GET/POST | `/nutricion/pacientes/:id/registros-ejercicio` | Mis registros de ejercicio |
| GET | `/nutricion/pacientes/:id/ejercicios` | Ver ejercicios asignados |
| GET/POST | `/nutricion/pacientes/:id/r24h` | Recordatorio 24h |
| GET/POST | `/nutricion/pacientes/:id/r24h/:r24hId/items` | Ítems del R24H |
| GET/POST | `/nutricion/pacientes/:id/sintomas` | Mis síntomas |
| GET/POST/DELETE | `/nutricion/pacientes/:id/preferencias` | Mis preferencias |
| GET/POST | `/nutricion/pacientes/:id/progreso` | Mi progreso físico |
| GET | `/nutricion/pacientes/:id/xp` | Mi XP |
| GET | `/nutricion/pacientes/:id/logros` | Mis logros |
| GET | `/nutricion/alimentos` | Catálogo de alimentos |
| GET | `/nutricion/ejercicios-catalogo` | Catálogo de ejercicios |
| POST | `/nutricion/formulas` | Calcular IMC/TMB/GET |
