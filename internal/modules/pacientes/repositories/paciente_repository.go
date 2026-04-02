package repositories

import (
	"log"
	"saas-medico/internal/modules/pacientes/models"

	"gorm.io/gorm"
)

// DoctorInfo — datos básicos del médico devueltos al paciente tras el login
type DoctorInfo struct {
	ID           uint
	Nombre       string
	Apellidos    string
	Especialidad string
}

type PacienteRepository struct {
	db *gorm.DB
}

func NewPacienteRepository(db *gorm.DB) *PacienteRepository {
	return &PacienteRepository{db: db}
}

func (r *PacienteRepository) Create(p *models.Paciente) error {
	return r.db.Create(p).Error
}

func (r *PacienteRepository) FindByID(id uint) (*models.Paciente, error) {
	var p models.Paciente
	if err := r.db.First(&p, "id = ? AND state = 'A'", id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PacienteRepository) FindAll(search string, page, pageSize int, clinicaId int, usuarioid int) ([]models.Paciente, int64, error) {
	var list []models.Paciente
	var total int64
	q := r.db.Model(&models.Paciente{}).Where("state = 'A'")
	if clinicaId > 0 {
		q = q.Where("clinica_id = ?", clinicaId)
	}
	if usuarioid > 0 {
		q = q.Where("created_by = ?", usuarioid)
	}
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("nombres LIKE ? OR apellidos LIKE ? OR numero_documento LIKE ? OR telefono LIKE ?",
			like, like, like, like)
	}

	q.Count(&total)
	offset := (page - 1) * pageSize
	err := q.Offset(offset).Limit(pageSize).Order("apellidos ASC").Find(&list).Error
	return list, total, err
}

func (r *PacienteRepository) Update(p *models.Paciente) error {
	return r.db.Save(p).Error
}

func (r *PacienteRepository) SoftDelete(id uint) error {
	return r.db.Model(&models.Paciente{}).Where("id = ?", id).Update("state", "I").Error
}

// ── PrePaciente ───────────────────────────────────────────────────────────────

func (r *PacienteRepository) CreatePrePaciente(pp *models.PrePaciente) error {
	return r.db.Create(pp).Error
}

func (r *PacienteRepository) FindPrePacienteByID(id uint) (*models.PrePaciente, error) {
	var pp models.PrePaciente
	if err := r.db.First(&pp, "id = ? AND state = 'A'", id).Error; err != nil {
		return nil, err
	}
	return &pp, nil
}

func (r *PacienteRepository) FindPrePacientesByClinica(clinicaID uint, page, pageSize int) ([]models.PrePaciente, int64, error) {
	var list []models.PrePaciente
	var total int64
	r.db.Model(&models.PrePaciente{}).Where("clinica_id = ? AND state = 'A'", clinicaID).Count(&total)
	offset := (page - 1) * pageSize
	err := r.db.Where("clinica_id = ? AND state = 'A'", clinicaID).Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *PacienteRepository) UpdatePrePaciente(pp *models.PrePaciente) error {
	return r.db.Save(pp).Error
}

func (r *PacienteRepository) DeletePrePaciente(id uint) error {
	return r.db.Model(&models.PrePaciente{}).Where("id = ?", id).Update("state", "I").Error
}

// ── Aplicaciones ──────────────────────────────────────────────────────────────

func (r *PacienteRepository) FindAplicaciones(clinicaID uint) ([]models.Aplicacion, error) {
	log.Println(clinicaID)
	var list []models.Aplicacion

	err := r.db.Where("clinica_id = ? AND state = 'A'", clinicaID).Order("nombre ASC").Find(&list).Error
	log.Println("[INFO APPS]", list)
	return list, err

}

func (r *PacienteRepository) CreateAplicacion(a *models.Aplicacion) error {
	return r.db.Create(a).Error
}

func (r *PacienteRepository) FindAplicacionesByPaciente(pacienteID, clinicaID uint) ([]models.PacienteAplicacion, error) {
	var list []models.PacienteAplicacion
	err := r.db.Where("paciente_id = ? AND clinica_id = ?", pacienteID, clinicaID).Find(&list).Error
	return list, err
}

func (r *PacienteRepository) FindPacienteAplicacion(pacienteID, aplicacionID, clinicaID uint) (*models.PacienteAplicacion, error) {
	var pa models.PacienteAplicacion
	err := r.db.Preload("Aplicacion").Preload("Aplicacion.Medico").Where("paciente_id = ? AND aplicacion_id = ? AND clinica_id = ?",
		pacienteID, aplicacionID, clinicaID).First(&pa).Error
	return &pa, err
}

func (r *PacienteRepository) CreatePacienteAplicacion(pa *models.PacienteAplicacion) error {
	return r.db.Create(pa).Error
}

func (r *PacienteRepository) UpdatePacienteAplicacion(pa *models.PacienteAplicacion) error {
	return r.db.Save(pa).Error
}

// ── PacienteUsuario ───────────────────────────────────────────────────────────

func (r *PacienteRepository) CreatePacienteUsuario(pu *models.PacienteUsuario) error {
	return r.db.Create(pu).Error
}

func (r *PacienteRepository) FindPacienteUsuario(username string, clinicaID uint) (*models.PacienteUsuario, error) {
	var pu models.PacienteUsuario
	if err := r.db.Preload("Paciente").Where("username = ? AND clinica_id = ? AND state = 'A'", username, clinicaID).First(&pu).Error; err != nil {
		return nil, err
	}
	return &pu, nil
}

func (r *PacienteRepository) ExistsPacienteUsuario(pacienteID, clinicaID uint) bool {
	var count int64
	r.db.Model(&models.PacienteUsuario{}).Where("paciente_id = ? AND clinica_id = ?", pacienteID, clinicaID).Count(&count)
	return count > 0
}

// FindAplicacionByID devuelve la aplicación con su medico_id para enriquecer el login
func (r *PacienteRepository) FindAplicacionByID(id uint) (*models.Aplicacion, error) {
	var a models.Aplicacion
	err := r.db.First(&a, "id = ? AND state = 'A'", id).Error
	return &a, err
}

// FindDoctorByID busca los datos básicos del médico/profesional en la tabla de usuarios
func (r *PacienteRepository) FindDoctorByID(userID uint) (*DoctorInfo, error) {
	var result struct {
		ID        uint   `gorm:"column:id"`
		Nombre    string `gorm:"column:nombre"`
		Apellidos string `gorm:"column:apellidos"`
	}
	err := r.db.Table("usuarios").
		Select("id, nombre, apellidos").
		Where("id = ? AND state = 'A'", userID).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &DoctorInfo{
		ID:        result.ID,
		Nombre:    result.Nombre,
		Apellidos: result.Apellidos,
	}, nil
}

// ── PacienteAccesoApp ─────────────────────────────────────────────────────────

func (r *PacienteRepository) RegistrarAccesoApp(pacienteID, clinicaID, aplicacionID uint, tipo string) error {
	acceso := &models.PacienteAccesoApp{
		PacienteID:   pacienteID,
		ClinicaID:    clinicaID,
		AplicacionID: aplicacionID,
		Tipo:         tipo,
	}
	return r.db.Create(acceso).Error
}

func (r *PacienteRepository) FindAccesosByPaciente(pacienteID uint, desde, hasta string) ([]models.PacienteAccesoApp, error) {
	var list []models.PacienteAccesoApp
	q := r.db.Where("paciente_id = ?", pacienteID)
	if desde != "" {
		q = q.Where("DATE(creado_en) >= ?", desde)
	}
	if hasta != "" {
		q = q.Where("DATE(creado_en) <= ?", hasta)
	}
	err := q.Order("creado_en DESC").Find(&list).Error
	return list, err
}

// CountAccesosByPaciente cuenta accesos totales del paciente en un rango de fechas
func (r *PacienteRepository) CountAccesosByPaciente(pacienteID uint, desde, hasta string) int64 {
	var count int64
	q := r.db.Model(&models.PacienteAccesoApp{}).Where("paciente_id = ?", pacienteID)
	if desde != "" {
		q = q.Where("DATE(creado_en) >= ?", desde)
	}
	if hasta != "" {
		q = q.Where("DATE(creado_en) <= ?", hasta)
	}
	q.Count(&count)
	return count
}

// FindUltimoAccesoPaciente devuelve la fecha del último acceso
func (r *PacienteRepository) FindUltimoAccesoPaciente(pacienteID uint) string {
	var acceso models.PacienteAccesoApp
	if err := r.db.Where("paciente_id = ?", pacienteID).Order("creado_en DESC").First(&acceso).Error; err != nil {
		return ""
	}
	return acceso.CreadoEn.Format("2006-01-02T15:04:05Z07:00")
}
