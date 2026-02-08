package api

import (
	"net/http"

	"github.com/buenvecino/dmarc-analyzer/internal/report"
	"github.com/gin-gonic/gin"
)

// Server holds the report and optional refresh callback.
type Server struct {
	Report *report.Report
	Refresh func() (*report.Report, error)
}

// Summary returns the full report (by domain, by IP, by recipient).
func (s *Server) Summary(c *gin.Context) {
	if s.Report == nil {
		c.JSON(http.StatusOK, gin.H{"by_domain": nil, "by_ip": nil, "by_recipient": nil})
		return
	}
	c.JSON(http.StatusOK, s.Report)
}

// ByDomain returns report grouped by domain.
func (s *Server) ByDomain(c *gin.Context) {
	if s.Report == nil {
		c.JSON(http.StatusOK, gin.H{"by_domain": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"by_domain": s.Report.ByDomain})
}

// ByIP returns report grouped by source IP (with known_sender flag).
func (s *Server) ByIP(c *gin.Context) {
	if s.Report == nil {
		c.JSON(http.StatusOK, gin.H{"by_ip": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"by_ip": s.Report.ByIP})
}

// ByRecipient returns report grouped by envelope_to.
func (s *Server) ByRecipient(c *gin.Context) {
	if s.Report == nil {
		c.JSON(http.StatusOK, gin.H{"by_recipient": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"by_recipient": s.Report.ByRecipient})
}

// RefreshReport re-scans DMARC_DIR and rebuilds the report.
func (s *Server) RefreshReport(c *gin.Context) {
	if s.Refresh == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "refresh not configured"})
		return
	}
	newReport, err := s.Refresh()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	s.Report = newReport
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Health returns 200 for Docker/Compose healthchecks.
func (s *Server) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
