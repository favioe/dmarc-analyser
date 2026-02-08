package dmarc

// Feedback represents the root DMARC aggregate report (feedback).
type Feedback struct {
	ReportMetadata   ReportMetadata   `xml:"report_metadata"`
	PolicyPublished  PolicyPublished  `xml:"policy_published"`
	Records         []Record         `xml:"record"`
}

// ReportMetadata contains report identification and date range.
type ReportMetadata struct {
	OrgName   string    `xml:"org_name"`
	Email     string    `xml:"email"`
	ReportID  string    `xml:"report_id"`
	DateRange DateRange `xml:"date_range"`
}

// DateRange is the period covered by the report.
type DateRange struct {
	Begin string `xml:"begin"`
	End   string `xml:"end"`
}

// PolicyPublished contains the domain and policy requested by the domain owner.
type PolicyPublished struct {
	Domain string `xml:"domain"`
	P      string `xml:"p"`   // Policy: none, quarantine, reject
	Pct    int    `xml:"pct"` // Percentage of messages the policy applies to (often 100)
}

// Record is a single DMARC record (one row + identifiers + auth_results).
type Record struct {
	Row          Row          `xml:"row"`
	Identifiers  Identifiers  `xml:"identifiers"`
}

// Row holds source_ip, count, and policy_evaluated.
type Row struct {
	SourceIP         string          `xml:"source_ip"`
	Count            int             `xml:"count"`
	PolicyEvaluated  PolicyEvaluated `xml:"policy_evaluated"`
}

// PolicyEvaluated holds disposition, dkim, spf results.
type PolicyEvaluated struct {
	Disposition string `xml:"disposition"`
	DKIM        string `xml:"dkim"`
	SPF         string `xml:"spf"`
}

// Identifiers holds envelope_to, envelope_from, header_from (envelope_to may be missing).
type Identifiers struct {
	EnvelopeTo   string `xml:"envelope_to"`
	EnvelopeFrom string `xml:"envelope_from"`
	HeaderFrom   string `xml:"header_from"`
}
