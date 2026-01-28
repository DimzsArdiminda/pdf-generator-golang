# PDF Generator API

RESTful API untuk generate dan manage file PDF menggunakan Golang, PostgreSQL, dan library GoFPDF.

## 📋 Daftar Isi

- [Fitur](#-fitur)
- [Tech Stack](#-tech-stack)
- [Arsitektur Projek](#-arsitektur-projek)
- [Prasyarat](#-prasyarat)
- [Instalasi](#-instalasi)
- [Konfigurasi](#-konfigurasi)
- [Menjalankan Aplikasi](#-menjalankan-aplikasi)
- [API Endpoints](#-api-endpoints)
- [Postman Collection](#-postman-collection)
- [Database Schema](#-database-schema)
- [Contoh Request](#-contoh-request)

## ✨ Fitur

- ✅ Generate PDF dari JSON data dengan header institusi
- ✅ Upload file PDF
- ✅ List semua file PDF
- ✅ Soft delete PDF
- ✅ Environment configuration dengan .env
- ✅ PostgreSQL database integration
- ✅ RESTful API design

## 🛠 Tech Stack

- **Language:** Go 1.21.7
- **Web Framework:** Gorilla Mux
- **Database:** PostgreSQL
- **PDF Library:** GoFPDF
- **Environment:** godotenv

### Dependencies

```go
github.com/gorilla/mux v1.8.1
github.com/lib/pq v1.10.9
github.com/jung-kurt/gofpdf v1.16.2
github.com/joho/godotenv v1.5.1
```

## 📁 Arsitektur Projek

```
pdf-generate/
├── cmd/
│   └── main.go              # Entry point aplikasi
├── config/
│   └── database.go          # Database configuration & connection
├── handlers/
│   └── pdf_handler.go       # HTTP request handlers
├── models/
│   └── pdf_file.go          # Data models
├── routes/
│   └── routes.go            # Route definitions
├── uploads/
│   └── pdf/                 # Directory untuk menyimpan PDF
├── docs/                    # Dokumentasi tambahan
├── .env                     # Environment variables (jangan commit!)
├── .gitignore              # Git ignore rules
├── dump.sql                # Database schema
├── go.mod                  # Go modules
├── go.sum                  # Dependency checksums
├── README.md               # Dokumentasi ini
└── req.rest                # REST Client test file
```

## 📦 Prasyarat

Sebelum menjalankan aplikasi, pastikan Anda telah menginstall:

- [Go](https://golang.org/dl/) (versi 1.21 atau lebih baru)
- [PostgreSQL](https://www.postgresql.org/download/) (versi 12 atau lebih baru)
- [Git](https://git-scm.com/)

## 🚀 Instalasi

### 1. Clone Repository

```bash
git clone <repository-url>
cd pdf-generate
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Setup Database

Buat database PostgreSQL baru:

```bash
psql -U postgres
CREATE DATABASE pdf-generator;
\c pdf-generator
```

Jalankan schema dari file `dump.sql`:

```bash
psql -U postgres -d pdf-generator -f dump.sql
```

### 4. Setup Environment Variables

Copy file `.env` dan sesuaikan dengan konfigurasi Anda:

```bash
# Sudah ada file .env, sesuaikan nilai-nilainya
```

## ⚙️ Konfigurasi

Edit file `.env` dan sesuaikan dengan environment Anda:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=pdf-generator
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Application
APP_ENV=development
APP_DEBUG=true

# JWT Secret (if using authentication)
JWT_SECRET=your-secret-key-here-change-in-production

# Upload Configuration
MAX_UPLOAD_SIZE=10485760
UPLOAD_PATH=uploads/pdf/

# CORS (if needed)
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

**⚠️ PENTING:** Jangan lupa ganti `DB_PASSWORD` dan `JWT_SECRET` dengan nilai yang aman!

## 🏃 Menjalankan Aplikasi

### Development Mode

```bash
go run cmd/main.go
```

### Build & Run

```bash
go build -o pdf-api.exe cmd/main.go
./pdf-api.exe
```

Server akan berjalan di `http://localhost:8080`

## 📡 API Endpoints

| Method | Endpoint            | Description                 |
| ------ | ------------------- | --------------------------- |
| POST   | `/api/pdf/generate` | Generate PDF dari JSON data |
| POST   | `/api/pdf/upload`   | Upload file PDF             |
| GET    | `/api/pdf/list`     | List semua file PDF         |
| DELETE | `/api/pdf/{id}`     | Soft delete PDF by ID       |

### Detail Endpoints

#### 1. Generate PDF

**Endpoint:** `POST /api/pdf/generate`

**Request Body:**

```json
{
  "title": "Laporan Kunjungan Pasien",
  "institution_name": "RS Sehat Sentosa",
  "address": "Jl. Kesehatan No. 123, Jakarta",
  "phone": "(021) 123-4567",
  "logo_url": "http://localhost:8080/logo.png",
  "content": "Berikut adalah laporan kunjungan pasien bulan Januari 2025. Total kunjungan mencapai 1,234 pasien dengan berbagai keluhan kesehatan. Mayoritas pasien datang untuk pemeriksaan rutin dan konsultasi kesehatan. Tim medis telah memberikan pelayanan terbaik dengan standar kesehatan yang tinggi. Laporan ini dibuat untuk dokumentasi dan evaluasi kinerja rumah sakit."
}
```

**Response:**

```json
{
  "success": true,
  "message": "PDF generated successfully"
}
```

#### 2. Upload PDF

**Endpoint:** `POST /api/pdf/upload`

**Request:** Multipart form-data

- `file`: File PDF (max 10MB)

**Response:**

```json
{
  "success": true,
  "message": "PDF uploaded successfully"
}
```

#### 3. List PDF

**Endpoint:** `GET /api/pdf/list`

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "filename": "report_1738047600.pdf",
      "original_name": "document.pdf",
      "size": 25600,
      "status": "CREATED",
      "created_at": "2026-01-28T10:00:00Z"
    }
  ]
}
```

#### 4. Delete PDF

**Endpoint:** `DELETE /api/pdf/{id}`

**Response:**

```json
{
  "success": true,
  "message": "PDF deleted successfully"
}
```

## 📮 Postman Collection

Untuk testing yang lebih mudah, import Postman collection berikut:

**[PDF Generator API - Postman Collection](https://www.postman.com/mobile-engineer-team/workspace/public-workspace/collection/42465258-b1b86b0f-8301-4c64-935d-a4ab76eb4812?action=share&source=copy-link&creator=42465258ocumentation&active-environment=42465258-3cf4ec51-8aa3-42e4-90c9-35abf62fa3d6)**

Collection ini sudah termasuk:

- Semua endpoint dengan contoh request
- Pre-configured environment variables
- Test scripts untuk validasi response

## 🗄️ Database Schema

```sql
CREATE TABLE pdf_files (
    id BIGSERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    original_name VARCHAR(255),
    logo_url VARCHAR(50),
    filepath VARCHAR(500) NOT NULL,
    size BIGINT,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

### Status Values:

- `UPLOADED` - File di-upload oleh user
- `CREATED` - File di-generate oleh sistem
- `DELETED` - File sudah dihapus (soft delete)

## 📝 Contoh Request

### Menggunakan cURL

#### Generate PDF

```bash
curl -X POST http://localhost:8080/api/pdf/generate \
  -H "Content-Type: application/json" \
  -d '{
  "title": "Laporan Kunjungan Pasien",
  "institution_name": "RS Sehat Sentosa",
  "address": "Jl. Kesehatan No. 123, Jakarta",
  "phone": "(021) 123-4567",
  "logo_url": "http://localhost:8080/logo.png",
  "content": "Berikut adalah laporan kunjungan pasien bulan Januari 2025. Total kunjungan mencapai 1,234 pasien dengan berbagai keluhan kesehatan. Mayoritas pasien datang untuk pemeriksaan rutin dan konsultasi kesehatan. Tim medis telah memberikan pelayanan terbaik dengan standar kesehatan yang tinggi. Laporan ini dibuat untuk dokumentasi dan evaluasi kinerja rumah sakit."
}'
```

#### List PDF

```bash
curl http://localhost:8080/api/pdf/list
```

#### Upload PDF

```bash
curl -X POST http://localhost:8080/api/pdf/upload \
  -F "file=@/path/to/document.pdf"
```

#### Delete PDF

```bash
curl -X DELETE http://localhost:8080/api/pdf/1
```

### Menggunakan REST Client (VSCode Extension)

Install extension "REST Client" di VSCode, lalu gunakan file `req.rest` yang sudah tersedia di project.

## 🔧 Development

### Running Tests

```bash
go test ./...
```

### Format Code

```bash
go fmt ./...
```

### Check Dependencies

```bash
go mod tidy
```

## 📄 License

This project is licensed under the MIT License.

## 👤 Author

Created with ❤️ 

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!

## 📞 Support

Jika ada pertanyaan atau masalah, silakan buat issue di repository ini.

---

**Happy Coding! 🚀**
