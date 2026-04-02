package repositories

import (
	"gorm.io/gorm"
)

type PsicologiaRepository struct {
	db *gorm.DB
}

func NewPsicologiaRepository(db *gorm.DB) *PsicologiaRepository {
	return &PsicologiaRepository{db: db}
}

func (r *PsicologiaRepository) GetDB() *gorm.DB {
	return r.db
}

// Aquí irán los métodos de acceso a datos específicos de psicología
// Ejemplo:
// func (r *PsicologiaRepository) FindPacienteByID(id uint) (*models.Paciente, error) {
// 	var paciente models.Paciente
// 	if err := r.db.First(&paciente, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &paciente, nil
// }
