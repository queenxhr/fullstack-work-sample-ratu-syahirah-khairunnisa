import React, { useState, useEffect } from 'react';
import { Flame, ShieldCheck, Zap, Lock, RefreshCw, KeyRound, CheckCircle, XCircle, AlertTriangle } from 'lucide-react';

export default function FlashSaleConcurrencyDemo() {
  const [method, setMethod] = useState('ATOMIC_UPDATE');
  const [productStock, setProductStock] = useState(1);
  const [loading, setLoading] = useState(false);
  const [benchmarkResult, setBenchmarkResult] = useState(null);
  const [orders, setOrders] = useState([]);

  // Auth / JWT Revocation Demo State
  const [jwtToken, setJwtToken] = useState('');
  const [userProfile, setUserProfile] = useState(null);
  const [authMessage, setAuthMessage] = useState('');
  const [authLoading, setAuthLoading] = useState(false);
  const [blacklistedTokens, setBlacklistedTokens] = useState([]);

  const fetchStock = async () => {
    try {
      const res = await fetch('/api/v1/flashsale/stock/1');
      if (res.ok) {
        const json = await res.json();
        setProductStock(json.data.stock);
      }
    } catch (e) {
      console.error(e);
    }
  };

  const fetchOrders = async () => {
    try {
      const res = await fetch('/api/v1/flashsale/orders');
      if (res.ok) {
        const json = await res.json();
        setOrders(json.data || []);
      }
    } catch (e) {
      console.error(e);
    }
  };

  const fetchBlacklist = async () => {
    try {
      const res = await fetch('/api/v1/auth/blacklist');
      if (res.ok) {
        const json = await res.json();
        setBlacklistedTokens(json.data?.revoked_jtis || []);
      }
    } catch (e) {
      console.error(e);
    }
  };

  useEffect(() => {
    fetchStock();
    fetchOrders();
    fetchBlacklist();
  }, []);

  const handleResetStock = async (stock = 1) => {
    try {
      const res = await fetch('/api/v1/flashsale/reset-stock', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ product_id: 1, stock }),
      });
      if (res.ok) {
        await fetchStock();
        setBenchmarkResult(null);
      }
    } catch (e) {
      console.error(e);
    }
  };

  const handleRun100Benchmark = async () => {
    setLoading(true);
    setBenchmarkResult(null);
    try {
      const res = await fetch('/api/v1/flashsale/benchmark-100', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ product_id: 1, method }),
      });
      const json = await res.json();
      if (res.ok) {
        setBenchmarkResult(json.data);
        await fetchStock();
        await fetchOrders();
      }
    } catch (e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  };

  // JWT Flow: Login
  const handleLogin = async () => {
    setAuthLoading(true);
    setAuthMessage('');
    try {
      const res = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: 'ratu.syahirah@example.com', password: 'password123' }),
      });
      const json = await res.json();
      if (res.ok) {
        setJwtToken(json.data.token);
        setUserProfile(json.data.claims);
        setAuthMessage(`JWT Terbit! JTI: ${json.data.claims.jti}`);
      } else {
        setAuthMessage(json.message || 'Login failed');
      }
    } catch (e) {
      setAuthMessage('Connection error');
    } finally {
      setAuthLoading(false);
    }
  };

  // JWT Flow: Check Protected Profile
  const handleCheckProfile = async () => {
    setAuthLoading(true);
    try {
      const res = await fetch('/api/v1/auth/profile', {
        headers: { Authorization: `Bearer ${jwtToken}` },
      });
      const json = await res.json();
      if (res.ok) {
        setAuthMessage(`Akses Diterima! Token valid dan belum masuk blacklist.`);
      } else {
        setAuthMessage(`Akses Ditolak (401 Unauthorized): ${json.message}`);
      }
    } catch (e) {
      setAuthMessage('Profile check error');
    } finally {
      setAuthLoading(false);
    }
  };

  // JWT Flow: Logout & Revoke Token
  const handleLogout = async () => {
    setAuthLoading(true);
    try {
      const res = await fetch('/api/v1/auth/logout', {
        method: 'POST',
        headers: { Authorization: `Bearer ${jwtToken}` },
      });
      const json = await res.json();
      if (res.ok) {
        setAuthMessage('Logout Berhasil! Token JTI telah masuk ke Redis Blacklist.');
        await fetchBlacklist();
      } else {
        setAuthMessage(json.message);
      }
    } catch (e) {
      setAuthMessage('Logout error');
    } finally {
      setAuthLoading(false);
    }
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
      {/* Top Banner: Case Overview */}
      <div
        className="glass-panel"
        title="Simulasi 100 permintaan pembelian paralel untuk satu stok. Backend menggunakan Atomic SQL Update atau Pessimistic Lock untuk mencegah overselling."
        style={{ padding: '1.5rem', background: 'linear-gradient(135deg, rgba(15,23,42,0.9), rgba(30,41,59,0.7))' }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', marginBottom: '0.75rem' }}>
          <Flame size={24} color="var(--accent-cyan)" />
          <h2 style={{ fontSize: '1.25rem', fontWeight: 700 }}>
            Soal #2: High-Concurrency Scaling & Race Condition Handling (Golang + DB Locking)
          </h2>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1.5rem' }}>
        {/* Panel 1: Race Condition Simulator */}
        <div
          className="glass-panel"
          style={{ padding: '1.5rem' }}
        >
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.25rem' }}>
            <h3 style={{ fontSize: '1.05rem', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              <Zap size={18} color="var(--accent-amber)" /> Simulasi Flash Sale 100 Concurrency
            </h3>
            <span className="badge badge-accent">Produk #1</span>
          </div>

          <div style={{ padding: '1rem', background: 'var(--bg-secondary)', borderRadius: 'var(--radius-md)', marginBottom: '1.25rem' }}>
            <div style={{ fontSize: '0.8rem', color: 'var(--text-secondary)', marginBottom: '0.25rem' }}>Nama Produk Flash Sale:</div>
            <div style={{ fontWeight: 600, color: 'var(--accent-cyan)', marginBottom: '0.5rem' }}>
              🔥 Sony PlayStation 5 Pro Limited (Promo Flash Sale)
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <span style={{ fontSize: '0.875rem' }}>Sisa Stok Saat Ini:</span>
              <span className="mono" style={{ fontSize: '1.25rem', fontWeight: 800, color: productStock > 0 ? 'var(--accent-emerald)' : 'var(--accent-rose)' }}>
                {productStock} unit
              </span>
            </div>
          </div>

          {/* Strategy Selector */}
          <div style={{ marginBottom: '1.25rem' }}>
            <label style={{ fontSize: '0.8rem', color: 'var(--text-secondary)', display: 'block', marginBottom: '0.5rem' }}>
              Pilih Strategi Database Locking:
            </label>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.5rem' }}>
              <button
                onClick={() => setMethod('ATOMIC_UPDATE')}
                className={`btn ${method === 'ATOMIC_UPDATE' ? 'btn-primary' : 'btn-secondary'}`}
                style={{ fontSize: '0.775rem' }}
              >
                <ShieldCheck size={14} /> Atomic SQL Update
              </button>
              <button
                onClick={() => setMethod('PESSIMISTIC_LOCK')}
                className={`btn ${method === 'PESSIMISTIC_LOCK' ? 'btn-primary' : 'btn-secondary'}`}
                style={{ fontSize: '0.775rem' }}
              >
                <Lock size={14} /> Pessimistic Lock
              </button>
            </div>
          </div>

          {/* Actions */}
          <div style={{ display: 'flex', gap: '0.75rem', marginBottom: '1.5rem' }}>
            <button
              onClick={handleRun100Benchmark}
              disabled={loading}
              className="btn btn-primary"
              style={{ flex: 1, padding: '0.75rem' }}
            >
              {loading ? (
                <>Menjalankan 100 Goroutines...</>
              ) : (
                <>
                  <Flame size={16} /> Jalankan 100 Request Bersamaan
                </>
              )}
            </button>
            <button
              onClick={() => handleResetStock(1)}
              className="btn btn-secondary"
              title="Reset stok menjadi 1"
            >
              <RefreshCw size={14} /> Reset Stok = 1
            </button>
          </div>

          {/* Benchmark Results */}
          {benchmarkResult && (
            <div style={{ padding: '1.25rem', background: 'rgba(16, 185, 129, 0.08)', border: '1px solid rgba(16, 185, 129, 0.3)', borderRadius: 'var(--radius-md)' }}>
              <div style={{ fontWeight: 700, color: 'var(--accent-emerald)', display: 'flex', alignItems: 'center', gap: '0.5rem', marginBottom: '0.75rem' }}>
                <CheckCircle size={18} /> Hasil Pengujian Concurrency (Zero Overselling!)
              </div>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.5rem', fontSize: '0.85rem' }}>
                <div>Total Request: <strong>{benchmarkResult.total_requests} users</strong></div>
                <div>Metode: <strong>{benchmarkResult.method_used}</strong></div>
                <div style={{ color: 'var(--accent-emerald)' }}>Berhasil Beli: <strong>{benchmarkResult.successful_orders} order</strong></div>
                <div style={{ color: 'var(--accent-rose)' }}>Ditolak (409 Out of Stock): <strong>{benchmarkResult.rejected_out_of_stock}</strong></div>
                <div>Sisa Stok Akhir: <strong>{benchmarkResult.remaining_stock}</strong></div>
                <div>Waktu Eksekusi: <strong>{benchmarkResult.duration_ms} ms</strong></div>
              </div>
              {benchmarkResult.successful_user_ids?.length > 0 && (
                <div style={{ marginTop: '0.75rem', fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
                  Pemenang Pembelian: <span className="mono" style={{ color: 'var(--accent-cyan)' }}>{benchmarkResult.successful_user_ids[0]}</span>
                </div>
              )}
            </div>
          )}
        </div>

        {/* Panel 2: JWT Stateless Token Revocation Demo */}
        <div
          className="glass-panel"
          title="JWT tetap stateless, sementara token dapat dicabut segera dengan menyimpan JTI pada blacklist hingga token kedaluwarsa."
          style={{ padding: '1.5rem' }}
        >
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.25rem' }}>
            <h3 style={{ fontSize: '1.05rem', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              <KeyRound size={18} color="var(--accent-cyan)" /> Stateless JWT & Redis Blacklist Revocation
            </h3>
            <span className="badge badge-accent">Security</span>
          </div>

          <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '1rem' }}>
            <button onClick={handleLogin} disabled={authLoading} className="btn btn-primary" style={{ flex: 1 }}>
              1. Login & Terbitkan JWT
            </button>
            <button onClick={handleCheckProfile} disabled={!jwtToken || authLoading} className="btn btn-secondary">
              2. Tes Akses Profile
            </button>
            <button onClick={handleLogout} disabled={!jwtToken || authLoading} className="btn btn-danger">
              3. Logout (Revoke)
            </button>
          </div>

          {authMessage && (
            <div style={{ padding: '0.75rem 1rem', background: 'var(--bg-secondary)', borderRadius: 'var(--radius-md)', fontSize: '0.8rem', border: '1px solid var(--border-color)', marginBottom: '1rem' }}>
              {authMessage}
            </div>
          )}

          {jwtToken && (
            <div style={{ marginBottom: '1rem' }}>
              <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '0.25rem' }}>Active JWT Token (Bearer):</div>
              <div className="mono" style={{ fontSize: '0.7rem', wordBreak: 'break-all', background: 'var(--bg-primary)', padding: '0.5rem', borderRadius: '4px', color: 'var(--accent-cyan)' }}>
                {jwtToken}
              </div>
            </div>
          )}

          {/* Blacklisted JTIs */}
          <div style={{ marginTop: '1rem' }}>
            <div style={{ fontSize: '0.8rem', fontWeight: 600, marginBottom: '0.5rem', display: 'flex', justifyContent: 'space-between' }}>
              <span>Redis Blacklist Cache (Active Revoked JTIs):</span>
              <span className="mono" style={{ color: 'var(--accent-rose)' }}>{blacklistedTokens.length} active</span>
            </div>
            <div style={{ maxHeight: '120px', overflowY: 'auto', background: 'var(--bg-primary)', padding: '0.5rem', borderRadius: 'var(--radius-md)' }}>
              {blacklistedTokens.length > 0 ? (
                blacklistedTokens.map((jti, i) => (
                  <div key={i} className="mono" style={{ fontSize: '0.7rem', color: 'var(--accent-rose)', padding: '2px 0' }}>
                    ⛔ jti: {jti}
                  </div>
                ))
              ) : (
                <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>Belum ada token yang di-revoke.</div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
