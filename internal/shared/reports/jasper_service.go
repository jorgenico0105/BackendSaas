package reports

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	ReportsDirectory  = "storage/reports"
	TemplatesDirectory = "resources/jasper_templates"
)

type ReportFormat string

const (
	FormatPDF   ReportFormat = "pdf"
	FormatXLS   ReportFormat = "xls"
	FormatXLSX  ReportFormat = "xlsx"
	FormatCSV   ReportFormat = "csv"
	FormatHTML  ReportFormat = "html"
)

type ReportParams struct {
	TemplateName string
	OutputName   string
	Format       ReportFormat
	Parameters   map[string]interface{}
}

type JasperService struct {
	JasperStarterPath string
	JDBCDriverPath    string
	DBHost            string
	DBPort            string
	DBName            string
	DBUser            string
	DBPassword        string
}

func NewJasperService(jasperPath, jdbcPath, host, port, dbName, user, password string) *JasperService {
	return &JasperService{
		JasperStarterPath: jasperPath,
		JDBCDriverPath:    jdbcPath,
		DBHost:            host,
		DBPort:            port,
		DBName:            dbName,
		DBUser:            user,
		DBPassword:        password,
	}
}

func (j *JasperService) GenerateReport(params ReportParams) (string, error) {
	if err := os.MkdirAll(ReportsDirectory, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create reports directory: %w", err)
	}

	templatePath := filepath.Join(TemplatesDirectory, params.TemplateName+".jasper")
	outputPath := filepath.Join(ReportsDirectory, params.OutputName)

	jdbcURL := fmt.Sprintf("jdbc:mysql://%s:%s/%s", j.DBHost, j.DBPort, j.DBName)

	args := []string{
		"process",
		templatePath,
		"-o", outputPath,
		"-f", string(params.Format),
		"-t", "mysql",
		"-H", j.DBHost,
		"-n", j.DBName,
		"-u", j.DBUser,
		"-p", j.DBPassword,
		"--jdbc-url", jdbcURL,
		"--jdbc-driver", j.JDBCDriverPath,
	}

	for key, value := range params.Parameters {
		args = append(args, "-P", fmt.Sprintf("%s=%v", key, value))
	}

	cmd := exec.Command(j.JasperStarterPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("jasper error: %s - %w", string(output), err)
	}

	return outputPath + "." + string(params.Format), nil
}

func (j *JasperService) GetReportPath(fileName string) string {
	return filepath.Join(ReportsDirectory, fileName)
}

func DeleteReport(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(filePath)
}
