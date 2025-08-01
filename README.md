# 🧠 Backend Insider

JWT tabanlı kimlik doğrulama, kullanıcı yönetimi, bakiye ve işlem takibi gibi temel özelliklere sahip, Go dilinde geliştirilmiş ölçeklenebilir bir RESTful backend uygulamasıdır. PostgreSQL ve Docker desteğiyle birlikte gelir.

## 🚀 Özellikler

- 🔐 JWT ile kimlik doğrulama (access & refresh token desteği)
- 👤 Kullanıcı kayıt, giriş ve güncelleme
- 🏦 Bakiye görüntüleme ve para yatırma
- 💸 Transaction (kredi/debit/transfer) işlemleri
- 🧱 GORM destekli PostgreSQL veritabanı
- 🐳 Docker ve docker-compose desteği
- ⚙️ Katmanlı mimari: `handlers`, `repositories`, `services`
- 🧪 Token doğrulama ve role-based erişim (admin/user)

- ## 📁 Proje Klasör Yapısı

```bash
backend-insider/
├── cmd/server/              # Uygulama giriş noktası (main.go)
├── internal/
│   ├── config/              # Ortam değişkeni yönetimi
│   ├── db/                  # GORM veritabanı bağlantısı
│   ├── models/              # Veritabanı tabloları
│   ├── repositories/        # DB işlemleri
│   └── server/
│       ├── handlers/        # Endpoint işlemleri
│       └── routes/          # Route tanımlamaları
├── pkg/logger/              # Zerolog ile loglama
├── utils/                   # Yardımcı fonksiyonlar
├── schema.sql               # Veritabanı şeması
├── Dockerfile               # Docker yapılandırması
├── docker-compose.yml       # Docker Compose
├── go.mod / go.sum          # Go modül yönetimi
└── README.md                # Proje dökümantasyonu
```

### 1. `.env` Dosyası Oluşturma

Projeyi çalıştırmadan önce proje dizininde `.env` adında bir dosya oluşturun ve aşağıdaki ortam değişkenlerini içine yazın:

```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=verdiğiniz_şifre
DB_NAME=verdiğiniz_veritabanı_ismi
JWT_SECRET=defaultsecret


### 2. Docker ile Başlatma

```bash
docker-compose up --build
```

## 🧪 API Örnekleri

### 🔐 Auth

```http
POST /api/v1/auth/register
{
  "email": "user@example.com",
  "password": "123456"
}

POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "123456"
}
```

### 👤 Kullanıcı

```http
GET /api/v1/users/me
Authorization: Bearer <access_token>

PUT /api/v1/users/me
{
  "username": "yeni-isim",
  "email": "yeni@example.com"
}
```
### 🏦 Bakiye & İşlem

```http
GET /api/v1/balances/{id}

POST /api/v1/balances/deposit
{
  "user_id": 1,
  "amount": 100.0
}

GET  /api/v1/balances/{id}
POST /api/v1/balances/deposit
POST /api/v1/transactions/debit
POST /api/v1/transactions
GET  /api/v1/transactions
GET  /api/v1/transactions/{id}
GET  /api/v1/transactions/user/{id}
GET  /api/v1/transactions/history

---

## 📦 Bağımlılıklar

| Paket                      | Açıklama                   |
| -------------------------- | -------------------------- |
| `github.com/go-chi/chi/v5` | HTTP router                |
| `gorm.io/gorm`             | ORM aracı                  |
| `github.com/lib/pq`        | PostgreSQL sürücüsü        |
| `github.com/joho/godotenv` | .env dosyası okuyucu       |
| `github.com/rs/zerolog`    | JSON tabanlı loglama aracı |


## 👤 Geliştirici

**Damla Arpa**  
📧 damlarpa@gmail.com
🔗 [[github.com/Severustark](https://github.com/Severustark)](https://github.com/Severustark)
🔗https://www.linkedin.com/in/damla-arpa/

---
