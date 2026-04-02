package repositories

import (
	"saas-medico/internal/modules/historia/models"
	"time"

	"gorm.io/gorm"
)

type HistoriaRepository struct {
	db *gorm.DB
}

func NewHistoriaRepository(db *gorm.DB) *HistoriaRepository {
	return &HistoriaRepository{db: db}
}

// ─── Formularios ──────────────────────────────────────────────────────────────

func (r *HistoriaRepository) FindFormularioByID(id uint) (*models.Formulario, error) {
	var f models.Formulario
	err := r.db.First(&f, "id = ? AND state = 'A'", id).Error
	return &f, err
}

func (r *HistoriaRepository) FindFormularios(tipoID, usuarioID, clinicaID uint) ([]models.Formulario, error) {
	var list []models.Formulario
	q := r.db.Where("state = 'A'")
	if tipoID > 0 {
		q = q.Where("tipo_formulario_id = ?", tipoID)
	}
	if usuarioID > 0 {
		q = q.Where("usuario_id = ?", usuarioID)
	}
	if clinicaID > 0 {
		q = q.Where("clinica_id = ?", clinicaID)
	}
	err := q.Order("creado_en DESC").Find(&list).Error
	return list, err
}

func (r *HistoriaRepository) CreateFormulario(f *models.Formulario) error {
	return r.db.Create(f).Error
}

func (r *HistoriaRepository) UpdateFormulario(id uint, nombre, descripcion string) error {
	return r.db.Model(&models.Formulario{}).Where("id = ? AND state = 'A'", id).
		Updates(map[string]interface{}{"nombre": nombre, "descripcion": descripcion}).Error
}

func (r *HistoriaRepository) DeleteFormulario(id uint) error {
	return r.db.Model(&models.Formulario{}).Where("id = ?", id).Update("state", "I").Error
}

func (r *HistoriaRepository) FindTiposFormulario() ([]models.TipoFormulario, error) {
	var list []models.TipoFormulario
	err := r.db.Where("state = 'A'").Find(&list).Error
	return list, err
}

func (r *HistoriaRepository) CreatePregunta(p *models.FormularioPregunta) error {
	return r.db.Create(p).Error
}

func (r *HistoriaRepository) CreateOpcion(o *models.FormularioOpcion) error {
	return r.db.Create(o).Error
}

func (r *HistoriaRepository) DeletePreguntasByFormulario(formularioID uint) error {
	return r.db.Model(&models.FormularioPregunta{}).Where("formulario_id = ?", formularioID).
		Update("state", "I").Error
}

func (r *HistoriaRepository) DeleteOpcionesByFormulario(formularioID uint) error {
	return r.db.Exec(
		"UPDATE formulario_opciones fo JOIN formulario_preguntas fp ON fo.pregunta_id = fp.id SET fo.state = 'I' WHERE fp.formulario_id = ?",
		formularioID,
	).Error
}

func (r *HistoriaRepository) FindPreguntasByFormulario(formularioID uint) ([]models.FormularioPregunta, error) {
	var list []models.FormularioPregunta
	err := r.db.Preload("Opciones").Where("formulario_id = ? AND state = 'A'", formularioID).
		Order("orden ASC").Find(&list).Error
	return list, err
}

func (r *HistoriaRepository) FindOpcionesByPregunta(preguntaID uint) ([]models.FormularioOpcion, error) {
	var list []models.FormularioOpcion
	err := r.db.Where("pregunta_id = ? AND state = 'A'", preguntaID).
		Order("orden ASC").Find(&list).Error
	return list, err
}

// ─── Historia Clínica ─────────────────────────────────────────────────────────

func (r *HistoriaRepository) CreateHistoria(req *models.CreateHistoriaClinicaRequest, pacienteId int) (*models.HistoriaClinica, error) {

	var respuestas []models.HistoriaRespuesta

	fecha, err := time.Parse("2006-01-02", req.Fecha)
	if err != nil {
		return nil, err
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	historia := models.HistoriaClinica{
		PacienteID:         uint(pacienteId),
		MedicoID:           req.MedicoID,
		FormularioID:       req.FormularioID,
		Fecha:              fecha,
		ObservacionGeneral: req.ObservacionGeneral,
		State:              "A",
	}

	if err := tx.Create(&historia).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, val := range req.Preguntas {
		var respuestaFecha *time.Time

		if val.RespuestaFecha != nil && *val.RespuestaFecha != "" {
			f, err := time.Parse("2006-01-02", *val.RespuestaFecha)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			respuestaFecha = &f
		}

		resp := models.HistoriaRespuesta{
			HistoriaID:      historia.ID,
			PreguntaID:      val.PreguntaID,
			RespuestaTexto:  val.RespuestaTexto,
			RespuestaFecha:  respuestaFecha,
			RespuestaNumero: val.RespuestaNumero,
		}

		respuestas = append(respuestas, resp)
	}

	if len(respuestas) > 0 {
		if err := tx.Create(&respuestas).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &historia, nil
}
func resNumber(number *float64) *float64 {
	if *number != 0.0 {
		return number
	}
	a := 0.0
	return &a
}
func (r *HistoriaRepository) GetHistoriaClinicaByUser(userId, clinicaId, tipoForm int) ([]models.Formulario, error) {
	var list []models.Formulario
	err := r.db.
		Model(&models.Formulario{}).
		Preload("Preguntas").
		Where("clinica_id = ? AND usuario_id = ? AND tipo_formulario_id = ?", clinicaId, userId, tipoForm).
		Find(&list).Error
	return list, err
}

func (r *HistoriaRepository) FindHistoriasByPaciente(pacienteID uint) ([]models.HistoriaClinica, error) {
	var list []models.HistoriaClinica
	err := r.db.Where("paciente_id = ? AND state = 'A'", pacienteID).
		Order("fecha DESC").Find(&list).Error
	return list, err
}
func (r *HistoriaRepository) FindHistoriasClinicasByPaciente(pacienteID uint) ([]models.ResultadoHistoria, error) {
	var resultados []models.ResultadoHistoria
	err := r.db.
		Table("historia_respuestas hr").
		Select(`
		fp.pregunta AS pregunta,
		fp.orden AS orden_pregunta,
		fp.permite_multi AS multiple,
		hr.respuesta_texto AS respuesta_text,
		hr.respuesta_numero AS respuesta_numero,
		hc.id AS id_historia_clinica,
		hc.fecha AS fecha_registro,
		f.nombre AS nombre_formulario
	`).
		Joins("INNER JOIN formulario_preguntas fp ON fp.id = hr.pregunta_id").
		Joins("INNER JOIN historias_clinicas hc ON hc.id = hr.historia_id").
		Joins("INNER JOIN formularios f ON f.id = hc.formulario_id").
		Where("hc.paciente_id = ?", pacienteID).
		Scan(&resultados).Error
	if err != nil {
		return nil, err
	}
	return resultados, nil
}

func (r *HistoriaRepository) FindHistoriaByID(id uint) (*models.HistoriaClinica, error) {
	var h models.HistoriaClinica
	err := r.db.First(&h, "id = ? AND state = 'A'", id).Error
	return &h, err
}

func (r *HistoriaRepository) CreateRespuestas(respuestas []models.HistoriaRespuesta) error {
	return r.db.Create(&respuestas).Error
}

// ─── Alergias ─────────────────────────────────────────────────────────────────

func (r *HistoriaRepository) FindAlergiasByPaciente(pacienteID uint) ([]models.PacienteAlergia, error) {
	var list []models.PacienteAlergia
	err := r.db.Where("paciente_id = ? AND state = 'A'", pacienteID).Find(&list).Error
	return list, err
}

func (r *HistoriaRepository) CreateAlergia(a *models.PacienteAlergia) error {
	return r.db.Create(a).Error
}

func (r *HistoriaRepository) SoftDeleteAlergia(id, pacienteID uint) error {
	return r.db.Model(&models.PacienteAlergia{}).
		Where("id = ? AND paciente_id = ?", id, pacienteID).
		Update("state", "I").Error
}

// ─── Antecedentes ─────────────────────────────────────────────────────────────

func (r *HistoriaRepository) FindAntecedentesByPaciente(pacienteID uint) ([]models.PacienteAntecedente, error) {
	var list []models.PacienteAntecedente
	err := r.db.Where("paciente_id = ? AND state = 'A'", pacienteID).Find(&list).Error
	return list, err
}

func (r *HistoriaRepository) CreateAntecedente(a *models.PacienteAntecedente) error {
	return r.db.Create(a).Error
}

// ─── Hábitos ──────────────────────────────────────────────────────────────────

func (r *HistoriaRepository) FindHabitosByPaciente(pacienteID uint) ([]models.PacienteHabito, error) {
	var list []models.PacienteHabito
	err := r.db.Where("paciente_id = ? AND state = 'A'", pacienteID).Find(&list).Error
	return list, err
}

func (r *HistoriaRepository) SaveHabito(h *models.PacienteHabito) error {
	return r.db.Save(h).Error
}

// ─── Diagnósticos ─────────────────────────────────────────────────────────────

func (r *HistoriaRepository) FindDiagnosticosByPaciente(pacienteID uint) ([]models.PacienteDiagnostico, error) {
	var list []models.PacienteDiagnostico
	err := r.db.Where("paciente_id = ? AND state = 'A'", pacienteID).
		Order("fecha_diagnostico DESC").Find(&list).Error
	return list, err
}

func (r *HistoriaRepository) CreateDiagnostico(d *models.PacienteDiagnostico) error {
	return r.db.Create(d).Error
}

func (r *HistoriaRepository) UpdateDiagnostico(d *models.PacienteDiagnostico) error {
	return r.db.Save(d).Error
}

// ─── Imágenes del paciente ────────────────────────────────────────────────────

func (r *HistoriaRepository) CreatePacienteImagen(img *models.PacienteImagen) error {
	return r.db.Create(img).Error
}

func (r *HistoriaRepository) FindImagenesByPaciente(pacienteID uint) ([]models.PacienteImagen, error) {
	var list []models.PacienteImagen
	err := r.db.Where("paciente_id = ?", pacienteID).Order("creado_en DESC").Find(&list).Error
	return list, err
}
