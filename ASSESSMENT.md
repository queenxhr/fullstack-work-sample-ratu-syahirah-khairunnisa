# WORK SAMPLE ASSESSMENT: SENIOR FULLSTACK REACT & BACKEND SYSTEMS

> **Assessment Title:** Senior Fullstack JavaScript / React Engineer Role | Technical & Domain Assessment  
> **Candidate Name:** Ratu Syahirah Khairunnisa  
> **Date:** 26 Oktober 2026  
> **Stack:** React 18 (Vite, react-window, memoization) & Golang 1.25 (Clean Layered Architecture, GORM, PostgreSQL)

---

## 📑 Daftar Isi
1. [Arsitektur Sistem (System Architecture)](#-arsitektur-sistem-system-architecture)
2. [Soal #1: Optimasi Performa & State Management React (Live Coding / Implementation)](#-soal-1-optimasi-performa--state-management-react-live-coding--implementation)
3. [Soal #2: Arsitektur Backend, Autentikasi JWT, & High-Concurrency Scaling](#-soal-2-arsitektur-backend-autentikasi-jwt--high-concurrency-scaling)
4. [Soal #3: Strategi Integrasi API, Webhook Reliability, & Resiliency](#-soal-3-strategi-integrasi-api-webhook-reliability--resiliency)
5. [Soal #4: Evaluasi TCO (Total Cost of Ownership) & Trade-off Monolith vs Microservices](#-soal-4-evaluasi-tco-total-cost-of-ownership--trade-off-monolith-vs-microservices)
6. [Soal #5: Prioritisasi Feature, Metrik Software (LTV/CAC), & Handling Tech Debt](#-soal-5-prioritisasi-feature-metrik-software-ltvcac--handling-tech-debt)
7. [Struktur Direktori Repositori](#-struktur-direktori-repositori)

---

## 🏛 Arsitektur Sistem (System Architecture)

<img width="4425" height="7938" alt="Golang Product Service Architecture-2026-10-26-080808" src="https://github.com/user-attachments/assets/271b5cbe-a689-406b-a9b9-8f407647b334" />

---

## ⚡ Soal #1: Optimasi Performa & State Management React (Live Coding / Implementation)
**Kategori:** Frontend React JS  

### 1.1 Diagnosa Teknis: Penyebab Utama Unnecessary Re-Render pada Dataset Besar (5.000+ Items)
Ketika merender 5.000 data produk dalam satu tabel e-commerce secara real-time, terjadi degradasi performa ekstrem (lagging, dropped frames hingga < 10 FPS, main-thread blocking). Diagnosa teknis dari akar permasalahannya adalah:

1. **Reconciliation & Fiber Tree Diffing Overhead:**
   React bekerja menggunakan Virtual DOM. Ketika satu item mengalami update (misal stok berubah dari 10 menjadi 9), React secara default mengeksekusi rekonsiliasi ke seluruh subtree. Membandingkan (diffing) 5.000 fiber nodes satu per satu memakan waktu CPU > 100ms, melampaui frame budget browser (16.6ms untuk 60 FPS), menyebabkan browser freeze (jank).

2. **Parent Re-render Cascading (Reference Invalidation):**
   Jika 5.000 data disimpan dalam satu state monolitik di parent (`const [products, setProducts] = useState(list)`), memperbarui satu item dengan `setProducts([...products])` menghasilkan array baru. Tanpa memoization referensial yang tepat, seluruh 5.000 child components menerima props baru dan dipaksa me-render ulang.

3. **Inline Anonymous Functions & Prop Instability:**
   Mengoper callback inline seperti `onClick={() => onUpdate(product.id)}` atau objek gaya inline `{ style: ... }` membuat referensi memori baru di setiap render cycle. Akibatnya, `React.memo` bawaan gagal melakukan shallow comparison (`prevProps !== nextProps` selalu bernilai `false`).

4. **DOM Bloat & Layout Thrashing:**
   Merender 5.000 baris tabel secara native menghasilkan lebih dari **40.000 DOM nodes** sekaligus (5.000 baris × ~8 elemen per baris). Browser engine menghabiskan waktu masif untuk proses *Recalculate Styles*, *Layout/Reflow*, dan *Compositing layers* setiap kali ada mutasi DOM kecil.

---

### 1.2 Implementasi Komponen React (Virtualization + Memoization)
Solusi lengkap telah diimplementasikan secara fungsional pada [`ProductTableOptimized.jsx`](./frontend/src/components/ProductTableOptimized.jsx).

Berikut adalah arsitektur kode kuncinya:

```jsx
import React, { memo, useCallback, useMemo, useRef } from 'react';
import { FixedSizeList as List, areEqual } from 'react-window';

/**
 * 1. Komponen ProductRow di-memoize dengan areEqual dari react-window.
 * Row ini HANYA akan re-render jika data row spesifik ini yang berubah!
 */
const ProductRow = memo(({ data, index, style }) => {
  const { products, onUpdateStock } = data;
  const product = products[index];

  // Visual counter untuk membuktikan row lain tidak ikut re-render
  const renderCountRef = useRef(0);
  renderCountRef.current += 1;

  if (!product) return null;

  return (
    <div style={style} className="product-row">
      <span className="sku">{product.sku}</span>
      <span className="name">{product.name}</span>
      <span className="price">Rp {product.price.toLocaleString('id-ID')}</span>
      <span className="stock">{product.stock} pcs</span>
      <span className="status-badge">{product.status}</span>
      <div className="actions">
        <button onClick={() => onUpdateStock(product.id, product.stock + 1)}>
          +1
        </button>
        <button onClick={() => onUpdateStock(product.id, Math.max(0, product.stock - 1))}>
          -1
        </button>
        <span className="render-indicator">Render: {renderCountRef.current}x</span>
      </div>
    </div>
  );
}, areEqual);

/**
 * 2. Virtualized Table Component
 * Hanya merender ~12-15 elemen DOM yang berada di viewport aktif,
 * mengabaikan sisa 4.985 item dari DOM tree!
 */
export default function ProductTableOptimized({ products, onUpdateStock, searchQuery }) {
  // Memoize data yang difilter agar tidak dihitung ulang setiap render
  const filteredProducts = useMemo(() => {
    if (!searchQuery) return products;
    const q = searchQuery.toLowerCase();
    return products.filter((p) => p.name.toLowerCase().includes(q) || p.sku.toLowerCase().includes(q));
  }, [products, searchQuery]);

  // Memoize itemData agar referensi ke react-window stabil
  const itemData = useMemo(() => ({
    products: filteredProducts,
    onUpdateStock,
  }), [filteredProducts, onUpdateStock]);

  return (
    <div className="virtual-table-container">
      <List
        height={600}
        itemCount={filteredProducts.length}
        itemSize={54}
        width="100%"
        itemData={itemData}
      >
        {ProductRow}
      </List>
    </div>
  );
}
```

---

### 1.3 Strategi State Management agar Update Item Individual Tidak Memicu Re-Render Global
Untuk memastikan update data real-time pada 1 item tidak memicu re-render pada 4.999 baris lainnya:

1. **Normalized State Structure (Dictionary by ID):**
   Ubah struktur state dari array flat menjadi skema normalized entity (seperti Normalizr / Redux Toolkit):
   ```ts
   type ProductState = {
     byId: { [id: number]: Product };
     allIds: number[];
   };
   ```
2. **Selective / Atomic Subscription (Zustand / Recoil / Signals):**
   Dengan Zustand atau Context Selector, `ProductRow` hanya men-subscribe perubahan slice ID miliknya:
   ```ts
   // Di dalam ProductRow:
   const product = useProductStore(useCallback((state) => state.byId[id], [id]));
   // Baris 42 TIDAK PERNAH re-render ketika baris 100 berubah!
   ```
3. **Web Worker & Micro-Task Batching untuk Real-Time Stream:**
   Saat lonjakan ratusan event perubahan stok per detik masuk melalui WebSocket/SSE, deserialization dilakukan di Web Worker, kemudian di-batch menggunakan `requestAnimationFrame` sebelum di-dispatch ke React agar tidak membombardir event loop.

---

## 🛡 Soal #2: Arsitektur Backend, Autentikasi JWT, & High-Concurrency Scaling
**Kategori:** Backend Architecture & Security  

### 2.1 Pilihan Tech Stack Backend & Justifikasi
**Pilihan:** **Golang 1.21+ dengan Gin Web Framework & GORM (Clean Layered Architecture)**

**Alasan & Justifikasi:**
1. **Lightweight Concurrency Model (Goroutines vs OS Threads):**
   Satu Goroutine di Go hanya membutuhkan memory stack awal sebesar **~2 KB** (dibandingkan OS thread Java atau PHP-FPM worker yang membutuhkan 1 MB – 2 MB per koneksi). Go runtime scheduler (M:N scheduler) mampu menangani **puluhan ribu concurrent HTTP requests per detik (RPS)** saat promo Flash Sale dengan CPU dan memory overhead yang sangat rendah.
2. **Compiled, Predictable Low Latency:**
   Go dikompilasi langsung ke native machine code tanpa perantara Virtual Machine (seperti JVM) atau Single-Threaded Event Loop (seperti Node.js V8 yang rentan terhadap CPU-blocking tasks).
3. **Clean Architecture Separation of Concerns:**
   Mengikuti pola domain-driven layer `entity -> repository -> usecase -> handler` seperti yang diimplementasikan pada repositori ini. Komponen storage dan transport terisolasi secara murni, memungkinkan pengujian unit yang terisolasi dan fleksibilitas integrasi database.

---

### 2.2 Mekanisme Token Revocation (Logout / Invalidation) pada Stateless JWT
**Tantangan:** JWT secara arsitektur bersifat *self-contained* dan *stateless*. Setelah token ditandatangani, token tersebut valid secara matematis sampai batas waktu kadaluarsa (`exp`).

**Solusi Arsitektur yang Diimplementasikan:**
1. **Redis Revocation Blacklist dengan Automatic TTL:**
   - Setiap kali JWT diterbitkan, sematkan klaim standar `jti` (JWT ID / UUID unik).
   - Ketika user menekan tombol **Logout**, backend membaca `jti` dan menghitung sisa waktu aktif token:
     $$\text{TTL} = \text{Claims.ExpiresAt} - \text{CurrentTime}$$
   - Simpan `jti` ke dalam Redis:
     ```redis
     SETEX "jwt_blacklist:<jti>" <TTL> "revoked"
     ```
   - Di HTTP Authentication Middleware, lakukan pengecekan:
     ```go
     if isBlacklisted(claims.JTI) {
         return 401 Unauthorized ("Token telah dicabut / logged out")
     }
     ```
   - **Efisiensi:** Karena menggunakan TTL otomatis yang sama dengan sisa masa berlaku token, memori Redis tidak akan pernah bocor (no memory leak) karena token yang sudah expired secara natural akan dihapus otomatis oleh Redis.
2. **Short-Lived Access Token (15 Menit) + Refresh Token Rotation (RTR):**
   Access token hanya aktif selama 10-15 menit. Token yang terlanjur diterbitkan akan hilang kekuatannya dalam waktu singkat tanpa membebani penyimpanan blacklist dalam durasi lama.

---

### 2.3 Strategi Penanganan 'Race Condition' (100 User Membeli 1 Sisa Barang)
Jika 100 user checkout secara simultan terhadap 1 sisa barang, terjadi bahaya *Read-Modify-Write Race Condition* (semua 100 goroutine membaca `stock = 1`, semua menganggap stok cukup, lalu semua melakukan decrement sehingga stok menjadi -99).

**Solusi yang Diterapkan di Repositori:**

#### 1. Atomic Database Update (Best Practice untuk SQL / PostgreSQL)
```sql
UPDATE products 
SET stock = stock - 1, version = version + 1, updated_at = CURRENT_TIMESTAMP 
WHERE id = ? AND stock >= 1;
```
Di Golang (GORM):
```go
res := tx.Model(&ProductModel{}).
    Where("id = ? AND stock >= ?", productID, quantity).
    Update("stock", gorm.Expr("stock - ?", quantity))

if res.RowsAffected == 0 {
    // 99 request akan jatuh ke blok ini secara atomic!
    return ErrConflict("Stok habis atau pesanan melebihi sisa stok")
}
```
**Mengapa ini unggul?** Engine PostgreSQL menggunakan *Row-Level Lock* internal saat mengeksekusi `UPDATE`. Hanya 1 dari 100 transaksi yang mendapatkan `RowsAffected == 1`. 99 request lainnya langsung mendapatkan `RowsAffected == 0` dan ditolak (HTTP 409 Conflict) secara instan tanpa deadlocks.

#### 2. Pessimistic Locking (`SELECT ... FOR UPDATE`)
```go
tx := db.Begin()
var prod ProductModel
// Kunci baris produk secara eksklusif dalam transaksi
tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&prod, "id = ?", productID)

if prod.Stock < quantity {
    tx.Rollback()
    return ErrConflict("Stok tidak mencukupi")
}

tx.Model(&prod).Update("stock", prod.Stock - quantity)
tx.Commit()
```

#### 3. Redis In-Memory Lua Scripting (Untuk Skala 10.000+ RPS)
Untuk traffic spike ekstrem, stok didecrement di Redis sebelum menyentuh PostgreSQL menggunakan Atomic Lua Script:
```lua
local stock = tonumber(redis.call('get', KEYS[1]))
if stock and stock >= tonumber(ARGV[1]) then
    redis.call('decrby', KEYS[1], ARGV[1])
    return 1 -- Sukses
else
    return 0 -- Habis
end
```
User yang berhasil didecrement di Redis kemudian dimasukkan ke Message Queue (Kafka/RabbitMQ) untuk pencatatan order asynchronous ke PostgreSQL.

---

### 2.4 Diagram Alur Arsitektur Backend (Mermaid.js)

<img width="8192" height="3381" alt="Golang Product Service Architecture-2026-10-26-081155" src="https://github.com/user-attachments/assets/b1eb97d7-a620-4183-be68-7f7a68a5f554" />

---

## 📡 Soal #3: Strategi Integrasi API, Webhook Reliability, & Resiliency
**Kategori:** System Integration & Fullstack Communication  

### 3.1 Rancang Mekanisme 'Idempotency' pada Endpoint Webhook Backend
**Definisi:** Operasi idempotent menjamin bahwa menerima notifikasi webhook yang sama $N$ kali menghasilkan efek samping yang persis sama dengan menerima 1 kali ($f(f(x)) = f(x)$), mencegah pengkreditan saldo atau pengiriman barang ganda.

**Rancangan Implementasi:**
1. **Idempotency Key Extraction:**
   Ambil header `Idempotency-Key` atau ID transaksi unik dari Payment Gateway (`payment_id` / `event_id`).
2. **Database Unique Constraint & Atomic Locking:**
   Tabel `webhook_events` memiliki unique index pada `idempotency_key` dan `event_id`:
   ```sql
   CREATE TABLE webhook_events (
       id SERIAL PRIMARY KEY,
       event_id VARCHAR(100) UNIQUE NOT NULL,
       idempotency_key VARCHAR(100) UNIQUE NOT NULL,
       payment_id VARCHAR(100) NOT NULL,
       status VARCHAR(50) NOT NULL,
       processed_count INT DEFAULT 1,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   ```
3. **Alur Eksekusi Idempotent:**
   - Cek apakah `idempotency_key` sudah ada di database.
   - **Jika Ditemukan:** Increment `processed_count = processed_count + 1`. Return HTTP 200 OK dengan payload sebelumnya. **Skip seluruh side-effect** (jangan update order lagi, jangan kirim email lagi).
   - **Jika Baru:** Buka transaksi database, simpan event, ubah status order dari `PENDING` menjadi `PAID`, dan simpan perubahan secara atomik.

---

### 3.2 Strategi Fallback / Retry Mechanism untuk Pembayaran yang Menggantung
Jika koneksi internet putus atau terjadi timeout pada webhook:
1. **Exponential Backoff Retry via Dead Letter Queue (DLQ):**
   - Webhook yang gagal diproses diarahkan ke message broker (RabbitMQ / Kafka / Redis Asynq).
   - Sistem melakukan retry dengan jeda exponential: 10s, 30s, 2m, 15m, 1h.
   - Jika tetap gagal setelah 5x percobaan, pesan dipindahkan ke **Dead Letter Queue (DLQ)** dan alert otomatis dikirimkan ke tim DevOps (Slack / PagerDuty) untuk investigasi.
2. **Automated Active Reconciliation Cron Job:**
   - Worker terjadwal (misal: setiap 3 menit) melakukan query mencari order yang masih berstatus `PENDING` dan telah berusia $> 5$ menit.
   - Worker menembak **Outbound Check Status API** Payment Gateway:
     `GET /v1/payments/{payment_id}/status`
   - Jika status di Payment Gateway sudah `PAID`, worker merekonsiliasi status lokal secara otomatis. Jika sudah `EXPIRED`, batalkan order dan kembalikan stok produk.

---

### 3.3 Komunikasi Status Pembayaran ke Frontend React (SSE vs Polling vs WebSocket)

| Kriteria | Short Polling | WebSocket | Server-Sent Events (SSE) ✅ (Pilihan Terbaik) |
|---|---|---|---|
| **Arah Komunikasi** | Unidirectional (Client -> Server) | Full-Duplex (Bidirectional) | Unidirectional (Server -> Client) |
| **Beban Server & Koneksi** | **Sangat Boros:** Membuka/menutup ribuan HTTP socket kosong tiap 2 detik. | **Tinggi:** Perlu stateful TCP socket connection management di backend. | **Sangat Ringan:** Berjalan di atas HTTP standard dengan keep-alive chunked transfer. |
| **Auto-Reconnection** | Harus dikoding manual di React. | Harus dikoding manual dengan heartbeat / ping-pong. | **Otomatis Native Browser:** `EventSource` otomatis reconnect jika koneksi putus! |
| **Kesesuaian Kasus Pembayaran** | ❌ Boros resource server | ⚠️ Overkill (user hanya menunggu sinyal konfirmasi dari server) | **✅ Sempurna:** User menunggu 1 arah sinyal `event: payment_success`. |

**Implementasi di Frontend React:**
```javascript
useEffect(() => {
  const eventSource = new EventSource(`/api/v1/webhook/stream/${paymentId}`);

  eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data);
    if (data.status === 'SUCCESS') {
      setPaymentStatus('PAID');
      eventSource.close(); // Tutup stream segera setelah pembayaran selesai
    }
  };

  eventSource.onerror = () => {
    // Jika SSE gagal atau terhalang proxy lawas, fallback otomatis ke smart exponential polling!
  };

  return () => eventSource.close();
}, [paymentId]);
```

---

## 💰 Soal #4: Evaluasi TCO (Total Cost of Ownership) & Trade-off Monolith vs Microservices
**Kategori:** Business & Technical Decision Making  

### 4.1 Analisis Total Cost of Ownership (TCO): Monolith ($150/bln) vs Cloud-Native Microservices
**Konteks:** Aplikasi e-commerce saat ini berjalan dengan biaya server $150/bulan. Tim mengusulkan migrasi langsung ke Microservices (Kubernetes/EKS, AWS Lambda, Managed DBs).

#### Perbandingan Biaya Riil Finansial (Real Cost Breakdown):

| Komponen Biaya | Monolith Saat Ini (VPS Single/HA) | Cloud-Native Microservices (AWS EKS) |
|---|---|---|
| **Compute / Control Plane** | 1 VPS High-Spec: **$100/bln** | EKS Cluster Fee ($73) + 3 Worker Nodes Multi-AZ ($240): **$313/bln** |
| **Database Storage** | 1 Managed PostgreSQL: **$50/bln** | 3 Managed DBs per Domain (User, Catalog, Order): **$285/bln** |
| **Networking & Routing** | Gratis / Standard Gateway | ALB Ingress + NAT Gateways + Cross-AZ Traffic: **$180/bln** |
| **Observability & Tracing** | Standard CloudWatch / Loki: **$10/bln** | Datadog / OpenSearch Distributed Traces: **$250/bln** |
| **Subtotal Biaya Infrastruktur** | **$160 / bulan** | **$1.028 / bulan (Membengkak 6.4x lipat!)** |
| **Biaya SDM (Operational / DevOps)** | 0 Engineer Dedicated (Dikelola dev fullstack) | Minimal 1 Dedicated Platform/DevOps Engineer: **$2.500 - $4.000/bln** |
| **TOTAL TCO BULANAN** | **~$160 / bulan** | **~$3.528 - $5.000+ / bulan (Lonjakan >20x Lipat!)** |

**Kesimpulan Finansial:**
Secara bisnis, migrasi penuh ke Microservices untuk startup skala awal-menengah adalah **"Premature Optimization" yang berbahaya**. Langkah ini membakar runway modal perusahaan tanpa menghasilkan pertumbuhan revenue yang sebanding.

---

### 4.2 Hidden Costs (Biaya Tersembunyi) dari Arsitektur Microservices
1. **Biaya Egress Jaringan & NAT Gateways:** Komunikasi RPC/HTTP antar puluhan kontainer di subnet privat melintasi availability zones (AZ) menghasilkan tagihan data transfer tersembunyi yang mahal.
2. **Cognitive Load & Debugging Latency:** Melacak bug transaksi membutuhkan distributed tracing (Jaeger/OpenTelemetry) melintasi 5-10 services, memperlambat Mean Time to Resolution (MTTR).
3. **Kompleksitas Distributed Transaction (Saga Pattern):** Menghilangkan kemudahan atomic database transaction (`db.Begin()`). Kegagalan salah satu service membutuhkan orchestrator transaksi kompensasi yang sangat rawan bug inkonsistensi data.
4. **Contract Testing & Release Orchestration:** Butuh pipeline CI/CD yang jauh lebih kompleks, memakan waktu build lama, dan memperlambat delivery fitur bisnis harian.

---

### 4.3 Rekomendasi Arsitektur Kompromi: Modular Monolith
**Arsitektur Rekomendasi:** **Modular Monolith (Clean Architecture) dalam Golang** (seperti yang diimplementasikan pada repositori ini).

**Karakteristik & Keunggulan:**
- **Single Deployable Unit:** Satu binary aplikasi yang dapat di-deploy ke satu VPS murah ($100-$150/bln) via Docker Compose.
- **Strict Domain Boundaries:** Kode dipisahkan secara tegas ke dalam modul domain mandiri (`catalog`, `order`, `payment`, `user`) melalui Go interface.
- **In-Memory Zero-Latency Communication:** Komunikasi antar modul menggunakan pemanggilan fungsi Go native in-memory (0ms latency, zero network failure).
- **Future-Proof Scalability:** Jika suatu saat domain Flash Sale mengalami traffic jutaan user dan memerlukan scaling mandiri, modul tersebut sudah terisolasi dan **siap di-extract menjadi microservice independen** hanya dalam hitungan hari tanpa perlu menulis ulang seluruh sistem.

---

## 📈 Soal #5: Prioritisasi Feature, Metrik Software (LTV/CAC), & Handling Tech Debt
**Kategori:** Product Domain & Software Business Metrics  

### 5.1 Mengkomunikasikan Dampak Technical Debt kepada PM & Stakeholder Non-Teknis
Sebagai Senior/Lead Fullstack Developer, hindari terminologi teknis abstrak seperti *"kode kita spaghetti"* atau *"ada race condition"*. Komunikasikan dampak teknis dalam **Bahasa Finansial dan Risiko Bisnis (Financial & Revenue Risk)**:

1. **Kalkulasi Kerugian Finansial Akibat Downtime (Cost of Outage):**
   > *"Jika server down selama 2 jam saat promo flash sale akibat technical debt ini, estimasi transaksi yang gagal adalah Rp 120.000.000, melampaui seluruh target pendapatan dari 3 fitur baru yang sedang kita kejar."*
2. **Pemborosan Biaya Marketing Akuisisi Pelanggan (Wasted CAC):**
   > *"Kita mengeluarkan biaya iklan CAC sebesar Rp 200.000 per user baru. Jika aplikasi lagging atau checkout freeze, 45% dari user baru tersebut langsung uninstall aplikasi. Kita membuang anggaran marketing secara percuma."*
3. **Korelasi Data Riil:**
   Tunjukkan grafik APM (Application Performance Monitoring): lonjakan error rate 5xx, latensi checkout p99 yang menyentuh 6 detik, dan utilisasi database connection pool yang mencapai 98%.

---

### 5.2 Hubungan Reliabilitas Sistem dengan Churn Rate & Customer Lifetime Value (LTV)
Rumus Customer Lifetime Value (LTV):
$$\text{LTV} = \frac{\text{Average Revenue Per User (ARPU)} \times \text{Gross Margin}}{\text{Churn Rate}}$$

**Analisis Hubungan Kausal:**
1. **Churn Rate berada di Posisi Penyebut (Denominator):**
   Peningkatan kecil pada Churn Rate akan menghancurkan LTV secara signifikan.
   - Jika ARPU = Rp 250.000 dan Churn Rate normal = 5%, maka:
     $$\text{LTV} = \frac{\text{Rp 250.000}}{0.05} = \text{Rp 5.000.000}$$
   - Jika sistem sering down atau lagging parah sehingga Churn Rate naik menjadi 12%, maka:
     $$\text{LTV} = \frac{\text{Rp 250.000}}{0.12} = \text{Rp 2.083.333} \quad (\text{LTV anjlok sebesar } 58\%!)$$
2. **Rasio Kesehatan Bisnis ($\text{LTV} / \text{CAC}$):**
   Standar bisnis teknologi yang sehat menuntut rasio $\text{LTV} / \text{CAC} \ge 3.0$. Jika reliabilitas sistem buruk, LTV akan turun hingga mendekati CAC ($\text{LTV} / \text{CAC} \le 1.0$), menyebabkan perusahaan merugi pada setiap pelanggan baru yang diakuisisi!

---

### 5.3 Framework Prioritisasi Konkret: Membagi Kapasitas Sprint

#### 1. The 70 / 20 / 10 Sprint Capacity Allocation Rule
Dalam setiap sprint berdurasi 2 minggu, alokasikan kapasitas total story points engineering secara disiplin:
- **70% Kapasitas:** **Fitur Bisnis Baru & Revenue Growth** (Fokus Product Manager untuk mengejar CAC).
- **20% Kapasitas:** **Refactoring Technical Debt & Core Reliability** (Fokus Engineering untuk menstabilkan database connection, index query, dan optimasi virtual table).
- **10% Kapasitas:** **Buffer Maintenance, Security Patches, & Unplanned Bugfixes**.

#### 2. Adjusted RICE Scoring untuk Technical Debt
Gunakan formula RICE yang disesuaikan untuk membandingkan prioritas fitur bisnis vs tech debt secara objektif:
$$\text{Score} = \frac{\text{Reach} \times \text{Impact} \times \text{Confidence}}{\text{Effort}}$$

Untuk Technical Debt kritikal yang berisiko membuat server down:
- **Reach:** 100% dari seluruh pengguna aktif aplikasi.
- **Impact:** 3x (Massive - Mencegah catastrophic server outage).
- **Confidence:** 90%.
Skor ini akan secara objektif menempati prioritas teratas dalam backlog sprint tanpa menimbulkan perdebatan emosional antara tim PM dan Engineering.

#### 3. Kebijakan Error Budget / Service Level Objective (SLO) Circuit Breaker
- Tetapkan SLO Uptime sebesar **99.9%** (Error Budget toleransi downtime: ~43 menit per bulan).
- **Aturan Tegas:** Jika dalam suatu periode Error Budget terlampaui (karena insiden downtime atau lonjakan error rate di atas 1%), maka **seluruh pengerjaan fitur baru dihentikan sementara (Feature Freeze)**, dan kapasitas sprint berikutnya dialihkan **100% untuk stabilitas sistem dan pembersihan Technical Debt** sampai metrik reliabilitas pulih kembali.

---

## 📂 Struktur Direktori Repositori

```
fullstack-work-sample-ratu-syahirah-khairunnisa/
├── README.md                               # Setup, cara menjalankan, dan pengujian aplikasi
├── ASSESSMENT.md                           # Jawaban lengkap Soal #1 - #5
├── package.json                            # Root convenience runner scripts
├── .gitignore                              # Git ignore rules
│
├── backend/                                # Service Backend Golang (Clean Layered Architecture)
│   ├── config/
│   │   └── config.go                       # Environment configuration loader
│   ├── internal/
│   │   ├── entity/                         # Domain Structs (Product, Order, Token, Webhook)
│   │   │   ├── product.go
│   │   │   ├── order.go
│   │   │   ├── token.go
│   │   │   └── webhook.go
│   │   ├── repository/                     # Repository Interfaces (Domain contracts)
│   │   │   ├── product.go
│   │   │   ├── order.go
│   │   │   ├── token.go
│   │   │   └── webhook.go
│   │   ├── usecase/                        # Business Logic Layer
│   │   │   ├── product_usecase.go          # 5k items dataset, SSE streaming broadcaster
│   │   │   ├── flashsale_usecase.go        # 100-concurrency benchmark & race condition prevention
│   │   │   ├── auth_usecase.go             # Stateless JWT & Redis-style JTI Blacklist revocation
│   │   │   └── webhook_usecase.go          # Idempotency checks, DLQ fallback, and payment SSE
│   │   └── errors/
│   │       └── errors.go                   # Standard HTTP dynamic error handler
│   ├── infrastructure/
│   │   ├── database/
│   │   │   └── db.go                       # PostgreSQL connection with SQLite fallback & 5k product seeder
│   │   └── repository/                     # GORM Postgres Implementations
│   │       ├── product/
│   │       │   ├── model.go
│   │       │   └── postgres.go
│   │       ├── order/
│   │       │   ├── model.go
│   │       │   └── postgres.go
│   │       ├── webhook/
│   │       │   ├── model.go
│   │       │   └── postgres.go
│   │       └── token/
│   │           └── cache.go                # In-memory Redis-simulating TTL revocation blacklist
│   ├── handler/
│   │   └── v1/
│   │       ├── router.go                   # Gin Router, CORS middleware, and route registrations
│   │       ├── product/http/               # Product HTTP handlers & DTOs
│   │       ├── flashsale/http/             # Flash Sale HTTP handlers & Concurrency Benchmark
│   │       ├── auth/http/                  # Auth login, profile, and JWT middleware
│   │       └── webhook/http/               # Idempotent Webhook & SSE payment status stream
│   ├── utils/
│   │   └── http.go                         # Standard JSON API response envelope
│   ├── .env.example
│   ├── .env
│   ├── go.mod
│   ├── go.sum
│   └── main.go                             # Application entrypoint
│
└── frontend/                               # Mini App React JS (Vite + Modern Vanilla CSS)
    ├── package.json
    ├── vite.config.js
    ├── index.html
    └── src/
        ├── main.jsx                        # React DOM mounting
        ├── index.css                       # Modern dark glassmorphic design system
        ├── App.jsx                         # Main Dashboard layout, metrics bar, and tab navigation
        └── components/
            ├── ProductTableOptimized.jsx   # Virtualized (react-window) table for 5,000+ products
            ├── ProductTableUnoptimized.jsx # Comparison table demonstrating layout thrashing & lag
            ├── FlashSaleConcurrencyDemo.jsx# Interactive test panel for 100 concurrent requests against 1 stock
            ├── WebhookPaymentResilience.jsx# Interactive panel for Idempotency verification & SSE stream
            └── ArchitectureDiagrams.jsx    # Interactive TCO Calculator & LTV/CAC Sprint Planner
```

---

*Disusun dengan standar teknis tingkat tinggi untuk evaluasi posisi **Senior Fullstack React & Systems Engineer** oleh **Ratu Syahirah Khairunnisa**.*
