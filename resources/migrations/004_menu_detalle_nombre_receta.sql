-- Agrega columna nombre_receta a nutricion_menu_detalle
ALTER TABLE nutricion_menu_detalle
    ADD COLUMN IF NOT EXISTS nombre_receta VARCHAR(150) NULL DEFAULT NULL;
