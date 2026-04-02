package repositories

import (
	"gorm.io/gorm"
)

type OdontologiaRepository struct {
	db *gorm.DB
}

func NewOdontologiaRepository(db *gorm.DB) *OdontologiaRepository {
	return &OdontologiaRepository{db: db}
}

func (r *OdontologiaRepository) GetDB() *gorm.DB {
	return r.db
}

// Aquí irán los métodos de acceso a datos específicos de odontología
// Ejemplo:
// func (r *OdontologiaRepository) FindHistorialByID(id uint) (*models.HistorialDental, error) {
// 	var historial models.HistorialDental
// 	if err := r.db.First(&historial, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &historial, nil
// }
