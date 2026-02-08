package api

import (
	"github.com/gin-gonic/gin"
)

// Router sets up routes and CORS for the API.
func Router(s *Server) *gin.Engine {
	r := gin.Default()
	r.Use(corsMiddleware())

	api := r.Group("/api")
	{
		api.GET("/report", s.Summary)
		api.GET("/report/by-domain", s.ByDomain)
		api.GET("/report/by-ip", s.ByIP)
		api.GET("/report/by-recipient", s.ByRecipient)
		api.POST("/refresh", s.RefreshReport)
	}
	r.GET("/health", s.Health)

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
