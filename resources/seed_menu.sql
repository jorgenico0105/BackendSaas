-- ============================================================
-- Seed: Menú del sistema (transacciones + rol_transaccion)
-- general = true  → visible para todas las clínicas
-- visible = true  → aparece en el sidebar
-- Ejecutar una sola vez contra la BD de desarrollo/producción
-- ============================================================

-- ─── 1. TRANSACCIONES (menú jerárquico) ───────────────────

-- Módulo: Pacientes
INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Pacientes', NULL, 1, NULL, 'ri-group-line', 'MENU', true, true, 'A');
SET @id_pacientes = LAST_INSERT_ID();

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Lista de Pacientes', @id_pacientes, 1, '/dashboard/pacientes', 'ri-user-line', 'ITEM', true, true, 'A');

-- Módulo: Nutrición
INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Nutrición', NULL, 2, NULL, 'ri-restaurant-line', 'MENU', true, true, 'A');
SET @id_nutricion = LAST_INSERT_ID();

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Dashboard Nutrición', @id_nutricion, 1, '/dashboard/nutricion', 'ri-pie-chart-line', 'ITEM', true, true, 'A');

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Pacientes Nutrición', @id_nutricion, 2, '/dashboard/nutricion/pacientes', 'ri-user-heart-line', 'ITEM', true, true, 'A');

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Alimentos', @id_nutricion, 3, '/dashboard/nutricion/alimentos', 'ri-leaf-line', 'ITEM', true, true, 'A');

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Ejercicios', @id_nutricion, 4, '/dashboard/nutricion/ejercicios', 'ri-run-line', 'ITEM', true, true, 'A');

-- Módulo: Historia / Formularios
INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Historia Clínica', NULL, 3, NULL, 'ri-file-list-3-line', 'MENU', true, true, 'A');
SET @id_historia = LAST_INSERT_ID();

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Formularios', @id_historia, 1, '/dashboard/formularios', 'ri-survey-line', 'ITEM', true, true, 'A');

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Recursos', @id_historia, 2, '/dashboard/recursos', 'ri-folder-2-line', 'ITEM', true, true, 'A');

-- Módulo: Administración
INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Administración', NULL, 4, NULL, 'ri-settings-3-line', 'MENU', true, true, 'A');
SET @id_admin = LAST_INSERT_ID();

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Clínicas', @id_admin, 1, '/dashboard/admin/clinicas', 'ri-hospital-line', 'ITEM', true, true, 'A');

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Roles y Permisos', @id_admin, 2, '/dashboard/admin/roles', 'ri-shield-user-line', 'ITEM', true, true, 'A');

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Profesiones', @id_admin, 3, '/dashboard/admin/profesiones', 'ri-stethoscope-line', 'ITEM', true, true, 'A');

INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Planes SaaS', @id_admin, 4, '/dashboard/admin/planes-saas', 'ri-price-tag-3-line', 'ITEM', true, true, 'A');

-- Item global: Mi Perfil
INSERT INTO transacciones (nombre, padre_id, orden, ruta, icono, tipo, visible, general, state)
VALUES ('Mi Perfil', NULL, 99, '/dashboard/mi-perfil', 'ri-account-circle-line', 'ITEM', true, true, 'A');

-- ─── 2. ROL_TRANSACCION ───────────────────────────────────
-- Asignar todas las transacciones recién insertadas a cada rol.
-- Se usa subquery por nombre para no depender de IDs hardcodeados.

-- super_admin y admin_clinica → acceso completo
INSERT INTO rol_transaccion (rol_id, transaccion_id, state)
SELECT r.id, t.id, 'A'
FROM roles r
CROSS JOIN transacciones t
WHERE r.nombre IN ('super_admin', 'admin_clinica')
  AND t.nombre IN (
    'Pacientes', 'Lista de Pacientes',
    'Nutrición', 'Dashboard Nutrición', 'Pacientes Nutrición', 'Alimentos', 'Ejercicios',
    'Historia Clínica', 'Formularios', 'Recursos',
    'Administración', 'Clínicas', 'Roles y Permisos', 'Profesiones', 'Planes SaaS',
    'Mi Perfil'
  )
  AND t.state = 'A'
ON DUPLICATE KEY UPDATE state = 'A';

-- nutriologo → pacientes + nutrición + formularios (sin administración)
INSERT INTO rol_transaccion (rol_id, transaccion_id, state)
SELECT r.id, t.id, 'A'
FROM roles r
CROSS JOIN transacciones t
WHERE r.nombre = 'nutriologo'
  AND t.nombre IN (
    'Pacientes', 'Lista de Pacientes',
    'Nutrición', 'Dashboard Nutrición', 'Pacientes Nutrición', 'Alimentos', 'Ejercicios',
    'Historia Clínica', 'Formularios', 'Recursos',
    'Mi Perfil'
  )
  AND t.state = 'A'
ON DUPLICATE KEY UPDATE state = 'A';

-- medico, psicologo, odontologo → pacientes + formularios + mi perfil
INSERT INTO rol_transaccion (rol_id, transaccion_id, state)
SELECT r.id, t.id, 'A'
FROM roles r
CROSS JOIN transacciones t
WHERE r.nombre IN ('medico', 'psicologo', 'odontologo')
  AND t.nombre IN (
    'Pacientes', 'Lista de Pacientes',
    'Historia Clínica', 'Formularios', 'Recursos',
    'Mi Perfil'
  )
  AND t.state = 'A'
ON DUPLICATE KEY UPDATE state = 'A';

-- recepcionista → solo pacientes + mi perfil
INSERT INTO rol_transaccion (rol_id, transaccion_id, state)
SELECT r.id, t.id, 'A'
FROM roles r
CROSS JOIN transacciones t
WHERE r.nombre = 'recepcionista'
  AND t.nombre IN (
    'Pacientes', 'Lista de Pacientes',
    'Mi Perfil'
  )
  AND t.state = 'A'
ON DUPLICATE KEY UPDATE state = 'A';
