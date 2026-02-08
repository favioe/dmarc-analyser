package report

// DomainStats holds aggregated counts for a single domain.
type DomainStats struct {
	Domain          string          `json:"domain"`
	PolicyRequested string          `json:"policy_requested"` // none, quarantine, reject (what the domain asked for)
	Pct             int             `json:"pct"`               // percentage of messages policy applies to (from reports)
	Total           int             `json:"total"`
	ByDisposition   map[string]int  `json:"by_disposition"`   // none, quarantine, reject (what receivers did)
	ByDKIM          map[string]int  `json:"by_dkim"`
	BySPF           map[string]int   `json:"by_spf"`
}

// IPStats holds aggregated counts for a source IP, plus whether it's a known sender IP.
type IPStats struct {
	IP           string          `json:"ip"`
	KnownSender  bool            `json:"known_sender"`
	Total        int             `json:"total"`
	ByDisposition map[string]int `json:"by_disposition"`
	ByDKIM       map[string]int  `json:"by_dkim"`
	BySPF        map[string]int  `json:"by_spf"`
}

// RecipientStats holds aggregated counts per envelope_to (recipient domain), with optional per-IP breakdown.
type RecipientStats struct {
	Recipient    string              `json:"recipient"`
	Total        int                 `json:"total"`
	ByIP         map[string]IPCounts `json:"by_ip,omitempty"`
	ByDisposition map[string]int     `json:"by_disposition"`
	ByDKIM       map[string]int      `json:"by_dkim"`
	BySPF        map[string]int     `json:"by_spf"`
}

// IPCounts is a small struct for per-IP counts inside a recipient.
type IPCounts struct {
	Total        int            `json:"total"`
	ByDisposition map[string]int `json:"by_disposition"`
}

// Report holds the full in-memory aggregate for the API.
type Report struct {
	ByDomain    map[string]*DomainStats    `json:"by_domain"`
	ByIP        map[string]*IPStats       `json:"by_ip"`
	ByRecipient map[string]*RecipientStats `json:"by_recipient"`
}

// NewReport creates an empty Report ready for aggregation.
func NewReport() *Report {
	return &Report{
		ByDomain:    make(map[string]*DomainStats),
		ByIP:        make(map[string]*IPStats),
		ByRecipient: make(map[string]*RecipientStats),
	}
}

func inc(m map[string]int, key string, delta int) {
	if key == "" {
		key = "(empty)"
	}
	m[key] += delta
}
