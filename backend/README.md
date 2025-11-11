# WaiHus Go Backend (minimal)

This is a minimal Go backend that stores anonymous votes in a local SQLite database and exposes the following endpoints:

- POST /api/vote  { characterId, name, image_url }
- GET  /api/report -> aggregated counts per character (JSON)
- GET  /api/votes  -> raw votes (JSON)
- GET  /api/export -> aggregated CSV (download)

Notes:
- Uses `modernc.org/sqlite` driver (pure-Go SQLite) to avoid cgo on Windows.
- Default listen port: 3001 (set via PORT env var to change)

Run locally:

```powershell
cd backend-go
go mod tidy
go run main.go
```

Then point the frontend to the Go backend API (default http://localhost:3001). For example, update `Chart.vue` `axios.post`/`axios.get` URLs if you prefer the Go backend port.

If you want the server to listen on port 3000, set `PORT=3000` before starting.

Windows note: `modernc.org/sqlite` is pure Go and should work without installing extra C toolchains.
