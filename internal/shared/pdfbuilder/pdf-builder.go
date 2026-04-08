package pdfbuilder

import (
	"log"
	"os"
	"path/filepath"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/core/entity"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type UseCase interface {
	CreatePdf() error
}

type PdfService struct {
	WaterMark string
}

func NewPdfBuilder(watermark string) *PdfService {
	return &PdfService{WaterMark: watermark}
}

type PdfBuilder interface {
	GeneratePdfBuilder() (core.Maroto, error)
}

func (pdfs *PdfService) GeneratePdfBuilder() (core.Maroto, error) {
	m := GetMaroto(pdfs.WaterMark)

	return m, nil
}

func GetMaroto(watermark string) core.Maroto {
	regularPath := resourcePath("Poppins", "Poppins-Regular.ttf")
	boldPath := resourcePath("Poppins", "Poppins-Bold.ttf")
	italicPath := resourcePath("Poppins", "Poppins-Italic.ttf")

	customFonts := []*entity.CustomFont{
		{
			Family: "Poppins",
			Style:  fontstyle.Normal,
			File:   regularPath,
			Bytes:  mustReadFont(regularPath),
		},
		{
			Family: "Poppins",
			Style:  fontstyle.Bold,
			File:   boldPath,
			Bytes:  mustReadFont(boldPath),
		},
		{
			Family: "Poppins",
			Style:  fontstyle.Italic,
			File:   italicPath,
			Bytes:  mustReadFont(italicPath),
		},
	}

	cfg := config.NewBuilder().
		WithPageNumber().
		WithCustomFonts(customFonts).
		WithLeftMargin(10).
		WithTopMargin(15).
		WithRightMargin(10).
		WithBackgroundImage(mustReadFont(watermark), extension.Png).
		WithDefaultFont(&props.Font{Family: "Poppins"}).
		Build()

	mrt := maroto.New(cfg)
	return maroto.NewMetricsDecorator(mrt)
}

func mustReadFont(path string) []byte {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read font %s: %v", path, err)
	}
	return b
}

func resourcePath(parts ...string) string {
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}

	allParts := append([]string{baseDir, "resources"}, parts...)
	return filepath.Join(allParts...)
}
