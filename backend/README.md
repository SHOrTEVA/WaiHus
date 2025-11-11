# WaiHus Backend (Python)

Simple FastAPI + SQLite backend for anonymous voting.

Endpoints:
- POST /api/vote { characterId, name, image_url }
- GET  /api/report -> aggregated counts
- GET  /api/votes  -> raw votes
- GET  /api/export -> CSV download of aggregated counts

## Run locally (Python)

```powershell
cd backend
pip install -r requirements.txt
python -m uvicorn main:app --reload --port 3000
```

Or:

```bash
cd backend
pip install -r requirements.txt
uvicorn main:app --reload --port 3000
```

The server listens on port 3000 by default.

## Database

SQLite is auto-initialized when the app starts. The database file `votes.db` is created in the `backend/` folder.

## Development

- FastAPI automatically generates API docs at `http://localhost:3000/docs` (Swagger UI)
- Use `--reload` flag to auto-restart on code changes
