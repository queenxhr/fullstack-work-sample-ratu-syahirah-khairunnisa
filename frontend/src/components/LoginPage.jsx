import React, { useState } from 'react';
import { Lock, Mail, KeyRound, Sparkles, ShieldCheck, ArrowRight, CheckCircle2, AlertCircle, LogOut } from 'lucide-react';

export default function LoginPage({ onLoginSuccess, currentToken, currentUser, onLogout }) {
  const [email, setEmail] = useState('ratu.syahirah@example.com');
  const [password, setPassword] = useState('password123');
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState('');
  const [successMsg, setSuccessMsg] = useState('');

  // Protected API Test State
  const [testResult, setTestResult] = useState(null);
  const [testLoading, setTestLoading] = useState(false);

  const handleLoginSubmit = async (e) => {
    e?.preventDefault();
    setLoading(true);
    setErrorMsg('');
    setSuccessMsg('');

    try {
      const res = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });

      const json = await res.json();
      if (res.ok) {
        setSuccessMsg(`Login Berhasil! JWT Token diterbitkan dengan JTI: ${json.data.claims.jti}`);
        onLoginSuccess(json.data.token, json.data.claims);
      } else {
        setErrorMsg(json.message || 'Login gagal. Periksa kembali email dan password.');
      }
    } catch (err) {
      setErrorMsg('Gagal terhubung ke backend server (Pastikan Go server aktif di port 8080).');
    } finally {
      setLoading(false);
    }
  };

  const fillDemoAccount = () => {
    setEmail('ratu.syahirah@example.com');
    setPassword('password123');
    setErrorMsg('');
  };

  // Test calling protected profile endpoint with token
  const handleTestProtectedAPI = async () => {
    setTestLoading(true);
    setTestResult(null);
    try {
      const res = await fetch('/api/v1/auth/profile', {
        headers: {
          Authorization: `Bearer ${currentToken}`,
        },
      });
      const json = await res.json();
      setTestResult({
        status: res.status,
        ok: res.ok,
        data: json,
        testedAt: new Date().toLocaleTimeString(),
      });
    } catch (err) {
      setTestResult({
        status: 500,
        ok: false,
        data: { message: 'Network connection error' },
        testedAt: new Date().toLocaleTimeString(),
      });
    } finally {
      setTestLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: '840px', margin: '0 auto', display: 'flex', flexDirection: 'column', gap: '2rem' }}>
      {/* If User is already logged in, show Auth State & Token Testing Playground */}
      {currentToken ? (
        <div
          className="glass-panel"
          title="Menampilkan klaim JWT yang aktif dan menyediakan pengujian endpoint terproteksi."
          style={{ padding: '2rem', border: '1px solid rgba(16, 185, 129, 0.4)' }}
        >
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem', borderBottom: '1px solid var(--border-color)', paddingBottom: '1rem' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
              <div style={{ width: '48px', height: '48px', borderRadius: '50%', background: 'linear-gradient(135deg, var(--accent-emerald), var(--accent-cyan))', display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#fff', fontWeight: 800 }}>
                RS
              </div>
              <div>
                <h2 style={{ fontSize: '1.25rem', fontWeight: 700 }}>Akun Terautentikasi Aktif</h2>
                <div style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>
                  {currentUser?.email} • <span className="badge badge-in-stock">{currentUser?.role}</span>
                </div>
              </div>
            </div>
            <button onClick={onLogout} className="btn btn-danger" style={{ fontSize: '0.825rem' }}>
              <LogOut size={16} /> Logout &amp; Revoke Token (Redis Blacklist)
            </button>
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1.5rem', marginBottom: '1.5rem' }}>
            <div style={{ background: 'var(--bg-secondary)', padding: '1.25rem', borderRadius: 'var(--radius-md)' }}>
              <h3 style={{ fontSize: '0.95rem', fontWeight: 600, color: 'var(--accent-cyan)', marginBottom: '0.75rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <KeyRound size={16} /> Informasi JWT Claims:
              </h3>
              <div style={{ fontSize: '0.8rem', display: 'flex', flexDirection: 'column', gap: '0.35rem' }}>
                <div>User ID: <span className="mono" style={{ color: 'var(--text-primary)' }}>{currentUser?.user_id}</span></div>
                <div>Unique JTI: <span className="mono" style={{ color: 'var(--accent-cyan)' }}>{currentUser?.jti}</span></div>
                <div>Role: <span className="mono" style={{ color: 'var(--accent-emerald)' }}>{currentUser?.role}</span></div>
                <div>Expires At: <span className="mono">{currentUser?.exp ? new Date(currentUser.exp * 1000).toLocaleTimeString() : '-'}</span></div>
              </div>
            </div>

            <div
              title="Memanggil GET /api/v1/auth/profile menggunakan access token pada header Authorization."
              style={{ background: 'var(--bg-secondary)', padding: '1.25rem', borderRadius: 'var(--radius-md)' }}
            >
              <h3 style={{ fontSize: '0.95rem', fontWeight: 600, color: 'var(--accent-purple)', marginBottom: '0.75rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <ShieldCheck size={16} /> Pengujian Endpoint Terproteksi:
              </h3>
              <p style={{ fontSize: '0.775rem', color: 'var(--text-secondary)', marginBottom: '0.75rem' }}>
                Uji akses ke endpoint <code>GET /api/v1/auth/profile</code> dengan Authorization Bearer header.
              </p>
              <button
                onClick={handleTestProtectedAPI}
                disabled={testLoading}
                className="btn btn-primary"
                style={{ width: '100%', fontSize: '0.8rem' }}
              >
                {testLoading ? 'Menguji API...' : '🚀 Tes Akses Endpoint Terproteksi'}
              </button>
            </div>
          </div>

          {/* Test Result Display */}
          {testResult && (
            <div
              style={{
                padding: '1rem',
                borderRadius: 'var(--radius-md)',
                background: testResult.ok ? 'rgba(16, 185, 129, 0.1)' : 'rgba(244, 63, 94, 0.1)',
                border: `1px solid ${testResult.ok ? 'rgba(16, 185, 129, 0.3)' : 'rgba(244, 63, 94, 0.3)'}`,
              }}
            >
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem' }}>
                <span className="mono" style={{ fontWeight: 700, color: testResult.ok ? 'var(--accent-emerald)' : 'var(--accent-rose)' }}>
                  HTTP {testResult.status} {testResult.ok ? '• Akses Diterima (Authorized)' : '• Akses Ditolak (Unauthorized)'}
                </span>
                <span className="mono" style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>{testResult.testedAt}</span>
              </div>
              <pre className="mono" style={{ fontSize: '0.75rem', background: 'var(--bg-primary)', padding: '0.5rem', borderRadius: '4px', overflowX: 'auto', color: testResult.ok ? '#86efac' : '#fda4af' }}>
                {JSON.stringify(testResult.data, null, 2)}
              </pre>
            </div>
          )}
        </div>
      ) : (
        /* Login Form */
        <div className="glass-panel" style={{ padding: '2.5rem', position: 'relative', overflow: 'hidden' }}>
          <div style={{ textAlign: 'center', marginBottom: '2rem' }}>
            <div style={{ width: '56px', height: '56px', margin: '0 auto 1rem auto', borderRadius: 'var(--radius-md)', background: 'linear-gradient(135deg, var(--accent-blue), var(--accent-cyan))', display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#fff', boxShadow: '0 0 20px rgba(56, 189, 248, 0.4)' }}>
              <Lock size={28} />
            </div>
            <h2
              title="Login menerbitkan token JWT dengan JTI unik untuk pengujian autentikasi dan pencabutan token."
              style={{ fontSize: '1.5rem', fontWeight: 700, letterSpacing: '-0.02em', marginBottom: '0.35rem' }}
            >
              Autentikasi Stateless JWT &amp; Revocation
            </h2>
            <p style={{ color: 'var(--text-secondary)', fontSize: '0.875rem' }}>
              Masuk untuk mendapatkan token JWT ber-JTI unik untuk pengujian Soal #2
            </p>
          </div>

          <form onSubmit={handleLoginSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '1.25rem' }}>
            <div>
              <label style={{ display: 'block', fontSize: '0.8rem', fontWeight: 600, color: 'var(--text-secondary)', marginBottom: '0.4rem' }}>
                Email Akun:
              </label>
              <div style={{ position: 'relative' }}>
                <Mail size={16} style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: 'var(--text-muted)' }} />
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                  placeholder="nama@perusahaan.com"
                  style={{ width: '100%', padding: '0.65rem 0.75rem 0.65rem 2.5rem', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-md)', color: '#fff', fontSize: '0.875rem', outline: 'none' }}
                />
              </div>
            </div>

            <div>
              <label style={{ display: 'block', fontSize: '0.8rem', fontWeight: 600, color: 'var(--text-secondary)', marginBottom: '0.4rem' }}>
                Password:
              </label>
              <div style={{ position: 'relative' }}>
                <Lock size={16} style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: 'var(--text-muted)' }} />
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  placeholder="••••••••"
                  style={{ width: '100%', padding: '0.65rem 0.75rem 0.65rem 2.5rem', background: 'var(--bg-secondary)', border: '1px solid var(--border-color)', borderRadius: 'var(--radius-md)', color: '#fff', fontSize: '0.875rem', outline: 'none' }}
                />
              </div>
            </div>

            {errorMsg && (
              <div style={{ padding: '0.75rem', background: 'rgba(244, 63, 94, 0.1)', border: '1px solid rgba(244, 63, 94, 0.3)', borderRadius: 'var(--radius-md)', color: '#fb7185', fontSize: '0.8rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <AlertCircle size={16} /> {errorMsg}
              </div>
            )}

            {successMsg && (
              <div style={{ padding: '0.75rem', background: 'rgba(16, 185, 129, 0.1)', border: '1px solid rgba(16, 185, 129, 0.3)', borderRadius: 'var(--radius-md)', color: '#86efac', fontSize: '0.8rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <CheckCircle2 size={16} /> {successMsg}
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="btn btn-primary"
              style={{ width: '100%', padding: '0.75rem', marginTop: '0.5rem' }}
            >
              {loading ? 'Mengautentikasi...' : (
                <>
                  Masuk &amp; Dapatkan JWT Token <ArrowRight size={16} />
                </>
              )}
            </button>
          </form>

          {/* 1-Click Demo Fill */}
          <div style={{ marginTop: '1.5rem', paddingTop: '1.25rem', borderTop: '1px solid var(--border-color)', textAlign: 'center' }}>
            <button
              onClick={fillDemoAccount}
              className="btn btn-secondary"
              style={{ fontSize: '0.75rem', padding: '0.4rem 0.85rem' }}
            >
              <Sparkles size={13} color="var(--accent-cyan)" /> Gunakan Akun Demo (Ratu Syahirah Khairunnisa)
            </button>
            <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)', marginTop: '0.5rem' }}>
              Demo: <code>ratu.syahirah@example.com</code> / <code>password123</code>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
