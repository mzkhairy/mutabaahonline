# 📘 Mutabaah Online

**Mutabaah Online** adalah aplikasi web untuk **manajemen dan monitoring aktivitas akademik siswa (Mutabaah Yaumiyah)** yang dirancang khusus untuk lembaga pendidikan Islam/Tahfidz.

Aplikasi ini memudahkan pencatatan **hafalan**, **kehadiran**, dan **kegiatan harian siswa** secara digital, terintegrasi antara **Admin**, **Guru**, dan **Murid** dalam satu sistem.

## User Manual & Demo 

> User manual & tampilan web
- 📘 [View PDF](./MANUAL-OPERASI-SISTEM.pdf)

- ▶️ [Watch Demo (Google Drive)](https://drive.google.com/file/d/1e6atv2E5x0rt1ZRmg3UFSbAN_Umvo6Kr/view?usp=sharing)
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
- Wajib versi **1.24.3** atau lebih baru.
Download: https://go.dev/dl/
- Cek versi:
  ```
  go version
  ```

### 2. Node.js & npm
- Wajib Versi 22+
 - Download: https://nodejs.org/en/download/
- Cek versi:
  ```
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

## 🚀 Panduan Instalasi (Lokal)

Untuk menjalankan aplikasi ini di komputer Anda (Localhost), Anda bisa memilih salah satu dari dua metode di bawah ini:

1.  **Via Docker** (⭐️ **Sangat Disarankan**: Lebih mudah, bersih, dan tidak perlu install banyak *tools*).
2.  **Via Manual** (Untuk keperluan *development* mendalam).

---

### 🐳 Opsi 1: Instalasi via Docker (Recommended)

Metode ini akan menjalankan Database, Backend, dan Frontend sekaligus dalam container terisolasi.

#### 1. Persiapan (Prerequisites)
Pastikan komputer Anda sudah terinstall **Docker** dan **Docker Compose**.

* **🖥️ Windows:**
    * Download & Install **[Docker Desktop for Windows](https://www.docker.com/products/docker-desktop/)**.
    * ⚠️ **PENTING:** Pastikan fitur **WSL 2 (Windows Subsystem for Linux)** sudah diaktifkan di Windows Anda agar Docker berjalan stabil dan cepat.
* **🍎 Mac:**
    * Download & Install **[Docker Desktop for Mac](https://www.docker.com/products/docker-desktop/)**.
* **🐧 Linux:**
    * Install via terminal: `curl -fsSL https://get.docker.com -o get-docker.sh && sudo sh get-docker.sh`

#### 2. Konfigurasi Environment
Buat file baru bernama `.env` di folder root proyek (sejajar dengan `docker-compose.prod.yml`).
Salin konfigurasi berikut:

```ini
# --- DATABASE CONFIG ---
# Docker akan otomatis membuat database & user ini
DB_USER=docker_user
DB_PASSWORD=docker_password
DB_NAME=mutabaah_db_docker

# --- SECURITY ---
JWT_SECRET=rahasia_lokal_saja

# --- FRONTEND CONFIG ---
# Karena dijalankan via Docker Lokal, arahkan ke localhost
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

#### 3. Jalankan Aplikasi
Buka terminal di folder root proyek, lalu jalankan perintah:
```
docker compose -f docker-compose.prod.yml up --build
```
(Tunggu hingga proses download image dan build selesai. Proses pertama kali mungkin memakan waktu beberapa menit).

#### 4. Akses Aplikasi
Setelah log terminal berhenti bergerak dan tidak ada error:

Web Frontend: Buka browser ke http://localhost (Tidak perlu port :5173, karena Docker menjalankannya di Port 80).
API Backend: http://localhost:8080

#### 5. Menghentikan Aplikasi
Tekan Ctrl+C di terminal, atau jalankan perintah:
```
docker compose -f docker-compose.prod.yml down
```
---

⚙️ Opsi 2: Instalasi Manual (Tanpa Docker)
Gunakan cara ini jika Anda ingin menginstall Go, Node.js, dan PostgreSQL secara manual di sistem operasi Anda.

1️⃣ Clone Repository

git clone [https://github.com/mzkhairy/mutabaahonline.git](https://github.com/mzkhairy/mutabaahonline.git)
cd mutabaahonline
2️⃣ Setup Database (PostgreSQL)
Pastikan PostgreSQL sudah berjalan di komputer Anda. Masuk ke console database dan buat database baru:

psql -U postgres

-- Masuk ke psql terminal, lalu jalankan:
CREATE DATABASE mutabaah_db;
3️⃣ Konfigurasi Environment (.env)
Duplikasi file .env.example menjadi .env, lalu sesuaikan isinya dengan konfigurasi komputer Anda:

```
# --- SERVER CONFIG ---
APP_ENV=dev
PORT=8080

# --- DATABASE CONFIG ---
# Format: postgres://USER:PASSWORD@HOST:PORT/DBNAME
# Ganti 'password_anda' dengan password postgres di laptop Anda
DATABASE_URL=postgres://postgres:password_anda@localhost:5432/mutabaah_db?sslmode=disable

# --- MIGRATIONS ---
RUN_MIGRATIONS=true
MIGRATIONS_DIR=migrations

# --- AUTH ---
JWT_SECRET=ganti_dengan_teks_acak_panjang
JWT_TTL_MINUTES=43200

# --- FRONTEND ---
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

4️⃣ Menjalankan Backend (Go)
Pastikan Go (Golang) versi terbaru sudah terinstall.


#### 1. Download dependency (Go Modules)
```
go mod tidy
```
#### 2. Jalankan Server
```
go run cmd/api/main.go
```
✅ Berhasil: Muncul pesan Server starting on port 8080 dan log migrasi database sukses.

5️⃣ Menjalankan Frontend (Vue.js)
Pastikan Node.js dan npm sudah terinstall. Buka terminal baru (terminal backend jangan ditutup):


#### 1. Masuk ke folder frontend (jika struktur folder terpisah) atau root
```
cd frontend
```
#### 2. Install dependency
```
npm install
```
#### 3. Jalankan mode development
```
npm run dev
```
✅ Berhasil: Aplikasi dapat diakses di http://localhost:5173/

---

## 🚀 Cara Penggunaan Pertama Kali

### 1️⃣ Akses Aplikasi
Buka browser:
```
http://localhost:5173
```
atau

```
http://localhost
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

