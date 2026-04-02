package repositories

import (
	"saas-medico/internal/modules/agenda/models"
	"time"

	"gorm.io/gorm"
)

type AgendaRepository struct {
	db *gorm.DB
}

func NewAgendaRepository(db *gorm.DB) *AgendaRepository {
	return &AgendaRepository{db: db}
}

// ── Citas ─────────────────────────────────────────────────────────────────────

func (r *AgendaRepository) CreateCita(c *models.Cita) error {
	return r.db.Create(c).Error
}

func (r *AgendaRepository) FindCitaByID(id uint) (*models.Cita, error) {
	var c models.Cita
	if err := r.db.Preload("TipoCita").Preload("EstadoCita").Preload("Paciente").
		First(&c, "id = ? AND state = 'A'", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *AgendaRepository) FindCitas(medicoID, clinicaID uint, fecha *time.Time, page, pageSize int) ([]models.Cita, int64, error) {
	var list []models.Cita
	var total int64

	q := r.db.Model(&models.Cita{}).Where("state = 'A'")
	if medicoID > 0 {
		q = q.Where("id_medico = ?", medicoID)
	}
	if clinicaID > 0 {
		q = q.Where("id_clinica = ?", clinicaID)
	}
	if fecha != nil {
		q = q.Where("DATE(fecha) = DATE(?)", *fecha)
	}

	q.Count(&total)
	offset := (page - 1) * pageSize
	err := q.Preload("TipoCita").Preload("EstadoCita").Preload("Paciente").Preload("PrePaciente").
		Offset(offset).Limit(pageSize).Order("fecha ASC, hora ASC").Find(&list).Error
	return list, total, err
}

func (r *AgendaRepository) UpdateCitaPaciente(id, pacienteID uint) error {
	return r.db.Model(&models.Cita{}).Where("id = ?", id).Update("id_paciente", pacienteID).Error
}

func (r *AgendaRepository) UpdateCita(c *models.Cita) error {
	return r.db.Save(c).Error
}

func (r *AgendaRepository) UpdateEstadoCita(id, estadoID uint) error {
	return r.db.Model(&models.Cita{}).Where("id = ?", id).Update("estado_cita_id", estadoID).Error
}

func (r *AgendaRepository) DeleteCita(id uint) error {
	return r.db.Model(&models.Cita{}).Where("id = ?", id).Update("state", "I").Error
}

// ── Sesiones ──────────────────────────────────────────────────────────────────

func (r *AgendaRepository) CreateSesion(s *models.Sesion) error {
	return r.db.Create(s).Error
}

func (r *AgendaRepository) FindSesionByID(id uint) (*models.Sesion, error) {
	var s models.Sesion
	if err := r.db.Preload("Cita").First(&s, "id = ? AND state = 'A'", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *AgendaRepository) FindSesionByCitaID(citaID uint) (*models.Sesion, error) {
	var s models.Sesion
	if err := r.db.Where("cita_id = ? AND state = 'A'", citaID).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *AgendaRepository) UpdateSesion(s *models.Sesion) error {
	return r.db.Save(s).Error
}

// ── Horarios ──────────────────────────────────────────────────────────────────

func (r *AgendaRepository) CreateHorario(h *models.HorarioMedico) error {
	return r.db.Create(h).Error
}

func (r *AgendaRepository) FindHorariosByMedico(medicoID uint) ([]models.HorarioMedico, error) {
	var list []models.HorarioMedico
	err := r.db.Where("medico_id = ? AND state = 'A'", medicoID).Order("dia_semana ASC").Find(&list).Error
	return list, err
}

func (r *AgendaRepository) FindHorarioByID(id uint) (*models.HorarioMedico, error) {
	var h models.HorarioMedico
	if err := r.db.First(&h, "id = ? AND state = 'A'", id).Error; err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *AgendaRepository) UpdateHorario(h *models.HorarioMedico) error {
	return r.db.Save(h).Error
}

func (r *AgendaRepository) DeleteHorario(id uint) error {
	return r.db.Model(&models.HorarioMedico{}).Where("id = ?", id).Update("state", "I").Error
}

// ── Bloqueos ──────────────────────────────────────────────────────────────────

func (r *AgendaRepository) CreateBloqueo(b *models.BloqueoAgenda) error {
	return r.db.Create(b).Error
}

func (r *AgendaRepository) FindBloqueosByMedico(medicoID uint) ([]models.BloqueoAgenda, error) {
	var list []models.BloqueoAgenda
	err := r.db.Where("medico_id = ? AND state = 'A'", medicoID).Find(&list).Error
	return list, err
}

func (r *AgendaRepository) DeleteBloqueo(id uint) error {
	return r.db.Model(&models.BloqueoAgenda{}).Where("id = ?", id).Update("state", "I").Error
}

// ── Catálogos ─────────────────────────────────────────────────────────────────

func (r *AgendaRepository) FindAllTiposCita() ([]models.TipoCita, error) {
	var list []models.TipoCita
	err := r.db.Where("state = 'A'").Find(&list).Error
	return list, err
}

func (r *AgendaRepository) FindTiposCitaByRol(rolID uint) ([]models.TipoCita, error) {
	var list []models.TipoCita
	err := r.db.Where("state = 'A' AND id_rol = ?", rolID).Find(&list).Error
	return list, err
}

func (r *AgendaRepository) FindEstadoCitaByCodigo(codigo string) (*models.EstadoCita, error) {
	var e models.EstadoCita
	if err := r.db.Where("codigo = ?", codigo).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *AgendaRepository) FindAllEstadosCita() ([]models.EstadoCita, error) {
	var list []models.EstadoCita
	err := r.db.Find(&list).Error
	return list, err
}
