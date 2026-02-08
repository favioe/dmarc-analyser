import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';
import ReportByDomain from './components/ReportByDomain';
import ReportByIP from './components/ReportByIP';
import ReportByRecipient from './components/ReportByRecipient';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

function App() {
  const [activeTab, setActiveTab] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [report, setReport] = useState(null);

  const fetchReport = async () => {
    try {
      setLoading(true);
      const res = await axios.get(`${API_URL}/api/report`);
      setReport(res.data);
      setError(null);
    } catch (err) {
      setError(err.message || 'Error loading report. Make sure the API is running.');
      setReport(null);
    } finally {
      setLoading(false);
    }
  };

  const handleRefresh = async () => {
    try {
      await axios.post(`${API_URL}/api/refresh`);
      await fetchReport();
    } catch (err) {
      setError(err.message || 'Error refreshing.');
    }
  };

  useEffect(() => {
    fetchReport();
  }, []);

  if (loading && !report) {
    return (
      <div className="App">
        <div className="loading">
          <p>Loading DMARC report...</p>
        </div>
      </div>
    );
  }

  if (error && !report) {
    return (
      <div className="App">
        <div className="error">
          <h2>Error</h2>
          <p>{error}</p>
          <button onClick={fetchReport}>Retry</button>
        </div>
      </div>
    );
  }

  const tabs = [
    { label: 'By domain', panel: <ReportByDomain report={report} /> },
    { label: 'By IP', panel: <ReportByIP report={report} /> },
    { label: 'By recipient', panel: <ReportByRecipient report={report} /> },
  ];

  return (
    <div className="App">
      <header className="App-header">
        <h1>DMARC Analyzer</h1>
        <p className="subtitle">Domain trust from DMARC reports</p>
        <button type="button" className="refresh-btn" onClick={handleRefresh}>
          Refresh report
        </button>
      </header>

      <section className="policy-glossary" style={{ margin: '0 24px 24px', padding: '16px 20px', background: '#f8f9fa', borderRadius: 8, border: '1px solid #e0e0e0' }}>
        <h3 style={{ margin: '0 0 12px 0', fontSize: '1.1rem' }}>DMARC policy glossary</h3>
        <p style={{ margin: '0 0 8px 0', fontSize: '0.95rem', color: '#333' }}>
          <strong>Disposition</strong> is what receivers do with messages that fail DMARC (e.g. DKIM/SPF fail). The domain can request one of these policies:
        </p>
        <ul style={{ margin: '0 0 0 20px', padding: 0, fontSize: '0.9rem', color: '#444' }}>
          <li><strong>none</strong> (monitor) — No action; mail is delivered. Used to collect data without affecting delivery.</li>
          <li><strong>quarantine</strong> — Failing messages are treated as suspicious (e.g. placed in spam/junk).</li>
          <li><strong>reject</strong> — Failing messages are rejected and not delivered.</li>
        </ul>
        <p style={{ margin: '8px 0 0 0', fontSize: '0.85rem', color: '#666' }}>
          Reports also show <strong>p</strong> (policy for the domain) and optionally <strong>sp</strong> (policy for subdomains). Moving from <em>none</em> → <em>quarantine</em> → <em>reject</em> is the usual path to stricter protection.
        </p>
      </section>

      <main className="tabs">
        <div className="tab-list">
          {tabs.map((t, i) => (
            <button
              key={i}
              type="button"
              className={activeTab === i ? 'active' : ''}
              onClick={() => setActiveTab(i)}
            >
              {t.label}
            </button>
          ))}
        </div>
        <div className="tab-panel">{tabs[activeTab].panel}</div>
      </main>
    </div>
  );
}

export default App;
