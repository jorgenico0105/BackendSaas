# SaasMedico — Esquema de Base de Datos (Contexto para desarrollo)

> **Stack:** NestJS + TypeORM · MySQL 8 · Base de datos: `SaasMedico`
> **Propósito:** SaaS médico multi-clínica. Un usuario (profesional de salud) puede pertenecer a una o varias clínicas. Los pacientes pertenecen a la clínica.

---

## MÓDULO: ADMINISTRACIÓN Y AUTENTICACIÓN

### `Usuarios`
Profesionales de salud que usan el sistema (médicos, nutricionistas, psicólogos, etc.).

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| nombre | VARCHAR(100) | |
| apellidos | VARCHAR(100) | |
| username | VARCHAR(20) | |
| sexo | VARCHAR(10) | |
| codigo_profesion | INT | Código interno de profesión |
| universidad | VARCHAR(50) | |
| celular | VARCHAR(20) | |
| correo | VARCHAR(100) | UNIQUE |
| password | VARCHAR(255) | Hash bcrypt |
| state | VARCHAR(10) | 'A' activo, 'I' inactivo |
| created_at | DATE | |
| created_by | INT | |
| foto | VARCHAR(250) | URL foto de perfil |

---

### `profesiones`
Catálogo de profesiones (Medicina, Nutrición, Odontología, Psicología, etc.).

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| nombre | VARCHAR(100) | UNIQUE |
| descripcion | VARCHAR(255) | |
| state | VARCHAR(10) | |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |

---

### `roles`
Roles dentro de una clínica (Admin, Doctor, Recepcionista, etc.).

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| nombre | VARCHAR(50) | UNIQUE |
| descripcion | VARCHAR(150) | |
| state | VARCHAR(10) | |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |

---

### `clinicas`
Clínica u organización médica. Unidad principal del SaaS.

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| nombre | VARCHAR(150) | |
| ruc | VARCHAR(13) | UNIQUE |
| razon_social | VARCHAR(200) | |
| direccion | VARCHAR(250) | |
| ciudad | VARCHAR(100) | |
| provincia | VARCHAR(100) | |
| pais | VARCHAR(100) | Default 'Ecuador' |
| telefono | VARCHAR(20) | |
| correo | VARCHAR(150) | |
| sitio_web | VARCHAR(150) | |
| representante_legal | VARCHAR(150) | |
| tipo_clinica | VARCHAR(100) | |
| state | VARCHAR(5) | |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |

---

### `sucursales`
Sedes o sucursales de una clínica.

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| id_clinica | INT FK→clinicas | CASCADE DELETE |
| nombre | VARCHAR(150) | |
| codigo | VARCHAR(30) | |
| direccion | VARCHAR(255) | |
| ciudad | VARCHAR(120) | |
| provincia | VARCHAR(120) | |
| telefono | VARCHAR(30) | |
| correo | VARCHAR(150) | |
| state | CHAR(1) | |
| creado_en | TIMESTAMP | |
| actualizado_en | TIMESTAMP | |

---

### `consultorios`
Consultorios físicos dentro de una sucursal.

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| sucursal_id | INT FK→sucursales | CASCADE DELETE |
| nombre | VARCHAR(120) | |
| codigo | VARCHAR(30) | |
| piso | VARCHAR(30) | |
| descripcion | VARCHAR(255) | |
| state | CHAR(1) | |
| creado_en | TIMESTAMP | |

---

### `usuarios_clinicas`
Relación N:M entre usuarios y clínicas. Define en qué clínica/sucursal trabaja cada usuario y con qué rol.

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| usuario_id | INT FK→Usuarios | |
| clinica_id | INT FK→clinicas | |
| sucursal_id | INT FK→sucursales | Nullable |
| rol_id | INT FK→roles | Nullable |
| state | CHAR(1) | |
| creado_en | TIMESTAMP | |

> **UNIQUE:** (usuario_id, clinica_id, sucursal_id)

---

### `usuarios_consultorios`
Relación usuario ↔ consultorio (qué consultorios tiene asignados).

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| usuario_id | INT FK→Usuarios | |
| consultorio_id | INT FK→consultorios | |
| state | CHAR(1) | |
| creado_en | TIMESTAMP | |

---

### `planes_saas`
Planes de suscripción del SaaS.

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| codigo | VARCHAR(30) | UNIQUE |
| nombre | VARCHAR(120) | |
| descripcion | VARCHAR(255) | |
| precio_mensual | DECIMAL(10,2) | |
| precio_anual | DECIMAL(10,2) | |
| max_usuarios | INT | Nullable = ilimitado |
| max_pacientes | INT | Nullable = ilimitado |
| state | CHAR(1) | |
| creado_en | TIMESTAMP | |

---

### `estados_suscripcion`
| código | nombre |
|---|---|
| PRUEBA | Período de prueba |
| ACTIVA | Suscripción activa |
| PAUSADA | Suscripción pausada |
| VENCIDA | Suscripción vencida |
| CANCELADA | Suscripción cancelada |

### `suscripciones`
Suscripción de una clínica o usuario a un plan SaaS.

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| plan_id | INT FK→planes_saas | |
| clinica_id | INT FK→clinicas | Nullable |
| usuario_id | INT FK→Usuarios | Nullable |
| estado | VARCHAR(20) | |
| estado_id | INT FK→estados_suscripcion | |
| inicio | DATE | |
| fin | DATE | Nullable |
| proximo_cobro | DATE | |
| gracia_hasta | DATE | |
| state | CHAR(1) | |
| creado_en / actualizado_en | TIMESTAMP | |

---

### `estilos_clinica`
Logo y colores personalizados de la clínica.

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| clinica_id | INT FK→clinicas | |
| usuario_id | INT FK→Usuarios | |
| nombre_archivo | VARCHAR(255) | |
| url_archivo | VARCHAR(500) | |
| tipo_logo | VARCHAR(30) | 'PRINCIPAL', 'FAVICON', etc. |
| primary_color | VARCHAR(100) | HEX color |
| secondary_color | VARCHAR(100) | |
| third_color | VARCHAR(100) | |
| dark_mode | TINYINT(1) | |
| es_activo | TINYINT(1) | |
| state | CHAR(1) | |
| creado_en | TIMESTAMP | |

---

### `transacciones`
Árbol de permisos/menú del sistema (estructura jerárquica).

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| nombre | VARCHAR(120) | |
| padre_id | INT FK→transacciones | Self-reference, nullable |
| orden | INT | Orden de aparición en menú |
| ruta | VARCHAR(200) | Ruta Angular/frontend |
| icono | VARCHAR(80) | |
| tipo | CHAR(10) | 'MENU', 'ITEM', 'ACTION' |
| visible | TINYINT(1) | |
| general | TINYINT(1) | Si aplica a todos |
| clinica_id | INT FK→clinicas | Nullable |
| usuario_id | INT | Nullable |
| state | CHAR(1) | |
| created_at / updated_at | TIMESTAMP | |

---

### `bloqueos_acceso`
Bloqueos temporales de acceso a una clínica o usuario.

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| clinica_id | INT FK→clinicas | Nullable |
| usuario_id | INT FK→Usuarios | Nullable |
| motivo | VARCHAR(255) | |
| bloqueado_desde | DATETIME | |
| bloqueado_hasta | DATETIME | Nullable = indefinido |
| state | CHAR(1) | |

---

## MÓDULO: PACIENTES

### `Paciente`
Paciente registrado en la clínica.

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| nombres | VARCHAR(100) | |
| apellidos | VARCHAR(100) | |
| sexo | VARCHAR(5) | |
| fecha_nacimiento | DATE | |
| lugar_nacimiento | VARCHAR(100) | |
| nacionalidad | INT | |
| direccion | VARCHAR(200) | |
| telefono | VARCHAR(20) | |
| correo | VARCHAR(100) | Nullable |
| contacto_emergencia | VARCHAR(100) | |
| telefono_emergencia | VARCHAR(100) | |
| tipo_documento | INT | Nullable |
| numero_documento | VARCHAR(100) | Nullable |
| tipo_sangre | VARCHAR(10) | Nullable |
| tipo_paciente | INT | |
| foto | VARCHAR(250) | Nullable |
| state | INT | |
| creado | DATE | |
| created_by | INT | |
| updated_at | DATE | |

---

### `pre_pacientes`
Pacientes registrados antes de confirmarse (leads/prospectos desde formulario web o WhatsApp).

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| id_clinica | INT FK→clinicas | |
| nombres | VARCHAR(120) | |
| apellidos | VARCHAR(120) | |
| telefono | VARCHAR(30) | |
| correo | VARCHAR(150) | Nullable |
| identificacion | VARCHAR(30) | Nullable |
| fecha_nacimiento | DATE | Nullable |
| sexo | CHAR(1) | Nullable |
| origen | VARCHAR(50) | 'WEB', 'WHATSAPP', 'MANUAL' |
| notas | VARCHAR(255) | |
| state | CHAR(1) | |
| creado_en | TIMESTAMP | |

---

## MÓDULO: AGENDA Y CITAS

### `tipo_citas` · `estado_citas`
Catálogos de tipos y estados de cita.

**estados_citas comunes:** PE (Pendiente), CF (Confirmada), AT (Atendida), CA (Cancelada), NA (No asistió)

### `citas`
Cita médica agendada.

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| fecha | DATE | |
| hora | TIME | |
| duracion_min | INT | Default 30 |
| id_medico | INT FK→Usuarios | |
| id_paciente | INT FK→Paciente | |
| id_clinica | INT FK→clinicas | |
| tipo_cita_id | INT FK→tipo_citas | |
| estado_cita_id | INT FK→estado_citas | |
| sucursal_id | INT FK→sucursales | Nullable |
| consultorio_id | INT FK→consultorios | Nullable |
| pre_paciente_id | INT FK→pre_pacientes | Nullable |
| motivo | VARCHAR(255) | |
| url_sesion | VARCHAR(250) | Para teleconsulta |
| notificado | TINYINT(1) | |
| state | CHAR(1) | |
| creado_en / actualizado_en | TIMESTAMP | |

---

### `sesiones`
Sesión clínica que ocurre cuando se atiende una cita. Relación 1:1 con cita.

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| cita_id | INT FK→citas | UNIQUE, CASCADE DELETE |
| inicio | DATETIME | |
| fin | DATETIME | |
| resumen | TEXT | |
| conclusiones | TEXT | |
| state | CHAR(1) | |
| creado_en / actualizado_en | TIMESTAMP | |

---

### `horarios_medico`
Disponibilidad horaria semanal del médico.

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| medico_id | INT FK→Usuarios | |
| clinica_id | INT FK→clinicas | Nullable |
| consultorio_id | INT FK→consultorios | Nullable |
| dia_semana | TINYINT | 0=Lunes … 6=Domingo |
| hora_inicio | TIME | |
| hora_fin | TIME | |
| intervalo_min | INT | Duración de cada slot |
| state | CHAR(1) | |

---

### `bloqueos_agenda`
Bloqueos de horario en la agenda (vacaciones, feriados, etc.).

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| clinica_id | INT FK→clinicas | |
| sucursal_id | INT FK→sucursales | |
| consultorio_id | INT FK→consultorios | |
| medico_id | INT FK→Usuarios | |
| fecha_inicio | DATETIME | |
| fecha_fin | DATETIME | |
| motivo | VARCHAR(255) | |
| tipo_bloqueo | VARCHAR(30) | |
| state | CHAR(1) | |

---

## MÓDULO: COBROS Y PAGOS

### `estados_cobro`
| código | nombre |
|---|---|
| PE | Pendiente |
| PA | Parcial |
| CO | Cobrado |
| AN | Anulado |

### `medios_pago`
EFE (Efectivo), TRA (Transferencia), TAR (Tarjeta), DEP (Depósito), DNA (DeUna), AHO (Ahorita)

### `cobros_sesion`
Cobro generado por una sesión clínica.

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| sesion_id | INT FK→sesiones | UNIQUE |
| id_paciente | INT FK→Paciente | |
| id_medico | INT FK→Usuarios | |
| id_clinica | INT FK→clinicas | |
| monto_cobrar | DECIMAL(10,2) | |
| descuento | DECIMAL(10,2) | |
| recargo | DECIMAL(10,2) | |
| monto_total | DECIMAL(10,2) | |
| estado_cobro_id | INT FK→estados_cobro | |
| observacion | VARCHAR(255) | |
| state | CHAR(1) | |

### `pagos`
Pagos realizados contra un cobro (puede haber varios pagos parciales).

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| cobro_id | INT FK→cobros_sesion | |
| id_paciente | INT FK→Paciente | |
| fecha_pago | DATETIME | |
| monto_pagado | DECIMAL(10,2) | |
| medio_pago_id | INT FK→medios_pago | |
| referencia | VARCHAR(100) | Número comprobante |
| observacion | VARCHAR(255) | |
| state | CHAR(1) | |

### `egresos`
Gastos operativos de la clínica.

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| id_clinica | INT FK→clinicas | |
| tipo_egreso_id | INT FK→tipo_egreso | |
| fecha | DATETIME | |
| monto | DECIMAL(10,2) | |
| descripcion | VARCHAR(255) | |
| proveedor | VARCHAR(150) | |
| referencia | VARCHAR(100) | |
| state | CHAR(1) | |

**tipo_egreso:** MAT (Materiales), SER (Servicios), ALQ (Alquiler), CON (Consultas), OTR (Otros)

---

## MÓDULO: HISTORIA CLÍNICA

### `formularios` → `formulario_preguntas` → `formulario_opciones`
Sistema de formularios dinámicos. Un formulario tiene preguntas; las preguntas de tipo selección tienen opciones.

- `tipo_formulario`: catálogo de tipos (historia clínica, anamnesis, seguimiento, test psicológico, etc.)
- `formulario_preguntas.tipo_respuesta`: TEXT, NUMBER, DATE, SELECT, MULTISELECT, BOOLEAN
- `formulario_preguntas.puntua`: si la pregunta suma puntos para un test

### `historias_clinicas`
Historia clínica de un paciente basada en un formulario dinámico.

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| id_paciente | INT FK→Paciente | |
| id_medico | INT FK→Usuarios | |
| formulario_id | INT FK→formularios | |
| fecha | DATETIME | |
| observacion_general | TEXT | |
| state | CHAR(1) | |

### `historia_respuestas`
Respuestas individuales de una historia clínica.

---

### `paciente_alergias`
Alergias del paciente (FK→alergias_catalogo).

### `paciente_antecedentes`
Antecedentes médicos (tipos: PER, FAM, QUI, PAT, FAR, OTR).

### `paciente_habitos`
Hábitos del paciente: TAB, ALC, SUE, EJE, DIE, CAF.

### `paciente_diagnosticos`
Diagnósticos activos/resueltos del paciente. FK→diagnosticos_catalogo (creado por médico), FK→sesiones, FK→citas.

### `paciente_examenes_resultados`
Archivos de resultados de exámenes del paciente.

---

## MÓDULO: TESTS PSICOLÓGICOS

### `test_reglas` → `test_reglas_detalle`
Reglas de puntuación de un formulario. Define rangos de puntuación y su resultado (ej: 0-10 = Normal, 11-20 = Moderado).

### `tests`
Test aplicado a un paciente. Guarda puntaje total y resultado según reglas.

### `test_respuestas`
Respuestas individuales del test.

---

## MÓDULO: PLANES DE TRATAMIENTO

### `planes_tratamiento` (cabecera)
Plan de tratamiento personalizado para el paciente.

| Columna | Tipo | Notas |
|---|---|---|
| id | INT PK AI | |
| id_paciente | INT FK→Paciente | |
| id_profesional | INT FK→Usuarios | |
| profesion_id | INT FK→profesiones | |
| id_clinica | INT FK→clinicas | |
| titulo | VARCHAR(180) | |
| fecha_inicio | DATE | |
| duracion | INT | En días |
| estado_plan | VARCHAR(20) | ACTIVO, COMPLETADO, CANCELADO |
| costo_total_estimado | DECIMAL(10,2) | |

### `plan_tratamiento_items`
Acciones/procedimientos del plan. FK→acciones_catalogo_tratamientos.
Estado por ítem: PENDIENTE, EN_PROCESO, COMPLETADO.

### `plan_tratamiento_cobros` → `plan_tratamiento_pagos`
Cobros y pagos específicos del plan de tratamiento (independiente de cobros_sesion).

---

## MÓDULO: ODONTOLOGÍA

### `odontogramas`
Odontograma de un paciente. Tiene versiones.

### `odontograma_piezas`
Estado de cada pieza dental (32 piezas posibles).

### `odontograma_piezas_caras`
Estado de cada cara de la pieza (V, L, M, D, O).

### `odontograma_eventos`
Historial de procedimientos realizados en cada pieza. FK→procedimientos_odontologia, FK→citas, FK→sesiones.

---

## MÓDULO: CONSENTIMIENTOS Y PRESCRIPCIONES

### `consentimientos` → `consentimiento_archivos` + `consentimiento_firmas`
Consentimientos informados. Tienen firma digital (base64) del paciente y/o representante.

### `prescripciones` → `prescripcion_items`
Receta médica por sesión. Cada ítem es un medicamento con dosis, frecuencia, duración, vía.

---

## MÓDULO: TAREAS Y SEGUIMIENTO

### `tareas_paciente`
Tareas asignadas al paciente en una sesión. FK→estados_tarea (por profesión), FK→sesiones (asignación y cierre).

### `tarea_progreso`
Registro diario de progreso de una tarea (porcentaje, notas).

### `sesion_observaciones`
Observaciones adicionales registradas durante una sesión.

### `sesion_grabaciones`
Archivos de grabación de sesiones de teleconsulta.

---

## MÓDULO: NOTIFICACIONES (WhatsApp / SMS)

### `tipo_notificacion`
| código | nombre |
|---|---|
| CCON | Confirmación de cita |
| CREM | Recordatorio 24h |
| CRE2 | Recordatorio 2h |
| CREP | Reprogramación |
| CCAN | Cancelación |
| CLNK | Enlace sesión virtual |
| CSEG | Seguimiento post sesión |
| CDEU | Aviso de deuda |
| CPAG | Confirmación de pago |

### `plantillas_notificacion`
Plantillas de mensajes por tipo y canal. Variables dinámicas en el contenido (ej: {{nombre_paciente}}).

### `notificaciones`
Cola de notificaciones. Estado: PENDIENTE, ENVIADA, ERROR, CANCELADA.
FK a citas, sesiones, pagos, pacientes.

---

## MÓDULO: RECURSOS PSICOEDUCATIVOS

### `tematicas_psicoeducativas`
Categorías de recursos (Ansiedad, Depresión, Nutrición, etc.).

### `recursos_psicoeducativos` → `recursos_psicoeducativos_archivos`
Recursos (PDFs, videos, links) organizados por temática. tipo_recurso: FILE, LINK, VIDEO.

### `plantillas_intervencion` → `plantilla_intervencion_archivos`
Plantillas de intervención clínica reutilizables.

---

## MÓDULO: NUTRICIÓN (app móvil paciente)

### Catálogos
- `nutricion_alimentos` — macros por gramos de porción (calorias, proteinas, carbohidratos, grasas, fibra, sodio...)
- `nutricion_dietas_catalogo` — plantillas de dietas con macros diarios estimados y perfil de paciente objetivo
- `nutricion_ejercicios_catalogo` — ejercicios con categoría, grupo muscular, kcal/hora, nivel
- `nutricion_tipo_comida` — Desayuno(1), Media Mañana(2), Almuerzo(3), Media Tarde(4), Merienda/Cena(5)

### Plan del paciente
- `nutricion_dieta_paciente` — cabecera del plan personalizado (duración_dias, fecha_inicio, fecha_fin, resultado_esperado, macros objetivo)
- `nutricion_dieta_detalle` — cada comida de cada día: UNIQUE(dieta_paciente_id, dia_numero, tipo_comida_id)
- `nutricion_dieta_alimentos` — alimentos asignados a cada comida con gramos y macros pre-calculados
- `nutricion_ejercicios_paciente` — ejercicios prescritos (dia_numero, series, repeticiones, duracion_min)

### Registro desde app móvil
- `nutricion_registro_comidas` — el paciente registra qué comió. FK→dieta_detalle_id (cumplimiento del plan) o fuera_de_plan=1
- `nutricion_registro_alimentos` — alimentos individuales de cada registro
- `nutricion_registro_ejercicios` — ejercicio realizado. FK→ejercicio_paciente_id (del plan) o libre

### Progreso y gamificación
- `nutricion_progreso_paciente` — peso, IMC, grasa corporal, medidas (cintura, cadera, etc.), cumplimiento %
- `nutricion_logros_catalogo` — insignias con condicion_tipo (RACHA_DIAS, PESO_META, etc.) y puntos_xp
- `nutricion_logros_paciente` — logros obtenidos. UNIQUE(id_paciente, logro_id)
- `nutricion_paciente_xp` — XP total, nivel, racha_actual, racha_maxima. UNIQUE por paciente

---

## RELACIONES CLAVE (resumen)

```
clinicas ──< sucursales ──< consultorios
clinicas ──< usuarios_clinicas >── Usuarios
Usuarios ──< horarios_medico
Paciente ──< citas >── Usuarios (médico)
citas ──1:1── sesiones
sesiones ──< cobros_sesion ──< pagos
Paciente ──< nutricion_dieta_paciente ──< nutricion_dieta_detalle ──< nutricion_dieta_alimentos
nutricion_dieta_detalle.tipo_comida_id ──> nutricion_tipo_comida
nutricion_dieta_detalle ──> nutricion_registro_comidas (cumplimiento desde app)
Paciente ──< nutricion_logros_paciente >── nutricion_logros_catalogo
Paciente ──1:1── nutricion_paciente_xp
```

---

## CONVENCIONES GENERALES

| Campo | Significado |
|---|---|
| `state = 'A'` | Registro activo |
| `state = 'I'` | Registro inactivo (soft delete) |
| `creado_en` | Timestamp de creación |
| `actualizado_en` | Timestamp última modificación |
| Todos los FK | Con índice explícito para performance |
| Tablas catálogo | Prefijo del módulo (nutricion_, tipo_, estado_) |
| Tablas principales | Nombre en plural o descriptivo del concepto |
