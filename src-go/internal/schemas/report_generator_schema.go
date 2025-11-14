package schemas

import "time"

// ReportRequest é o schema para a requisição de geração de relatório em PDF.
type ReportRequest struct {
	ReportType string    `json:"report_type" binding:"required"`
	TargetID   uint      `json:"target_id" binding:"required"`
	DateFrom   time.Time `json:"date_from" binding:"required"`
	DateTo     time.Time `json:"date_to" binding:"required"`
}
