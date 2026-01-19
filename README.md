# 📘 Mutabaah Online

**Mutabaah Online** adalah aplikasi web untuk **manajemen dan monitoring aktivitas akademik siswa (Mutabaah Yaumiyah)** yang dirancang khusus untuk lembaga pendidikan Islam/Tahfidz.

Aplikasi ini memudahkan pencatatan **hafalan**, **kehadiran**, dan **kegiatan harian siswa** secara digital, terintegrasi antara **Admin**, **Guru**, dan **Murid** dalam satu sistem.

---

## ✨ Fitur Utama
- Manajemen data Sekolah/Lembaga
- Pencatatan mutabaah harian (hafalan, kehadiran, catatan)
- Role-based access (Admin, Guru, Murid)
- Manajemen kelas & wali kelas
- Template penilaian fleksibel
- Monitoring progres siswa secara terpusat

---

## 📋 Prerequisites (Prasyarat)

Sebelum memulai, pastikan perangkat Anda telah terinstal software berikut:

### 1. Go (Golang)
- **Wajib versi 1.20 atau lebih baru**
- Cek versi:
  ```bash
  go version
  ```

### 2. Node.js & npm
- Gunakan versi **LTS** (v18 atau v20 direkomendasikan)
- Cek versi:
  ```bash
  node -v
  npm -v
  ```

### 3. PostgreSQL
- Database Management System
- Pastikan service PostgreSQL sudah berjalan
- Perintah `psql` dapat diakses dari terminal

### 4. Git
- Digunakan untuk clone repository

---

## 🛠️ Instalasi & Konfigurasi (Dari Nol)

Ikuti langkah-langkah berikut secara berurutan untuk menjalankan aplikasi di **Localhost**.

---

### 1️⃣ Clone Repository

Buka terminal, lalu jalankan:

```bash
git clone https://github.com/mzkhairy/mutabaahonline.git
cd mutabaahonline
```

---

### 2️⃣ Setup Database (Via Terminal)

Pastikan PostgreSQL sudah berjalan.

Masuk ke console PostgreSQL:
```bash
psql -U postgres
```

Buat database baru:
```sql
CREATE DATABASE mutabaah_db;
```

(Opsional) Cek daftar database:
```sql
\l
```

Keluar dari console:
```sql
\q
```

---

### 3️⃣ Konfigurasi Environment (`.env`)

Duplikasi file `.env.example` menjadi `.env`:

- **Windows**
  ```bat
  copy .env.example .env
  ```
- **Mac / Linux**
  ```bash
  cp .env.example .env
  ```

Buka file `.env`, lalu **edit bagian DATABASE_URL**
(ganti `password123` dengan password PostgreSQL Anda).

Contoh konfigurasi `.env`:

```env
# --- SERVER CONFIG ---
APP_ENV=dev
PORT=8080

# --- DATABASE CONFIG ---
# Format: postgres://USERNAME:PASSWORD@HOST:PORT/DBNAME?sslmode=disable
DATABASE_URL=postgres://postgres:rahasia@localhost:5432/mutabaah_db?sslmode=disable

# --- MIGRATIONS ---
RUN_MIGRATIONS=true
MIGRATIONS_DIR=migrations

# --- AUTH & SECURITY ---
JWT_SECRET=ganti_dengan_teks_acak_rahasia_anda_disini
JWT_TTL_MINUTES=43200

# --- FRONTEND CONFIG ---
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

⚠️ **PENTING:**  
Pastikan `JWT_SECRET` diganti dengan teks acak yang panjang demi keamanan.

---

### 4️⃣ Menjalankan Backend (Go)

Proyek ini menggunakan **Go Modules + Vendor**.

Jalankan perintah berikut di root proyek:

```bash
go mod tidy
go mod vendor
go run -mod=vendor cmd/api/main.go
```

✅ **Tanda berhasil:**
- Muncul log migrasi database
- Muncul pesan:
  ```
  Server starting on port 8080
  ```

⚠️ Jangan tutup terminal ini.

---

### 5️⃣ Menjalankan Frontend (Vue.js)

Buka **terminal baru**, lalu jalankan:

```bash
npm install
npm run dev
```

✅ **Tanda berhasil:**
```
Local: http://localhost:5173/
```

---

## 🚀 Cara Penggunaan Pertama Kali

### 1️⃣ Akses Aplikasi
Buka browser:
```
http://localhost:5173
```

---

### 2️⃣ Daftarkan Sekolah
- Klik **"Daftarkan Sekolah Baru"**
- Isi data sekolah dan akun **Admin**
- Simpan **Kode Lembaga** yang dihasilkan

---

### 3️⃣ Login sebagai Admin
Gunakan:
- Kode Lembaga
- Username
- Password Admin

---

### 4️⃣ Setup Data Master (WAJIB URUT)

1. **Tahun Ajaran**
   - Buat tahun ajaran baru
   - Aktifkan tahun ajaran

2. **Kelola Akun**
   - Tambah akun Guru
   - Tambah akun Murid

3. **Manajemen Kelas**
   - Buat kelas
   - Tentukan wali kelas
   - Masukkan murid ke kelas

4. **Template Penilaian**
   - Buat template penilaian  
     (contoh: Hafalan, Kehadiran, Catatan)

---

## 🧩 Teknologi yang Digunakan
- **Frontend**: Vue.js + Vite
- **Backend**: Go (Gin Framework)
- **Database**: PostgreSQL
- **Auth**: JWT
- **Environment**: Localhost / WSL / Docker-ready

---

## 📄 Lisensi
Proyek ini menggunakan lisensi **MIT** (atau sesuaikan jika berbeda).

---

## 🤝 Kontribusi
Kontribusi sangat terbuka.  
Silakan fork repository ini dan ajukan pull request.

---

## 📬 Kontak
Jika ada pertanyaan atau saran, silakan buat **Issue** di repository ini.

