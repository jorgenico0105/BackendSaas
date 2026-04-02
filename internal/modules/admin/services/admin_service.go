package services

import (
	"errors"

	"saas-medico/internal/modules/admin/models"
	"saas-medico/internal/modules/admin/repositories"
)

var (
	ErrClinicaNotFound     = errors.New("clínica no encontrada")
	ErrConsultorioNotFound = errors.New("consultorio no encontrado")
	ErrProfesionNotFound   = errors.New("profesión no encontrada")
	ErrPlanNotFound        = errors.New("plan no encontrado")
)

type AdminService struct {
	repo *repositories.AdminRepository
}

func NewAdminService(repo *repositories.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

// ── Clinica ───────────────────────────────────────────────────────────────────

func (s *AdminService) CreateClinica(req models.CreateClinicaRequest) (*models.Clinica, error) {
	c := &models.Clinica{
		Nombre:             req.Nombre,
		Ruc:                req.Ruc,
		RazonSocial:        req.RazonSocial,
		Direccion:          req.Direccion,
		Ciudad:             req.Ciudad,
		Provincia:          req.Provincia,
		Pais:               req.Pais,
		Telefono:           req.Telefono,
		Correo:             req.Correo,
		SitioWeb:           req.SitioWeb,
		RepresentanteLegal: req.RepresentanteLegal,
		TipoClinica:        req.TipoClinica,
		State:              "A",
	}
	if c.Pais == "" {
		c.Pais = "Ecuador"
	}
	return c, s.repo.CreateClinica(c)
}

func (s *AdminService) GetClinica(id uint) (*models.Clinica, error) {
	c, err := s.repo.FindClinicaByID(id)
	if err != nil {
		return nil, ErrClinicaNotFound
	}
	return c, nil
}
func (s *AdminService) GetClinicasByUser(idUsuario int) ([]models.Clinica, int64, error) {
	return s.repo.GetClinicasByUser(idUsuario)
}
func (s *AdminService) ListClinicas(page, pageSize int) ([]models.Clinica, int64, error) {
	return s.repo.FindAllClinicas(page, pageSize)
}

func (s *AdminService) UpdateClinica(id uint, req models.CreateClinicaRequest) (*models.Clinica, error) {
	c, err := s.repo.FindClinicaByID(id)
	if err != nil {
		return nil, ErrClinicaNotFound
	}
	c.Nombre = req.Nombre
	c.Ruc = req.Ruc
	c.RazonSocial = req.RazonSocial
	c.Direccion = req.Direccion
	c.Ciudad = req.Ciudad
	c.Provincia = req.Provincia
	if req.Pais != "" {
		c.Pais = req.Pais
	}
	c.Telefono = req.Telefono
	c.Correo = req.Correo
	c.SitioWeb = req.SitioWeb
	c.RepresentanteLegal = req.RepresentanteLegal
	c.TipoClinica = req.TipoClinica
	return c, s.repo.UpdateClinica(c)
}

func (s *AdminService) DeleteClinica(id uint) error {
	if _, err := s.repo.FindClinicaByID(id); err != nil {
		return ErrClinicaNotFound
	}
	return s.repo.DeleteClinica(id)
}

// ── Consultorio ───────────────────────────────────────────────────────────────

func (s *AdminService) CreateConsultorio(clinicaID uint, req models.CreateConsultorioRequest) (*models.Consultorio, error) {
	if _, err := s.repo.FindClinicaByID(clinicaID); err != nil {
		return nil, ErrClinicaNotFound
	}
	c := &models.Consultorio{
		ClinicaID:   clinicaID,
		Nombre:      req.Nombre,
		Codigo:      req.Codigo,
		Piso:        req.Piso,
		Descripcion: req.Descripcion,
		State:       "A",
	}
	return c, s.repo.CreateConsultorio(c)
}

func (s *AdminService) GetConsultorio(id uint) (*models.Consultorio, error) {
	c, err := s.repo.FindConsultorioByID(id)
	if err != nil {
		return nil, ErrConsultorioNotFound
	}
	return c, nil
}

func (s *AdminService) ListConsultorios(clinicaID uint) ([]models.Consultorio, error) {
	return s.repo.FindConsultoriosByClinica(clinicaID)
}

func (s *AdminService) UpdateConsultorio(id uint, req models.CreateConsultorioRequest) (*models.Consultorio, error) {
	c, err := s.repo.FindConsultorioByID(id)
	if err != nil {
		return nil, ErrConsultorioNotFound
	}
	c.Nombre = req.Nombre
	c.Codigo = req.Codigo
	c.Piso = req.Piso
	c.Descripcion = req.Descripcion
	return c, s.repo.UpdateConsultorio(c)
}

func (s *AdminService) DeleteConsultorio(id uint) error {
	if _, err := s.repo.FindConsultorioByID(id); err != nil {
		return ErrConsultorioNotFound
	}
	return s.repo.DeleteConsultorio(id)
}

func (s *AdminService) AsignarUsuarioAConsultorio(consultorioID, usuarioID uint) error {
	return s.repo.AsignarUsuarioAConsultorio(&models.UsuarioConsultorio{
		ConsultorioID: consultorioID,
		UsuarioID:     usuarioID,
		State:         "A",
	})
}

func (s *AdminService) ListUsuariosByConsultorio(consultorioID uint) ([]models.UsuarioConsultorio, error) {
	return s.repo.FindUsuariosByConsultorio(consultorioID)
}

func (s *AdminService) RemoverUsuarioDeConsultorio(consultorioID, usuarioID uint) error {
	return s.repo.RemoveUsuarioDeConsultorio(usuarioID, consultorioID)
}

// ── Profesion ─────────────────────────────────────────────────────────────────

func (s *AdminService) CreateProfesion(req models.CreateProfesionRequest) (*models.Profesion, error) {
	p := &models.Profesion{Nombre: req.Nombre, Descripcion: req.Descripcion, State: "A"}
	return p, s.repo.CreateProfesion(p)
}

func (s *AdminService) ListProfesiones() ([]models.Profesion, error) {
	return s.repo.FindAllProfesiones()
}

func (s *AdminService) UpdateProfesion(id uint, req models.CreateProfesionRequest) (*models.Profesion, error) {
	p, err := s.repo.FindProfesionByID(id)
	if err != nil {
		return nil, ErrProfesionNotFound
	}
	p.Nombre = req.Nombre
	p.Descripcion = req.Descripcion
	return p, s.repo.UpdateProfesion(p)
}

func (s *AdminService) DeleteProfesion(id uint) error {
	if _, err := s.repo.FindProfesionByID(id); err != nil {
		return ErrProfesionNotFound
	}
	return s.repo.DeleteProfesion(id)
}

// ── PlanSaas ──────────────────────────────────────────────────────────────────

func (s *AdminService) CreatePlan(req models.CreatePlanSaasRequest) (*models.PlanSaas, error) {
	p := &models.PlanSaas{
		Codigo:        req.Codigo,
		Nombre:        req.Nombre,
		Descripcion:   req.Descripcion,
		PrecioMensual: req.PrecioMensual,
		PrecioAnual:   req.PrecioAnual,
		MaxUsuarios:   req.MaxUsuarios,
		MaxPacientes:  req.MaxPacientes,
		State:         "A",
	}
	return p, s.repo.CreatePlan(p)
}

func (s *AdminService) ListPlanes() ([]models.PlanSaas, error) {
	return s.repo.FindAllPlanes()
}

// ── UsuarioClinica ────────────────────────────────────────────────────────────

func (s *AdminService) AsignarUsuario(clinicaID uint, req models.AsignarUsuarioClinicaRequest) (*models.UsuarioClinica, error) {
	uc := &models.UsuarioClinica{
		UsuarioID: req.UsuarioID,
		ClinicaID: clinicaID,
		RolID:     req.RolID,
		State:     "A",
	}
	return uc, s.repo.AsignarUsuarioAClinica(uc)
}

func (s *AdminService) ListUsuariosClinica(clinicaID uint) ([]models.UsuarioClinica, error) {
	return s.repo.FindUsuariosByClinica(clinicaID)
}

func (s *AdminService) RemoverUsuario(clinicaID, usuarioID uint) error {
	return s.repo.RemoveUsuarioDeClinica(usuarioID, clinicaID)
}
