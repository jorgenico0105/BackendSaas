package services

import (
	"saas-medico/internal/modules/psicologia/repositories"
)

type PsicologiaService struct {
	repo *repositories.PsicologiaRepository
}

func NewPsicologiaService(repo *repositories.PsicologiaRepository) *PsicologiaService {
	return &PsicologiaService{repo: repo}
}

func (s *PsicologiaService) Ping() string {
	return "pong from psicologia"
}

// Aquí irán los métodos de lógica de negocio específicos de psicología
// Ejemplo:
// func (s *PsicologiaService) GetPaciente(id uint) (*models.Paciente, error) {
// 	return s.repo.FindPacienteByID(id)
// }
