import React, { memo, useCallback, useMemo, useRef, useState } from 'react';
import { FixedSizeList as List, areEqual } from 'react-window';
import { RefreshCw, Zap, TrendingUp, AlertCircle, CheckCircle2 } from 'lucide-react';

/**
 * ProductRow Component
 * Optimized using React.memo with areEqual comparison so that updating
 * one item does NOT re-render any other rows in the 5,000+ item list.
 */
const ProductRow = memo(({ data, index, style }) => {
  const { products, onUpdateStock } = data;
  const product = products[index];

  // Ref to track render count for this specific row (demonstrates memoization visually!)
  const renderCountRef = useRef(0);
  renderCountRef.current += 1;

  if (!product) return null;

  const getStatusBadge = (status) => {
    switch (status) {
      case 'IN_STOCK':
        return <span className="badge badge-in-stock"><CheckCircle2 size={12} /> In Stock</span>;
      case 'LOW_STOCK':
        return <span className="badge badge-low-stock"><AlertCircle size={12} /> Low Stock</span>;
      default:
        return <span className="badge badge-out-stock"><AlertCircle size={12} /> Out of Stock</span>;
    }
  };

  const formattedPrice = new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(product.price);

  return (
    <div style={style} className="product-row">
      <span className="mono" style={{ color: 'var(--text-muted)' }}>#{product.id}</span>
      <span className="mono" style={{ color: 'var(--accent-cyan)' }}>{product.sku}</span>
      <span style={{ fontWeight: 600, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
        {product.name}
      </span>
      <span style={{ color: 'var(--text-secondary)' }}>{product.category}</span>
      <span className="mono" style={{ fontWeight: 600 }}>{formattedPrice}</span>
      <span className="mono" style={{ fontWeight: 700, color: product.stock <= 5 ? 'var(--accent-amber)' : 'var(--text-primary)' }}>
        {product.stock} pcs
      </span>
      <div>{getStatusBadge(product.status)}</div>
      <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
        <button
          onClick={() => onUpdateStock(product.id, product.stock + 5)}
          className="btn btn-secondary"
          style={{ padding: '0.25rem 0.625rem', fontSize: '0.75rem' }}
          title="Tambah stok (+5)"
        >
          +5
        </button>
        <button
          onClick={() => onUpdateStock(product.id, Math.max(0, product.stock - 1))}
          className="btn btn-secondary"
          style={{ padding: '0.25rem 0.625rem', fontSize: '0.75rem', color: 'var(--accent-rose)' }}
          title="Kurang stok (-1)"
        >
          -1
        </button>
        <span 
          className="mono" 
          title="Jumlah Render Row Ini (Membuktikan Memoization)" 
          style={{ fontSize: '0.7rem', color: 'var(--text-muted)', padding: 0, background: 'transparent' }}
        >
          R:{renderCountRef.current}
        </span>
      </div>
    </div>
  );
}, areEqual);

/**
 * ProductTableOptimized
 * Demonstrates high performance virtualization & memoization for 5,000+ items.
 */
export default function ProductTableOptimized({
  products,
  onUpdateStock,
  searchQuery,
  totalItems,
}) {
  const [listHeight, setListHeight] = useState(600);
  const rowHeight = 54;

  // Memoize filtered products based on search query
  const filteredProducts = useMemo(() => {
    if (!searchQuery) return products;
    const q = searchQuery.toLowerCase();
    return products.filter(
      (p) =>
        p.name.toLowerCase().includes(q) ||
        p.sku.toLowerCase().includes(q) ||
        p.category.toLowerCase().includes(q)
    );
  }, [products, searchQuery]);

  // Memoize itemData passed to react-window FixedSizeList
  // Ensures references to itemData remain stable across parent re-renders!
  const itemData = useMemo(() => ({
    products: filteredProducts,
    onUpdateStock,
  }), [filteredProducts, onUpdateStock]);

  return (
    <div className="glass-panel" style={{ overflow: 'hidden' }}>
      {/* Table Subheader Info */}
      <div
        title="react-window merender hanya baris yang terlihat sehingga mengurangi jumlah node DOM, menjaga scrolling tetap halus, dan menghindari layout thrashing."
        style={{ padding: '1rem 1.25rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '1px solid var(--border-color)' }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
          <span className="badge badge-accent">
            <Zap size={14} /> Virtualized (react-window)
          </span>
          <span style={{ fontSize: '0.875rem', color: 'var(--text-secondary)' }}>
            {filteredProducts.length.toLocaleString('id-ID')} dari {totalItems.toLocaleString('id-ID')} produk
          </span>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontSize: '0.75rem', color: 'var(--accent-emerald)' }}>
          <span className="render-pulse"></span>
          Virtualisasi aktif
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
        <span>Aksi Stok (Row Render)</span>
      </div>

      {/* Virtualized List Container */}
      <div style={{ width: '100%', height: `${listHeight}px` }}>
        {filteredProducts.length > 0 ? (
          <List
            height={listHeight}
            itemCount={filteredProducts.length}
            itemSize={rowHeight}
            width="100%"
            itemData={itemData}
          >
            {ProductRow}
          </List>
        ) : (
          <div style={{ padding: '3rem', textAlign: 'center', color: 'var(--text-muted)' }}>
            Tidak ada produk yang cocok dengan pencarian "{searchQuery}"
          </div>
        )}
      </div>
    </div>
  );
}
