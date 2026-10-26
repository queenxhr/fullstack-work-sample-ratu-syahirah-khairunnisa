# Fullstack Work Sample

Repository aplikasi asesmen Fullstack React dan Golang. Dokumen jawaban teknis untuk Soal #1-#5 ada di [ASSESSMENT.md](./ASSESSMENT.md); README ini khusus untuk setup, menjalankan, dan memvalidasi aplikasi.

## Prasyarat

- Node.js 18 atau lebih baru
- Go 1.25 (sesuai `backend/go.mod`)
- PostgreSQL opsional. Backend memakai `DB_DRIVER=auto` dan dapat fallback ke SQLite jika PostgreSQL tidak tersedia.

## Menjalankan aplikasi

### Backend

Buka terminal di root repository:

```powershell
if (!(Test-Path backend\.env)) { Copy-Item backend\.env.example backend\.env }
Set-Location backend
go run .
```

API berjalan pada `http://localhost:8080`. Untuk menggunakan PostgreSQL, sesuaikan konfigurasi koneksi di `backend\.env` terlebih dahulu.

### Frontend

Buka terminal kedua di root repository:

```powershell
Set-Location frontend
npm install
npm run dev
```

Dashboard tersedia pada `http://localhost:3000`. Vite meneruskan request `/api` ke backend di port 8080. Untuk instalasi dari root, jalankan `npm run install:all`, kemudian `npm start`.

## Build dan pengujian otomatis

Jalankan backend tests dari folder `backend`:

```powershell
go test ./...
```

Validasi build frontend dari folder `frontend`:

```powershell
npm run build
```

Frontend belum memiliki skrip unit test otomatis. Jalankan langkah pengujian manual di bawah setelah backend dan frontend aktif.

## Pengujian manual per soal

Buka dashboard di `http://localhost:3000` dan pastikan backend aktif di `http://localhost:8080` sebelum memulai. Hasil yang diharapkan dicantumkan pada tiap langkah agar mudah dibandingkan.

### Soal #1 — Performa React dan pengelolaan stok

1. Buka tab **Soal #1: React Optimization** dan tunggu sampai tabel produk termuat.
2. Pada mode **Optimized (react-window)**, scroll tabel dan amati kolom `R` pada baris: hanya baris yang terlihat yang dirender, dan penghitung render berubah pada baris yang terdampak.
3. Cari produk menggunakan nama, kategori, atau SKU. Pastikan daftar tersaring dan jumlah hasil berubah.
4. Gunakan tombol `+5` atau `-1`. Stok dan status berubah; stok tidak boleh kurang dari nol. Jika backend aktif, perubahan juga dikirim ke API dan disiarkan lewat SSE.
5. Beralih ke mode **Unoptimized**. Mode ini merender paling banyak 1.000 baris sekaligus untuk menjaga browser tetap responsif. Bandingkan perilaku scroll dan pembaruan dengan mode optimized.

**Passed IF:** pencarian bekerja, stok/status konsisten, perubahan hanya memengaruhi baris terkait dalam mode optimized, dan dashboard tetap dapat digunakan saat membandingkan mode.

### Soal #2 — Concurrency, stok, dan autentikasi JWT

1. Buka tab **Soal #2: Concurrency & JWT**.
2. Untuk metode **Atomic SQL Update**, tekan **Reset Stok = 1**, lalu **Jalankan 100 Request Bersamaan**. Hasil harus menunjukkan 100 request, tepat 1 pesanan berhasil, 99 ditolak karena stok habis, dan stok akhir 0.
3. Tekan **Reset Stok = 1** lagi dan ulangi dengan **Pessimistic Lock**. Hasil yang diharapkan sama: tidak ada overselling.
4. Di panel JWT pada tab yang sama, tekan **Login & Terbitkan JWT**, lalu **Tes Akses Profile**. Pastikan profil diterima dan token/JTI ditampilkan.
5. Tekan **Logout (Revoke)**, lalu uji kembali akses profil. Token yang sudah dicabut harus ditolak; daftar blacklist menampilkan JTI yang dicabut.
6. Alternatif pengujian halaman login: buka tab **Halaman Login**, gunakan akun demo yang tertera di layar, lalu uji endpoint terproteksi dan logout.

**Passed IF:** satu stok hanya menghasilkan satu pesanan, kedua strategi memberi hasil aman, token aktif dapat mengakses profil, dan token yang dicabut tidak dapat digunakan kembali.

### Soal #3 — Idempotency webhook dan SSE

1. Buka tab **Soal #3: Webhook & Resilience**. Pastikan indikator SSE menunjukkan **SSE Connected**; jika belum, pastikan backend berjalan lalu tunggu koneksi tersambung.
2. Pertahankan nilai **Idempotency-Key** dan **Payment ID** yang sama, kemudian tekan **Kirim Webhook** dua kali.
3. Periksa log: pengiriman pertama diproses, sedangkan pengiriman ulang ditandai sebagai duplikat. Event pembayaran juga muncul pada panel SSE untuk Payment ID tersebut.
4. Gunakan **Generate Baru** untuk membuat idempotency key baru, lalu kirim lagi. Permintaan ini harus diproses sebagai event baru.
5. Tekan **Jalankan Simulasi Fallback Cron Reconciler** dan periksa notifikasi jumlah event yang direkonsiliasi.

**Passed IF:** key yang sama tidak memproses pembayaran dua kali, key baru diproses, event diterima pada stream pembayaran, dan aksi rekonsiliasi menghasilkan respons backend.

### Soal #4 — Perbandingan TCO

1. Buka tab **Soal #4 & #5: TCO & Metrics** dan temukan kalkulator TCO pada bagian Soal #4.
2. Catat nilai awal biaya monolith, infrastruktur microservices, dan total TCO.
3. Geser jumlah **EKS Worker Nodes**, **Database Per Service**, dan biaya **DevOps/Platform Engineer** satu per satu.
4. Pastikan biaya dan rasio terhadap monolith diperbarui setelah setiap perubahan. Kembalikan slider ke nilai awal untuk membandingkan hasilnya.

**Passed IF:** perubahan jumlah node/DB memengaruhi biaya infrastruktur, perubahan biaya tenaga kerja memengaruhi total TCO, dan kalkulator tetap menampilkan rekomendasi berdasarkan perbandingan tersebut.

### Soal #5 — LTV/CAC dan alokasi sprint

1. Pada tab **Soal #4 & #5: TCO & Metrics**, buka bagian Soal #5.
2. Geser **Customer Monthly Churn Rate** dari nilai rendah ke tinggi. Amati nilai LTV dan rasio LTV/CAC.
3. Pastikan churn yang lebih tinggi menurunkan LTV dan rasionya; status rasio berubah mengikuti ambang sehat 3x.
4. Periksa pembagian sprint **70% fitur bisnis, 20% reliability/technical debt, dan 10% maintenance**. Bandingkan penjelasan tersebut dengan rekomendasi prioritisasi dan Error Budget di [ASSESSMENT.md](./ASSESSMENT.md).

**Passed IF:** LTV/CAC berubah konsisten dengan churn dan penjelasan alokasi 70/20/10 sesuai dengan dokumen asesmen.

## Database

Backend melakukan migrasi dan seeding data saat startup. Jika ingin menyiapkan PostgreSQL secara manual, gunakan [db_worksample_ratu.sql](./db_worksample_ratu.sql) dan pastikan nilai koneksi pada `backend\.env` sesuai dengan database lokal.

## Dokumen

- [Jawaban asesmen teknis Soal #1-#5](./ASSESSMENT.md)
- [Skema database](./db_worksample_ratu.sql)
