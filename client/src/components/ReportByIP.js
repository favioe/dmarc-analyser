import React from 'react';

function ReportByIP({ report }) {
  if (!report?.by_ip) return <p>No IP data.</p>;
  const ips = Object.values(report.by_ip);
  if (ips.length === 0) return <p>No IP data.</p>;

  return (
    <>
      <h2 className="section-title">Breakdown by IP (source_ip)</h2>
      <p style={{ marginBottom: 16 }}>IPs marked as &quot;known sender&quot; when listed in KNOWN_SERVER_IPS (e.g. Mailgun, SendGrid).</p>
      <div className="table-wrap">
        <table>
          <thead>
            <tr>
              <th>IP</th>
              <th>Known sender</th>
              <th>Total</th>
              <th>Disposition (none / quarantine / reject)</th>
              <th>DKIM (pass / fail)</th>
              <th>SPF (pass / fail)</th>
            </tr>
          </thead>
          <tbody>
            {ips.map((ip) => (
              <tr key={ip.ip} className={ip.known_sender ? 'known-sender' : ''}>
                <td><code>{ip.ip}</code></td>
                <td>{ip.known_sender ? 'Yes' : '—'}</td>
                <td>{ip.total}</td>
                <td>{formatMap(ip.by_disposition)}</td>
                <td>{formatMap(ip.by_dkim)}</td>
                <td>{formatMap(ip.by_spf)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  );
}

function formatMap(m) {
  if (!m || typeof m !== 'object') return '—';
  return Object.entries(m)
    .map(([k, v]) => `${k}: ${v}`)
    .join(' · ');
}

export default ReportByIP;
