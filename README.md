# GRADIFY

Proyek e-learning **Nusantara Kreatif** — portofolio studi kasus pengembangan platform pendaftaran siswa.

## Struktur

```
GRADIFY/
├── server/              # Go backend + frontend
│   ├── handlers.go      # HTTP handler + validasi input
│   ├── handlers_test.go # Integration test
│   ├── main.go          # Entry point server
│   └── static/          # HTML frontend (index.html)
├── formtest/            # Unit test validasi form
│   ├── main.go
│   └── main_test.go
├── analisis/            # Analisis data pendaftar
│   ├── main.go          # Baca CSV + ranking kota
│   ├── visualisasi.html # Choropleth map (Leaflet.js)
│   └── data/
│       └── pendaftaran.csv
└── form-preview.html    # Preview form statis
```

## Cara Menjalankan

```bash
# Jalankan server
cd server
go run .

# Unit test validasi form
cd formtest
go test -v

# Integration test API
cd server
go test -v

# Analisis data
cd analisis
go run .
```

Server berjalan di `http://localhost:8080`.

## Teknologi

- **Backend**: Go (net/http), bcrypt, regexp
- **Frontend**: HTML, CSS, Vanilla JS
- **Visualisasi**: Leaflet.js, GeoJSON
- **Testing**: Go testing + httptest
