package repositories

import (
	"log"
	"saas-medico/internal/modules/admin/models"

	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// ── Clinica ──────────────────────────────────────────────────────────────────

func (r *AdminRepository) CreateClinica(c *models.Clinica) error {
	return r.db.Create(c).Error
}

func (r *AdminRepository) FindClinicaByID(id uint) (*models.Clinica, error) {
	var c models.Clinica
	if err := r.db.First(&c, "id = ? AND state = 'A'", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *AdminRepository) FindAllClinicas(page, pageSize int) ([]models.Clinica, int64, error) {
	var list []models.Clinica
	var total int64
	r.db.Model(&models.Clinica{}).Where("state = 'A'").Count(&total)
	offset := (page - 1) * pageSize
	err := r.db.Where("state = 'A'").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *AdminRepository) GetClinicasByUser(id int) ([]models.Clinica, int64, error) {
	var list []models.Clinica
	var total int64
	query := r.db.
		Model(&models.Clinica{}).
		Joins("JOIN usuarios_clinicas uc ON uc.clinica_id = clinicas.id").
		Where("uc.usuario_id = ? AND uc.state = ? AND clinicas.state = ?", id, "A", "A")

	//query := r.db.
	log.Println("[Info del query]", query)

	if err := query.Count(&total).Error; err != nil {
		return list, total, err
	}
	if err := query.Preload("Estilo").Find(&list).Error; err != nil {
		return list, total, err
	}
	return list, total, nil
}
func (r *AdminRepository) FindEstiloByClinica(clinicaID uint) (*models.EstiloClinica, error) {
	var estilo models.EstiloClinica
	err := r.db.
		Where("clinica_id = ? AND es_activo = true AND state = 'A'", clinicaID).
		First(&estilo).Error
	if err != nil {
		return nil, err
	}
	return &estilo, nil
}

func (r *AdminRepository) GetMenuByRolAndClinica(rolID, clinicaID uint) ([]models.Transaccion, error) {
	var list []models.Transaccion
	err := r.db.
		Joins("INNER JOIN rol_transaccion rt ON rt.transaccion_id = transacciones.id AND rt.state = 'A'").
		Where("rt.rol_id = ?", rolID).
		Where("transacciones.state = 'A' AND transacciones.visible = true").
		Where("transacciones.clinica_id = ? OR transacciones.general = true", clinicaID).
		Order("transacciones.orden ASC").
		Find(&list).Error
	return list, err
}

func (r *AdminRepository) UpdateClinica(c *models.Clinica) error {
	return r.db.Save(c).Error
}

func (r *AdminRepository) DeleteClinica(id uint) error {
	return r.db.Model(&models.Clinica{}).Where("id = ?", id).Update("state", "I").Error
}

// ── Sucursal ─────────────────────────────────────────────────────────────────

func (r *AdminRepository) CreateSucursal(s *models.Sucursal) error {
	return r.db.Create(s).Error
}

func (r *AdminRepository) FindSucursalByID(id uint) (*models.Sucursal, error) {
	var s models.Sucursal
	if err := r.db.Preload("Clinica").First(&s, "id = ? AND state = 'A'", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *AdminRepository) FindSucursalesByClinica(clinicaID uint) ([]models.Sucursal, error) {
	var list []models.Sucursal
	err := r.db.Where("clinica_id = ? AND state = 'A'", clinicaID).Find(&list).Error
	return list, err
}

func (r *AdminRepository) UpdateSucursal(s *models.Sucursal) error {
	return r.db.Save(s).Error
}

func (r *AdminRepository) DeleteSucursal(id uint) error {
	return r.db.Model(&models.Sucursal{}).Where("id = ?", id).Update("state", "I").Error
}

// ── Consultorio ───────────────────────────────────────────────────────────────

func (r *AdminRepository) CreateConsultorio(c *models.Consultorio) error {
	return r.db.Create(c).Error
}

func (r *AdminRepository) FindConsultorioByID(id uint) (*models.Consultorio, error) {
	var c models.Consultorio
	if err := r.db.First(&c, "id = ? AND state = 'A'", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *AdminRepository) FindConsultoriosByClinica(clinicaID uint) ([]models.Consultorio, error) {
	var list []models.Consultorio
	err := r.db.Where("clinica_id = ? AND state = 'A'", clinicaID).Find(&list).Error
	return list, err
}

func (r *AdminRepository) UpdateConsultorio(c *models.Consultorio) error {
	return r.db.Save(c).Error
}

func (r *AdminRepository) DeleteConsultorio(id uint) error {
	return r.db.Model(&models.Consultorio{}).Where("id = ?", id).Update("state", "I").Error
}

// ── UsuarioConsultorio ────────────────────────────────────────────────────────

func (r *AdminRepository) AsignarUsuarioAConsultorio(uc *models.UsuarioConsultorio) error {
	return r.db.Create(uc).Error
}

func (r *AdminRepository) FindUsuariosByConsultorio(consultorioID uint) ([]models.UsuarioConsultorio, error) {
	var list []models.UsuarioConsultorio
	err := r.db.Where("consultorio_id = ? AND state = 'A'", consultorioID).Find(&list).Error
	return list, err
}

func (r *AdminRepository) RemoveUsuarioDeConsultorio(usuarioID, consultorioID uint) error {
	return r.db.Model(&models.UsuarioConsultorio{}).
		Where("usuario_id = ? AND consultorio_id = ?", usuarioID, consultorioID).
		Update("state", "I").Error
}

// ── Profesion ─────────────────────────────────────────────────────────────────

func (r *AdminRepository) CreateProfesion(p *models.Profesion) error {
	return r.db.Create(p).Error
}

func (r *AdminRepository) FindAllProfesiones() ([]models.Profesion, error) {
	var list []models.Profesion
	err := r.db.Where("state = 'A'").Find(&list).Error
	return list, err
}

func (r *AdminRepository) FindProfesionByID(id uint) (*models.Profesion, error) {
	var p models.Profesion
	if err := r.db.First(&p, "id = ? AND state = 'A'", id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *AdminRepository) UpdateProfesion(p *models.Profesion) error {
	return r.db.Save(p).Error
}

func (r *AdminRepository) DeleteProfesion(id uint) error {
	return r.db.Model(&models.Profesion{}).Where("id = ?", id).Update("state", "I").Error
}

// ── PlanSaas ──────────────────────────────────────────────────────────────────

func (r *AdminRepository) CreatePlan(p *models.PlanSaas) error {
	return r.db.Create(p).Error
}

func (r *AdminRepository) FindAllPlanes() ([]models.PlanSaas, error) {
	var list []models.PlanSaas
	err := r.db.Where("state = 'A'").Find(&list).Error
	return list, err
}

func (r *AdminRepository) FindPlanByID(id uint) (*models.PlanSaas, error) {
	var p models.PlanSaas
	if err := r.db.First(&p, "id = ? AND state = 'A'", id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *AdminRepository) UpdatePlan(p *models.PlanSaas) error {
	return r.db.Save(p).Error
}

// ── UsuarioClinica ────────────────────────────────────────────────────────────

func (r *AdminRepository) AsignarUsuarioAClinica(uc *models.UsuarioClinica) error {
	return r.db.Create(uc).Error
}

func (r *AdminRepository) FindUsuariosByClinica(clinicaID uint) ([]models.UsuarioClinica, error) {
	var list []models.UsuarioClinica
	err := r.db.Where("clinica_id = ? AND state = 'A'", clinicaID).Find(&list).Error
	return list, err
}

func (r *AdminRepository) RemoveUsuarioDeClinica(usuarioID, clinicaID uint) error {
	return r.db.Model(&models.UsuarioClinica{}).
		Where("usuario_id = ? AND clinica_id = ?", usuarioID, clinicaID).
		Update("state", "I").Error
}

// ── RolTransaccion ────────────────────────────────────────────────────────────

func (r *AdminRepository) FindTransaccionesByRol(rolID uint) ([]models.RolTransaccion, error) {
	var list []models.RolTransaccion
	err := r.db.Where("rol_id = ? AND state = 'A'", rolID).Find(&list).Error
	return list, err
}

func (r *AdminRepository) AsignarTransaccionArol(rt *models.RolTransaccion) error {
	// Upsert: si ya existe (inactiva) la reactiva, si no la crea
	return r.db.Where(models.RolTransaccion{RolID: rt.RolID, TransaccionID: rt.TransaccionID}).
		Assign(map[string]interface{}{"state": "A"}).
		FirstOrCreate(rt).Error
}

func (r *AdminRepository) RevocarTransaccionDeRol(rolID, transaccionID uint) error {
	return r.db.Model(&models.RolTransaccion{}).
		Where("rol_id = ? AND transaccion_id = ?", rolID, transaccionID).
		Update("state", "I").Error
}

func (r *AdminRepository) FindTransacciones(clinicaID *uint) ([]models.Transaccion, error) {
	var list []models.Transaccion
	q := r.db.Where("state = 'A' AND visible = true")
	if clinicaID != nil {
		q = q.Where("clinica_id = ? OR general = true", *clinicaID)
	} else {
		q = q.Where("general = true")
	}
	err := q.Order("orden ASC").Find(&list).Error
	return list, err
}
