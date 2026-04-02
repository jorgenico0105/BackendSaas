package repositories

import (
	"saas-medico/internal/modules/cobros/models"

	"gorm.io/gorm"
)

type CobroRepository struct {
	db *gorm.DB
}

func NewCobroRepository(db *gorm.DB) *CobroRepository {
	return &CobroRepository{db: db}
}

func (r *CobroRepository) CreateCobro(c *models.CobroSesion) error {
	return r.db.Create(c).Error
}

func (r *CobroRepository) FindCobroByID(id uint) (*models.CobroSesion, error) {
	var c models.CobroSesion
	if err := r.db.Preload("EstadoCobro").Preload("Pagos.MedioPago").
		First(&c, "id = ? AND state = 'A'", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CobroRepository) FindCobrosBySesion(sesionID uint) (*models.CobroSesion, error) {
	var c models.CobroSesion
	if err := r.db.Preload("EstadoCobro").Where("sesion_id = ? AND state = 'A'", sesionID).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CobroRepository) FindCobrosByPaciente(pacienteID uint, page, size int) ([]models.CobroSesion, int64, error) {
	var list []models.CobroSesion
	var total int64
	r.db.Model(&models.CobroSesion{}).Where("id_paciente = ? AND state = 'A'", pacienteID).Count(&total)
	offset := (page - 1) * size
	err := r.db.Preload("EstadoCobro").Where("id_paciente = ? AND state = 'A'", pacienteID).
		Offset(offset).Limit(size).Order("creado_en DESC").Find(&list).Error
	return list, total, err
}

func (r *CobroRepository) UpdateCobro(c *models.CobroSesion) error {
	return r.db.Save(c).Error
}

func (r *CobroRepository) UpdateEstadoCobro(id, estadoID uint) error {
	return r.db.Model(&models.CobroSesion{}).Where("id = ?", id).Update("estado_cobro_id", estadoID).Error
}

func (r *CobroRepository) FindEstadoCobro(codigo string) (*models.EstadoCobro, error) {
	var e models.EstadoCobro
	if err := r.db.Where("codigo = ?", codigo).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *CobroRepository) CreatePago(p *models.Pago) error {
	return r.db.Create(p).Error
}

func (r *CobroRepository) FindPagosByCobroID(cobroID uint) ([]models.Pago, error) {
	var list []models.Pago
	err := r.db.Preload("MedioPago").Where("cobro_id = ? AND state = 'A'", cobroID).Find(&list).Error
	return list, err
}

func (r *CobroRepository) FindAllMediosPago() ([]models.MedioPago, error) {
	var list []models.MedioPago
	err := r.db.Find(&list).Error
	return list, err
}

func (r *CobroRepository) FindAllEstadosCobro() ([]models.EstadoCobro, error) {
	var list []models.EstadoCobro
	err := r.db.Find(&list).Error
	return list, err
}

func (r *CobroRepository) CreateEgreso(e *models.Egreso) error {
	return r.db.Create(e).Error
}

func (r *CobroRepository) FindEgresosByClinica(clinicaID uint, page, size int) ([]models.Egreso, int64, error) {
	var list []models.Egreso
	var total int64
	r.db.Model(&models.Egreso{}).Where("id_clinica = ? AND state = 'A'", clinicaID).Count(&total)
	offset := (page - 1) * size
	err := r.db.Preload("TipoEgreso").Where("id_clinica = ? AND state = 'A'", clinicaID).
		Offset(offset).Limit(size).Order("fecha DESC").Find(&list).Error
	return list, total, err
}

func (r *CobroRepository) FindAllTiposEgreso() ([]models.TipoEgreso, error) {
	var list []models.TipoEgreso
	err := r.db.Find(&list).Error
	return list, err
}
