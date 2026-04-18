package repositories

import (
	"log"
	"saas-medico/internal/modules/nutricion/models"
	"time"

	"gorm.io/gorm"
)

type NutricionRepository struct {
	db *gorm.DB
}

func NewNutricionRepository(db *gorm.DB) *NutricionRepository {
	return &NutricionRepository{db: db}
}

// ─── Grupos de alimento ───────────────────────────────────────────────────────

func (r *NutricionRepository) FindGruposAlimento() ([]models.NutricionGrupoAlimento, error) {
	var list []models.NutricionGrupoAlimento
	err := r.db.Where("state = 'A'").Order("orden ASC, nombre ASC").Find(&list).Error
	return list, err
}

func (r *NutricionRepository) FindGrupoAlimentoByID(id uint) (*models.NutricionGrupoAlimento, error) {
	var g models.NutricionGrupoAlimento
	err := r.db.First(&g, "id = ? AND state = 'A'", id).Error
	return &g, err
}

func (r *NutricionRepository) CreateGrupoAlimento(g *models.NutricionGrupoAlimento) error {
	return r.db.Create(g).Error
}

func (r *NutricionRepository) UpdateGrupoAlimento(g *models.NutricionGrupoAlimento) error {
	return r.db.Save(g).Error
}

func (r *NutricionRepository) DeleteGrupoAlimento(id uint) error {
	return r.db.Model(&models.NutricionGrupoAlimento{}).
		Where("id = ?", id).
		Update("state", "I").Error
}

// ─── Alimentos ────────────────────────────────────────────────────────────────
type DatosRequerimeintos struct {
	ID              uint
	TipoComidaID    uint
	GrupoAlimentoID uint
}

func (r *NutricionRepository) DesactivarMenusAnitguos() {
	today := time.Now().Format("2006-01-02")
	err := r.db.
		Model(&models.NutricionMenu{}).
		Where("fecha_fin = ? AND state = ?", today, "A").
		Updates(map[string]interface{}{
			"state":  "I",
			"estado": "FIN",
		}).Error
	if err != nil {
		log.Println("error desactivando menús:", err)
		return
	}

	log.Println("menús desactivados correctamente")
}

func (r *NutricionRepository) GetRequerimientosPorComida() (map[uint][]uint, error) {
	var data []models.NutricionTipoComidaGrupo

	err := r.db.Find(&data).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uint][]uint)

	for _, d := range data {
		result[d.TipoComidaID] = append(result[d.TipoComidaID], d.GrupoAlimentoID)
	}

	return result, nil
}
func (r *NutricionRepository) FindAlimentos(categoria string) ([]models.NutricionAlimento, error) {
	var list []models.NutricionAlimento
	q := r.db.Where("state = 'A'")
	if categoria != "" {
		q = q.Where("categoria = ?", categoria)
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *NutricionRepository) FindAlimentoByID(id uint) (*models.NutricionAlimento, error) {
	var a models.NutricionAlimento
	err := r.db.First(&a, "id = ? AND state = 'A'", id).Error
	return &a, err
}

func (r *NutricionRepository) CreateAlimento(a *models.NutricionAlimento) error {
	return r.db.Create(a).Error
}

func (r *NutricionRepository) UpdateAlimento(a *models.NutricionAlimento) error {
	return r.db.Save(a).Error
}

func (r *NutricionRepository) FindTipoComida() ([]models.NutricionTipoComida, error) {
	var list []models.NutricionTipoComida
	q := r.db.Where("state = 'A'")
	err := q.Find(&list).Error
	return list, err
}

// ─── Dietas catálogo ──────────────────────────────────────────────────────────

func (r *NutricionRepository) FindDietasCatalogo() ([]models.NutricionDietaCatalogo, error) {
	var list []models.NutricionDietaCatalogo
	err := r.db.Where("state = 'A'").Find(&list).Error
	return list, err
}

// ─── Plan de dieta del paciente ───────────────────────────────────────────────

func (r *NutricionRepository) FindDietasByPaciente(pacienteID uint) ([]models.NutricionDietaPaciente, error) {
	var list []models.NutricionDietaPaciente
	err := r.db.Preload("Paciente").Where("paciente_id = ? AND state = 'A'", pacienteID).
		Order("fecha_inicio DESC").Find(&list).Error
	return list, err
}

func (r *NutricionRepository) FindDietaByID(id uint) (*models.NutricionDietaPaciente, error) {
	var d models.NutricionDietaPaciente
	err := r.db.Preload("Paciente").First(&d, "id = ? AND state = 'A'", id).Error
	return &d, err
}

func (r *NutricionRepository) CreateDieta(d *models.NutricionDietaPaciente) error {
	return r.db.Create(d).Error
}

func (r *NutricionRepository) UpdateDieta(d *models.NutricionDietaPaciente) error {
	return r.db.Save(d).Error
}

func (r *NutricionRepository) CreateDetalle(d *models.NutricionMenuDetalle) error {
	return r.db.Create(d).Error
}

func (r *NutricionRepository) UpdateDetalleReceta(id uint, receta string) (*models.NutricionMenuDetalle, error) {
	var d models.NutricionMenuDetalle
	if err := r.db.First(&d, "id = ? AND state = 'A'", id).Error; err != nil {
		return nil, err
	}
	d.NombreReceta = receta
	if err := r.db.Save(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *NutricionRepository) DeleteDetalle(id uint) error {
	return r.db.Model(&models.NutricionMenuDetalle{}).
		Where("id = ?", id).Update("state", "I").Error
}

func (r *NutricionRepository) FindAlimentosByMenuDetalle(detalleID uint) ([]models.NutricionMenuAlimento, error) {
	var list []models.NutricionMenuAlimento
	err := r.db.Where("menu_detalle_id = ? AND state = 'A'", detalleID).Find(&list).Error
	return list, err
}

func (r *NutricionRepository) CreateMenuAlimento(a *models.NutricionMenuAlimento) error {
	return r.db.Create(a).Error
}

func (r *NutricionRepository) DeleteMenuAlimento(id uint) error {
	return r.db.Model(&models.NutricionMenuAlimento{}).
		Where("id = ?", id).Update("state", "I").Error
}

func (r *NutricionRepository) UpdateMenuAlimentoGramos(id uint, gramos, cal, prot, carb, gras float64) error {
	return r.db.Model(&models.NutricionMenuAlimento{}).Where("id = ?", id).Updates(map[string]interface{}{
		"gramos_asignados":     gramos,
		"calorias_calc":        cal,
		"proteinas_g_calc":     prot,
		"carbohidratos_g_calc": carb,
		"grasas_g_calc":        gras,
	}).Error
}

func (r *NutricionRepository) FindMenuAlimentoByID(id uint) (*models.NutricionMenuAlimento, error) {
	var a models.NutricionMenuAlimento
	err := r.db.Preload("Alimento").Where("id = ? AND state = 'A'", id).First(&a).Error
	return &a, err
}

// ─── Menú semanal ─────────────────────────────────────────────────────────────

func (r *NutricionRepository) CreateMenu(m *models.NutricionMenu) error {
	return r.db.Create(m).Error
}
func (r *NutricionRepository) CreateMenuPlantilla(mp *models.NutricionMenuPlantilla) error {
	return r.db.Create(mp).Error
}
func (r *NutricionRepository) GetMenuPlantilla(dietaID uint) (*models.NutricionMenuPlantilla, error) {
	var mp models.NutricionMenuPlantilla
	err := r.db.Where("dieta_paciente_id = ?", dietaID).Find(&mp).Error
	if err != nil {
		return nil, err
	}
	return &mp, nil
}
func (r *NutricionRepository) UpdateMenuPlantilla(dietaPacienteID, newMenuID uint) error {
	return r.db.Model(&models.NutricionMenuPlantilla{}).
		Where("dieta_paciente_id = ?", dietaPacienteID).
		Update("menu_id", newMenuID).Error
}
func (r *NutricionRepository) FindMenusByDieta(dietaID uint) ([]models.NutricionMenu, error) {
	var list []models.NutricionMenu
	err := r.db.Where("dieta_paciente_id = ? AND state = 'A'", dietaID).
		Order("semana_numero ASC").Find(&list).Error
	return list, err
}

func (r *NutricionRepository) FindMenuByID(id uint) (*models.NutricionMenu, error) {
	var m models.NutricionMenu
	err := r.db.
		Preload("Detalles", "state IN ('A','C')").
		Preload("Detalles.Alimentos", "state = 'A'").
		Preload("Detalles.Alimentos.Alimento").
		First(&m, "id = ? AND state = 'A'", id).Error
	return &m, err
}

func (r *NutricionRepository) UpdateMenu(m *models.NutricionMenu) error {
	return r.db.Save(m).Error
}

func (r *NutricionRepository) FindMenusRequierenCambio() ([]models.NutricionMenu, error) {
	var list []models.NutricionMenu
	err := r.db.
		Preload("Dieta").
		Preload("Dieta.Paciente").
		Where("DATE(fecha_fin) = CURDATE() AND state = 'A'").
		Find(&list).Error
	return list, err
}

func (r *NutricionRepository) FindDetallesByMenu(menuID uint) ([]models.NutricionMenuDetalle, error) {
	var list []models.NutricionMenuDetalle
	err := r.db.Where("menu_id = ? AND state IN ('A','C')", menuID).Preload("Alimentos").Preload("Alimentos.Alimento").
		Order("dia_numero ASC, tipo_comida_id ASC").Find(&list).Error
	return list, err
}

func (r *NutricionRepository) FindMenuDetalleByID(id uint) (*models.NutricionMenuDetalle, error) {
	var d models.NutricionMenuDetalle
	err := r.db.First(&d, "id = ? AND state IN ('A','C')", id).Error
	return &d, err
}

func (r *NutricionRepository) MarcarMenuDetalleConsumido(id uint) error {
	return r.db.Model(&models.NutricionMenuDetalle{}).
		Where("id = ?", id).Update("state", "C").Error
}

func (r *NutricionRepository) CreateMenuDetalles(comidas []*models.NutricionMenuDetalle) ([]*models.NutricionMenuDetalle, error) {
	err := r.db.CreateInBatches(comidas, 100).Error
	if err != nil {
		return nil, err
	}
	return comidas, nil
}

func (r *NutricionRepository) AddAlimentosToComidas(alimentos []*models.NutricionMenuAlimento) ([]*models.NutricionMenuAlimento, error) {
	err := r.db.CreateInBatches(alimentos, 100).Error
	if err != nil {
		return nil, err
	}
	return alimentos, nil
}

// func (r *NutricionRepository) AddAlimentosToComida(alimentos []*models.NutricionAlimento) error{

// }

// ─── Recordatorio 24 horas ────────────────────────────────────────────────────

func (r *NutricionRepository) CreateR24H(rec *models.NutricionR24H) error {
	return r.db.Create(rec).Error
}

func (r *NutricionRepository) FindR24HByPaciente(pacienteID uint) ([]models.NutricionR24H, error) {
	var list []models.NutricionR24H
	err := r.db.Where("paciente_id = ? AND state = 'A'", pacienteID).
		Order("fecha DESC").Find(&list).Error
	return list, err
}

func (r *NutricionRepository) FindR24HByID(id uint) (*models.NutricionR24H, error) {
	var rec models.NutricionR24H
	err := r.db.First(&rec, "id = ? AND state = 'A'", id).Error
	return &rec, err
}

func (r *NutricionRepository) CreateR24HItem(item *models.NutricionR24HItem) error {
	return r.db.Create(item).Error
}

func (r *NutricionRepository) FindR24HItems(r24hID uint) ([]models.NutricionR24HItem, error) {
	var list []models.NutricionR24HItem
	err := r.db.Where("r24h_id = ? AND state = 'A'", r24hID).Find(&list).Error
	return list, err
}

// ─── Preferencias alimentarias ────────────────────────────────────────────────

func (r *NutricionRepository) CreatePreferencia(p *models.NutricionPreferenciaAlimento) error {
	return r.db.Create(p).Error
}

func (r *NutricionRepository) FindPreferenciasByPaciente(pacienteID uint) ([]models.NutricionPreferenciaAlimento, error) {
	var list []models.NutricionPreferenciaAlimento
	err := r.db.Preload("Alimento").Where("paciente_id = ? AND state = 'A'", pacienteID).Find(&list).Error
	return list, err
}

func (r *NutricionRepository) DeletePreferencia(id uint) error {
	return r.db.Model(&models.NutricionPreferenciaAlimento{}).
		Where("id = ?", id).Update("state", "I").Error
}

// ─── Síntomas ─────────────────────────────────────────────────────────────────

func (r *NutricionRepository) CreateSintoma(s *models.NutricionSintoma) error {
	return r.db.Create(s).Error
}

func (r *NutricionRepository) FindSintomasByPaciente(pacienteID uint, fechaDesde, fechaHasta string) ([]models.NutricionSintoma, error) {
	var list []models.NutricionSintoma
	q := r.db.Where("paciente_id = ? AND state = 'A'", pacienteID)
	if fechaDesde != "" {
		q = q.Where("fecha >= ?", fechaDesde)
	}
	if fechaHasta != "" {
		q = q.Where("fecha <= ?", fechaHasta)
	}
	err := q.Order("fecha DESC").Find(&list).Error
	return list, err
}

// ─── Tipo de Recurso ──────────────────────────────────────────────────────────

func (r *NutricionRepository) FindTipoRecursos() ([]models.NutricionTipoRecurso, error) {
	var list []models.NutricionTipoRecurso
	err := r.db.Where("state = 'A'").Order("nombre ASC").Find(&list).Error
	return list, err
}

func (r *NutricionRepository) CreateTipoRecurso(t *models.NutricionTipoRecurso) error {
	return r.db.Create(t).Error
}

func (r *NutricionRepository) UpdateTipoRecurso(id uint, nombre string) (*models.NutricionTipoRecurso, error) {
	var t models.NutricionTipoRecurso
	if err := r.db.Where("id = ? AND state = 'A'", id).First(&t).Error; err != nil {
		return nil, err
	}
	t.Nombre = nombre
	if err := r.db.Save(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *NutricionRepository) DeleteTipoRecurso(id uint) error {
	return r.db.Model(&models.NutricionTipoRecurso{}).Where("id = ?", id).Update("state", "I").Error
}

// ─── Archivos PDF ─────────────────────────────────────────────────────────────

func (r *NutricionRepository) CreateArchivoPDF(a *models.NutricionArchivoPDF) error {
	if err := r.db.Create(a).Error; err != nil {
		return err
	}
	return r.db.Preload("TipoRecurso").First(a, a.ID).Error
}

func (r *NutricionRepository) FindArchivosPDF(clinicaID uint, pacienteID *uint, tipoRecursoID *uint) ([]models.NutricionArchivoPDF, error) {
	var list []models.NutricionArchivoPDF
	q := r.db.Preload("TipoRecurso").Where("nutricion_archivos_pdf.clinica_id = ? AND nutricion_archivos_pdf.state = 'A'", clinicaID)

	if tipoRecursoID != nil {
		q = q.Where("nutricion_archivos_pdf.tipo = ?", *tipoRecursoID)
	}

	err := q.Order("nutricion_archivos_pdf.creado_en DESC").Find(&list).Error
	return list, err
}
func (r *NutricionRepository) FindArchivosPDFByUser(clinicaID uint, pacienteID uint) ([]models.NutricionArchivoPDF, error) {
	var list []models.NutricionArchivoPDF
	err := r.db.Preload("TipoRecurso").
		Where("nutricion_archivos_pdf.clinica_id = ? AND nutricion_archivos_pdf.state = 'A' AND (paciente_id IS NULL OR paciente_id = ?)", clinicaID, pacienteID).
		Order("nutricion_archivos_pdf.creado_en DESC").
		Find(&list).Error
	return list, err
}

func (r *NutricionRepository) DeleteArchivoPDF(id uint) error {
	return r.db.Model(&models.NutricionArchivoPDF{}).
		Where("id = ?", id).Update("state", "I").Error
}

// ─── Ejercicios ───────────────────────────────────────────────────────────────

func (r *NutricionRepository) CreateEjercicioCatalogo(e *models.NutricionEjercicioCatalogo) error {
	return r.db.Create(e).Error
}

func (r *NutricionRepository) FindEjerciciosCatalogo(categoria string) ([]models.NutricionEjercicioCatalogo, error) {
	var list []models.NutricionEjercicioCatalogo
	q := r.db.Where("state = 'A'")
	if categoria != "" {
		q = q.Where("categoria = ?", categoria)
	}
	err := q.Find(&list).Error
	return list, err
}

func (r *NutricionRepository) FindEjercicioCatalogoByID(id uint) (*models.NutricionEjercicioCatalogo, error) {
	var e models.NutricionEjercicioCatalogo
	err := r.db.Where("id = ? AND state = 'A'", id).First(&e).Error
	return &e, err
}

func (r *NutricionRepository) FindEjerciciosByPaciente(pacienteID uint) ([]models.NutricionEjercicioPaciente, error) {
	var list []models.NutricionEjercicioPaciente
	err := r.db.Where("paciente_id = ? AND state = 'A'", pacienteID).Find(&list).Error
	return list, err
}

func (r *NutricionRepository) CreateEjercicioPaciente(e *models.NutricionEjercicioPaciente) error {
	return r.db.Create(e).Error
}

func (r *NutricionRepository) UpdateEjercicioPaciente(e *models.NutricionEjercicioPaciente) error {
	return r.db.Save(e).Error
}

// ─── Registros ────────────────────────────────────────────────────────────────

func (r *NutricionRepository) FindRegistrosComida(pacienteID uint, fecha, desde, hasta string) ([]models.NutricionRegistroComida, error) {
	var list []models.NutricionRegistroComida
	q := r.db.Where("paciente_id = ? AND state = 'A'", pacienteID)
	if fecha != "" {
		q = q.Where("DATE(fecha) = ?", fecha)
	} else {
		if desde != "" {
			q = q.Where("DATE(fecha) >= ?", desde)
		}
		if hasta != "" {
			q = q.Where("DATE(fecha) <= ?", hasta)
		}
	}
	err := q.Order("fecha DESC").Find(&list).Error
	return list, err
}

func (r *NutricionRepository) CreateRegistroComida(rc *models.NutricionRegistroComida) error {
	return r.db.Create(rc).Error
}

func (r *NutricionRepository) FindRegistroComidaByID(id uint) (*models.NutricionRegistroComida, error) {
	var rc models.NutricionRegistroComida
	err := r.db.First(&rc, "id = ? AND state = 'A'", id).Error
	return &rc, err
}

func (r *NutricionRepository) UpdateFotoComida(id uint, rutaFoto string) (*models.NutricionRegistroComida, string, error) {
	rc, err := r.FindRegistroComidaByID(id)
	if err != nil {
		return nil, "", err
	}
	oldFoto := rc.FotoComida
	rc.FotoComida = rutaFoto
	if err := r.db.Save(rc).Error; err != nil {
		return nil, "", err
	}
	return rc, oldFoto, nil
}

func (r *NutricionRepository) MarcarRegistroComidaConsumida(id uint) error {
	return r.db.Model(&models.NutricionRegistroComida{}).
		Where("id = ? AND state = 'A'", id).
		Update("estado", models.EstadoRegistroComidaConsumida).Error
}

// FindRegistroComidaByMenuDetalle busca si ya existe un registro consumido para un detalle de menú en una fecha
func (r *NutricionRepository) FindRegistroComidaByMenuDetalle(pacienteID, menuDetalleID uint, fecha string) (*models.NutricionRegistroComida, error) {
	var rc models.NutricionRegistroComida
	err := r.db.Where("paciente_id = ? AND menu_detalle_id = ? AND DATE(fecha) = ? AND estado = 'C' AND state = 'A'",
		pacienteID, menuDetalleID, fecha).First(&rc).Error
	return &rc, err
}

func (r *NutricionRepository) CreateRegistroAlimento(ra *models.NutricionRegistroAlimento) error {
	return r.db.Create(ra).Error
}

func (r *NutricionRepository) FindRegistroAlimentosByRegistro(registroComidaID uint) ([]models.NutricionRegistroAlimento, error) {
	var list []models.NutricionRegistroAlimento
	err := r.db.Where("registro_comida_id = ? AND state = 'A'", registroComidaID).Find(&list).Error
	return list, err
}

func (r *NutricionRepository) UpdateRegistroComidaMacros(registroID uint, cal, prot, carb, gras float64) error {
	return r.db.Model(&models.NutricionRegistroComida{}).
		Where("id = ?", registroID).
		Updates(map[string]interface{}{
			"calorias_consumidas": cal,
			"proteinas_g":         prot,
			"carbohidratos_g":     carb,
			"grasas_g":            gras,
		}).Error
}

// FindRegistroAlimentosByRegistros carga los alimentos de múltiples registros de comida
func (r *NutricionRepository) FindRegistroAlimentosByRegistros(ids []uint) ([]models.NutricionRegistroAlimento, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []models.NutricionRegistroAlimento
	err := r.db.Where("registro_comida_id IN ? AND state = 'A'", ids).Find(&list).Error
	return list, err
}

func (r *NutricionRepository) FindRegistrosEjercicio(pacienteID uint, fecha, desde, hasta string) ([]models.NutricionRegistroEjercicio, error) {
	var list []models.NutricionRegistroEjercicio
	q := r.db.Where("paciente_id = ? AND state = 'A'", pacienteID)
	if fecha != "" {
		q = q.Where("DATE(fecha) = ?", fecha)
	} else {
		if desde != "" {
			q = q.Where("DATE(fecha) >= ?", desde)
		}
		if hasta != "" {
			q = q.Where("DATE(fecha) <= ?", hasta)
		}
	}
	err := q.Order("fecha DESC").Find(&list).Error
	return list, err
}

func (r *NutricionRepository) CreateRegistroEjercicio(re *models.NutricionRegistroEjercicio) error {
	return r.db.Create(re).Error
}

// ─── Progreso ─────────────────────────────────────────────────────────────────

func (r *NutricionRepository) FindProgresoByPaciente(pacienteID uint) ([]models.NutricionProgresoPaciente, error) {
	var list []models.NutricionProgresoPaciente
	err := r.db.Where("paciente_id = ? AND state = 'A'", pacienteID).
		Order("fecha DESC").Find(&list).Error
	return list, err
}

func (r *NutricionRepository) CreateProgreso(p *models.NutricionProgresoPaciente) error {
	return r.db.Create(p).Error
}

func (r *NutricionRepository) UpdateProgreso(p *models.NutricionProgresoPaciente) error {
	return r.db.Save(p).Error
}

// FindProgresoPorFecha devuelve el registro de progreso de un paciente en una fecha específica
func (r *NutricionRepository) FindProgresoPorFecha(pacienteID uint, fecha string) (*models.NutricionProgresoPaciente, error) {
	var p models.NutricionProgresoPaciente
	err := r.db.Where("paciente_id = ? AND DATE(fecha) = ? AND state = 'A'", pacienteID, fecha).First(&p).Error
	return &p, err
}

// FindUltimoProgreso devuelve el registro de progreso más reciente del paciente
func (r *NutricionRepository) FindUltimoProgreso(pacienteID uint) (*models.NutricionProgresoPaciente, error) {
	var p models.NutricionProgresoPaciente
	err := r.db.Where("paciente_id = ? AND state = 'A'", pacienteID).
		Order("fecha DESC").First(&p).Error
	return &p, err
}

// ─── XP y logros ──────────────────────────────────────────────────────────────

func (r *NutricionRepository) FindOrCreateXP(pacienteID uint) (*models.NutricionPacienteXP, error) {
	var xp models.NutricionPacienteXP
	err := r.db.Where("paciente_id = ?", pacienteID).FirstOrCreate(&xp, models.NutricionPacienteXP{
		PacienteID: pacienteID,
	}).Error
	return &xp, err
}

func (r *NutricionRepository) SaveXP(xp *models.NutricionPacienteXP) error {
	return r.db.Save(xp).Error
}

func (r *NutricionRepository) FindLogrosByPaciente(pacienteID uint) ([]models.NutricionLogroPaciente, error) {
	var list []models.NutricionLogroPaciente
	err := r.db.Where("paciente_id = ? AND state = 'A'", pacienteID).Find(&list).Error
	return list, err
}

func (r *NutricionRepository) GrantLogro(l *models.NutricionLogroPaciente) error {
	return r.db.Create(l).Error
}

func (r *NutricionRepository) FindLogrosCatalogo() ([]models.NutricionLogroCatalogo, error) {
	var list []models.NutricionLogroCatalogo
	err := r.db.Where("state = 'A'").Find(&list).Error
	return list, err
}

// FindDietaActivaByPaciente devuelve la dieta ACTIVA más reciente del paciente
func (r *NutricionRepository) FindDietaActivaByPaciente(pacienteID uint) (*models.NutricionDietaPaciente, error) {
	var d models.NutricionDietaPaciente
	err := r.db.Where("paciente_id = ? AND estado = 'ACTIVA' AND state = 'A'", pacienteID).
		Order("fecha_inicio DESC").First(&d).Error
	return &d, err
}

// ─── Plantillas de menú semanal ───────────────────────────────────────────────

func (r *NutricionRepository) CreatePlantillaSemana(p *models.NutricionMenuPlantillaSemana) error {
	return r.db.Create(p).Error
}

func (r *NutricionRepository) FindPlantillas(clinicaID uint, numComidas *int, semanaNumero *int) ([]models.NutricionMenuPlantillaSemana, error) {
	var list []models.NutricionMenuPlantillaSemana
	q := r.db.Where("clinica_id = ? AND state = 'A'", clinicaID)
	if numComidas != nil {
		q = q.Where("num_comidas = ?", *numComidas)
	}
	if semanaNumero != nil {
		q = q.Where("semana_numero = ?", *semanaNumero)
	}
	err := q.Order("semana_numero ASC, num_comidas ASC").Find(&list).Error
	return list, err
}

func (r *NutricionRepository) FindPlantillaSemanaByID(id uint) (*models.NutricionMenuPlantillaSemana, error) {
	var p models.NutricionMenuPlantillaSemana
	err := r.db.Preload("Detalles", "state = 'A'").
		Preload("Detalles.Alimentos", "state = 'A'").
		Preload("Detalles.Alimentos.Alimento").
		First(&p, "id = ? AND state = 'A'", id).Error
	return &p, err
}

func (r *NutricionRepository) UpdatePlantillaSemana(p *models.NutricionMenuPlantillaSemana) error {
	return r.db.Save(p).Error
}

func (r *NutricionRepository) DeletePlantillaSemana(id uint) error {
	return r.db.Model(&models.NutricionMenuPlantillaSemana{}).Where("id = ?", id).Update("state", "I").Error
}

// ─── Detalles de plantilla ────────────────────────────────────────────────────

func (r *NutricionRepository) CreateDetallePlantilla(d *models.NutricionMenuDetallePlantilla) error {
	return r.db.Create(d).Error
}

func (r *NutricionRepository) FindDetallesByPlantilla(plantillaID uint) ([]models.NutricionMenuDetallePlantilla, error) {
	var list []models.NutricionMenuDetallePlantilla
	err := r.db.Preload("Alimentos", "state = 'A'").
		Preload("Alimentos.Alimento").
		Where("menu_id = ? AND state = 'A'", plantillaID).
		Order("dia_numero ASC, tipo_comida_id ASC").Find(&list).Error
	return list, err
}

func (r *NutricionRepository) FindDetallePlantillaByID(id uint) (*models.NutricionMenuDetallePlantilla, error) {
	var d models.NutricionMenuDetallePlantilla
	err := r.db.Preload("Alimentos", "state = 'A'").
		First(&d, "id = ? AND state = 'A'", id).Error
	return &d, err
}

func (r *NutricionRepository) UpdateDetallePlantilla(d *models.NutricionMenuDetallePlantilla) error {
	return r.db.Save(d).Error
}

func (r *NutricionRepository) DeleteDetallePlantilla(id uint) error {
	return r.db.Model(&models.NutricionMenuDetallePlantilla{}).Where("id = ?", id).Update("state", "I").Error
}

// ─── Alimentos de plantilla ───────────────────────────────────────────────────

func (r *NutricionRepository) CreateAlimentoPlantilla(a *models.NutricionMenuAlimentoPlantilla) error {
	return r.db.Create(a).Error
}

func (r *NutricionRepository) FindAlimentosByDetallePlantilla(detalleID uint) ([]models.NutricionMenuAlimentoPlantilla, error) {
	var list []models.NutricionMenuAlimentoPlantilla
	err := r.db.Preload("Alimento").
		Where("menu_detalle_id = ? AND state = 'A'", detalleID).Find(&list).Error
	return list, err
}

func (r *NutricionRepository) FindAlimentoPlantillaByID(id uint) (*models.NutricionMenuAlimentoPlantilla, error) {
	var a models.NutricionMenuAlimentoPlantilla
	err := r.db.Preload("Alimento").First(&a, "id = ? AND state = 'A'", id).Error
	return &a, err
}

func (r *NutricionRepository) UpdateAlimentoPlantillaGramos(id uint, gramos, cal, prot, carb, gras float64) error {
	return r.db.Model(&models.NutricionMenuAlimentoPlantilla{}).Where("id = ?", id).Updates(map[string]interface{}{
		"gramos_asignados":     gramos,
		"calorias_calc":        cal,
		"proteinas_g_calc":     prot,
		"carbohidratos_g_calc": carb,
		"grasas_g_calc":        gras,
	}).Error
}

func (r *NutricionRepository) DeleteAlimentoPlantilla(id uint) error {
	return r.db.Model(&models.NutricionMenuAlimentoPlantilla{}).Where("id = ?", id).Update("state", "I").Error
}

func (r *NutricionRepository) CreateMenuDetallesPlantilla(detalles []*models.NutricionMenuDetallePlantilla) ([]*models.NutricionMenuDetallePlantilla, error) {
	if err := r.db.Create(&detalles).Error; err != nil {
		return nil, err
	}
	return detalles, nil
}

func (r *NutricionRepository) AddAlimentosToPlantillaComidas(alimentos []*models.NutricionMenuAlimentoPlantilla) ([]*models.NutricionMenuAlimentoPlantilla, error) {
	if err := r.db.Create(&alimentos).Error; err != nil {
		return nil, err
	}
	return alimentos, nil
}
