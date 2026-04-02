-- ============================================================
-- Migración 002 — Reestructuración detalle de menú nutricional
-- Fecha: 2026-03-19
--
-- Cambios:
--   • Elimina nutricion_dieta_detalle  (reemplazada por nutricion_menu_detalle)
--   • Elimina nutricion_dieta_alimentos (reemplazada por nutricion_menu_alimentos)
--   • Crea nutricion_menu_detalle       (sin dieta_paciente_id redundante)
--   • Crea nutricion_menu_alimentos     (FK → menu_detalle_id)
--   • Renombra columna en nutricion_registro_comidas:
--       dieta_detalle_id → menu_detalle_id
--
-- EJECUTAR en este orden. Si las tablas no existen aún, salta al bloque CREATE.
-- ============================================================

-- ─── 1. Eliminar tablas viejas (si existen) ────────────────────────────────
DROP TABLE IF EXISTS `nutricion_dieta_alimentos`;
DROP TABLE IF EXISTS `nutricion_dieta_detalle`;

-- ─── 2. Crear nutricion_menu_detalle ──────────────────────────────────────
-- Jerarquía: DietaPaciente → Menu (semana) → MenuDetalle (día+comida) → MenuAlimento
CREATE TABLE IF NOT EXISTS `nutricion_menu_detalle` (
  `id`                    BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `menu_id`               BIGINT UNSIGNED NOT NULL,
  `tipo_comida_id`        BIGINT UNSIGNED NOT NULL,
  `dia_numero`            TINYINT         NOT NULL COMMENT '1=Lun … 7=Dom',
  `nombre_comida`         VARCHAR(150)    DEFAULT NULL,
  `instrucciones`         TEXT            DEFAULT NULL,
  `calorias_total`        DECIMAL(8,2)    DEFAULT NULL,
  `proteinas_g_total`     DECIMAL(8,2)    DEFAULT NULL,
  `carbohidratos_g_total` DECIMAL(8,2)    DEFAULT NULL,
  `grasas_g_total`        DECIMAL(8,2)    DEFAULT NULL,
  `state`                 CHAR(1)         NOT NULL DEFAULT 'A',
  `creado_en`             DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `actualizado_en`        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `udx_menu_dia_comida` (`menu_id`, `dia_numero`, `tipo_comida_id`),
  KEY `idx_menu_id` (`menu_id`),
  CONSTRAINT `fk_mdet_menu` FOREIGN KEY (`menu_id`)
      REFERENCES `nutricion_menu` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ─── 3. Crear nutricion_menu_alimentos ───────────────────────────────────
CREATE TABLE IF NOT EXISTS `nutricion_menu_alimentos` (
  `id`                   BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `menu_detalle_id`      BIGINT UNSIGNED NOT NULL,
  `alimento_id`          BIGINT UNSIGNED NOT NULL,
  `gramos_asignados`     DECIMAL(8,2)    NOT NULL,
  `calorias_calc`        DECIMAL(8,2)    DEFAULT NULL,
  `proteinas_g_calc`     DECIMAL(8,2)    DEFAULT NULL,
  `carbohidratos_g_calc` DECIMAL(8,2)    DEFAULT NULL,
  `grasas_g_calc`        DECIMAL(8,2)    DEFAULT NULL,
  `observacion`          VARCHAR(255)    DEFAULT NULL,
  `state`                CHAR(1)         NOT NULL DEFAULT 'A',
  `creado_en`            DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_menu_detalle_id` (`menu_detalle_id`),
  KEY `idx_alimento_id` (`alimento_id`),
  CONSTRAINT `fk_malim_detalle` FOREIGN KEY (`menu_detalle_id`)
      REFERENCES `nutricion_menu_detalle` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_malim_alimento` FOREIGN KEY (`alimento_id`)
      REFERENCES `nutricion_alimentos` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ─── 4. Renombrar columna en nutricion_registro_comidas ──────────────────
-- Si la columna dieta_detalle_id existe, renombrarla a menu_detalle_id.
-- MySQL 8.0+: usar RENAME COLUMN. MySQL < 8.0: usar CHANGE COLUMN.
ALTER TABLE `nutricion_registro_comidas`
  RENAME COLUMN `dieta_detalle_id` TO `menu_detalle_id`;

-- Si la tabla aún no existe (primera migración), AutoMigrate la creará
-- directamente con menu_detalle_id. En ese caso el ALTER fallará inofensivamente;
-- puedes comentarlo según tu situación.
