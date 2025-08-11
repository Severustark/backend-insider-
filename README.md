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
