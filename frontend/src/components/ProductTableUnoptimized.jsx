import React, { useRef } from 'react';
import { AlertTriangle, Clock } from 'lucide-react';

/**
 * ProductTableUnoptimized
 * Demonstrates the severe lag, layout thrashing, and unnecessary re-render cascade
 * when rendering 5,000+ items directly into the DOM without virtualization or memoization.
 */
export default function ProductTableUnoptimized({
  products,
  onUpdateStock,
  searchQuery,
}) {
  const filteredProducts = products.filter(
    (p) =>
      !searchQuery ||
      p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.sku.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="glass-panel" style={{ overflow: 'hidden', border: '1px solid rgba(244, 63, 94, 0.4)' }}>
      {/* Warning Banner */}
      <div
        title={`Mode tanpa virtualisasi: perubahan stok dapat memicu render ulang hingga ${filteredProducts.length.toLocaleString('id-ID')} baris.`}
        style={{ padding: '1rem 1.25rem', background: 'rgba(244, 63, 94, 0.1)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', color: 'var(--accent-rose)' }}>
          <AlertTriangle size={18} />
          <span style={{ fontSize: '0.875rem', fontWeight: 600 }}>
            Mode Unoptimized: Merender {filteredProducts.length.toLocaleString('id-ID')} DOM Element Sekaligus!
          </span>
        </div>
        <div style={{ fontSize: '0.75rem', color: 'var(--accent-rose)', display: 'flex', alignItems: 'center', gap: '0.35rem' }}>
          <span className="render-pulse lag"></span>
          Render massal aktif
        </div>
      </div>

      {/* Header Row */}
      <div className="table-header">
        <span>ID</span>
        <span>SKU</span>
        <span>Nama Produk</span>
        <span>Kategori</span>
        <span>Harga</span>
        <span>Stok</span>
        <span>Status</span>
        <span>Aksi Stok (Un-memoized)</span>
      </div>

      {/* Unvirtualized Native Scroll Container */}
      <div style={{ maxHeight: '600px', overflowY: 'auto' }}>
        {filteredProducts.slice(0, 1000).map((product) => (
          // Inline callback without memoization causes full cascade!
          <div key={product.id} className="product-row">
            <span className="mono" style={{ color: 'var(--text-muted)' }}>#{product.id}</span>
            <span className="mono" style={{ color: 'var(--accent-rose)' }}>{product.sku}</span>
            <span style={{ fontWeight: 500, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
              {product.name}
            </span>
            <span style={{ color: 'var(--text-secondary)' }}>{product.category}</span>
            <span className="mono">
              Rp {product.price.toLocaleString('id-ID')}
            </span>
            <span className="mono" style={{ fontWeight: 700 }}>
              {product.stock} pcs
            </span>
            <div>
              <span className={`badge ${product.status === 'IN_STOCK' ? 'badge-in-stock' : product.status === 'LOW_STOCK' ? 'badge-low-stock' : 'badge-out-stock'}`}>
                {product.status.replaceAll('_', ' ')}
              </span>
            </div>
            <div style={{ display: 'flex', gap: '0.5rem' }}>
              <button
                onClick={() => onUpdateStock(product.id, product.stock + 1)}
                className="btn btn-secondary"
                style={{ padding: '0.2rem 0.5rem', fontSize: '0.75rem' }}
              >
              +1
              </button>
            </div>
          </div>
        ))}
      </div>
      {filteredProducts.length > 1000 && (
        <div style={{ padding: '0.75rem 1rem', background: 'var(--bg-secondary)', textAlign: 'center', fontSize: '0.8rem', color: 'var(--text-muted)' }}>
          (Dibatasi ke 1.000 render pertama agar browser tidak crash/freeze total akibat ~40.000 DOM nodes).
        </div>
      )}
    </div>
  );
}
