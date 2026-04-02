package services

import (
	"errors"

	"saas-medico/internal/modules/cobros/models"
	"saas-medico/internal/modules/cobros/repositories"
)

var (
	ErrCobroNotFound = errors.New("cobro no encontrado")
	ErrPagoNotFound  = errors.New("pago no encontrado")
)

type CobroService struct {
	repo *repositories.CobroRepository
}

func NewCobroService(repo *repositories.CobroRepository) *CobroService {
	return &CobroService{repo: repo}
}

func (s *CobroService) CreateCobro(req models.CreateCobroRequest) (*models.CobroSesion, error) {
	estado, err := s.repo.FindEstadoCobro(models.CobroPendiente)
	if err != nil {
		return nil, errors.New("estado de cobro inicial no configurado")
	}

	montoTotal := req.MontoCobrar - req.Descuento + req.Recargo

	c := &models.CobroSesion{
		SesionID:      req.SesionID,
		PacienteID:    req.PacienteID,
		MedicoID:      req.MedicoID,
		ClinicaID:     req.ClinicaID,
		MontoCobrar:   req.MontoCobrar,
		Descuento:     req.Descuento,
		Recargo:       req.Recargo,
		MontoTotal:    montoTotal,
		EstadoCobroID: estado.ID,
		Observacion:   req.Observacion,
		State:         "A",
	}
	return c, s.repo.CreateCobro(c)
}

func (s *CobroService) GetCobro(id uint) (*models.CobroSesion, error) {
	c, err := s.repo.FindCobroByID(id)
	if err != nil {
		return nil, ErrCobroNotFound
	}
	return c, nil
}

func (s *CobroService) ListCobrosPaciente(pacienteID uint, page, size int) ([]models.CobroSesion, int64, error) {
	return s.repo.FindCobrosByPaciente(pacienteID, page, size)
}

func (s *CobroService) RegistrarPago(cobroID uint, req models.RegistrarPagoRequest, pacienteID uint) (*models.Pago, error) {
	cobro, err := s.repo.FindCobroByID(cobroID)
	if err != nil {
		return nil, ErrCobroNotFound
	}

	pago := &models.Pago{
		CobroID:     cobroID,
		PacienteID:  pacienteID,
		FechaPago:   req.FechaPago,
		MontoPagado: req.MontoPagado,
		MedioPagoID: req.MedioPagoID,
		Referencia:  req.Referencia,
		Observacion: req.Observacion,
		State:       "A",
	}

	if err := s.repo.CreatePago(pago); err != nil {
		return nil, err
	}

	// Calcular total pagado y actualizar estado del cobro
	pagos, _ := s.repo.FindPagosByCobroID(cobroID)
	var totalPagado float64
	for _, p := range pagos {
		totalPagado += p.MontoPagado
	}

	var nuevoCodigo string
	switch {
	case totalPagado >= cobro.MontoTotal:
		nuevoCodigo = models.CobradoCobrado
	case totalPagado > 0:
		nuevoCodigo = models.CobroParcial
	default:
		nuevoCodigo = models.CobroPendiente
	}

	if estado, err := s.repo.FindEstadoCobro(nuevoCodigo); err == nil {
		_ = s.repo.UpdateEstadoCobro(cobroID, estado.ID)
	}

	return pago, nil
}

func (s *CobroService) CreateEgreso(req models.CreateEgresoRequest) (*models.Egreso, error) {
	e := &models.Egreso{
		ClinicaID:    req.ClinicaID,
		TipoEgresoID: req.TipoEgresoID,
		Fecha:        req.Fecha,
		Monto:        req.Monto,
		Descripcion:  req.Descripcion,
		Proveedor:    req.Proveedor,
		Referencia:   req.Referencia,
		State:        "A",
	}
	return e, s.repo.CreateEgreso(e)
}

func (s *CobroService) ListEgresos(clinicaID uint, page, size int) ([]models.Egreso, int64, error) {
	return s.repo.FindEgresosByClinica(clinicaID, page, size)
}

func (s *CobroService) ListMediosPago() ([]models.MedioPago, error) {
	return s.repo.FindAllMediosPago()
}

func (s *CobroService) ListEstadosCobro() ([]models.EstadoCobro, error) {
	return s.repo.FindAllEstadosCobro()
}

func (s *CobroService) ListTiposEgreso() ([]models.TipoEgreso, error) {
	return s.repo.FindAllTiposEgreso()
}
