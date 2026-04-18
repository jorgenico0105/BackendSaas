package usecases

import (
	"fmt"
	"math"
	"saas-medico/internal/modules/nutricion/models"
	"strconv"
	"strings"
	"time"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type MenuPdf struct {
	Dieta       *models.NutricionDietaPaciente
	Menu        *models.NutricionMenu
	MarotoI     core.Maroto
	LogoPath    string
	OutputhPath string
}
type ComidaAlimentos struct {
	NombreComida string
	NombreReceta string
	Alimentos    []models.NutricionMenuAlimento
}

func NewMenuPdfUseCase(dieta *models.NutricionDietaPaciente, menu *models.NutricionMenu, m core.Maroto, logoPath, outputhPath string) *MenuPdf {
	return &MenuPdf{Dieta: dieta, Menu: menu, MarotoI: m, LogoPath: logoPath, OutputhPath: outputhPath}
}

func (mpdf *MenuPdf) CreatePdf() error {
	semana := mpdf.Menu.SemanaNumero

	// Level 1: group by day number → ordered list of meals
	// Each NutricionMenuDetalle is already one meal (unique by dia+tipo_comida),
	// so we just bucket detalles by DiaNúmero.
	diaMap := make(map[int8][]*ComidaAlimentos)
	for _, detalle := range mpdf.Menu.Detalles {
		diaMap[detalle.DiaNúmero] = append(diaMap[detalle.DiaNúmero], &ComidaAlimentos{
			NombreComida: detalle.NombreComida,
			NombreReceta: detalle.NombreReceta,
			Alimentos:    detalle.Alimentos,
		})
	}

	mpdf.MarotoI.AddRows(getPageHeader(mpdf.LogoPath))

	mpdf.MarotoI.AddRows(text.NewRow(20, "Plan Nutricional Semana "+strconv.Itoa(semana), props.Text{
		Top:    5,
		Size:   15,
		Bottom: 3,
		Style:  fontstyle.Bold,
		Family: "Poppins",
		Align:  align.Center,
	}))
	mpdf.MarotoI.AddRows(getPacienteTable(mpdf.Dieta)...)

	// Iterate days 1-7 in order
	for dia := int8(1); dia <= 7; dia++ {
		comidas, ok := diaMap[dia]
		if !ok {
			continue
		}

		// Day header
		mpdf.MarotoI.AddRows(text.NewRow(8, "Dia "+strconv.Itoa(int(dia)), props.Text{
			Top:    2,
			Size:   11,
			Style:  fontstyle.Bold,
			Family: "Poppins",
			Left:   3,
		}))

		// Table header once per day
		mpdf.MarotoI.AddRows(getAlimentosHeader())

		// Level 2: iterate each meal of the day (data rows only)
		for _, comida := range comidas {
			mpdf.MarotoI.AddRows(getAlimetosContent(comida)...)
		}
	}

	document, err := mpdf.MarotoI.Generate()
	if err != nil {
		return err
	}

	return document.Save(mpdf.OutputhPath)
}
func getPageHeader(logoPath string) core.Row {
	return row.New(40).Add(
		image.NewFromFileCol(12, logoPath, props.Rect{
			Center:  true,
			Percent: 100,
		}),
	)
}

func getPacienteTable(dieta *models.NutricionDietaPaciente) []core.Row {
	rows := []core.Row{
		row.New(6).Add(
			col.New(1),
			text.NewCol(2, "Nombre : ", props.Text{
				Size:  9,
				Style: fontstyle.Bold,
				Top:   1,
				Left:  2,
			}).WithStyle(&props.Cell{BorderType: border.Full}),
			text.NewCol(8, dieta.Paciente.Nombres+dieta.Paciente.Apellidos, props.Text{
				Size:  9,
				Style: fontstyle.Normal,
				Top:   1,
				Left:  2,
			}).WithStyle(&props.Cell{BorderType: border.Full}),
		),
		row.New(6).Add(
			col.New(1),
			text.NewCol(2, "Edad : ", props.Text{
				Size:  9,
				Style: fontstyle.Bold,
				Top:   1,
				Left:  2,
			}).WithStyle(&props.Cell{BorderType: border.Full}),
			text.NewCol(8, calcularEdad(dieta.Paciente.FechaNacimiento), props.Text{
				Size:  9,
				Left:  2,
				Top:   1,
				Style: fontstyle.Normal,
			}).WithStyle(&props.Cell{BorderType: border.Full}),
		),
		row.New(6).Add(
			col.New(1),
			text.NewCol(2, "Objetivo : ", props.Text{
				Size:  10,
				Style: fontstyle.Bold,
				Top:   1,
				Left:  2,
			}).WithStyle(&props.Cell{BorderType: border.Full}),
			text.NewCol(8, dieta.Objetivo, props.Text{
				Size:  10,
				Left:  2,
				Top:   1,
				Style: fontstyle.Normal,
			}).WithStyle(&props.Cell{BorderType: border.Full}),
		),
		row.New(6).Add(
			col.New(1),
			text.NewCol(2, "Calorías : ", props.Text{
				Size:  10,
				Style: fontstyle.Bold,
				Top:   1,
				Left:  2,
			}).WithStyle(&props.Cell{BorderType: border.Full}),
			text.NewCol(8, fmt.Sprintf("%.0f", *dieta.CaloriasDiaObjetivo), props.Text{
				Size:  10,
				Left:  2,
				Top:   1,
				Style: fontstyle.Normal,
			}).WithStyle(&props.Cell{BorderType: border.Full}),
		),
	}
	return rows
}

func calcularEdad(fechaNacimiento *time.Time) string {
	if fechaNacimiento == nil {
		return "N/A"
	}
	hoy := time.Now()
	años := hoy.Year() - fechaNacimiento.Year()
	if hoy.Month() < fechaNacimiento.Month() ||
		(hoy.Month() == fechaNacimiento.Month() && hoy.Day() < fechaNacimiento.Day()) {
		años--
	}
	return strconv.Itoa(años) + " años"
}
func getAlimentosHeader() core.Row {
	cellStyle := &props.Cell{BorderType: border.Full}
	return row.New(7).Add(
		col.New(1),
		text.NewCol(3, "CANTIDAD", props.Text{
			Size:  9,
			Style: fontstyle.Bold,
			Top:   2,
			Align: align.Center,
		}).WithStyle(cellStyle),
		text.NewCol(8, "ALIMENTOS", props.Text{
			Size:  9,
			Style: fontstyle.Bold,
			Top:   2,
			Align: align.Center,
		}).WithStyle(cellStyle),
	)
}

func getAlimetosContent(ca *ComidaAlimentos) []core.Row {
	cellStyle := &props.Cell{BorderType: border.Full}

	// Meal name as subheader row
	rows := []core.Row{
		row.New(7).Add(
			col.New(1),
			text.NewCol(11, strings.ToUpper(ca.NombreComida), props.Text{
				Size:  9,
				Style: fontstyle.Bold,
				Top:   2,
				Left:  3,
			}).WithStyle(&props.Cell{BorderType: border.Full}),
		),
	}

	// Recipe name row (only if set)
	if ca.NombreReceta != "" {
		rows = append(rows, row.New(5).Add(
			col.New(1),
			text.NewCol(11, "Receta: "+ca.NombreReceta, props.Text{
				Size:  8,
				Style: fontstyle.Italic,
				Top:   1,
				Left:  3,
			}).WithStyle(&props.Cell{BorderType: border.Full}),
		))
	}

	for _, alimento := range ca.Alimentos {
		var cantidad string
		if alimento.Observacion != "" {
			cantidad = alimento.Observacion
		} else if alimento.Alimento.NeedUnidad && alimento.Alimento.GramosUnidad != nil && *alimento.Alimento.GramosUnidad > 0 {
			unidades := math.Ceil(alimento.GramosAsignados / *alimento.Alimento.GramosUnidad)
			medida := alimento.Alimento.Medida
			if medida == "" {
				medida = "unidad(es)"
			}
			cantidad = fmt.Sprintf("%.0f %s", unidades, medida)
		} else {
			cantidad = fmt.Sprintf("%.0f gr", alimento.GramosAsignados)
		}

		rows = append(rows, row.New(6).Add(
			col.New(1),
			text.NewCol(3, cantidad, props.Text{
				Size:  9,
				Style: fontstyle.Normal,
				Top:   1,
				Align: align.Center,
			}).WithStyle(cellStyle),
			text.NewCol(8, alimento.Alimento.Nombre, props.Text{
				Size:  9,
				Style: fontstyle.Normal,
				Top:   1,
				Left:  2,
			}).WithStyle(cellStyle),
		))
	}

	return rows
}
