package repositories

import (
	"saas-medico/internal/modules/tests/models"

	"gorm.io/gorm"
)

type TestRepository struct {
	db *gorm.DB
}

func NewTestRepository(db *gorm.DB) *TestRepository {
	return &TestRepository{db: db}
}

// ─── Reglas ───────────────────────────────────────────────────────────────────

func (r *TestRepository) FindReglasByFormulario(formularioID uint) ([]models.TestRegla, error) {
	var list []models.TestRegla
	err := r.db.Where("formulario_id = ? AND state = 'A'", formularioID).Find(&list).Error
	return list, err
}

func (r *TestRepository) FindReglaDetalles(reglaID uint) ([]models.TestReglaDetalle, error) {
	var list []models.TestReglaDetalle
	err := r.db.Where("regla_id = ? AND state = 'A'", reglaID).
		Order("orden ASC").Find(&list).Error
	return list, err
}

func (r *TestRepository) CreateRegla(regla *models.TestRegla) error {
	return r.db.Create(regla).Error
}

func (r *TestRepository) CreateReglaDetalle(d *models.TestReglaDetalle) error {
	return r.db.Create(d).Error
}

// ─── Tests ────────────────────────────────────────────────────────────────────

func (r *TestRepository) FindTestsByPaciente(pacienteID uint) ([]models.Test, error) {
	var list []models.Test
	err := r.db.Where("paciente_id = ? AND state = 'A'", pacienteID).
		Order("fecha DESC").Find(&list).Error
	return list, err
}

func (r *TestRepository) FindTestByID(id uint) (*models.Test, error) {
	var t models.Test
	err := r.db.First(&t, "id = ? AND state = 'A'", id).Error
	return &t, err
}

func (r *TestRepository) CreateTest(t *models.Test) error {
	return r.db.Create(t).Error
}

func (r *TestRepository) UpdateTest(t *models.Test) error {
	return r.db.Save(t).Error
}

func (r *TestRepository) CreateRespuestas(respuestas []models.TestRespuesta) error {
	return r.db.Create(&respuestas).Error
}

func (r *TestRepository) FindRespuestasByTest(testID uint) ([]models.TestRespuesta, error) {
	var list []models.TestRespuesta
	err := r.db.Where("test_id = ?", testID).Find(&list).Error
	return list, err
}

// ─── Sesión ↔ Tests ───────────────────────────────────────────────────────────

func (r *TestRepository) LinkTestToSesion(st *models.SesionTest) error {
	return r.db.Create(st).Error
}

func (r *TestRepository) FindTestsBySesion(sesionID uint) ([]models.SesionTest, error) {
	var list []models.SesionTest
	err := r.db.Where("sesion_id = ? AND state = 'A'", sesionID).Find(&list).Error
	return list, err
}
