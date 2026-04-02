package services

import (
	"errors"
	"time"

	"saas-medico/internal/modules/agenda/models"
	"saas-medico/internal/modules/agenda/repositories"
)

var (
	ErrCitaNotFound    = errors.New("cita no encontrada")
	ErrSesionNotFound  = errors.New("sesión no encontrada")
	ErrHorarioNotFound = errors.New("horario no encontrado")
	ErrSesionExiste    = errors.New("ya existe una sesión para esta cita")
)

type AgendaService struct {
	repo *repositories.AgendaRepository
}

func NewAgendaService(repo *repositories.AgendaRepository) *AgendaService {
	return &AgendaService{repo: repo}
}

// ── Citas ─────────────────────────────────────────────────────────────────────

func (s *AgendaService) CreateCita(req models.CreateCitaRequest) (*models.Cita, error) {
	fecha, err := models.ParseFecha(req.Fecha)
	if err != nil {
		return nil, errors.New("formato de fecha inválido, use YYYY-MM-DD")
	}

	// Obtener el estado inicial PE (Pendiente)
	estado, err := s.repo.FindEstadoCitaByCodigo(models.CitaPendiente)
	if err != nil {
		return nil, errors.New("estado de cita inicial no configurado")
	}

	duracion := req.DuracionMin
	if duracion == 0 {
		duracion = 30
	}

	c := &models.Cita{
		Fecha:         fecha,
		Hora:          req.Hora,
		DuracionMin:   duracion,
		MedicoID:      req.MedicoID,
		PacienteID:    req.PacienteID,
		ClinicaID:     req.ClinicaID,
		TipoCitaID:    req.TipoCitaID,
		EstadoCitaID:  estado.ID,
		SucursalID:    req.SucursalID,
		ConsultorioID: req.ConsultorioID,
		PrePacienteID: req.PrePacienteID,
		Motivo:        req.Motivo,
		UrlSesion:     req.UrlSesion,
		State:         "A",
	}
	return c, s.repo.CreateCita(c)
}

func (s *AgendaService) GetCita(id uint) (*models.Cita, error) {
	c, err := s.repo.FindCitaByID(id)
	if err != nil {
		return nil, ErrCitaNotFound
	}
	return c, nil
}

func (s *AgendaService) ListCitas(medicoID, clinicaID uint, fechaStr string, page, size int) ([]models.Cita, int64, error) {
	var fecha *time.Time
	if fechaStr != "" {
		t, err := time.ParseInLocation("2006-01-02", fechaStr, time.Local)
		if err == nil {
			fecha = &t
		}
	}
	return s.repo.FindCitas(medicoID, clinicaID, fecha, page, size)
}

func (s *AgendaService) UpdateEstadoCita(id uint, estadoCodigo string) (*models.Cita, error) {
	c, err := s.repo.FindCitaByID(id)
	if err != nil {
		return nil, ErrCitaNotFound
	}
	estado, err := s.repo.FindEstadoCitaByCodigo(estadoCodigo)
	if err != nil {
		return nil, errors.New("estado de cita no válido")
	}
	if err := s.repo.UpdateEstadoCita(id, estado.ID); err != nil {
		return nil, err
	}
	c.EstadoCitaID = estado.ID
	c.EstadoCita = estado
	return c, nil
}

func (s *AgendaService) UpdateCitaPaciente(id, pacienteID uint) (*models.Cita, error) {
	if _, err := s.repo.FindCitaByID(id); err != nil {
		return nil, ErrCitaNotFound
	}
	if err := s.repo.UpdateCitaPaciente(id, pacienteID); err != nil {
		return nil, err
	}
	return s.repo.FindCitaByID(id)
}

func (s *AgendaService) DeleteCita(id uint) error {
	if _, err := s.repo.FindCitaByID(id); err != nil {
		return ErrCitaNotFound
	}
	return s.repo.DeleteCita(id)
}

// ── Sesiones ──────────────────────────────────────────────────────────────────

func (s *AgendaService) CreateSesion(citaID uint, req models.CreateSesionRequest) (*models.Sesion, error) {
	// Verificar que la cita existe
	cita, err := s.repo.FindCitaByID(citaID)
	if err != nil {
		return nil, ErrCitaNotFound
	}

	// Verificar que no exista ya una sesión para esta cita
	if _, err := s.repo.FindSesionByCitaID(citaID); err == nil {
		return nil, ErrSesionExiste
	}

	sesion := &models.Sesion{
		CitaID:       citaID,
		Inicio:       req.Inicio,
		Fin:          req.Fin,
		Resumen:      req.Resumen,
		Conclusiones: req.Conclusiones,
		State:        "A",
	}

	if err := s.repo.CreateSesion(sesion); err != nil {
		return nil, err
	}

	// Actualizar estado de la cita a Atendida
	estadoAT, _ := s.repo.FindEstadoCitaByCodigo(models.CitaAtendida)
	if estadoAT != nil {
		_ = s.repo.UpdateEstadoCita(cita.ID, estadoAT.ID)
	}

	return sesion, nil
}

func (s *AgendaService) GetSesion(id uint) (*models.Sesion, error) {
	sesion, err := s.repo.FindSesionByID(id)
	if err != nil {
		return nil, ErrSesionNotFound
	}
	return sesion, nil
}

func (s *AgendaService) UpdateSesion(id uint, req models.UpdateSesionRequest) (*models.Sesion, error) {
	sesion, err := s.repo.FindSesionByID(id)
	if err != nil {
		return nil, ErrSesionNotFound
	}
	if req.Fin != nil {
		sesion.Fin = req.Fin
	}
	if req.Resumen != "" {
		sesion.Resumen = req.Resumen
	}
	if req.Conclusiones != "" {
		sesion.Conclusiones = req.Conclusiones
	}
	return sesion, s.repo.UpdateSesion(sesion)
}

// ── Horarios ──────────────────────────────────────────────────────────────────

func (s *AgendaService) CreateHorario(req models.CreateHorarioRequest) (*models.HorarioMedico, error) {
	intervalo := req.IntervaloMin
	if intervalo == 0 {
		intervalo = 30
	}
	h := &models.HorarioMedico{
		MedicoID:      req.MedicoID,
		ClinicaID:     req.ClinicaID,
		ConsultorioID: req.ConsultorioID,
		DiaSemana:     req.DiaSemana,
		HoraInicio:    req.HoraInicio,
		HoraFin:       req.HoraFin,
		IntervaloMin:  intervalo,
		State:         "A",
	}
	return h, s.repo.CreateHorario(h)
}

func (s *AgendaService) ListHorarios(medicoID uint) ([]models.HorarioMedico, error) {
	return s.repo.FindHorariosByMedico(medicoID)
}

func (s *AgendaService) DeleteHorario(id uint) error {
	if _, err := s.repo.FindHorarioByID(id); err != nil {
		return ErrHorarioNotFound
	}
	return s.repo.DeleteHorario(id)
}

// ── Bloqueos ──────────────────────────────────────────────────────────────────

func (s *AgendaService) CreateBloqueo(req models.CreateBloqueoRequest) (*models.BloqueoAgenda, error) {
	b := &models.BloqueoAgenda{
		ClinicaID:     req.ClinicaID,
		SucursalID:    req.SucursalID,
		ConsultorioID: req.ConsultorioID,
		MedicoID:      req.MedicoID,
		FechaInicio:   req.FechaInicio,
		FechaFin:      req.FechaFin,
		Motivo:        req.Motivo,
		TipoBloqueo:   req.TipoBloqueo,
		State:         "A",
	}
	return b, s.repo.CreateBloqueo(b)
}

func (s *AgendaService) ListBloqueos(medicoID uint) ([]models.BloqueoAgenda, error) {
	return s.repo.FindBloqueosByMedico(medicoID)
}

func (s *AgendaService) DeleteBloqueo(id uint) error {
	return s.repo.DeleteBloqueo(id)
}

// ── Catálogos ─────────────────────────────────────────────────────────────────

func (s *AgendaService) ListTiposCita(rolID uint) ([]models.TipoCita, error) {
	if rolID > 0 {
		return s.repo.FindTiposCitaByRol(rolID)
	}
	return s.repo.FindAllTiposCita()
}

func (s *AgendaService) ListEstadosCita() ([]models.EstadoCita, error) {
	return s.repo.FindAllEstadosCita()
}
