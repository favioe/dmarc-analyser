import React from 'react';
import { PieChart, Pie, Cell, Legend, Tooltip, ResponsiveContainer } from 'recharts';

const COLORS = { none: '#4caf50', quarantine: '#ff9800', reject: '#f44336', pass: '#4caf50', fail: '#f44336' };

const DISP_LABELS = {
  none: 'delivered (none)',
  quarantine: 'quarantine',
  reject: 'reject',
};

function ReportByDomain({ report }) {
  if (!report?.by_domain) return <p>No domain data.</p>;
  const domains = Object.values(report.by_domain);
  if (domains.length === 0) return <p>No domain data.</p>;

  return (
    <>
      <h2 className="section-title">Summary by domain</h2>
      {domains.map((d) => {
        const dispData = Object.entries(d.by_disposition || {}).map(([k, v]) => ({ name: k, value: v }));
        const dkimData = Object.entries(d.by_dkim || {}).map(([k, v]) => ({ name: `DKIM ${k}`, value: v }));
        const spfData = Object.entries(d.by_spf || {}).map(([k, v]) => ({ name: `SPF ${k}`, value: v }));
        const total = d.total || 1;
        const policyLabel = d.policy_requested === 'quarantine' ? 'quarantine' : d.policy_requested === 'reject' ? 'reject' : (d.policy_requested || 'none');
        const conclusionParts = Object.entries(d.by_disposition || {})
          .filter(([, v]) => v > 0)
          .map(([k, v]) => `${DISP_LABELS[k] || k} at ${((v / total) * 100).toFixed(1)}%`)
          .join(', ');

        return (
          <div key={d.domain} style={{ marginBottom: 32 }}>
            <h3 style={{ marginBottom: 16 }}>{d.domain}</h3>
            <p><strong>Total messages:</strong> {d.total}</p>

            <div className="conclusiones" style={{ marginTop: 16, marginBottom: 16, padding: '12px 16px', background: '#f0f7ff', borderRadius: 8, borderLeft: '4px solid #1976d2' }}>
              <h4 style={{ margin: '0 0 8px 0', fontSize: '1rem' }}>Summary</h4>
              <p style={{ margin: 0 }}>
                This domain published policy <strong>{policyLabel}</strong>
                {d.pct ? ` (applied to ${d.pct}% of messages)` : ''}. Based on that policy and authentication results (DKIM/SPF), receivers applied: {conclusionParts || '—'} to your emails.
              </p>
            </div>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: 24, marginTop: 16 }}>
              {dispData.length > 0 && (
                <div style={{ width: 280, height: 220 }}>
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie data={dispData} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={70} label={(e) => `${e.name}: ${e.value}`}>
                        {dispData.map((_, i) => (
                          <Cell key={i} fill={COLORS[dispData[i].name] || '#888'} />
                        ))}
                      </Pie>
                      <Tooltip />
                      <Legend />
                    </PieChart>
                  </ResponsiveContainer>
                  <p style={{ textAlign: 'center', margin: 0 }}>Disposition</p>
                </div>
              )}
              {dkimData.length > 0 && (
                <div style={{ width: 280, height: 220 }}>
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie data={dkimData} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={70} label={(e) => `${e.name}: ${e.value}`}>
                        {dkimData.map((_, i) => (
                          <Cell key={i} fill={COLORS[dkimData[i].name?.replace('DKIM ', '')] || '#888'} />
                        ))}
                      </Pie>
                      <Tooltip />
                      <Legend />
                    </PieChart>
                  </ResponsiveContainer>
                  <p style={{ textAlign: 'center', margin: 0 }}>DKIM</p>
                </div>
              )}
              {spfData.length > 0 && (
                <div style={{ width: 280, height: 220 }}>
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie data={spfData} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={70} label={(e) => `${e.name}: ${e.value}`}>
                        {spfData.map((_, i) => (
                          <Cell key={i} fill={COLORS[spfData[i].name?.replace('SPF ', '')] || '#888'} />
                        ))}
                      </Pie>
                      <Tooltip />
                      <Legend />
                    </PieChart>
                  </ResponsiveContainer>
                  <p style={{ textAlign: 'center', margin: 0 }}>SPF</p>
                </div>
              )}
            </div>
          </div>
        );
      })}
    </>
  );
}

export default ReportByDomain;
