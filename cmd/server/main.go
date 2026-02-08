package main

import (
	"log"
	"os"

	"github.com/buenvecino/dmarc-analyzer/internal/api"
	"github.com/buenvecino/dmarc-analyzer/internal/config"
	"github.com/buenvecino/dmarc-analyzer/internal/report"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// Build report at startup
	rep, err := report.Build(cfg)
	if err != nil {
		log.Fatalf("build report: %v", err)
	}
	log.Printf("report built: %d domains, %d IPs, %d recipients", len(rep.ByDomain), len(rep.ByIP), len(rep.ByRecipient))

	srv := &api.Server{
		Report: rep,
		Refresh: func() (*report.Report, error) {
			return report.Build(cfg)
		},
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("listening on %s", addr)
	if err := api.Router(srv).Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}
