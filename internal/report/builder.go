package report

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/buenvecino/dmarc-analyzer/internal/config"
	"github.com/buenvecino/dmarc-analyzer/internal/dmarc"
)

// Build scans dmarcDir, processes each .gz and .zip file, and returns the aggregated Report.
// Known sender IPs/ranges (e.g. KNOWN_SERVER_IPS) are read from cfg for KnownSender on IPStats.
func Build(cfg *config.Config) (*Report, error) {
	r := NewReport()

	entries, err := os.ReadDir(cfg.DMARCDir)
	if err != nil {
		return nil, fmt.Errorf("read dir %s: %w", cfg.DMARCDir, err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		lower := strings.ToLower(name)
		if !strings.HasSuffix(lower, ".gz") && !strings.HasSuffix(lower, ".zip") {
			continue
		}
		path := filepath.Join(cfg.DMARCDir, name)
		if err := processFile(path, r, cfg); err != nil {
			return nil, fmt.Errorf("process %s: %w", name, err)
		}
	}

	return r, nil
}

func processFile(path string, r *Report, cfg *config.Config) error {
	rc, err := dmarc.OpenReport(path)
	if err != nil {
		return err
	}
	defer rc.Close()

	fb, err := dmarc.Decode(rc)
	if err != nil {
		return err
	}

	domain := strings.TrimSpace(fb.PolicyPublished.Domain)
	if domain == "" {
		return nil
	}

	// Ensure domain stats exist and set policy from this report
	if r.ByDomain[domain] == nil {
		r.ByDomain[domain] = &DomainStats{
			Domain:        domain,
			ByDisposition: make(map[string]int),
			ByDKIM:        make(map[string]int),
			BySPF:         make(map[string]int),
		}
	}
	dom := r.ByDomain[domain]
	dom.PolicyRequested = strings.TrimSpace(fb.PolicyPublished.P)
	if dom.PolicyRequested == "" {
		dom.PolicyRequested = "none"
	}
	if fb.PolicyPublished.Pct > 0 {
		dom.Pct = fb.PolicyPublished.Pct
	}

	for _, rec := range fb.Records {
		ip := strings.TrimSpace(rec.Row.SourceIP)
		count := rec.Row.Count
		if count <= 0 {
			count = 1
		}
		disp := strings.TrimSpace(rec.Row.PolicyEvaluated.Disposition)
		dkim := strings.TrimSpace(rec.Row.PolicyEvaluated.DKIM)
		spf := strings.TrimSpace(rec.Row.PolicyEvaluated.SPF)

		// By domain
		dom.Total += count
		inc(dom.ByDisposition, disp, count)
		inc(dom.ByDKIM, dkim, count)
		inc(dom.BySPF, spf, count)

		// By IP
		if ip != "" {
			if r.ByIP[ip] == nil {
				r.ByIP[ip] = &IPStats{
					IP:            ip,
					KnownSender:   cfg.IsKnownSenderIP(ip),
					ByDisposition: make(map[string]int),
					ByDKIM:        make(map[string]int),
					BySPF:         make(map[string]int),
				}
			}
			ipStats := r.ByIP[ip]
			ipStats.Total += count
			inc(ipStats.ByDisposition, disp, count)
			inc(ipStats.ByDKIM, dkim, count)
			inc(ipStats.BySPF, spf, count)
		}

		// By recipient (envelope_to)
		recipient := strings.TrimSpace(rec.Identifiers.EnvelopeTo)
		if recipient == "" {
			recipient = "(empty)"
		}
		if r.ByRecipient[recipient] == nil {
			r.ByRecipient[recipient] = &RecipientStats{
				Recipient:     recipient,
				ByIP:          make(map[string]IPCounts),
				ByDisposition: make(map[string]int),
				ByDKIM:        make(map[string]int),
				BySPF:         make(map[string]int),
			}
		}
		recStats := r.ByRecipient[recipient]
		recStats.Total += count
		inc(recStats.ByDisposition, disp, count)
		inc(recStats.ByDKIM, dkim, count)
		inc(recStats.BySPF, spf, count)
		if ip != "" {
			pc := recStats.ByIP[ip]
			if pc.ByDisposition == nil {
				pc.ByDisposition = make(map[string]int)
			}
			pc.Total += count
			inc(pc.ByDisposition, disp, count)
			recStats.ByIP[ip] = pc
		}
	}

	return nil
}
