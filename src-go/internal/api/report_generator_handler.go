package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-api/internal/schemas"
	"go-api/internal/services"
)

// ReportGeneratorHandler gerencia as requisições HTTP para a geração de relatórios.
type ReportGeneratorHandler struct {
	generatorService services.ReportGeneratorService
	reportService    services.ReportService
}

// NewReportGeneratorHandler cria uma nova instância de ReportGeneratorHandler.
func NewReportGeneratorHandler(
	generatorService services.ReportGeneratorService,
	reportService services.ReportService,
) *ReportGeneratorHandler {
	return &ReportGeneratorHandler{
		generatorService: generatorService,
		reportService:    reportService,
	}
}

// GeneratePDF lida com a geração de um relatório em PDF.
func (h *ReportGeneratorHandler) GeneratePDF(c *gin.Context) {
	// orgID, exists := middleware.GetOrganizationID(c)
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Organização não identificada"})
	// 	return
	// }

	var req schemas.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var data interface{}
	var templatePath string
	var err error

	switch req.ReportType {
	case "activity_by_driver":
		// data, err = h.reportService.GetDriverActivityData(orgID, req.TargetID, req.DateFrom, req.DateTo)
		templatePath = "internal/templates/driver_activity_report.html"
		// Placeholder de dados
		data = gin.H{"DriverName": "João Silva", "TotalDistance": 1234.5}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de relatório inválido"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch report data"})
		return
	}
	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Não foram encontrados dados para este relatório."})
		return
	}

	pdfBytes, err := h.generatorService.GeneratePDFFromHTML(templatePath, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	filename := fmt.Sprintf("%s_report.pdf", req.ReportType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
