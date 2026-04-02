-- Migration 003: Create or update paciente_imagenes table
-- Run: mysql -u$DB_USER -p$DB_PASSWORD $DB_NAME < resources/migrations/003_paciente_imagenes.sql

CREATE TABLE IF NOT EXISTS `paciente_imagenes` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `paciente_id`    BIGINT UNSIGNED NOT NULL,
  `medico_id`      BIGINT UNSIGNED DEFAULT NULL,
  `nombre_archivo` VARCHAR(255)    NOT NULL,
  `url_archivo`    VARCHAR(500)    NOT NULL,
  `tipo_imagen`    INT             NOT NULL DEFAULT 1,
  `descripcion`    VARCHAR(255)    DEFAULT NULL,
  `creado_en`      DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  INDEX `idx_paciente_imagenes_paciente_id` (`paciente_id`),
  INDEX `idx_paciente_imagenes_medico_id`   (`medico_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- If the table already existed with fewer columns, add missing ones:
ALTER TABLE `paciente_imagenes`
  ADD COLUMN IF NOT EXISTS `medico_id`      BIGINT UNSIGNED DEFAULT NULL,
  ADD COLUMN IF NOT EXISTS `url_archivo`    VARCHAR(500)    NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS `tipo_imagen`    INT             NOT NULL DEFAULT 1,
  ADD COLUMN IF NOT EXISTS `descripcion`    VARCHAR(255)    DEFAULT NULL;
