# SaasMedico — Esquema BD: Nutrición · Tests · Formularios · Historia Clínica

> **Stack:** NestJS + TypeORM · MySQL 8 · DB: `appMedico`
> Este documento es contexto para implementar DTOs y estructuras (structs/entities) en Go o NestJS.

---

## 1. MÓDULO NUTRICIÓN

### `nutricion_tipo_comida`
Catálogo fijo de tiempos de comida del día.
```sql
id          INT PK AI
codigo      VARCHAR(10) UNIQUE   -- 'DES', 'MMA', 'ALM', 'MTA', 'MER'
nombre      VARCHAR(80)          -- 'Desayuno', 'Media Mañana', 'Almuerzo', 'Media Tarde', 'Merienda/Cena'
orden       INT                  -- 1 a 5 (orden cronológico)
hora_ref    TIME                 -- Hora sugerida: 07:30, 10:00, 13:00, 16:00, 19:00
state       CHAR(1) DEFAULT 'A'
```
**Datos fijos:**
| id | codigo | nombre | orden | hora_ref |
|---|---|---|---|---|
| 1 | DES | Desayuno | 1 | 07:30 |
| 2 | MMA | Media Mañana | 2 | 10:00 |
| 3 | ALM | Almuerzo | 3 | 13:00 |
| 4 | MTA | Media Tarde | 4 | 16:00 |
| 5 | MER | Merienda/Cena | 5 | 19:00 |

---

### `nutricion_alimentos`
Catálogo de alimentos con macros por porción de referencia.
```sql
id                  INT PK AI
nombre              VARCHAR(150) NOT NULL
descripcion         VARCHAR(255)
categoria           VARCHAR(80)          -- 'Fruta', 'Verdura', 'Cereal', 'Proteína', 'Lácteo', 'Grasa'
gramos_porcion      DECIMAL(8,2) DEFAULT 100.00  -- gramos a los que corresponden los macros
-- Macros obligatorios
calorias            DECIMAL(8,2) NOT NULL DEFAULT 0   -- kcal
proteinas_g         DECIMAL(8,2) NOT NULL DEFAULT 0
carbohidratos_g     DECIMAL(8,2) NOT NULL DEFAULT 0
grasas_g            DECIMAL(8,2) NOT NULL DEFAULT 0
-- Micronutrientes opcionales
fibra_g             DECIMAL(8,2)
azucares_g          DECIMAL(8,2)
sodio_mg            DECIMAL(8,2)
grasas_saturadas_g  DECIMAL(8,2)
grasas_trans_g      DECIMAL(8,2)
colesterol_mg       DECIMAL(8,2)
-- Control
state               CHAR(1) DEFAULT 'A'
creado_por          INT FK→Usuarios.id ON DELETE SET NULL
creado_en           TIMESTAMP DEFAULT CURRENT_TIMESTAMP
actualizado_en      TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
```

---

### `nutricion_dietas_catalogo`
Plantillas de dietas predefinidas y reutilizables.
```sql
id                    INT PK AI
nombre                VARCHAR(150) NOT NULL
descripcion           TEXT
tipo_paciente_perfil  VARCHAR(255)   -- 'Adulto sedentario', 'Deportista', 'Diabético', 'Embarazada'
objetivo              VARCHAR(100)   -- 'Pérdida de peso', 'Ganancia muscular', 'Mantenimiento'
-- Macros diarios objetivo
calorias_dia          DECIMAL(8,2)
proteinas_g_dia       DECIMAL(8,2)
carbohidratos_g_dia   DECIMAL(8,2)
grasas_g_dia          DECIMAL(8,2)
fibra_g_dia           DECIMAL(8,2)
-- Control
state                 CHAR(1) DEFAULT 'A'
creado_por            INT FK→Usuarios.id ON DELETE SET NULL
creado_en             TIMESTAMP
actualizado_en        TIMESTAMP
```

---

### `nutricion_dieta_paciente` ← CABECERA del plan
Plan de dieta personalizado asignado a un paciente concreto.
```sql
id                      INT PK AI
id_paciente             INT NOT NULL FK→Paciente.id CASCADE DELETE
id_medico               INT NOT NULL FK→Usuarios.id          -- nutricionista responsable
dieta_catalogo_id       INT FK→nutricion_dietas_catalogo.id SET NULL  -- base opcional
nombre                  VARCHAR(150) NOT NULL
descripcion             TEXT
objetivo                VARCHAR(150)   -- 'Pérdida de peso', 'Mantenimiento', 'Volumen'
resultado_esperado      TEXT           -- descripción del resultado esperado al finalizar
-- Período
fecha_inicio            DATE NOT NULL DEFAULT CURRENT_DATE
duracion_dias           INT NOT NULL DEFAULT 7
fecha_fin               DATE                                  -- calculado o manual
-- Macros diarios objetivo de este plan específico
calorias_dia_objetivo   DECIMAL(8,2)
proteinas_g_dia         DECIMAL(8,2)
carbohidratos_g_dia     DECIMAL(8,2)
grasas_g_dia            DECIMAL(8,2)
fibra_g_dia             DECIMAL(8,2)
-- Estado
estado                  VARCHAR(20) DEFAULT 'ACTIVA'  -- 'ACTIVA','COMPLETADA','CANCELADA','PAUSADA'
state                   CHAR(1) DEFAULT 'A'
creado_en               TIMESTAMP
actualizado_en          TIMESTAMP
```
**Índices:** idx en id_paciente, id_medico, estado, (fecha_inicio, fecha_fin)

---

### `nutricion_dieta_detalle`
Cada celda del plan: "Día N – Tipo de comida X".
Un registro = Día 2 – Almuerzo del plan del paciente.
```sql
id                      INT PK AI
dieta_paciente_id       INT NOT NULL FK→nutricion_dieta_paciente.id CASCADE DELETE
tipo_comida_id          INT NOT NULL FK→nutricion_tipo_comida.id
dia_numero              TINYINT NOT NULL        -- 1 a 7 (o más si duracion_dias > 7)
nombre_comida           VARCHAR(150)            -- nombre libre opcional ("Almuerzo mediterráneo")
instrucciones           TEXT                    -- notas de preparación generales
-- Macros totales calculados (suma de alimentos asignados)
calorias_total          DECIMAL(8,2)
proteinas_g_total       DECIMAL(8,2)
carbohidratos_g_total   DECIMAL(8,2)
grasas_g_total          DECIMAL(8,2)
state                   CHAR(1) DEFAULT 'A'
creado_en               TIMESTAMP
actualizado_en          TIMESTAMP

UNIQUE KEY (dieta_paciente_id, dia_numero, tipo_comida_id)
```

---

### `nutricion_dieta_alimentos`
Alimentos asignados a cada comida del plan con su cantidad en gramos.
```sql
id                    INT PK AI
dieta_detalle_id      INT NOT NULL FK→nutricion_dieta_detalle.id CASCADE DELETE
alimento_id           INT NOT NULL FK→nutricion_alimentos.id
gramos_asignados      DECIMAL(8,2) NOT NULL
-- Macros pre-calculados para esta porción: (gramos_asignados / gramos_porcion) * macro
calorias_calc         DECIMAL(8,2)
proteinas_g_calc      DECIMAL(8,2)
carbohidratos_g_calc  DECIMAL(8,2)
grasas_g_calc         DECIMAL(8,2)
observacion           VARCHAR(255)   -- 'cocido', 'crudo', 'sin piel'
state                 CHAR(1) DEFAULT 'A'
creado_en             TIMESTAMP
```

---

### `nutricion_ejercicios_catalogo`
Catálogo de ejercicios disponibles.
```sql
id                INT PK AI
nombre            VARCHAR(150) NOT NULL
descripcion       TEXT
categoria         VARCHAR(80)    -- 'Cardio', 'Fuerza', 'Flexibilidad', 'HIIT'
grupo_muscular    VARCHAR(120)   -- 'Piernas', 'Espalda', 'Pecho', 'Full Body'
calorias_por_hora DECIMAL(8,2)  -- kcal/hora estimado
unidad_medida     VARCHAR(30) DEFAULT 'minutos'  -- 'minutos','repeticiones','series','km'
nivel             VARCHAR(20)    -- 'Principiante', 'Intermedio', 'Avanzado'
url_referencia    VARCHAR(500)   -- video o imagen de referencia
state             CHAR(1) DEFAULT 'A'
creado_por        INT FK→Usuarios.id SET NULL
creado_en         TIMESTAMP
```

---

### `nutricion_ejercicios_paciente`
Ejercicios prescritos por el médico al paciente (plan de ejercicios).
```sql
id                    INT PK AI
id_paciente           INT NOT NULL FK→Paciente.id CASCADE DELETE
id_medico             INT NOT NULL FK→Usuarios.id
dieta_paciente_id     INT FK→nutricion_dieta_paciente.id SET NULL  -- asociado a un plan (opcional)
ejercicio_id          INT NOT NULL FK→nutricion_ejercicios_catalogo.id
dia_numero            TINYINT        -- día del plan (1-7); NULL = todos los días
dia_semana            VARCHAR(15)    -- 'Lunes', 'Martes'... (alternativa a dia_numero)
duracion_min          INT            -- duración prescrita en minutos
series                INT
repeticiones          INT
peso_kg               DECIMAL(6,2)   -- peso asignado si aplica
descanso_seg          INT            -- segundos de descanso entre series
calorias_estimadas    DECIMAL(8,2)
instrucciones         TEXT
estado                VARCHAR(20) DEFAULT 'PENDIENTE'  -- 'PENDIENTE','COMPLETADO','SALTADO'
state                 CHAR(1) DEFAULT 'A'
creado_en             TIMESTAMP
actualizado_en        TIMESTAMP
```

---

### `nutricion_registro_comidas`
El paciente registra desde la app qué comió cada día.
```sql
id                    INT PK AI
id_paciente           INT NOT NULL FK→Paciente.id CASCADE DELETE
fecha                 DATE NOT NULL
tipo_comida_id        INT NOT NULL FK→nutricion_tipo_comida.id
dieta_detalle_id      INT FK→nutricion_dieta_detalle.id SET NULL  -- comida del plan cumplida
fuera_de_plan         TINYINT(1) DEFAULT 0   -- 1 = comida no planificada
descripcion_libre     VARCHAR(255)            -- si es fuera del plan
-- Totales consumidos (calculados o estimados)
calorias_consumidas   DECIMAL(8,2)
proteinas_g           DECIMAL(8,2)
carbohidratos_g       DECIMAL(8,2)
grasas_g              DECIMAL(8,2)
porcentaje_cumplido   INT             -- 0-100, qué % de la comida del plan cumplió
foto_comida           VARCHAR(500)    -- URL foto tomada por el paciente
notas                 VARCHAR(255)
state                 CHAR(1) DEFAULT 'A'
creado_en             TIMESTAMP

INDEX (id_paciente, fecha)
```

---

### `nutricion_registro_alimentos`
Alimentos individuales de cada registro de comida.
```sql
id                    INT PK AI
registro_comida_id    INT NOT NULL FK→nutricion_registro_comidas.id CASCADE DELETE
alimento_id           INT FK→nutricion_alimentos.id SET NULL   -- del catálogo (opcional)
nombre_libre          VARCHAR(150)    -- si no está en catálogo
gramos_consumidos     DECIMAL(8,2) NOT NULL
calorias_calc         DECIMAL(8,2)
proteinas_g_calc      DECIMAL(8,2)
carbohidratos_g_calc  DECIMAL(8,2)
grasas_g_calc         DECIMAL(8,2)
state                 CHAR(1) DEFAULT 'A'
creado_en             TIMESTAMP
```

---

### `nutricion_registro_ejercicios`
El paciente registra ejercicios realizados desde la app.
```sql
id                      INT PK AI
id_paciente             INT NOT NULL FK→Paciente.id CASCADE DELETE
fecha                   DATE NOT NULL
ejercicio_paciente_id   INT FK→nutricion_ejercicios_paciente.id SET NULL  -- del plan prescrito
ejercicio_id            INT FK→nutricion_ejercicios_catalogo.id SET NULL  -- libre del catálogo
nombre_libre            VARCHAR(150)   -- si no está en catálogo
-- Lo que realmente realizó
duracion_min_real       INT
series_real             INT
repeticiones_real       INT
peso_kg_real            DECIMAL(6,2)
calorias_quemadas       DECIMAL(8,2)
frecuencia_cardiaca_max INT            -- ppm
nivel_esfuerzo          TINYINT        -- escala 1-10
notas                   VARCHAR(255)
state                   CHAR(1) DEFAULT 'A'
creado_en               TIMESTAMP

INDEX (id_paciente, fecha)
```

---

### `nutricion_progreso_paciente`
Registro periódico de métricas corporales y cumplimiento del plan.
```sql
id                          INT PK AI
id_paciente                 INT NOT NULL FK→Paciente.id CASCADE DELETE
id_medico                   INT FK→Usuarios.id SET NULL   -- si lo registra el nutricionista
dieta_paciente_id           INT FK→nutricion_dieta_paciente.id SET NULL
fecha                       DATE NOT NULL
-- Métricas corporales
peso_kg                     DECIMAL(6,2)
altura_cm                   DECIMAL(6,2)
imc                         DECIMAL(5,2)   -- calculado: peso/(altura_m^2)
grasa_corporal_pct          DECIMAL(5,2)
masa_muscular_kg            DECIMAL(6,2)
cintura_cm                  DECIMAL(6,2)
cadera_cm                   DECIMAL(6,2)
pecho_cm                    DECIMAL(6,2)
brazo_cm                    DECIMAL(6,2)
muslo_cm                    DECIMAL(6,2)
-- Cumplimiento
calorias_consumidas_dia     DECIMAL(8,2)
pct_cumplimiento_dieta      INT            -- 0-100%
pct_cumplimiento_ejercicio  INT            -- 0-100%
-- Bienestar subjetivo
energia_nivel               TINYINT        -- escala 1-10
sueno_horas                 DECIMAL(4,2)
hidratacion_litros          DECIMAL(4,2)
notas                       TEXT
foto_progreso               VARCHAR(500)
state                       CHAR(1) DEFAULT 'A'
creado_en                   TIMESTAMP

INDEX (id_paciente, fecha)
```

---

### `nutricion_logros_catalogo`
Catálogo de logros/insignias disponibles en el sistema.
```sql
id               INT PK AI
codigo           VARCHAR(30) UNIQUE NOT NULL
nombre           VARCHAR(120) NOT NULL
descripcion      VARCHAR(255)
icono            VARCHAR(100)    -- nombre del icono o URL
categoria        VARCHAR(50)     -- 'Dieta', 'Ejercicio', 'Progreso', 'Racha', 'Hito', 'Hábito'
condicion_tipo   VARCHAR(50)     -- 'RACHA_DIAS','PESO_META','DIAS_REGISTRADOS','EJERCICIOS_TOTAL','PLANES_COMPLETADOS','HIDRATACION_DIAS'
condicion_valor  INT             -- valor numérico umbral (ej: 7 para RACHA_DIAS=7)
puntos_xp        INT DEFAULT 0
state            CHAR(1) DEFAULT 'A'
creado_en        TIMESTAMP
```
**Datos iniciales:**
| codigo | nombre | condicion_tipo | condicion_valor | puntos_xp |
|---|---|---|---|---|
| PRIMER_DIA | Primer Paso | DIAS_REGISTRADOS | 1 | 50 |
| RACHA_3 | En Racha | RACHA_DIAS | 3 | 100 |
| RACHA_7 | Una Semana Perfecta | RACHA_DIAS | 7 | 250 |
| RACHA_30 | Mes Dedicado | RACHA_DIAS | 30 | 750 |
| META_PESO | Meta Alcanzada | PESO_META | 1 | 1000 |
| PRIMER_EJERCICIO | Primer Esfuerzo | EJERCICIOS_TOTAL | 1 | 50 |
| EJ_10 | En Movimiento | EJERCICIOS_TOTAL | 10 | 200 |
| EJ_50 | Atleta | EJERCICIOS_TOTAL | 50 | 500 |
| AGUA_7 | Hidratado | HIDRATACION_DIAS | 7 | 150 |
| PLAN_COMPLETO | Plan Completado | PLANES_COMPLETADOS | 1 | 500 |

---

### `nutricion_logros_paciente`
Logros obtenidos por cada paciente. Un logro se otorga una sola vez.
```sql
id               INT PK AI
id_paciente      INT NOT NULL FK→Paciente.id CASCADE DELETE
logro_id         INT NOT NULL FK→nutricion_logros_catalogo.id
fecha_obtenido   DATETIME DEFAULT CURRENT_TIMESTAMP
puntos_xp        INT DEFAULT 0   -- copia del valor al momento de obtenerlo
notas            VARCHAR(255)
state            CHAR(1) DEFAULT 'A'

UNIQUE KEY (id_paciente, logro_id)
```

---

### `nutricion_paciente_xp`
XP total, nivel y racha de actividad del paciente. Un registro por paciente.
```sql
id               INT PK AI
id_paciente      INT NOT NULL FK→Paciente.id CASCADE DELETE  UNIQUE
xp_total         INT DEFAULT 0
nivel            INT DEFAULT 1
racha_actual     INT DEFAULT 0   -- días consecutivos con actividad registrada
racha_maxima     INT DEFAULT 0
ultimo_registro  DATE            -- último día con actividad
state            CHAR(1) DEFAULT 'A'
actualizado_en   TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
```

---

## 2. MÓDULO FORMULARIOS DINÁMICOS

### `tipo_formulario`
Catálogo de tipos de formulario.
```sql
id          INT PK AI
codigo      CHAR(3) UNIQUE NOT NULL   -- ej: 'HCL', 'ANM', 'SEG', 'TST'
nombre      VARCHAR(100) NOT NULL      -- 'Historia Clínica', 'Anamnesis', 'Seguimiento', 'Test'
descripcion VARCHAR(255)
state       CHAR(1) DEFAULT 'A'
created_at  TIMESTAMP
created_by  INT
```

---

### `formularios`
Formulario creado por un profesional para una profesión.
```sql
id                  INT PK AI
nombre              VARCHAR(150) NOT NULL
descripcion         VARCHAR(255)
profesion_id        INT FK→profesiones.id SET NULL   -- a qué profesión aplica
usuario_id          INT NOT NULL FK→Usuarios.id       -- quién lo creó
tipo_formulario_id  INT NOT NULL FK→tipo_formulario.id
state               CHAR(1) DEFAULT 'A'
creado_en           TIMESTAMP
```

---

### `formulario_preguntas`
Preguntas de un formulario.
```sql
id              INT PK AI
formulario_id   INT NOT NULL FK→formularios.id CASCADE DELETE
pregunta        VARCHAR(255) NOT NULL
tipo_respuesta  VARCHAR(30) NOT NULL
                -- Valores posibles: 'TEXT', 'NUMBER', 'DATE', 'SELECT', 'MULTISELECT', 'BOOLEAN'
obligatorio     TINYINT(1) DEFAULT 0
orden           INT DEFAULT 0
state           CHAR(1) DEFAULT 'A'
puntua          TINYINT(1) DEFAULT 0    -- si esta pregunta suma puntos en un test
peso            DECIMAL(10,2) DEFAULT 1 -- factor multiplicador del puntaje
min_val         DECIMAL(10,2)           -- valor mínimo válido (para NUMBER)
max_val         DECIMAL(10,2)           -- valor máximo válido (para NUMBER)
permite_multi   TINYINT(1) DEFAULT 0    -- permite selección múltiple

INDEX (formulario_id, orden)
```

---

### `formulario_opciones`
Opciones de preguntas tipo SELECT o MULTISELECT.
```sql
id           INT PK AI
pregunta_id  INT NOT NULL FK→formulario_preguntas.id CASCADE DELETE
valor        VARCHAR(100) NOT NULL   -- valor interno
etiqueta     VARCHAR(150) NOT NULL   -- texto visible al usuario
orden        INT DEFAULT 0
puntos       DECIMAL(10,2) DEFAULT 0  -- puntos que aporta esta opción al puntaje total
state        CHAR(1) DEFAULT 'A'
```

---

## 3. MÓDULO TESTS PSICOLÓGICOS

### `test_reglas`
Conjunto de reglas de puntuación para un formulario. Versión permite evolucionar las reglas.
```sql
id            INT PK AI
formulario_id INT NOT NULL FK→formularios.id CASCADE DELETE
version       INT DEFAULT 1
nombre        VARCHAR(150) NOT NULL
descripcion   VARCHAR(255)
state         CHAR(1) DEFAULT 'A'
creado_en     TIMESTAMP

UNIQUE KEY (formulario_id, version)
```

---

### `test_reglas_detalle`
Rangos de puntuación y su resultado (ej: 0-10 = Normal, 11-20 = Moderado, 21+ = Severo).
```sql
id         INT PK AI
regla_id   INT NOT NULL FK→test_reglas.id CASCADE DELETE
min_val    DECIMAL(10,2) NOT NULL
max_val    DECIMAL(10,2) NOT NULL
resultado  VARCHAR(150) NOT NULL    -- ej: 'Normal', 'Leve', 'Moderado', 'Severo'
mensaje    VARCHAR(255)             -- mensaje explicativo para el paciente
orden      INT DEFAULT 0
state      CHAR(1) DEFAULT 'A'

INDEX (regla_id, orden)
```

---

### `tests`
Test aplicado a un paciente. Contiene el puntaje total y el resultado según las reglas.
```sql
id            INT PK AI
id_paciente   INT NOT NULL FK→Paciente.id CASCADE DELETE
id_medico     INT NOT NULL FK→Usuarios.id
formulario_id INT NOT NULL FK→formularios.id
regla_id      INT NOT NULL FK→test_reglas.id
fecha         DATETIME DEFAULT CURRENT_TIMESTAMP
puntaje_total DECIMAL(10,2)       -- suma calculada de respuestas que puntúan
resultado     VARCHAR(150)        -- resultado según test_reglas_detalle (ej: 'Moderado')
observacion   TEXT
state         CHAR(1) DEFAULT 'A'
creado_en     TIMESTAMP

INDEX (id_paciente, fecha)
```

---

### `test_respuestas`
Respuestas individuales de un test.
```sql
id                INT PK AI
test_id           INT NOT NULL FK→tests.id CASCADE DELETE
pregunta_id       INT NOT NULL FK→formulario_preguntas.id
opcion_id         INT FK→formulario_opciones.id SET NULL  -- si es SELECT/MULTISELECT
respuesta_texto   TEXT          -- si tipo_respuesta = TEXT
respuesta_numero  DECIMAL(10,2) -- si tipo_respuesta = NUMBER
creado_en         TIMESTAMP

INDEX (test_id)
```

---

### `test_archivos`
Archivos adjuntos a un test (imágenes, PDFs, resultados).
```sql
id             INT PK AI
test_id        INT NOT NULL FK→tests.id CASCADE DELETE
nombre_archivo VARCHAR(255) NOT NULL
tipo_archivo   VARCHAR(150)
state          CHAR(1) DEFAULT 'A'
creado_en      TIMESTAMP
```

---

### `sesion_tests`
Relación N:M entre sesiones clínicas y tests (tests aplicados en una sesión).
```sql
id         INT PK AI
sesion_id  INT NOT NULL FK→sesiones.id CASCADE DELETE
test_id    INT NOT NULL FK→tests.id CASCADE DELETE
state      CHAR(1) DEFAULT 'A'
creado_en  TIMESTAMP

UNIQUE KEY (sesion_id, test_id)
```

---

## 4. MÓDULO HISTORIA CLÍNICA DEL PACIENTE

### `historias_clinicas`
Historia clínica de un paciente basada en un formulario dinámico.
```sql
id                  INT PK AI
id_paciente         INT NOT NULL FK→Paciente.id CASCADE DELETE
id_medico           INT NOT NULL FK→Usuarios.id
formulario_id       INT NOT NULL FK→formularios.id
fecha               DATETIME DEFAULT CURRENT_TIMESTAMP
observacion_general TEXT
state               CHAR(1) DEFAULT 'A'
```

---

### `historia_respuestas`
Respuestas a cada pregunta del formulario en esa historia clínica.
```sql
id               INT PK AI
historia_id      INT NOT NULL FK→historias_clinicas.id CASCADE DELETE
pregunta_id      INT NOT NULL FK→formulario_preguntas.id
respuesta_texto  TEXT
respuesta_numero DECIMAL(10,2)
respuesta_fecha  DATE
creado_en        TIMESTAMP
```

---

### `paciente_alergias`
Alergias registradas del paciente.
```sql
id               INT PK AI
id_paciente      INT NOT NULL FK→Paciente.id CASCADE DELETE
alergia_id       INT NOT NULL FK→alergias_catalogo.id
severidad        VARCHAR(50)    -- 'Leve', 'Moderada', 'Severa', 'Anafiláctica'
reaccion         VARCHAR(255)
observacion      VARCHAR(255)
fecha_registro   DATETIME DEFAULT CURRENT_TIMESTAMP
id_medico        INT FK→Usuarios.id SET NULL
state            CHAR(1) DEFAULT 'A'

UNIQUE KEY (id_paciente, alergia_id)
```

---

### `paciente_antecedentes`
Antecedentes médicos del paciente.
```sql
id                  INT PK AI
id_paciente         INT NOT NULL FK→Paciente.id CASCADE DELETE
tipo_antecedente_id INT NOT NULL FK→tipos_antecedente.id
descripcion         TEXT NOT NULL
fecha_registro      DATETIME DEFAULT CURRENT_TIMESTAMP
id_medico           INT FK→Usuarios.id SET NULL
state               CHAR(1) DEFAULT 'A'
```
**tipos_antecedente:** PER (Personal), FAM (Familiar), QUI (Quirúrgico), PAT (Patológico), FAR (Farmacológico), OTR (Otro)

---

### `paciente_habitos`
Hábitos del paciente (uno por tipo).
```sql
id              INT PK AI
id_paciente     INT NOT NULL FK→Paciente.id CASCADE DELETE
habito_id       INT NOT NULL FK→habitos_catalogo.id
valor           VARCHAR(120)   -- ej: '10 cigarrillos/día', 'Social'
frecuencia      VARCHAR(120)   -- ej: 'Diario', 'Semanal', '3 veces/semana'
observacion     VARCHAR(255)
fecha_registro  DATETIME DEFAULT CURRENT_TIMESTAMP
id_medico       INT FK→Usuarios.id SET NULL
state           CHAR(1) DEFAULT 'A'

UNIQUE KEY (id_paciente, habito_id)
```
**habitos_catalogo:** TAB (Tabaco), ALC (Alcohol), SUE (Sueño), EJE (Ejercicio), DIE (Dieta), CAF (Cafeína)

---

### `paciente_diagnosticos`
Diagnósticos activos o resueltos del paciente.
```sql
id                  INT PK AI
id_paciente         INT NOT NULL FK→Paciente.id CASCADE DELETE
diagnostico_id      INT NOT NULL FK→diagnosticos_catalogo.id
id_medico           INT NOT NULL FK→Usuarios.id
sesion_id           INT FK→sesiones.id SET NULL
cita_id             INT FK→citas.id SET NULL
estado_clinico      VARCHAR(30) DEFAULT 'ACTIVO'   -- 'ACTIVO', 'RESUELTO', 'CRONICO', 'EN_SEGUIMIENTO'
fecha_diagnostico   DATE DEFAULT CURRENT_DATE
fecha_resolucion    DATE
observaciones       TEXT
state               CHAR(1) DEFAULT 'A'
creado_en           TIMESTAMP
actualizado_en      TIMESTAMP
```

---

### `paciente_examenes_resultados`
Archivos de resultados de exámenes del paciente.
```sql
id               INT PK AI
id_paciente      INT NOT NULL FK→Paciente.id CASCADE DELETE
id_tipo_examen   INT NOT NULL FK→tipo_examen.id
nombre_archivo   VARCHAR(255) NOT NULL   -- nombre del archivo subido
fecha_examen     DATE
creado_en        TIMESTAMP
```

---

### `paciente_imagenes`
Imágenes clínicas del paciente (radiografías, fotos, etc.).
```sql
id               INT PK AI
id_paciente      INT NOT NULL FK→Paciente.id CASCADE DELETE
nombre_archivo   VARCHAR(255) NOT NULL
creado_en        TIMESTAMP
```

---

### `paciente_certificados`
Certificados médicos generados para el paciente.
```sql
id               INT PK AI
id_paciente      INT NOT NULL FK→Paciente.id CASCADE DELETE
nombre_archivo   VARCHAR(255) NOT NULL
creado_en        TIMESTAMP
```

---

## RELACIONES ENTRE MÓDULOS

```
formularios ──< formulario_preguntas ──< formulario_opciones
formularios ──< test_reglas ──< test_reglas_detalle

historias_clinicas (usa formulario)
    └─< historia_respuestas (responde formulario_preguntas)

tests (usa formulario + test_reglas)
    └─< test_respuestas (responde formulario_preguntas, opcionalmente formulario_opciones)
    └─< test_archivos
sesiones >──< tests  (via sesion_tests)

Paciente ──< paciente_alergias >── alergias_catalogo
Paciente ──< paciente_antecedentes >── tipos_antecedente
Paciente ──< paciente_habitos >── habitos_catalogo
Paciente ──< paciente_diagnosticos >── diagnosticos_catalogo

Paciente ──< nutricion_dieta_paciente (cabecera plan)
    └─< nutricion_dieta_detalle [dia_numero + tipo_comida_id  UNIQUE]
        └─< nutricion_dieta_alimentos >── nutricion_alimentos

Paciente ──< nutricion_ejercicios_paciente >── nutricion_ejercicios_catalogo
Paciente ──< nutricion_registro_comidas ──> nutricion_dieta_detalle (cumplimiento)
    └─< nutricion_registro_alimentos
Paciente ──< nutricion_registro_ejercicios ──> nutricion_ejercicios_paciente
Paciente ──< nutricion_progreso_paciente
Paciente ──< nutricion_logros_paciente >── nutricion_logros_catalogo
Paciente ──1:1── nutricion_paciente_xp
```

---

## NOTAS PARA IMPLEMENTACIÓN (NestJS + TypeORM)

- **Soft delete:** todos usan `state CHAR(1)` ('A'=activo, 'I'=inactivo), no hay DELETE físico
- **Timestamps:** `creado_en` = created timestamp, `actualizado_en` = UpdateDateColumn
- **Macros calculados:** los campos `*_calc` en `nutricion_dieta_alimentos` y `nutricion_registro_alimentos` se calculan al guardar: `(gramos / gramos_porcion) * macro_base`
- **IMC:** calculado en backend o en el DTO, no se confía en el valor de la BD
- **XP y rachas:** la tabla `nutricion_paciente_xp` se actualiza cada vez que el paciente registra actividad en `nutricion_registro_comidas` o `nutricion_registro_ejercicios`
- **Logros:** se evalúan en un service/event listener después de cada actualización de XP
- **Unique constraints críticos:**
  - `nutricion_dieta_detalle`: (dieta_paciente_id, dia_numero, tipo_comida_id)
  - `nutricion_logros_paciente`: (id_paciente, logro_id)
  - `nutricion_paciente_xp`: (id_paciente)
  - `paciente_alergias`: (id_paciente, alergia_id)
  - `paciente_habitos`: (id_paciente, habito_id)
  - `sesion_tests`: (sesion_id, test_id)
