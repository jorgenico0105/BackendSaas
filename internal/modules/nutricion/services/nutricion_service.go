package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"saas-medico/internal/modules/nutricion/models"
	"saas-medico/internal/modules/nutricion/repositories"
	"saas-medico/internal/shared/pdfbuilder"
	"saas-medico/internal/shared/pdfbuilder/usecases"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrDietaNotFound    = errors.New("dieta no encontrada")
	ErrAlimentoNotFound = errors.New("alimento no encontrado")
	ErrMenuNotFound     = errors.New("menú no encontrado")
	ErrPerfilNotFound   = errors.New("perfil nutricional no encontrado")
	ErrR24HNotFound     = errors.New("recordatorio 24h no encontrado")
)

type NutricionService struct {
	repo  *repositories.NutricionRepository
	redis *redis.Client
}

func NewNutricionService(repo *repositories.NutricionRepository, redis *redis.Client) *NutricionService {
	return &NutricionService{repo: repo, redis: redis}
}

// Get Menu Report — returns the path of the generated PDF
func (s *NutricionService) GetMenuReport(dieta *models.NutricionDietaPaciente, menu *models.NutricionMenu) (string, error) {
	outputPath := fmt.Sprintf("storage/temp/menu_%d.pdf", menu.ID)
	logoPath := "storage/logos/logo_biohealth.png"
	waterMark := "storage/watermarks/watermark.png"
	pdfS := pdfbuilder.NewPdfBuilder(waterMark)
	m, err := pdfS.GeneratePdfBuilder()
	if err != nil {
		return "", err
	}
	useCase := usecases.NewMenuPdfUseCase(dieta, menu, m, logoPath, outputPath)
	if err = useCase.CreatePdf(); err != nil {
		return "", err
	}
	return outputPath, nil
}

// Desactivar menus
func (s *NutricionService) DeactivateOldMenus(ctx context.Context) {
	s.repo.DesactivarMenusAnitguos()
}

// ─── Alimentos ────────────────────────────────────────────────────────────────

func (s *NutricionService) ListAlimentos(categoria string) ([]models.NutricionAlimento, error) {
	return s.repo.FindAlimentos(categoria)
}

func (s *NutricionService) GetAlimento(id uint) (*models.NutricionAlimento, error) {
	return s.repo.FindAlimentoByID(id)
}

func (s *NutricionService) ListTipoComidas() ([]models.NutricionTipoComida, error) {
	return s.repo.FindTipoComida()
}

func (s *NutricionService) CreateAlimento(req models.CreateAlimentoRequest, creadoPor uint) (*models.NutricionAlimento, error) {
	porcion := req.GramosPorcion
	if porcion == 0 {
		porcion = 100
	}
	a := &models.NutricionAlimento{
		Nombre:         req.Nombre,
		Descripcion:    req.Descripcion,
		Categoria:      req.Categoria,
		GramosPorcion:  porcion,
		Calorias:       req.Calorias,
		ProteínasG:     req.ProteínasG,
		CarbohidratosG: req.CarbohidratosG,
		GrasasG:        req.GrasasG,
		FibraG:         req.FibraG,
		AzucaresG:      req.AzucaresG,
		SodioMg:        req.SodioMg,
		Desayuno:       req.Desayuno,
		MediaTardeMana: req.MediaTardeMana,
		Almuerzo:       req.Almuerzo,
		Merienda:       req.Merienda,
		State:          "A",
		CreadoPor:      &creadoPor,
	}
	if err := s.repo.CreateAlimento(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *NutricionService) UpdateAlimento(id uint, req models.UpdateAlimentoRequest) (*models.NutricionAlimento, error) {
	a, err := s.repo.FindAlimentoByID(id)
	if err != nil {
		return nil, err
	}
	porcion := req.GramosPorcion
	if porcion == 0 {
		porcion = 100
	}
	a.Nombre = req.Nombre
	a.Descripcion = req.Descripcion
	a.Categoria = req.Categoria
	a.GramosPorcion = porcion
	a.Calorias = req.Calorias
	a.ProteínasG = req.ProteínasG
	a.CarbohidratosG = req.CarbohidratosG
	a.GrasasG = req.GrasasG
	a.FibraG = req.FibraG
	a.AzucaresG = req.AzucaresG
	a.SodioMg = req.SodioMg
	a.Desayuno = req.Desayuno
	a.MediaTardeMana = req.MediaTardeMana
	a.Almuerzo = req.Almuerzo
	a.Merienda = req.Merienda
	if err := s.repo.UpdateAlimento(a); err != nil {
		return nil, err
	}
	return a, nil
}

// ─── Catálogo dietas ──────────────────────────────────────────────────────────

func (s *NutricionService) ListDietasCatalogo() ([]models.NutricionDietaCatalogo, error) {
	return s.repo.FindDietasCatalogo()
}

// ─── Plan de dieta del paciente ───────────────────────────────────────────────

func (s *NutricionService) ListDietasByPaciente(pacienteID uint) ([]models.NutricionDietaPaciente, error) {
	return s.repo.FindDietasByPaciente(pacienteID)
}

func (s *NutricionService) GetDieta(id uint) (*models.NutricionDietaPaciente, error) {
	return s.repo.FindDietaByID(id)
}

func (s *NutricionService) CreateDieta(pacienteID, medicoID uint, req models.CreateDietaRequest) (*models.NutricionDietaPaciente, error) {
	fechaInicio, err := time.ParseInLocation("2006-01-02", req.FechaInicio, time.Local)

	if err != nil {
		fechaInicio = time.Now()
	}

	duracion := req.DuracionDias
	if duracion == 0 {
		duracion = 7
	}

	numComidas := req.NumComidas
	if numComidas < 3 {
		numComidas = 3
	}
	if numComidas > 5 {
		numComidas = 5
	}

	fechaFin := fechaInicio.AddDate(0, 0, duracion)

	d := &models.NutricionDietaPaciente{
		PacienteID:          pacienteID,
		MedicoID:            medicoID,
		DietaCatalogoID:     req.DietaCatalogoID,
		Nombre:              req.Nombre,
		Descripcion:         req.Descripcion,
		Objetivo:            req.Objetivo,
		PesoObjetivo:        req.ResultadoEsperado,
		FechaInicio:         fechaInicio,
		DuracionDias:        duracion,
		FechaFin:            &fechaFin,
		NumComidas:          numComidas,
		CaloriasDiaObjetivo: req.CaloriasDiaObjetivo,
		ProteínasGDia:       req.ProteínasGDia,
		CarbohidratosGDia:   req.CarbohidratosGDia,
		GrasasGDia:          req.GrasasGDia,
		FibraGDia:           req.FibraGDia,
		Estado:              "ACTIVA",
		State:               "A",
	}

	if err := s.repo.CreateDieta(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *NutricionService) UpdateDieta(id uint, req models.UpdateDietaRequest) (*models.NutricionDietaPaciente, error) {
	d, err := s.repo.FindDietaByID(id)
	if err != nil {
		return nil, ErrDietaNotFound
	}

	if req.Nombre != "" {
		d.Nombre = req.Nombre
	}
	if req.Descripcion != "" {
		d.Descripcion = req.Descripcion
	}
	if req.Objetivo != "" {
		d.Objetivo = req.Objetivo
	}
	if req.ResultadoEsperado != nil {
		d.PesoObjetivo = req.ResultadoEsperado
	}
	if req.Estado != "" {
		d.Estado = req.Estado
	}
	if req.CaloriasDiaObjetivo != nil {
		d.CaloriasDiaObjetivo = req.CaloriasDiaObjetivo
	}
	if req.ProteínasGDia != nil {
		d.ProteínasGDia = req.ProteínasGDia
	}
	if req.CarbohidratosGDia != nil {
		d.CarbohidratosGDia = req.CarbohidratosGDia
	}
	if req.GrasasGDia != nil {
		d.GrasasGDia = req.GrasasGDia
	}
	if req.FibraGDia != nil {
		d.FibraGDia = req.FibraGDia
	}

	if err := s.repo.UpdateDieta(d); err != nil {
		return nil, err
	}
	return d, nil
}

// ─── Menús ────────────────────────────────────────────────────────────────────
func (s *NutricionService) CreateMenu(dietaID, pacienteID uint, req models.CreateMenuRequest) (*models.NutricionMenu, error) {
	fechaInicio, err := time.ParseInLocation("2006-01-02", req.FechaInicio, time.Local)
	if err != nil {
		fechaInicio = time.Now()
	}
	fechaFin := fechaInicio.AddDate(0, 0, 6)
	mpa, err := s.repo.GetMenuPlantilla(dietaID)
	if err != nil {
		return nil, err
	}
	if mpa.MenuID != 0 {
		menu, err := s.repo.FindMenuByID(mpa.MenuID)
		if err != nil {
			return nil, err
		}

		menuToSave := &models.NutricionMenu{
			DietaPacienteID: dietaID,
			SemanaNumero:    req.SemanaNumero,
			Nombre:          req.Nombre,
			Notas:           req.Notas,
			FechaInicio:     fechaInicio,
			FechaFin:        fechaFin,
			Estado:          "ACTIVO",
			State:           "A",
		}

		err = s.repo.CreateMenu(menuToSave)
		if err != nil {
			return nil, err
		}
		type detallePair struct {
			original models.NutricionMenuDetalle
			nuevo    *models.NutricionMenuDetalle
		}

		var pares []detallePair
		var detallesNuevos []*models.NutricionMenuDetalle

		for i := range menu.Detalles {
			src := menu.Detalles[i]

			nuevoDetalle := &models.NutricionMenuDetalle{
				MenuID:              menuToSave.ID,
				TipoComidaID:        src.TipoComidaID,
				DiaNúmero:           src.DiaNúmero,
				NombreComida:        src.NombreComida,
				Instrucciones:       src.Instrucciones,
				CaloriasTotal:       src.CaloriasTotal,
				ProteinasGTotal:     src.ProteinasGTotal,
				CarbohidratosGTotal: src.CarbohidratosGTotal,
				GrasasGTotal:        src.GrasasGTotal,
				State:               src.State,
			}

			detallesNuevos = append(detallesNuevos, nuevoDetalle)
			pares = append(pares, detallePair{
				original: src,
				nuevo:    nuevoDetalle,
			})
		}

		foods, err := s.repo.CreateMenuDetalles(detallesNuevos)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		if len(foods) != len(pares) {
			return nil, errors.New("no coincide la cantidad de detalles creados con la plantilla")
		}

		for i := range foods {
			pares[i].nuevo.ID = foods[i].ID
		}

		var alimentos []*models.NutricionMenuAlimento

		for _, par := range pares {
			for j := range par.original.Alimentos {
				srcAl := par.original.Alimentos[j]

				nuevoAlimento := &models.NutricionMenuAlimento{
					MenuDetalleID:      par.nuevo.ID,
					AlimentoID:         srcAl.AlimentoID,
					GramosAsignados:    srcAl.GramosAsignados,
					CaloriasCalc:       srcAl.CaloriasCalc,
					ProteinasGCalc:     srcAl.ProteinasGCalc,
					CarbohidratosGCalc: srcAl.CarbohidratosGCalc,
					GrasasGCalc:        srcAl.GrasasGCalc,
					Observacion:        srcAl.Observacion,
					State:              srcAl.State,
				}

				alimentos = append(alimentos, nuevoAlimento)
			}
		}

		if len(alimentos) > 0 {
			_, err = s.repo.AddAlimentosToComidas(alimentos)
			if err != nil {
				return nil, err
			}
		}

		err = s.repo.UpdateMenuPlantilla(mpa.DietaPacienteID, menuToSave.ID)
		if err != nil {
			return nil, err
		}

		return menuToSave, nil
	}
	var userAlimento []models.NutricionAlimento
	done := make(chan struct{})

	m := &models.NutricionMenu{
		DietaPacienteID: dietaID,
		SemanaNumero:    req.SemanaNumero,
		Nombre:          req.Nombre,
		Notas:           req.Notas,
		FechaInicio:     fechaInicio,
		FechaFin:        fechaFin,
		Estado:          "ACTIVO",
		State:           "A",
	}

	if err := s.repo.CreateMenu(m); err != nil {
		return nil, err
	}
	mp := &models.NutricionMenuPlantilla{
		MenuID:          m.ID,
		DietaPacienteID: dietaID,
	}
	if err := s.repo.CreateMenuPlantilla(mp); err != nil {
		return nil, err
	}
	var (
		preferences []models.NutricionPreferenciaAlimento
		alimentos   []models.NutricionAlimento
		dieta       *models.NutricionDietaPaciente
		tipoComida  []models.NutricionTipoComida
	)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error
	setErr := func(err error) {
		mu.Lock()
		defer mu.Unlock()
		if firstErr == nil {
			firstErr = err
		}
	}
	wg.Add(4)
	go func() {
		defer wg.Done()
		data, err := s.ListPreferencias(pacienteID)
		if err != nil {
			setErr(err)
			return
		}
		mu.Lock()
		preferences = data
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		data, err := s.ListAlimentos("")
		if err != nil {
			setErr(err)
			return
		}
		mu.Lock()
		alimentos = data
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		data, err := s.GetDieta(dietaID)
		if err != nil {
			setErr(err)
			return
		}
		mu.Lock()
		dieta = data
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		data, err := s.ListTipoComidas()
		if err != nil {
			setErr(err)
			return
		}
		mu.Lock()
		tipoComida = data
		mu.Unlock()
	}()
	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}

	prefMap := make(map[uint]bool)
	for _, pref := range preferences {
		if pref.AlimentoID != nil {
			prefMap[*pref.AlimentoID] = true
		}
	}

	for _, alimento := range alimentos {
		if !prefMap[alimento.ID] {
			userAlimento = append(userAlimento, alimento)
		}
	}

	go s.GenerateMenuFoods(userAlimento, dieta, m, tipoComida, done)
	<-done

	return m, nil
}

func (s *NutricionService) GenerateMenuFoods(
	alimentosUser []models.NutricionAlimento,
	dieta *models.NutricionDietaPaciente,
	menu *models.NutricionMenu,
	tiposComida []models.NutricionTipoComida,
	done chan struct{},
) (*models.NutricionMenu, error) {

	defer close(done)
	doneAl := make(chan struct{})

	var menuDetalles []*models.NutricionMenuDetalle

	menuID := menu.ID
	numeroDeComidas := dieta.NumComidas

	if dieta.CaloriasDiaObjetivo == nil || dieta.ProteínasGDia == nil || dieta.CarbohidratosGDia == nil || dieta.GrasasGDia == nil {
		return nil, errors.New("Error")
	}

	caloriasDieta := *dieta.CaloriasDiaObjetivo
	proteinasDieta := clampMacro(*dieta.ProteínasGDia, gramosProteMin, gramosProteMax)
	carbosGDia := clampMacro(*dieta.CarbohidratosGDia, gramosCarbMin, gramosCarbMax)
	grasasGDia := clampMacro(*dieta.GrasasGDia, gramosGrasaMin, gramosGrasaMax)

	// porcentaje por tipo_comida_id
	foodsMap := make(map[uint]float64)

	switch numeroDeComidas {
	case 3:
		foodsMap[1] = 0.30
		foodsMap[3] = 0.40
		foodsMap[5] = 0.30
	case 4:
		foodsMap[1] = 0.25
		foodsMap[2] = 0.15
		foodsMap[3] = 0.35
		foodsMap[5] = 0.25
	case 5:
		foodsMap[1] = 0.25
		foodsMap[2] = 0.10
		foodsMap[3] = 0.30
		foodsMap[4] = 0.10
		foodsMap[5] = 0.25
	default:
		return nil, nil
	}

	for _, comida := range tiposComida {
		porcion, ok := foodsMap[comida.ID]
		if !ok {
			continue
		}

		for dia := 1; dia <= 7; dia++ {
			proteinasComida := proteinasDieta * porcion
			carbosComida := carbosGDia * porcion
			grasasComida := grasasGDia * porcion
			totalCaloriasComida := caloriasDieta * porcion

			data := &models.NutricionMenuDetalle{
				MenuID:              menuID,
				TipoComidaID:        comida.ID,
				DiaNúmero:           int8(dia),
				NombreComida:        comida.Nombre,
				Instrucciones:       "(Opcional)",
				CaloriasTotal:       &totalCaloriasComida,
				ProteinasGTotal:     &proteinasComida,
				CarbohidratosGTotal: &carbosComida,
				GrasasGTotal:        &grasasComida,
				State:               "A",
			}

			menuDetalles = append(menuDetalles, data)
		}
	}

	foods, err := s.repo.CreateMenuDetalles(menuDetalles)
	if err != nil {
		return nil, err
	}

	go s.addAlimentosToComida(doneAl, foods, alimentosUser)

	<-doneAl

	return menu, nil
}

type AlimentosPorMomento struct {
	Desayuno map[uint][]*models.NutricionAlimento
	Media    map[uint][]*models.NutricionAlimento
	Almuerzo map[uint][]*models.NutricionAlimento
	Merienda map[uint][]*models.NutricionAlimento
}

type AlimentoConGramos struct {
	Alimento *models.NutricionAlimento
	Gramos   float64
	Tipo     string // "P", "C", "G", "F"
}

func (s *NutricionService) addAlimentosToComida(
	doneA chan struct{},
	foods []*models.NutricionMenuDetalle,
	alimentosUser []models.NutricionAlimento,
) error {
	defer close(doneA)
	var alimnetosToSave []*models.NutricionMenuAlimento
	requerimientos, err := s.repo.GetRequerimientosPorComida()
	if err != nil {
		return err
	}

	foodDay := make(map[int][]*models.NutricionMenuDetalle)
	alimentoMenuDetalle := make(map[int][]*models.NutricionAlimento)

	for _, food := range foods {
		dia := int(food.DiaNúmero)
		foodDay[dia] = append(foodDay[dia], food)
	}

	agrupados := agruparAlimentosPorMomento(alimentosUser)

	for _, comidas := range foodDay {
		for _, c := range comidas {
			if c.TipoComidaID == 1 {
				for _, v := range requerimientos[c.TipoComidaID] {
					val := agrupados.Desayuno[uint(v)]
					if len(val) == 0 {
						log.Printf("sin alimentos para desayuno, grupo %d, comida detalle %d", v, c.ID)
						continue
					}
					alimento := randomAlimento(val)
					if alimento != nil {
						alimentoMenuDetalle[int(c.ID)] = append(alimentoMenuDetalle[int(c.ID)], alimento)
					}
				}
			}

			if c.TipoComidaID == 2 || c.TipoComidaID == 4 {
				gruposUsados := make(map[uint]bool)
				for _, v := range requerimientos[c.TipoComidaID] {

					if len(alimentoMenuDetalle[int(c.ID)]) >= 2 {
						break
					}

					if gruposUsados[uint(v)] {
						continue
					}

					val := agrupados.Media[uint(v)]
					if len(val) == 0 {
						log.Printf("sin alimentos para media y tarde, grupo %d, comida detalle %d", v, c.ID)
						continue
					}
					alimento := randomAlimento(val)
					if alimento != nil {
						alimentoMenuDetalle[int(c.ID)] = append(alimentoMenuDetalle[int(c.ID)], alimento)
						gruposUsados[uint(v)] = true
					}
				}
			}

			if c.TipoComidaID == 3 {
				for _, v := range requerimientos[c.TipoComidaID] {
					val := agrupados.Almuerzo[uint(v)]
					if len(val) == 0 {
						log.Printf("sin alimentos para almuerzo, grupo %d, comida detalle %d", v, c.ID)
						continue
					}
					alimento := randomAlimento(val)
					if alimento != nil {
						alimentoMenuDetalle[int(c.ID)] = append(alimentoMenuDetalle[int(c.ID)], alimento)
					}
				}
			}

			if c.TipoComidaID == 5 {
				for _, v := range requerimientos[c.TipoComidaID] {
					val := agrupados.Merienda[uint(v)]
					if len(val) == 0 {
						log.Printf("sin alimentos para merienda, grupo %d, comida detalle %d", v, c.ID)
						continue
					}
					alimento := randomAlimento(val)
					if alimento != nil {
						alimentoMenuDetalle[int(c.ID)] = append(alimentoMenuDetalle[int(c.ID)], alimento)
					}
				}
			}
		}
	}
	factorToleranciaGramos := 1.5
	//var alimentosConGramos []AlimentoConGramos
	//alimentoConGramos := make(map[uint]AlimentoConGramos)
	for _, food := range foods {
		objProteinas := 0.0
		objGrasas := 0.0
		objCarbos := 0.0
		objCalorias := 0.0
		if food.CaloriasTotal != nil {
			objCalorias = *food.CaloriasTotal
		}
		if food.ProteinasGTotal != nil {
			objProteinas = *food.ProteinasGTotal
		}
		if food.GrasasGTotal != nil {
			objGrasas = *food.GrasasGTotal
		}
		if food.CarbohidratosGTotal != nil {
			objCarbos = *food.CarbohidratosGTotal
		}

		alimentos := alimentoMenuDetalle[int(food.ID)]
		//log.Println("[Detalle por comida]", food.ID, objProteinas, objCarbos, objGrasas)
		// Las calorías se derivan de los macros (P*4 + C*4 + G*9); pasar 0 omite
		// la fase de ajuste calórico que sobreescribía los macros ya convergidos.
		itemsConGramos := calcularGramosPorAlimento(
			alimentos,
			objProteinas,
			objCarbos,
			objGrasas,
			objCalorias,
			factorToleranciaGramos,
		)

		for _, item := range itemsConGramos {
			p, c, g, cal := macrosDeAlimento(item.Alimento, item.Gramos)

			alimnetosToSave = append(alimnetosToSave, &models.NutricionMenuAlimento{
				MenuDetalleID:      food.ID,
				AlimentoID:         item.Alimento.ID,
				GramosAsignados:    item.Gramos,
				CaloriasCalc:       &cal,
				ProteinasGCalc:     &p,
				CarbohidratosGCalc: &c,
				GrasasGCalc:        &g,
			})
		}
	}
	if len(alimnetosToSave) > 0 {
		if _, err := s.repo.AddAlimentosToComidas(alimnetosToSave); err != nil {
			return err
		}
	}
	return nil
}

func randomAlimento(alimentos []*models.NutricionAlimento) *models.NutricionAlimento {
	if len(alimentos) == 0 {
		return nil
	}

	index := rand.Intn(len(alimentos))
	return alimentos[index]
}

const (
	gramosPorcionRef = 100.0
	maxIteraciones   = 40

	// tolerancias por comida
	toleranciaProte = 4.0
	toleranciaCarb  = 6.0
	toleranciaGrasa = 2.0

	toleranciaCalorias = 25.0 // ±25 kcal por comida

	// rangos por tipo de alimento
	proteMin = 80.0
	proteMax = 220.0

	carbMin = 60.0
	carbMax = 260.0

	grasaMin = 5.0
	grasaMax = 30.0

	fibraMin = 30.0
	fibraMax = 300.0

	// límites diarios que ya usas al crear objetivos
	gramosProteMin = 80.0
	gramosProteMax = 200.0
	gramosCarbMin  = 100.0
	gramosCarbMax  = 300.0
	gramosGrasaMin = 5.0
	gramosGrasaMax = 30.0
)

func macrosDeAlimento(a *models.NutricionAlimento, gramos float64) (proteinas, carbos, grasas, calorias float64) {
	factor := gramos / gramosPorcionRef
	return a.ProteínasG * factor, a.CarbohidratosG * factor, a.GrasasG * factor, a.Calorias * factor
}

func sumarMacros(items []AlimentoConGramos) (proteinas, carbos, grasas, calorias float64) {
	for _, item := range items {
		p, c, g, cal := macrosDeAlimento(item.Alimento, item.Gramos)
		proteinas += p
		carbos += c
		grasas += g
		calorias += cal
	}
	return
}

func calcularGramosPorAlimento(
	alimentos []*models.NutricionAlimento,
	objProteinas, objCarbos, objGrasas float64,
	objCalorias float64,
	_ float64, // se ignora la tolerancia vieja para no romper la firma
) []AlimentoConGramos {
	if len(alimentos) == 0 {
		return nil
	}

	resultado := make([]AlimentoConGramos, 0, len(alimentos))

	// 1) asignación inicial por tipo real del alimento
	for _, a := range alimentos {
		tipo := detectarTipoAlimento(a)
		gramos := gramosInicialesPorTipo(a, objProteinas, objCarbos, objGrasas)

		resultado = append(resultado, AlimentoConGramos{
			Alimento: a,
			Gramos:   gramos,
			Tipo:     tipo,
		})
	}

	// 2) ajuste iterativo por macro, respetando el tipo de alimento
	for iter := 0; iter < maxIteraciones; iter++ {
		totalP, totalC, totalG, totalCal := sumarMacros(resultado)

		diffP := objProteinas - totalP
		diffC := objCarbos - totalC
		diffG := objGrasas - totalG
		diffCal := objCalorias - totalCal

		okP := math.Abs(diffP) <= toleranciaProte
		okC := math.Abs(diffC) <= toleranciaCarb
		okG := math.Abs(diffG) <= toleranciaGrasa
		okCal := objCalorias <= 0 || math.Abs(diffCal) <= toleranciaCalorias

		if okP && okC && okG && okCal {
			break
		}

		if !okP {
			ajustarPorTipo(&resultado, "P", diffP, func(a *models.NutricionAlimento) float64 {
				return a.ProteínasG
			})
		}

		if !okC {
			ajustarPorTipo(&resultado, "C", diffC, func(a *models.NutricionAlimento) float64 {
				return a.CarbohidratosG
			})
		}

		if !okG {
			ajustarPorTipo(&resultado, "G", diffG, func(a *models.NutricionAlimento) float64 {
				return a.GrasasG
			})
		}

		if !okCal && objCalorias > 0 {
			ajustarCalorias(&resultado, diffCal)
		}
	}

	return resultado
}
func clampMacro(v, minVal, maxVal float64) float64 {
	if v < minVal {
		return minVal
	}
	if v > maxVal {
		return maxVal
	}
	return v
}
func detectarTipoAlimento(a *models.NutricionAlimento) string {
	p := a.ProteínasG
	c := a.CarbohidratosG
	g := a.GrasasG

	if g >= p && g >= c && g >= 8 {
		return "G"
	}

	if p >= c && p >= g && p >= 8 {
		return "P"
	}

	if c >= p && c >= g && c >= 8 {
		return "C"
	}

	return "F"
}

func clampRango(v, minVal, maxVal float64) float64 {
	if v < minVal {
		return minVal
	}
	if v > maxVal {
		return maxVal
	}
	return v
}

func clampGramosPorTipo(tipo string, gramos float64) float64 {
	switch tipo {
	case "P":
		return clampRango(gramos, proteMin, proteMax)
	case "C":
		return clampRango(gramos, carbMin, carbMax)
	case "G":
		return clampRango(gramos, grasaMin, grasaMax)
	default:
		return clampRango(gramos, fibraMin, fibraMax)
	}
}

func gramosInicialesPorTipo(a *models.NutricionAlimento, objP, objC, objG float64) float64 {
	tipo := detectarTipoAlimento(a)

	switch tipo {
	case "P":
		if a.ProteínasG <= 0 {
			return proteMin
		}
		return clampGramosPorTipo(tipo, (objP/a.ProteínasG)*gramosPorcionRef)

	case "C":
		if a.CarbohidratosG <= 0 {
			return carbMin
		}
		return clampGramosPorTipo(tipo, (objC/a.CarbohidratosG)*gramosPorcionRef)

	case "G":
		if a.GrasasG <= 0 {
			return grasaMin
		}
		return clampGramosPorTipo(tipo, (objG/a.GrasasG)*gramosPorcionRef)

	default:
		// vegetales/frutas: porción moderada base
		return clampGramosPorTipo(tipo, 80.0)
	}
}

func ajustarPorTipo(
	items *[]AlimentoConGramos,
	tipo string,
	diff float64,
	aporteFn func(*models.NutricionAlimento) float64,
) {
	bestIdx := -1
	bestAporte := -1.0

	for i, item := range *items {
		if item.Tipo != tipo {
			continue
		}

		ap := aporteFn(item.Alimento)
		if ap > bestAporte {
			bestAporte = ap
			bestIdx = i
		}
	}

	if bestIdx < 0 || bestAporte <= 0 {
		return
	}

	delta := (diff / bestAporte) * gramosPorcionRef
	nuevo := (*items)[bestIdx].Gramos + delta
	(*items)[bestIdx].Gramos = clampGramosPorTipo(tipo, nuevo)
}

func existeTipo(items []AlimentoConGramos, tipo string) bool {
	for _, item := range items {
		if item.Tipo == tipo {
			return true
		}
	}
	return false
}

func ajustarCalorias(items *[]AlimentoConGramos, diffCal float64) {
	// primero intenta ajustar carbos
	if existeTipo(*items, "C") {
		bestIdx := -1
		bestCal := -1.0

		for i, item := range *items {
			if item.Tipo != "C" {
				continue
			}
			if item.Alimento.Calorias > bestCal {
				bestCal = item.Alimento.Calorias
				bestIdx = i
			}
		}

		if bestIdx >= 0 && bestCal > 0 {
			delta := (diffCal / bestCal) * gramosPorcionRef
			nuevo := (*items)[bestIdx].Gramos + delta
			(*items)[bestIdx].Gramos = clampGramosPorTipo("C", nuevo)
			return
		}
	}

	// si no hay carbos, intenta con grasas
	if existeTipo(*items, "G") {
		bestIdx := -1
		bestCal := -1.0

		for i, item := range *items {
			if item.Tipo != "G" {
				continue
			}
			if item.Alimento.Calorias > bestCal {
				bestCal = item.Alimento.Calorias
				bestIdx = i
			}
		}

		if bestIdx >= 0 && bestCal > 0 {
			delta := (diffCal / bestCal) * gramosPorcionRef
			nuevo := (*items)[bestIdx].Gramos + delta
			(*items)[bestIdx].Gramos = clampGramosPorTipo("G", nuevo)
		}
	}
}

func agruparAlimentosPorMomento(alimentos []models.NutricionAlimento) AlimentosPorMomento {
	result := AlimentosPorMomento{
		Desayuno: make(map[uint][]*models.NutricionAlimento),
		Media:    make(map[uint][]*models.NutricionAlimento),
		Almuerzo: make(map[uint][]*models.NutricionAlimento),
		Merienda: make(map[uint][]*models.NutricionAlimento),
	}

	for i := range alimentos {
		alimento := &alimentos[i]

		if alimento.GrupoID == nil {
			continue
		}

		grupoID := *alimento.GrupoID

		if alimento.Desayuno {
			result.Desayuno[grupoID] = append(result.Desayuno[grupoID], alimento)
		}
		if alimento.MediaTardeMana {
			result.Media[grupoID] = append(result.Media[grupoID], alimento)
		}
		if alimento.Almuerzo {
			result.Almuerzo[grupoID] = append(result.Almuerzo[grupoID], alimento)
		}
		if alimento.Merienda {
			result.Merienda[grupoID] = append(result.Merienda[grupoID], alimento)
		}
	}

	return result
}

func (s *NutricionService) ListMenusByDieta(dietaID uint) ([]models.NutricionMenu, error) {
	return s.repo.FindMenusByDieta(dietaID)
}

func (s *NutricionService) GetMenu(id uint) (*models.NutricionMenu, error) {
	m, err := s.repo.FindMenuByID(id)
	if err != nil {
		return nil, ErrMenuNotFound
	}
	return m, nil
}

func (s *NutricionService) GetDetallesMenu(menuID uint) ([]models.NutricionMenuDetalle, error) {
	return s.repo.FindDetallesByMenu(menuID)
}

func (s *NutricionService) AddDetalleMenu(menuID uint, req models.AddDetalleMenuRequest) (*models.NutricionMenuDetalle, error) {
	d := &models.NutricionMenuDetalle{
		MenuID:        menuID,
		TipoComidaID:  req.TipoComidaID,
		DiaNúmero:     req.DiaNúmero,
		NombreComida:  req.NombreComida,
		Instrucciones: req.Instrucciones,
		State:         "A",
	}
	if err := s.repo.CreateDetalle(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *NutricionService) UpdateDetalleReceta(detalleID uint, receta string) (*models.NutricionMenuDetalle, error) {
	return s.repo.UpdateDetalleReceta(detalleID, receta)
}

func (s *NutricionService) GetAlimentosMenuDetalle(detalleID uint) ([]models.NutricionMenuAlimento, error) {
	return s.repo.FindAlimentosByMenuDetalle(detalleID)
}

func (s *NutricionService) DeleteAlimentoMenuDetalle(id uint) error {
	return s.repo.DeleteMenuAlimento(id)
}

func (s *NutricionService) UpdateAlimentoMenuDetalle(id uint, gramos float64) (*models.NutricionMenuAlimento, error) {
	a, err := s.repo.FindMenuAlimentoByID(id)
	if err != nil {
		return nil, ErrAlimentoNotFound
	}
	alimento := a.Alimento
	porcion := alimento.GramosPorcion
	if porcion == 0 {
		porcion = 100
	}
	ratio := gramos / porcion
	cal := round2(alimento.Calorias * ratio)
	prot := round2(alimento.ProteínasG * ratio)
	carb := round2(alimento.CarbohidratosG * ratio)
	gras := round2(alimento.GrasasG * ratio)
	if err := s.repo.UpdateMenuAlimentoGramos(id, gramos, cal, prot, carb, gras); err != nil {
		return nil, err
	}
	a.GramosAsignados = gramos
	a.CaloriasCalc = &cal
	a.ProteinasGCalc = &prot
	a.CarbohidratosGCalc = &carb
	a.GrasasGCalc = &gras
	return a, nil
}

func (s *NutricionService) AddAlimentoMenuDetalle(detalleID, alimentoID uint, req models.AddAlimentoMenuRequest) (*models.NutricionMenuAlimento, error) {
	alimento, err := s.repo.FindAlimentoByID(alimentoID)
	if err != nil {
		return nil, ErrAlimentoNotFound
	}

	gramos := req.GramosAsignados
	porcion := alimento.GramosPorcion
	if porcion == 0 {
		porcion = 100
	}
	ratio := gramos / porcion

	cal := round2(alimento.Calorias * ratio)
	prot := round2(alimento.ProteínasG * ratio)
	carb := round2(alimento.CarbohidratosG * ratio)
	gras := round2(alimento.GrasasG * ratio)

	a := &models.NutricionMenuAlimento{
		MenuDetalleID:      detalleID,
		AlimentoID:         alimentoID,
		GramosAsignados:    gramos,
		CaloriasCalc:       &cal,
		ProteinasGCalc:     &prot,
		CarbohidratosGCalc: &carb,
		GrasasGCalc:        &gras,
		Observacion:        req.Observacion,
		State:              "A",
	}
	if err := s.repo.CreateMenuAlimento(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *NutricionService) ListDietasRequierenCambio() ([]models.NutricionMenu, error) {
	return s.repo.FindMenusRequierenCambio()
}

// ─── Fórmulas nutricionales (cálculo puro, sin persistencia) ─────────────────

// CalcularFormulas calcula IMC, ICC y Harris-Benedict a partir de los datos físicos
// del request. No persiste nada — los resultados se guardan en historia clínica.
func (s *NutricionService) CalcularFormulas(req models.CalcularFormulasRequest) *models.NutricionFormulasResult {
	result := &models.NutricionFormulasResult{}

	if req.AlturaCm != nil && req.PesoKg != nil {
		alturaM := *req.AlturaCm / 100.0
		imcVal := round2(*req.PesoKg / (alturaM * alturaM))
		result.IMC = &imcVal

		switch {
		case imcVal < 18.5:
			result.ClasificacionIMC = "Bajo peso"
		case imcVal < 25.0:
			result.ClasificacionIMC = "Normal"
		case imcVal < 30.0:
			result.ClasificacionIMC = "Sobrepeso"
		default:
			result.ClasificacionIMC = "Obesidad"
		}
	}

	if req.CinturaCm != nil && req.CaderaCm != nil && *req.CaderaCm > 0 {
		iccVal := round2(*req.CinturaCm / *req.CaderaCm)
		result.ICC = &iccVal

		var riesgo string
		if req.Sexo == "M" {
			switch {
			case iccVal < 0.9:
				riesgo = "BAJO"
			case iccVal <= 1.0:
				riesgo = "MODERADO"
			default:
				riesgo = "ALTO"
			}
		} else {
			switch {
			case iccVal < 0.8:
				riesgo = "BAJO"
			case iccVal <= 0.85:
				riesgo = "MODERADO"
			default:
				riesgo = "ALTO"
			}
		}
		result.RiesgoMetabolico = riesgo
	}

	if req.AlturaCm != nil && req.PesoKg != nil && req.EdadAnos != nil {
		peso := *req.PesoKg
		altura := *req.AlturaCm
		edad := float64(*req.EdadAnos)

		var tmb float64
		if req.Sexo == "M" {
			tmb = 66.5 + (13.75 * peso) + (5.003 * altura) - (6.75 * edad)
		} else {
			tmb = 655.1 + (9.563 * peso) + (1.850 * altura) - (4.676 * edad)
		}
		tmb = round2(tmb)
		result.TMB = &tmb
		result.GEB = &tmb

		factorActividad := 1.2
		if req.FactorActividad != nil {
			factorActividad = *req.FactorActividad
		}
		get := round2(tmb * factorActividad)
		result.GET = &get
	}

	return result
}

// ─── R24H ─────────────────────────────────────────────────────────────────────

func (s *NutricionService) CreateR24H(pacienteID, medicoID uint, req models.CreateR24HRequest) (*models.NutricionR24H, error) {
	fecha, err := time.ParseInLocation("2006-01-02", req.Fecha, time.Local)
	if err != nil {
		fecha = time.Now()
	}

	r := &models.NutricionR24H{
		PacienteID:    pacienteID,
		MedicoID:      medicoID,
		Fecha:         fecha,
		Observaciones: req.Observaciones,
		State:         "A",
	}
	if err := s.repo.CreateR24H(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *NutricionService) ListR24H(pacienteID uint) ([]models.NutricionR24H, error) {
	return s.repo.FindR24HByPaciente(pacienteID)
}

func (s *NutricionService) AddR24HItem(r24hID uint, req models.AddR24HItemRequest) (*models.NutricionR24HItem, error) {
	item := &models.NutricionR24HItem{
		R24HID:      r24hID,
		HoraAprox:   req.HoraAprox,
		TipoComida:  req.TipoComida,
		Alimento:    req.Alimento,
		Cantidad:    req.Cantidad,
		CaloriasEst: req.CaloriasEst,
		Notas:       req.Notas,
		State:       "A",
	}
	if err := s.repo.CreateR24HItem(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *NutricionService) ListR24HItems(r24hID uint) ([]models.NutricionR24HItem, error) {
	return s.repo.FindR24HItems(r24hID)
}

// ─── Preferencias ─────────────────────────────────────────────────────────────

func (s *NutricionService) AddPreferencia(pacienteID uint, req models.CreatePreferenciaRequest) (*models.NutricionPreferenciaAlimento, error) {
	p := &models.NutricionPreferenciaAlimento{
		PacienteID:  pacienteID,
		AlimentoID:  req.AlimentoID,
		NombreLibre: req.NombreLibre,
		Tipo:        req.Tipo,
		Notas:       req.Notas,
		State:       "A",
	}
	if err := s.repo.CreatePreferencia(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *NutricionService) ListPreferencias(pacienteID uint) ([]models.NutricionPreferenciaAlimento, error) {
	return s.repo.FindPreferenciasByPaciente(pacienteID)
}

func (s *NutricionService) DeletePreferencia(id uint) error {
	return s.repo.DeletePreferencia(id)
}

// ─── Síntomas ─────────────────────────────────────────────────────────────────

func (s *NutricionService) CreateSintoma(pacienteID uint, req models.CreateSintomaRequest) (*models.NutricionSintoma, error) {
	fecha, err := time.ParseInLocation("2006-01-02", req.Fecha, time.Local)
	if err != nil {
		fecha = time.Now()
	}

	sint := &models.NutricionSintoma{
		PacienteID:      pacienteID,
		Fecha:           fecha,
		Descripcion:     req.Descripcion,
		Tipo:            req.Tipo,
		Intensidad:      req.Intensidad,
		AlimentoPosible: req.AlimentoPosible,
		State:           "A",
	}
	if err := s.repo.CreateSintoma(sint); err != nil {
		return nil, err
	}
	return sint, nil
}

func (s *NutricionService) ListSintomas(pacienteID uint, fechaDesde, fechaHasta string) ([]models.NutricionSintoma, error) {
	return s.repo.FindSintomasByPaciente(pacienteID, fechaDesde, fechaHasta)
}

// ─── Ejercicios ───────────────────────────────────────────────────────────────

func (s *NutricionService) ListEjerciciosCatalogo(categoria string) ([]models.NutricionEjercicioCatalogo, error) {
	return s.repo.FindEjerciciosCatalogo(categoria)
}

func (s *NutricionService) CreateEjercicioCatalogo(req models.CreateEjercicioCatalogoRequest, creadoPor uint) (*models.NutricionEjercicioCatalogo, error) {
	e := &models.NutricionEjercicioCatalogo{
		Nombre:          req.Nombre,
		Descripcion:     req.Descripcion,
		Categoria:       req.Categoria,
		GrupoMuscular:   req.GrupoMuscular,
		CaloriasPorHora: req.CaloriasPorHora,
		UnidadMedida:    req.UnidadMedida,
		Nivel:           req.Nivel,
		State:           "A",
		CreadoPor:       &creadoPor,
	}
	if err := s.repo.CreateEjercicioCatalogo(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *NutricionService) ListEjerciciosByPaciente(pacienteID uint) ([]models.NutricionEjercicioPaciente, error) {
	return s.repo.FindEjerciciosByPaciente(pacienteID)
}

func (s *NutricionService) AddEjercicioPaciente(pacienteID, medicoID uint, req models.CreateEjercicioPacienteRequest) (*models.NutricionEjercicioPaciente, error) {
	e := &models.NutricionEjercicioPaciente{
		PacienteID:      pacienteID,
		MedicoID:        medicoID,
		DietaPacienteID: req.DietaPacienteID,
		EjercicioID:     req.EjercicioID,
		DiaNúmero:       req.DiaNúmero,
		DiaSemana:       req.DiaSemana,
		DuracionMin:     req.DuracionMin,
		Series:          req.Series,
		Repeticiones:    req.Repeticiones,
		Instrucciones:   req.Instrucciones,
		Estado:          "PENDIENTE",
		State:           "A",
	}
	if err := s.repo.CreateEjercicioPaciente(e); err != nil {
		return nil, err
	}
	return e, nil
}

// ─── Registros comida ─────────────────────────────────────────────────────────

func (s *NutricionService) ListRegistrosComida(pacienteID uint, fecha, desde, hasta string) ([]models.NutricionRegistroComida, error) {
	return s.repo.FindRegistrosComida(pacienteID, fecha, desde, hasta)
}

func (s *NutricionService) CreateRegistroComida(pacienteID uint, req models.CreateRegistroComidaRequest) (*models.NutricionRegistroComida, error) {
	fecha, err := time.ParseInLocation("2006-01-02", req.Fecha, time.Local)
	if err != nil {
		fecha = time.Now()
	}

	// Si viene de un detalle de menú, verificar que no xista ya un registro consumido
	if req.MenuDetalleID != nil {
		fechaStr := fecha.Format("2006-01-02")
		if existing, err := s.repo.FindRegistroComidaByMenuDetalle(pacienteID, *req.MenuDetalleID, fechaStr); err == nil && existing != nil {
			return existing, nil // ya registrado, devolver el existente
		}
	}

	calConsumidas := req.CaloriasConsumidas
	if calConsumidas == nil && req.MenuDetalleID != nil {
		if detalle, err := s.repo.FindMenuDetalleByID(*req.MenuDetalleID); err == nil && detalle.CaloriasTotal != nil {
			calConsumidas = detalle.CaloriasTotal
		}
	}

	rc := &models.NutricionRegistroComida{
		PacienteID:         pacienteID,
		Fecha:              fecha,
		TipoComidaID:       req.TipoComidaID,
		MenuDetalleID:      req.MenuDetalleID,
		FueraDePlan:        req.FueraDePlan,
		DescripcionLibre:   req.DescripcionLibre,
		CaloriasConsumidas: calConsumidas,
		FotoComida:         req.FotoComida,
		ProteínasG:         req.ProteinasConsumidas,
		GrasasG:            req.GrasasConsumidas,
		CarbohidratosG:     req.CarbosConsumidos,
		Notas:              req.Notas,
		Estado:             models.EstadoRegistroComidaConsumida, // siempre 'C' al crear
		State:              "A",
	}
	if err := s.repo.CreateRegistroComida(rc); err != nil {
		return nil, err
	}

	// Marcar el detalle del menú como consumido (state='C')
	if req.MenuDetalleID != nil {
		_ = s.repo.MarcarMenuDetalleConsumido(*req.MenuDetalleID)
	}

	return rc, nil
}

// MarcarConsumida cambia el Estado de un registro a 'C'
func (s *NutricionService) MarcarConsumida(registroID uint) error {
	return s.repo.MarcarRegistroComidaConsumida(registroID)
}

func (s *NutricionService) UpdateFotoComida(registroID uint, rutaFoto string) (*models.NutricionRegistroComida, string, error) {
	return s.repo.UpdateFotoComida(registroID, rutaFoto)
}

func (s *NutricionService) GetResumenDiario(pacienteID uint, fecha string) (*models.ResumenDiarioResponse, error) {
	if fecha == "" {
		fecha = time.Now().Format("2006-01-02")
	}
	var comidas []models.NutricionRegistroComida
	var ejercicios []models.NutricionRegistroEjercicio
	var progreso *models.NutricionProgresoPaciente
	var dieta *models.NutricionDietaPaciente
	var err error

	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error
	setErr := func(err error) {
		mu.Lock()
		defer mu.Unlock()
		if firstErr == nil {
			firstErr = err
		}
	}
	wg.Add(4)
	go func() {
		defer wg.Done()
		comidas, err = s.repo.FindRegistrosComida(pacienteID, fecha, "", "")
		if err != nil {
			setErr(err)
		}
	}()
	go func() {
		defer wg.Done()
		ejercicios, err = s.repo.FindRegistrosEjercicio(pacienteID, fecha, "", "")
		if err != nil {
			setErr(err)
		}

	}()
	go func() {
		defer wg.Done()
		progreso, _ = s.repo.FindProgresoPorFecha(pacienteID, fecha)
	}()

	var calObjetivo float64
	go func() {
		defer wg.Done()
		dieta, err = s.repo.FindDietaActivaByPaciente(pacienteID)
	}()
	wg.Wait()

	if dieta != nil && dieta.CaloriasDiaObjetivo != nil {
		calObjetivo = *dieta.CaloriasDiaObjetivo
	}
	var calCon, prot, carb, gras, calQuem float64

	var consumidoIDs []uint
	for _, r := range comidas {
		if r.Estado == models.EstadoRegistroComidaConsumida {
			consumidoIDs = append(consumidoIDs, r.ID)
		}
	}

	if calCon == 0 {
		for _, r := range comidas {
			if r.Estado == models.EstadoRegistroComidaConsumida {
				if r.CaloriasConsumidas != nil {
					calCon += *r.CaloriasConsumidas
				}
				if r.ProteínasG != nil {
					prot += *r.ProteínasG
				}
				if r.CarbohidratosG != nil {
					carb += *r.CarbohidratosG
				}
				if r.GrasasG != nil {
					gras += *r.GrasasG
				}
			}
		}
	}

	for _, e := range ejercicios {
		if e.CaloriasQuemadas != nil {
			calQuem += *e.CaloriasQuemadas
		}
	}

	pct := 0
	if calObjetivo > 0 && calCon > 0 {
		pct = int(math.Round((calCon / calObjetivo) * 100))
		if pct > 100 {
			pct = 100
		}
	}

	return &models.ResumenDiarioResponse{
		Fecha:                  fecha,
		CaloriasObjetivo:       calObjetivo,
		CaloriasConsumidas:     calCon,
		CaloriasQuemadas:       calQuem,
		ProteinasG:             prot,
		CarbohidratosG:         carb,
		GrasasG:                gras,
		PorcentajeCumplimiento: pct,
		RegistrosComida:        comidas,
		RegistrosEjercicio:     ejercicios,
		Progreso:               progreso,
	}, nil
}

func (s *NutricionService) AddRegistroAlimento(registroID uint, req models.AddRegistroAlimentoRequest) (*models.NutricionRegistroAlimento, error) {
	ra := &models.NutricionRegistroAlimento{
		RegistroComidaID: registroID,
		AlimentoID:       req.AlimentoID,
		NombreLibre:      req.NombreLibre,
		GramosConsumidos: req.GramosConsumidos,
		State:            "A",
	}

	// Calcular macros reales a partir del catálogo de alimentos
	if req.AlimentoID != nil {
		if alimento, err := s.repo.FindAlimentoByID(*req.AlimentoID); err == nil && alimento.GramosPorcion > 0 {
			factor := req.GramosConsumidos / alimento.GramosPorcion
			cal := alimento.Calorias * factor
			prot := alimento.ProteínasG * factor
			carb := alimento.CarbohidratosG * factor
			gras := alimento.GrasasG * factor
			ra.CaloriasCalc = &cal
			ra.ProteínasGCalc = &prot
			ra.CarbohidratosGCalc = &carb
			ra.GrasasGCalc = &gras
		}
	}

	if err := s.repo.CreateRegistroAlimento(ra); err != nil {
		return nil, err
	}

	// Actualizar los totales del registro_comida padre con la suma real de todos sus alimentos
	s.recalcRegistroComidaMacros(registroID)

	return ra, nil
}

// recalcRegistroComidaMacros suma los macros de todos los alimentos del registro y los persiste.
func (s *NutricionService) recalcRegistroComidaMacros(registroID uint) {
	alimentos, err := s.repo.FindRegistroAlimentosByRegistro(registroID)
	if err != nil {
		return
	}
	var cal, prot, carb, gras float64
	for _, a := range alimentos {
		if a.CaloriasCalc != nil {
			cal += *a.CaloriasCalc
		}
		if a.ProteínasGCalc != nil {
			prot += *a.ProteínasGCalc
		}
		if a.CarbohidratosGCalc != nil {
			carb += *a.CarbohidratosGCalc
		}
		if a.GrasasGCalc != nil {
			gras += *a.GrasasGCalc
		}
	}
	_ = s.repo.UpdateRegistroComidaMacros(registroID, math.Round(cal*100)/100, math.Round(prot*100)/100, math.Round(carb*100)/100, math.Round(gras*100)/100)
}

// ─── Registros ejercicio ──────────────────────────────────────────────────────

func (s *NutricionService) ListRegistrosEjercicio(pacienteID uint, fecha, desde, hasta string) ([]models.NutricionRegistroEjercicio, error) {
	return s.repo.FindRegistrosEjercicio(pacienteID, fecha, desde, hasta)
}

func (s *NutricionService) CreateRegistroEjercicio(pacienteID uint, req models.CreateRegistroEjercicioRequest) (*models.NutricionRegistroEjercicio, error) {
	fecha, err := time.ParseInLocation("2006-01-02", req.Fecha, time.Local)
	if err != nil {
		fecha = time.Now()
	}

	// Auto-calc calories from catalog when not provided
	calQuemadas := req.CaloriasQuemadas
	if calQuemadas == nil && req.EjercicioID != nil && req.DuracionMinReal != nil {
		if cat, err := s.repo.FindEjercicioCatalogoByID(*req.EjercicioID); err == nil && cat.CaloriasPorHora != nil {
			cal := *cat.CaloriasPorHora * float64(*req.DuracionMinReal) / 60.0
			cal = math.Round(cal*10) / 10
			calQuemadas = &cal
		}
	}

	re := &models.NutricionRegistroEjercicio{
		PacienteID:            pacienteID,
		Fecha:                 fecha,
		EjercicioPacienteID:   req.EjercicioPacienteID,
		EjercicioID:           req.EjercicioID,
		NombreLibre:           req.NombreLibre,
		DuracionMinReal:       req.DuracionMinReal,
		SeriesReal:            req.SeriesReal,
		RepeticionesReal:      req.RepeticionesReal,
		PesoKgReal:            req.PesoKgReal,
		CaloriasQuemadas:      calQuemadas,
		FrecuenciaCardiacaMax: req.FrecuenciaCardiacaMax,
		NivelEsfuerzo:         req.NivelEsfuerzo,
		Notas:                 req.Notas,
		State:                 "A",
	}
	if err := s.repo.CreateRegistroEjercicio(re); err != nil {
		return nil, err
	}
	return re, nil
}

// ─── Progreso ─────────────────────────────────────────────────────────────────

func (s *NutricionService) ListProgreso(pacienteID uint) ([]models.NutricionProgresoPaciente, error) {
	return s.repo.FindProgresoByPaciente(pacienteID)
}

func (s *NutricionService) AddProgreso(pacienteID, medicoID uint, req models.CreateProgresoRequest) (*models.NutricionProgresoPaciente, error) {
	fecha, err := time.ParseInLocation("2006-01-02", req.Fecha, time.Local)
	if err != nil {
		fecha = time.Now()
	}

	p := &models.NutricionProgresoPaciente{
		PacienteID:           pacienteID,
		MedicoID:             &medicoID,
		DietaPacienteID:      req.DietaPacienteID,
		Fecha:                fecha,
		PesoKg:               req.PesoKg,
		AlturaCm:             req.AlturaCm,
		GrasaCorporalPct:     req.GrasaCorporalPct,
		MasaMuscularKg:       req.MasaMuscularKg,
		CinturaCm:            req.CinturaCm,
		CaderaCm:             req.CaderaCm,
		PechoCm:              req.PechoCm,
		BrazoCm:              req.BrazoCm,
		MusloCm:              req.MusloCm,
		PctCumplimientoDieta: req.PctCumplimientoDieta,
		EnergiaNivel:         req.EnergiaNivel,
		SuenoHoras:           req.SuenoHoras,
		HidratacionLitros:    req.HidratacionLitros,
		Notas:                req.Notas,
		FotoProgreso:         req.FotoProgreso,
		State:                "A",
	}

	// Calcular IMC si hay peso y altura
	if p.PesoKg != nil && p.AlturaCm != nil && *p.AlturaCm > 0 {
		alturaM := *p.AlturaCm / 100.0
		imc := round2(*p.PesoKg / (alturaM * alturaM))
		p.IMC = &imc
	}

	if err := s.repo.CreateProgreso(p); err != nil {
		return nil, err
	}
	return p, nil
}

// ─── Tipo de Recurso ──────────────────────────────────────────────────────────

func (s *NutricionService) ListTipoRecursos() ([]models.NutricionTipoRecurso, error) {
	return s.repo.FindTipoRecursos()
}

func (s *NutricionService) CreateTipoRecurso(req models.CreateTipoRecursoRequest) (*models.NutricionTipoRecurso, error) {
	t := &models.NutricionTipoRecurso{Nombre: req.Nombre, State: "A"}
	return t, s.repo.CreateTipoRecurso(t)
}

func (s *NutricionService) UpdateTipoRecurso(id uint, req models.UpdateTipoRecursoRequest) (*models.NutricionTipoRecurso, error) {
	return s.repo.UpdateTipoRecurso(id, req.Nombre)
}

func (s *NutricionService) DeleteTipoRecurso(id uint) error {
	return s.repo.DeleteTipoRecurso(id)
}

// ─── Archivos PDF ─────────────────────────────────────────────────────────────

func (s *NutricionService) CreateArchivoPDF(clinicaID, medicoID uint, req models.CreateArchivoPDFRequest) (*models.NutricionArchivoPDF, error) {
	a := &models.NutricionArchivoPDF{
		ClinicaID:     clinicaID,
		MedicoID:      medicoID,
		PacienteID:    req.PacienteID,
		TipoRecursoID: req.TipoRecursoID,
		Titulo:        req.Titulo,
		RutaArchivo:   req.RutaArchivo,
		Descripcion:   req.Descripcion,
		State:         "A",
	}
	if err := s.repo.CreateArchivoPDF(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *NutricionService) ListArchivosPDF(clinicaID uint, pacienteID *uint, tipoRecursoID *uint) ([]models.NutricionArchivoPDF, error) {
	return s.repo.FindArchivosPDF(clinicaID, pacienteID, tipoRecursoID)
}

func (s *NutricionService) ListArchivosPDFByUser(clinicaID uint, pacienteID uint) ([]models.NutricionArchivoPDF, error) {
	return s.repo.FindArchivosPDFByUser(clinicaID, pacienteID)
}
func (s *NutricionService) DeleteArchivoPDF(id uint) error {
	return s.repo.DeleteArchivoPDF(id)
}

// ─── XP y logros ──────────────────────────────────────────────────────────────

func (s *NutricionService) GetXP(pacienteID uint) (*models.NutricionPacienteXP, error) {
	return s.repo.FindOrCreateXP(pacienteID)
}

func (s *NutricionService) ListLogros(pacienteID uint) ([]models.NutricionLogroPaciente, error) {
	return s.repo.FindLogrosByPaciente(pacienteID)
}

func (s *NutricionService) ListLogrosCatalogo() ([]models.NutricionLogroCatalogo, error) {
	return s.repo.FindLogrosCatalogo()
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func round2(x float64) float64 {
	return math.Round(x*100) / 100
}
