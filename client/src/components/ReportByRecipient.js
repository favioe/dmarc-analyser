import React from 'react';

function ReportByRecipient({ report }) {
  if (!report?.by_recipient) return <p>No recipient data.</p>;
  const recipients = Object.values(report.by_recipient);
  if (recipients.length === 0) return <p>No recipient data.</p>;

  return (
    <>
      <h2 className="section-title">Grouped by recipient (envelope_to)</h2>
      <p style={{ marginBottom: 16 }}>Totals and breakdown by IP to see IP–reject correlation per recipient.</p>
      {recipients.map((r) => (
        <div key={r.recipient} style={{ marginBottom: 24 }}>
          <h3 style={{ marginBottom: 8 }}>{r.recipient}</h3>
          <p><strong>Total:</strong> {r.total}</p>
          <p><strong>Disposition:</strong> {formatMap(r.by_disposition)}</p>
          <p><strong>DKIM:</strong> {formatMap(r.by_dkim)} · <strong>SPF:</strong> {formatMap(r.by_spf)}</p>
          {r.by_ip && Object.keys(r.by_ip).length > 0 && (
            <div className="table-wrap" style={{ marginTop: 12 }}>
              <table>
                <thead>
                  <tr>
                    <th>IP</th>
                    <th>Total</th>
                    <th>Disposition</th>
                  </tr>
                </thead>
                <tbody>
                  {Object.entries(r.by_ip).map(([ip, counts]) => (
                    <tr key={ip}>
                      <td><code>{ip}</code></td>
                      <td>{counts.total}</td>
                      <td>{formatMap(counts.by_disposition)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      ))}
    </>
  );
}

function formatMap(m) {
  if (!m || typeof m !== 'object') return '—';
  return Object.entries(m)
    .map(([k, v]) => `${k}: ${v}`)
    .join(' · ');
}

export default ReportByRecipient;
