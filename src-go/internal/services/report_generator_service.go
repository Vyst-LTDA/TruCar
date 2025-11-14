package services

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

// ReportGeneratorService define a interface para a geração de relatórios em PDF.
type ReportGeneratorService interface {
	GeneratePDFFromHTML(templatePath string, data interface{}) ([]byte, error)
}

type reportGeneratorService struct {
}

// NewReportGeneratorService cria uma nova instância de ReportGeneratorService.
func NewReportGeneratorService() ReportGeneratorService {
	return &reportGeneratorService{}
}

// GeneratePDFFromHTML gera um PDF a partir de um template HTML e dados.
func (s *reportGeneratorService) GeneratePDFFromHTML(templatePath string, data interface{}) ([]byte, error) {
	// 1. Renderizar o template HTML
	htmlContent, err := s.renderHTML(templatePath, data)
	if err != nil {
		return nil, fmt.Errorf("failed to render HTML template: %w", err)
	}

	// 2. Converter o HTML para PDF
	pdfBytes, err := s.convertHTMLToPDF(htmlContent)
	if err != nil {
		return nil, fmt.Errorf("failed to convert HTML to PDF: %w", err)
	}

	return pdfBytes, nil
}

func (s *reportGeneratorService) renderHTML(templatePath string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var renderedHTML bytes.Buffer
	err = tmpl.Execute(&renderedHTML, data)
	if err != nil {
		return "", err
	}

	return renderedHTML.String(), nil
}

func (s *reportGeneratorService) convertHTMLToPDF(htmlContent string) ([]byte, error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(htmlContent))))

	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
}
