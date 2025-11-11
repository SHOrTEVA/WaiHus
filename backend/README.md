# WaiHus Backend (minimal)

Simple Express + SQLite backend for anonymous voting.

Endpoints:
- POST /api/vote { characterId, name, image_url }
- GET  /api/report -> aggregated counts
- GET  /api/votes  -> raw votes
- GET  /api/export -> CSV download of aggregated counts

Run locally:

```powershell
cd backend
npm install
npm start
```

The server listens on port 3000 by default.
