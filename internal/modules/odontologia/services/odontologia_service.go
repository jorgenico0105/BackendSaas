package services

import (
	"saas-medico/internal/modules/odontologia/repositories"
)

type OdontologiaService struct {
	repo *repositories.OdontologiaRepository
}

func NewOdontologiaService(repo *repositories.OdontologiaRepository) *OdontologiaService {
	return &OdontologiaService{repo: repo}
}

func (s *OdontologiaService) Ping() string {
	return "pong from odontologia"
}

// Aquí irán los métodos de lógica de negocio específicos de odontología
// Ejemplo:
// func (s *OdontologiaService) GetHistorial(id uint) (*models.HistorialDental, error) {
// 	return s.repo.FindHistorialByID(id)
// }
