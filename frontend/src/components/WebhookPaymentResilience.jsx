import React, { useState, useEffect } from 'react';
import { Radio, RefreshCw, CheckCircle2, ShieldAlert, Cpu, Activity, Clock, Zap } from 'lucide-react';

export default function WebhookPaymentResilience() {
  const [idempotencyKey, setIdempotencyKey] = useState('IDEMP-PAY-2026-X889');
  const [paymentID, setPaymentID] = useState('pay-sample-001');
  const [amount, setAmount] = useState(1500000);
  const [webhookLog, setWebhookLog] = useState([]);
  const [isSending, setIsSending] = useState(false);

  // SSE State
  const [sseConnected, setSseConnected] = useState(false);
  const [sseEvents, setSseEvents] = useState([]);

  // Setup SSE Listener for the payment
  useEffect(() => {
    let eventSource = null;
    try {
      eventSource = new EventSource(`/api/v1/webhook/stream/${paymentID}`);

      eventSource.onopen = () => {
        setSseConnected(true);
      };

      eventSource.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          setSseEvents((prev) => [
            { time: new Date().toLocaleTimeString(), data },
            ...prev,
          ]);
        } catch (e) {
          console.error(e);
        }
      };

      eventSource.onerror = () => {
        setSseConnected(false);
      };
    } catch (e) {
      console.error(e);
    }

    return () => {
      if (eventSource) {
        eventSource.close();
      }
    };
  }, [paymentID]);

  // Trigger Payment Webhook
  const handleSendWebhook = async () => {
    setIsSending(true);
    try {
      const res = await fetch('/api/v1/webhook/payment', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Idempotency-Key': idempotencyKey,
        },
        body: JSON.stringify({
          event_id: idempotencyKey,
          payment_id: paymentID,
          order_id: `ord-${paymentID}`,
          amount: parseFloat(amount),
          currency: 'IDR',
          status: 'SUCCESS',
        }),
      });

      const json = await res.json();
      setWebhookLog((prev) => [
        {
          timestamp: new Date().toLocaleTimeString(),
          status: res.status,
          isDuplicate: json.data?.is_duplicate,
          processedCount: json.data?.event?.processed_count,
          message: json.message,
        },
        ...prev,
      ]);
    } catch (e) {
      console.error(e);
    } finally {
      setIsSending(false);
    }
  };

  // Trigger Fallback Cron Job Reconciler
  const handleTriggerCronReconciler = async () => {
    try {
      const res = await fetch('/api/v1/webhook/reconcile', { method: 'POST' });
      const json = await res.json();
      alert(`Fallback Cron Reconciler: ${json.message} (Diproses: ${json.data.reconciled_count} item)`);
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
      {/* Overview Banner */}
      <div
        className="glass-panel"
        title="Mencegah pemrosesan pembayaran ganda dengan Idempotency Key dan menyiarkan status pembayaran melalui Server-Sent Events."
        style={{ padding: '1.5rem', background: 'linear-gradient(135deg, rgba(15,23,42,0.9), rgba(30,41,59,0.7))' }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginBottom: '0.75rem' }}>
          <Radio size={24} color="#38bdf8" />
          <h2 style={{ fontSize: '1.25rem', fontWeight: 700 }}>
            Soal #3: Webhook Reliability, Idempotency, & Real-Time Sync (SSE vs Polling)
          </h2>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1.5rem' }}>
        {/* Panel 1: Idempotency Tester */}
        <div
          className="glass-panel"
          style={{ padding: '1.5rem' }}
        >
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.25rem' }}>
            <h3 style={{ fontSize: '1.05rem', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              <ShieldAlert size={18} color="var(--accent-cyan)" /> Uji Coba Idempotency Webhook
            </h3>
            <span className="badge badge-accent">Backend Guard</span>
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem', marginBottom: '1.25rem' }}>
            <div>
              <label style={{ fontSize: '0.75rem', color: 'var(--text-secondary)' }}>Header Idempotency-Key:</label>
              <input
                type="text"
                value={idempotencyKey}
                onChange={(e) => setIdempotencyKey(e.target.value)}
                className="mono"
                style={{ width: '100%', padding: '0.5rem 0.75rem', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff', fontSize: '0.85rem' }}
              />
            </div>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.75rem' }}>
              <div>
                <label style={{ fontSize: '0.75rem', color: 'var(--text-secondary)' }}>Payment ID:</label>
                <input
                  type="text"
                  value={paymentID}
                  onChange={(e) => setPaymentID(e.target.value)}
                  className="mono"
                  style={{ width: '100%', padding: '0.5rem 0.75rem', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff', fontSize: '0.85rem' }}
                />
              </div>
              <div>
                <label style={{ fontSize: '0.75rem', color: 'var(--text-secondary)' }}>Jumlah (IDR):</label>
                <input
                  type="number"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                  className="mono"
                  style={{ width: '100%', padding: '0.5rem 0.75rem', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-sm)', color: '#fff', fontSize: '0.85rem' }}
                />
              </div>
            </div>
          </div>

          <div style={{ display: 'flex', gap: '0.75rem', marginBottom: '1.5rem' }}>
            <button
              onClick={handleSendWebhook}
              disabled={isSending}
              className="btn btn-primary"
              style={{ flex: 1 }}
            >
              <Zap size={15} /> Kirim Webhook (Tes Retry/Ganda)
            </button>
            <button
              onClick={() => setIdempotencyKey(`IDEMP-PAY-${Date.now().toString().slice(-6)}`)}
              className="btn btn-secondary"
              title="Generate New Idempotency Key"
            >
              Generate Baru
            </button>
          </div>

          {/* Webhook Execution History */}
          <div>
            <div style={{ fontSize: '0.8rem', fontWeight: 600, marginBottom: '0.5rem', color: 'var(--text-secondary)' }}>
              Log Respons Webhook (Idempotency Logs):
            </div>
            <div style={{ maxHeight: '220px', overflowY: 'auto', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
              {webhookLog.length > 0 ? (
                webhookLog.map((log, i) => (
                  <div
                    key={i}
                    style={{
                      padding: '0.625rem 0.875rem',
                      borderRadius: 'var(--radius-sm)',
                      background: log.isDuplicate ? 'rgba(245, 158, 11, 0.1)' : 'rgba(16, 185, 129, 0.1)',
                      border: `1px solid ${log.isDuplicate ? 'rgba(245, 158, 11, 0.3)' : 'rgba(16, 185, 129, 0.3)'}`,
                      fontSize: '0.8rem',
                    }}
                  >
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.25rem' }}>
                      <span className="mono" style={{ fontWeight: 700, color: log.isDuplicate ? 'var(--accent-amber)' : 'var(--accent-emerald)' }}>
                        HTTP {log.status} {log.isDuplicate ? '• [DUPLIKAT KE-' + log.processedCount + ']' : '• [FIRST TIME PROCESSED]'}
                      </span>
                      <span className="mono" style={{ color: 'var(--text-muted)', fontSize: '0.75rem' }}>{log.timestamp}</span>
                    </div>
                    <div style={{ color: 'var(--text-secondary)', fontSize: '0.75rem' }}>{log.message}</div>
                  </div>
                ))
              ) : (
                <div style={{ padding: '1rem', textAlign: 'center', color: 'var(--text-muted)', fontSize: '0.8rem', background: 'var(--bg-secondary)', borderRadius: 'var(--radius-sm)' }}>
                  Klik tombol di atas untuk mengirim webhook pertama, lalu klik lagi untuk melihat mekanisme pencegahan transaksi ganda!
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Panel 2: SSE vs Polling & Fallback Cron */}
        <div
          className="glass-panel"
          title="SSE menggunakan satu koneksi HTTP satu arah dan mendukung koneksi ulang otomatis, berbeda dengan polling berkala."
          style={{ padding: '1.5rem' }}
        >
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.25rem' }}>
            <h3 style={{ fontSize: '1.05rem', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              <Activity size={18} color="var(--accent-emerald)" /> Server-Sent Events (SSE) Live Sync
            </h3>
            <span className={`badge ${sseConnected ? 'badge-in-stock' : 'badge-low-stock'}`}>
              {sseConnected ? 'SSE Connected' : 'Connecting...'}
            </span>
          </div>

          {/* SSE Events Feed */}
          <div style={{ marginBottom: '1.25rem' }}>
            <div style={{ fontSize: '0.8rem', fontWeight: 600, marginBottom: '0.5rem', display: 'flex', justifyContent: 'space-between' }}>
              <span>Live Inbound SSE Stream (Payment: {paymentID}):</span>
              <span className="mono" style={{ color: 'var(--accent-cyan)' }}>{sseEvents.length} events</span>
            </div>
            <div style={{ maxHeight: '150px', overflowY: 'auto', background: 'var(--bg-primary)', padding: '0.75rem', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)' }}>
              {sseEvents.length > 0 ? (
                sseEvents.map((evt, idx) => (
                  <div key={idx} className="mono" style={{ fontSize: '0.75rem', color: 'var(--accent-emerald)', marginBottom: '0.35rem' }}>
                    [{evt.time}] ⚡ Status: {evt.data.status} (Payment ID: {evt.data.payment_id})
                  </div>
                ))
              ) : (
                <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>Menunggu sinyal broadcast pembayaran dari backend...</div>
              )}
            </div>
          </div>

          {/* Fallback Cron Reconciler */}
          <div
            title="Jika webhook gagal mencapai backend, proses rekonsiliasi memeriksa status pembayaran tertunda melalui API payment gateway."
            style={{ padding: '1rem', background: 'var(--bg-secondary)', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)' }}
          >
            <div style={{ fontWeight: 600, fontSize: '0.875rem', marginBottom: '0.25rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              <Clock size={16} color="var(--accent-purple)" /> Fallback &amp; Retry Mechanism
            </div>
            <div style={{ fontSize: '0.775rem', color: 'var(--text-secondary)', marginBottom: '0.75rem' }}>
              Rekonsiliasi pembayaran tertunda
            </div>
            <button onClick={handleTriggerCronReconciler} className="btn btn-secondary" style={{ width: '100%', fontSize: '0.8rem' }}>
              <RefreshCw size={14} /> Jalankan Simulasi Fallback Cron Reconciler
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
