import React, { useState, useEffect, useCallback, useMemo, useRef } from 'react';
import ProductTableOptimized from './components/ProductTableOptimized';
import ProductTableUnoptimized from './components/ProductTableUnoptimized';
import FlashSaleConcurrencyDemo from './components/FlashSaleConcurrencyDemo';
import WebhookPaymentResilience from './components/WebhookPaymentResilience';
import ArchitectureDiagrams from './components/ArchitectureDiagrams';
import LoginPage from './components/LoginPage';
import {
  Zap,
  Flame,
  Radio,
  BarChart3,
  Search,
  RefreshCw,
  SlidersHorizontal,
  Package,
  Layers,
  Sparkles,
  Activity,
  AlertCircle,
  Lock,
  User,
  LogOut
} from 'lucide-react';

export default function App() {
  const [activeTab, setActiveTab] = useState('soal1');
  const [tableMode, setTableMode] = useState('optimized'); // 'optimized' or 'unoptimized'
  const [products, setProducts] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [searchQuery, setSearchQuery] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [sseConnected, setSseConnected] = useState(false);

  // Authentication & JWT Session State
  const [jwtToken, setJwtToken] = useState(() => localStorage.getItem('jwt_token') || '');
  const [currentUser, setCurrentUser] = useState(() => {
    try {
      const saved = localStorage.getItem('jwt_user');
      return saved ? JSON.parse(saved) : null;
    } catch {
      return null;
    }
  });

  const handleLoginSuccess = (token, claims) => {
    setJwtToken(token);
    setCurrentUser(claims);
    localStorage.setItem('jwt_token', token);
    localStorage.setItem('jwt_user', JSON.stringify(claims));
  };

  const handleLogout = async () => {
    if (jwtToken) {
      try {
        await fetch('/api/v1/auth/logout', {
          method: 'POST',
          headers: { Authorization: `Bearer ${jwtToken}` },
        });
      } catch (e) {}
    }
    setJwtToken('');
    setCurrentUser(null);
    localStorage.removeItem('jwt_token');
    localStorage.removeItem('jwt_user');
  };

  // FPS & Performance tracking
  const [fps, setFps] = useState(60);
  const fpsRef = useRef({ frames: 0, lastTime: performance.now() });

  // Measure FPS
  useEffect(() => {
    let animId;
    const calcFps = (time) => {
      fpsRef.current.frames++;
      if (time >= fpsRef.current.lastTime + 1000) {
        setFps(Math.round((fpsRef.current.frames * 1000) / (time - fpsRef.current.lastTime)));
        fpsRef.current.frames = 0;
        fpsRef.current.lastTime = time;
      }
      animId = requestAnimationFrame(calcFps);
    };
    animId = requestAnimationFrame(calcFps);
    return () => cancelAnimationFrame(animId);
  }, []);

  // Fetch initial 5,000 products from Golang Backend (or fallback mock if backend is still booting)
  const fetchProducts = useCallback(async () => {
    setIsLoading(true);
    try {
      const res = await fetch('/api/v1/products?limit=5000');
      if (res.ok) {
        const json = await res.json();
        setProducts(json.data.products || []);
        setTotalCount(json.data.total || json.data.products?.length || 0);
      } else {
        generateMockProducts();
      }
    } catch (err) {
      console.warn('Backend connecting...', err);
      generateMockProducts();
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Generate 5,000 client-side fallback products if backend is not started yet
  const generateMockProducts = () => {
    const list = [];
    const categories = ['Smartphones', 'Laptops', 'Audio & Wearables', 'Gaming & PC', 'Accessories', 'Home Appliances'];
    const brands = ['Apple', 'Sony', 'Logitech', 'Samsung', 'Asus', 'Dell', 'Bose', 'Xiaomi'];
    for (let i = 1; i <= 5000; i++) {
      const stock = i === 1 ? 1 : Math.floor(Math.random() * 150) + 1;
      let status = 'IN_STOCK';
      if (stock <= 0) status = 'OUT_OF_STOCK';
      else if (stock <= 5) status = 'LOW_STOCK';

      list.push({
        id: i,
        sku: `SKU-${String(i).padStart(5, '0')}`,
        name: i === 1 ? '🔥 [FLASH SALE] Sony PlayStation 5 Pro Limited' : `${brands[i % brands.length]} Model Pro #${i}`,
        category: categories[i % categories.length],
        price: (Math.floor(Math.random() * 1000) + 15) * 10000,
        stock: stock,
        status: status,
        version: 1,
      });
    }
    setProducts(list);
    setTotalCount(5000);
  };

  useEffect(() => {
    fetchProducts();
  }, [fetchProducts]);

  // Connect to SSE Server-Sent Events for real-time stock updates
  useEffect(() => {
    let eventSource = null;
    try {
      eventSource = new EventSource('/api/v1/products/stream');
      eventSource.onopen = () => setSseConnected(true);
      eventSource.onmessage = (event) => {
        try {
          const payload = JSON.parse(event.data);
          if (payload.type === 'STOCK_UPDATED' && payload.product) {
            const updated = payload.product;
            // Granular update: update only the matching item
            setProducts((prev) =>
              prev.map((item) => (item.id === updated.id ? { ...item, ...updated } : item))
            );
          }
        } catch (e) {
          console.error(e);
        }
      };
      eventSource.onerror = () => setSseConnected(false);
    } catch (e) {
      console.error(e);
    }

    return () => {
      if (eventSource) eventSource.close();
    };
  }, []);

  // Update Stock Callback: useCallback ensures referential stability!
  const handleUpdateStock = useCallback(async (id, newStock) => {
    // 1. Optimistic UI update in React state
    setProducts((prev) =>
      prev.map((item) => {
        if (item.id === id) {
          let status = 'IN_STOCK';
          if (newStock <= 0) status = 'OUT_OF_STOCK';
          else if (newStock <= 5) status = 'LOW_STOCK';
          return { ...item, stock: newStock, status };
        }
        return item;
      })
    );

    // 2. Persist to Backend API
    try {
      await fetch(`/api/v1/products/${id}/stock`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ stock: newStock }),
      });
    } catch (e) {
      console.warn('Backend update error, optimistic state preserved');
    }
  }, []);

  // Simulate real-time random stock fluctuation (e.g. 5 random items updated per second)
  const [isSimulatingSpike, setIsSimulatingSpike] = useState(false);
  useEffect(() => {
    if (!isSimulatingSpike) return;
    const interval = setInterval(() => {
      const randomId = Math.floor(Math.random() * Math.min(products.length, 500)) + 1;
      const randomStock = Math.floor(Math.random() * 80);
      handleUpdateStock(randomId, randomStock);
    }, 250);
    return () => clearInterval(interval);
  }, [isSimulatingSpike, products.length, handleUpdateStock]);

  // Derived stock statistics
  const stats = useMemo(() => {
    let inStock = 0;
    let lowStock = 0;
    let outStock = 0;
    for (const p of products) {
      if (p.stock <= 0) outStock++;
      else if (p.stock <= 5) lowStock++;
      else inStock++;
    }
    return { inStock, lowStock, outStock };
  }, [products]);

  return (
    <div className="app-container">
      {/* Sticky Top Header */}
      <header className="glass-panel app-header">
        <div className="logo-group">
          <div className="logo-badge">
            <Sparkles size={22} />
          </div>
          <div>
            <h1 className="app-title">Senior Fullstack Work Sample</h1>
            <p className="app-subtitle" title="React 18, Vite, react-window, dan Golang dengan Clean Architecture.">
              Kandidat: <strong>Ratu Syahirah Khairunnisa</strong>
            </p>
          </div>
        </div>

        {/* Navigation Tabs */}
        <nav className="tabs-nav">
          <button
            onClick={() => setActiveTab('soal1')}
            className={`tab-btn ${activeTab === 'soal1' ? 'active' : ''}`}
          >
            <Zap size={16} /> Soal #1: React Optimization (5k+)
          </button>
          <button
            onClick={() => setActiveTab('soal2')}
            className={`tab-btn ${activeTab === 'soal2' ? 'active' : ''}`}
          >
            <Flame size={16} /> Soal #2: Concurrency &amp; JWT
          </button>
          <button
            onClick={() => setActiveTab('soal3')}
            className={`tab-btn ${activeTab === 'soal3' ? 'active' : ''}`}
          >
            <Radio size={16} /> Soal #3: Webhook &amp; Resilience
          </button>
          <button
            onClick={() => setActiveTab('soal45')}
            className={`tab-btn ${activeTab === 'soal45' ? 'active' : ''}`}
          >
            <BarChart3 size={16} /> Soal #4 &amp; #5: TCO &amp; Metrics
          </button>
          <button
            onClick={() => setActiveTab('login')}
            className={`tab-btn ${activeTab === 'login' ? 'active' : ''}`}
            style={{ color: jwtToken ? 'var(--accent-emerald)' : 'var(--text-secondary)' }}
          >
            <Lock size={16} /> {jwtToken ? 'Akun & Token Lab' : 'Halaman Login'}
          </button>
        </nav>

        {/* User Auth Chip */}
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
          {jwtToken ? (
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', background: 'var(--bg-secondary)', padding: '0.35rem 0.75rem', borderRadius: 'var(--radius-full)', border: '1px solid var(--border-color)' }}>
              <span className="render-pulse"></span>
              <span style={{ fontSize: '0.8rem', fontWeight: 600 }}>{currentUser?.email?.split('@')[0]}</span>
              <button onClick={handleLogout} className="btn btn-secondary" style={{ padding: '2px 8px', fontSize: '0.7rem' }} title="Logout">
                <LogOut size={12} />
              </button>
            </div>
          ) : (
            <button onClick={() => setActiveTab('login')} className="btn btn-primary" style={{ padding: '0.4rem 0.85rem', fontSize: '0.8rem' }}>
              <User size={14} /> Masuk / Login
            </button>
          )}
        </div>
      </header>

      {/* Metrics Banner */}
      <div className="metrics-grid">
        <div className="glass-panel metric-card" title="Dataset produk e-commerce dengan pembaruan stok real-time.">
          <div className="metric-header">
            <span>Total Dataset Produk</span>
            <Package size={16} color="var(--accent-cyan)" />
          </div>
          <div className="metric-value mono">
            {totalCount.toLocaleString('id-ID')} items
          </div>
        </div>

        <div className="glass-panel metric-card" title={tableMode === 'optimized' ? 'Virtualisasi react-window menjaga jumlah baris DOM yang aktif tetap kecil.' : 'Mode ini merender banyak baris sekaligus dan dapat menyebabkan layout thrashing.'}>
          <div className="metric-header">
            <span>Real-time Rendering FPS</span>
            <Activity size={16} color={fps >= 50 ? 'var(--accent-emerald)' : 'var(--accent-rose)'} />
          </div>
          <div className="metric-value mono" style={{ color: fps >= 50 ? 'var(--accent-emerald)' : 'var(--accent-rose)' }}>
            {fps} FPS
          </div>
        </div>

        <div className="glass-panel metric-card" title="Jumlah produk menurut kondisi stok terkini.">
          <div className="metric-header">
            <span>Ketersediaan Stok</span>
            <Layers size={16} color="var(--accent-emerald)" />
          </div>
          <div style={{ display: 'flex', gap: '0.75rem', marginTop: '0.25rem' }}>
            <span className="badge badge-in-stock">{stats.inStock} In Stock</span>
            <span className="badge badge-low-stock">{stats.lowStock} Low</span>
            <span className="badge badge-out-stock">{stats.outStock} Out</span>
          </div>
        </div>

        <div className="glass-panel metric-card" title="Koneksi Server-Sent Events untuk sinkronisasi pembaruan produk.">
          <div className="metric-header">
            <span>Backend Real-time Sync (SSE)</span>
            <span className={`render-pulse ${sseConnected ? '' : 'lag'}`}></span>
          </div>
          <div className="metric-value mono" style={{ fontSize: '1.25rem', color: sseConnected ? 'var(--accent-emerald)' : 'var(--accent-amber)' }}>
            {sseConnected ? 'SSE Connected (Port 8080)' : 'Polling / Mock Mode'}
          </div>
        </div>
      </div>

      {/* Main Content Area based on Tab */}
      <main>
        {activeTab === 'soal1' && (
          <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
            {/* Control Bar: Mode Toggle, Search, Spike Simulation */}
            <div className="glass-panel" style={{ padding: '1rem 1.25rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '1rem' }}>
              {/* Search Box */}
              <div style={{ position: 'relative', width: '320px' }}>
                <Search size={16} style={{ position: 'absolute', left: '12px', top: '50%', transform: 'translateY(-50%)', color: 'var(--text-muted)' }} />
                <input
                  type="text"
                  placeholder="Cari dari 5.000+ produk (Nama / SKU)..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  style={{
                    width: '100%',
                    padding: '0.55rem 0.75rem 0.55rem 2.25rem',
                    background: 'var(--bg-secondary)',
                    border: '1px solid var(--border-color)',
                    borderRadius: 'var(--radius-md)',
                    color: '#ffffff',
                    fontSize: '0.85rem',
                    outline: 'none',
                  }}
                />
              </div>

              {/* Action Buttons */}
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
                <button
                  onClick={() => setIsSimulatingSpike(!isSimulatingSpike)}
                  className={`btn ${isSimulatingSpike ? 'btn-secondary' : 'btn-primary'}`}
                  style={{ fontSize: '0.8rem' }}
                >
                  <Activity size={14} />
                  {isSimulatingSpike ? 'Hentikan Simulasi Real-Time' : 'Simulasi Real-Time Stock Spike'}
                </button>

                {/* Comparison Mode Switcher */}
                <div style={{ display: 'flex', background: 'var(--bg-secondary)', padding: '3px', borderRadius: 'var(--radius-md)', border: '1px solid var(--border-color)' }}>
                  <button
                    onClick={() => setTableMode('optimized')}
                    className={`btn ${tableMode === 'optimized' ? 'btn-primary' : 'btn-secondary'}`}
                    style={{ padding: '0.4rem 0.85rem', fontSize: '0.775rem' }}
                  >
                    <Zap size={13} /> ⚡ Optimized (react-window)
                  </button>
                  <button
                    onClick={() => setTableMode('unoptimized')}
                    className={`btn ${tableMode === 'unoptimized' ? 'btn-primary' : 'btn-secondary'}`}
                    style={{ padding: '0.4rem 0.85rem', fontSize: '0.775rem' }}
                  >
                    <AlertCircle size={13} /> 🐢 Unoptimized (Direct DOM)
                  </button>
                </div>
              </div>
            </div>

            {/* Render Table */}
            {isLoading ? (
              <div className="glass-panel" style={{ padding: '4rem', textAlign: 'center', color: 'var(--text-muted)' }}>
                <RefreshCw size={24} className="spin" style={{ marginBottom: '1rem' }} />
                <div>Memuat 5.000+ data produk...</div>
              </div>
            ) : tableMode === 'optimized' ? (
              <ProductTableOptimized
                products={products}
                onUpdateStock={handleUpdateStock}
                searchQuery={searchQuery}
                totalItems={totalCount}
              />
            ) : (
              <ProductTableUnoptimized
                products={products}
                onUpdateStock={handleUpdateStock}
                searchQuery={searchQuery}
              />
            )}
          </div>
        )}

        {activeTab === 'soal2' && <FlashSaleConcurrencyDemo />}
        {activeTab === 'soal3' && <WebhookPaymentResilience />}
        {activeTab === 'soal45' && <ArchitectureDiagrams />}
        {activeTab === 'login' && (
          <LoginPage
            onLoginSuccess={handleLoginSuccess}
            currentToken={jwtToken}
            currentUser={currentUser}
            onLogout={handleLogout}
          />
        )}
      </main>
    </div>
  );
}
