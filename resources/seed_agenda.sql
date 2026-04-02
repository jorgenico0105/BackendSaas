-- ============================================================
-- Seed: Módulo Agenda/Citas
-- 1. Catálogos (tipo_citas, estado_citas)
-- 2. Transacciones de menú
-- 3. rol_transaccion
-- 4. Pre-pacientes anónimos
-- 5. Citas de ejemplo (2 con paciente real, 2 anónimas)
-- ============================================================

-- ─── 1. CATÁLOGOS ─────────────────────────────────────────

INSERT INTO tipo_citas (nombre, state) VALUES
    ('Consulta General',    'A'),
    ('Primera Consulta',    'A'),
    ('Control / Seguimiento', 'A'),
    ('Urgencia',            'A'),
    ('Evaluación Nutricional', 'A'),
    ('Sesión Psicológica',  'A'),
    ('Revisión Odontológica', 'A'),
    ('Teleconsulta',        'A')
ON DUPLICATE KEY UPDATE state = 'A';

-- Los códigos deben ser únicos (uniqueIndex en codigo)
INSERT INTO estado_citas (codigo, nombre) VALUES
    ('PE', 'Pendiente'),
    ('CF', 'Confirmada'),
    ('AT', 'Atendida'),
    ('CA', 'Cancelada'),
    ('NA', 'No Asistió')
ON DUPLICATE KEY UPDATE nombre = VALUES(nombre);

-- ─── 2. TRANSACCIONES (menú agenda) ───────────────────────
-- Este módulo es la página principal para todas las profesiones

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Agenda', NULL, 0, NULL, 'ri-calendar-2-line', 'MENU', true, true, 'A');
SET @id_agenda = LAST_INSERT_ID();

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Mis Citas', @id_agenda, 1, '/dashboard/citas', 'ri-calendar-check-line', 'ITEM', true, true, 'A');

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Calendario', @id_agenda, 2, '/dashboard/citas/calendario', 'ri-calendar-line', 'ITEM', true, true, 'A');

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Nueva Cita', @id_agenda, 3, '/dashboard/citas/nueva', 'ri-calendar-add-line', 'ITEM', true, true, 'A');

-- Ítem de sesión por profesión (general=false, solo aparece con la ruta correcta según rol)
-- Estos son navegados dinámicamente desde el dashboard de citas,
-- no aparecen en el menú lateral (visible=false)
INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Sesión Nutricionista',  @id_agenda, 10, '/dashboard/citas/sesion/nutricionista',  'ri-file-list-3-line', 'ITEM', false, true, 'A');

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Sesión Médico',         @id_agenda, 11, '/dashboard/citas/sesion/medico',         'ri-stethoscope-line',  'ITEM', false, true, 'A');

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Sesión Psicólogo',      @id_agenda, 12, '/dashboard/citas/sesion/psicologo',      'ri-mental-health-line','ITEM', false, true, 'A');

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Sesión Odontólogo',     @id_agenda, 13, '/dashboard/citas/sesion/odontologo',     'ri-tooth-line',        'ITEM', false, true, 'A');

-- ─── 3. ROL_TRANSACCION ───────────────────────────────────
-- Todos los roles profesionales acceden al módulo agenda
INSERT INTO rol_transaccion (rol_id, transaccion_id, state)
SELECT r.id, t.id, 'A'
FROM roles r
CROSS JOIN transacciones t
WHERE r.nombre IN ('super_admin', 'admin_clinica', 'medico', 'psicologo', 'nutriologo', 'odontologo', 'recepcionista')
  AND t.nombre IN (
    'Agenda', 'Mis Citas', 'Calendario', 'Nueva Cita',
    'Sesión Nutricionista', 'Sesión Médico', 'Sesión Psicólogo', 'Sesión Odontólogo'
  )
  AND t.state = 'A'
ON DUPLICATE KEY UPDATE state = 'A';

-- ─── 4. PRE-PACIENTES (anónimos / sin registrar) ──────────
-- Estas son personas que escribieron o llamaron pero no están en el sistema como pacientes

-- Asumimos clinica_id = 1 (ajusta según tu BD)
SET @clinica_id = 1;

INSERT INTO pre_pacientes (clinica_id, nombres, apellidos, telefono, correo, origen, notas, state)
VALUES
    (@clinica_id, 'Carlos',  'Mendoza Ríos',    '+593987654321', 'carlos.mendoza@gmail.com',  'MANUAL', 'Paciente nuevo por WhatsApp, consulta nutricional', 'A'),
    (@clinica_id, 'Patricia','Salazar Vega',    '+593976543210', '',                           'MANUAL', 'Llamó por teléfono, primera consulta médica',       'A');

SET @pre_pac1 = LAST_INSERT_ID();
SET @pre_pac2 = @pre_pac1 + 1;

-- ─── 5. CITAS ─────────────────────────────────────────────
-- Requiere: tipo_citas, estado_citas, users (médico), pacientes (existentes), clinica
--
-- ⚠️  Ajusta los IDs de medico y paciente a los reales de tu BD:
--     SET @medico_id   = (SELECT id FROM users WHERE state='A' LIMIT 1);
--     SET @paciente1   = (SELECT id FROM pacientes WHERE state='A' LIMIT 1);
--     SET @paciente2   = (SELECT id FROM pacientes WHERE state='A' LIMIT 1 OFFSET 1);
--
-- O reemplaza directamente con los IDs correctos.

SET @medico_id = (SELECT id FROM users WHERE state = 'A' LIMIT 1);
SET @tipo_control = (SELECT id FROM tipo_citas WHERE nombre = 'Control / Seguimiento' LIMIT 1);
SET @tipo_primera = (SELECT id FROM tipo_citas WHERE nombre = 'Primera Consulta' LIMIT 1);
SET @tipo_nutri   = (SELECT id FROM tipo_citas WHERE nombre = 'Evaluación Nutricional' LIMIT 1);
SET @estado_pe    = (SELECT id FROM estado_citas WHERE codigo = 'PE');
SET @estado_cf    = (SELECT id FROM estado_citas WHERE codigo = 'CF');
SET @estado_at    = (SELECT id FROM estado_citas WHERE codigo = 'AT');

-- ── Cita 1: Paciente real registrado — pasada (atendida)
SET @paciente1 = (SELECT id FROM pacientes WHERE state = 'A' ORDER BY id ASC LIMIT 1);
INSERT INTO citas (fecha, hora, duracion_min, id_medico, id_paciente, id_clinica, tipo_cita_id, estado_cita_id, pre_paciente_id, motivo, state)
VALUES (
    DATE_SUB(CURDATE(), INTERVAL 7 DAY),   -- hace 7 días
    '10:00',
    45,
    @medico_id, @paciente1, @clinica_id,
    @tipo_control, @estado_at,
    NULL,
    'Control mensual de peso y métricas nutricionales',
    'A'
);
SET @cita1_id = LAST_INSERT_ID();

-- Sesión de la cita 1 (ya completada)
INSERT INTO sesiones (cita_id, inicio, fin, resumen, conclusiones, state)
VALUES (
    @cita1_id,
    DATE_SUB(NOW(), INTERVAL 7 DAY),
    DATE_SUB(DATE_ADD(NOW(), INTERVAL 45 MINUTE), INTERVAL 7 DAY),
    'Paciente presenta buena adherencia al plan nutricional. Peso actual: 72 kg (-1.5 desde última consulta).',
    'Continuar con plan actual. Próxima cita en 4 semanas. Aumentar proteína en desayuno.',
    'A'
);

-- ── Cita 2: Paciente real registrado — futura (confirmada)
SET @paciente2 = (SELECT id FROM pacientes WHERE state = 'A' ORDER BY id ASC LIMIT 1 OFFSET 1);
INSERT INTO citas (fecha, hora, duracion_min, id_medico, id_paciente, id_clinica, tipo_cita_id, estado_cita_id, pre_paciente_id, motivo, state)
VALUES (
    DATE_ADD(CURDATE(), INTERVAL 2 DAY),   -- en 2 días
    '11:30',
    60,
    @medico_id, @paciente2, @clinica_id,
    @tipo_nutri, @estado_cf,
    NULL,
    'Primera evaluación nutricional completa con antropometría',
    'A'
);

-- ── Cita 3: Anónimo (pre_paciente) — pendiente
-- Para citas anónimas: id_paciente = 0 (sin FK física en GORM por defecto)
-- El pre_paciente_id es la referencia real al contacto
INSERT INTO citas (fecha, hora, duracion_min, id_medico, id_paciente, id_clinica, tipo_cita_id, estado_cita_id, pre_paciente_id, motivo, state)
VALUES (
    DATE_ADD(CURDATE(), INTERVAL 1 DAY),   -- mañana
    '09:00',
    45,
    @medico_id, 0, @clinica_id,
    @tipo_primera, @estado_pe,
    @pre_pac1,
    'Primera consulta nutricional — contacto por WhatsApp',
    'A'
);

-- ── Cita 4: Anónimo (pre_paciente) — pendiente
INSERT INTO citas (fecha, hora, duracion_min, id_medico, id_paciente, id_clinica, tipo_cita_id, estado_cita_id, pre_paciente_id, motivo, state)
VALUES (
    DATE_ADD(CURDATE(), INTERVAL 3 DAY),   -- en 3 días
    '15:00',
    30,
    @medico_id, 0, @clinica_id,
    @tipo_primera, @estado_pe,
    @pre_pac2,
    'Primera consulta médica general — contacto telefónico',
    'A'
);


-- ─── RESUMEN ──────────────────────────────────────────────
-- Cita 1: paciente real, ATENDIDA hace 7 días, con sesión registrada
-- Cita 2: paciente real, CONFIRMADA en 2 días
-- Cita 3: anónimo (Carlos Mendoza), PENDIENTE mañana
-- Cita 4: anónimo (Patricia Salazar), PENDIENTE en 3 días
-- ============================================================
