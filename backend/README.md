# WaiHus Backend (Rust)# WaiHus Backend (Python)



Fast, safe, and production-ready backend built with Actix-web and SQLite.Simple FastAPI + SQLite backend for anonymous voting.



Endpoints:Endpoints:

- POST /api/vote { character_id, name, image_url }- POST /api/vote { characterId, name, image_url }

- GET  /api/report -> aggregated vote counts (JSON)- GET  /api/report -> aggregated counts

- GET  /api/votes  -> all raw votes- GET  /api/votes  -> raw votes

- GET  /api/export -> CSV download of aggregated counts- GET  /api/export -> CSV download of aggregated counts



## Prerequisites## Run locally (Python)



- Rust 1.70+ ([install](https://rustup.rs/))```powershell

cd backend

## Run locallypip install -r requirements.txt

python -m uvicorn main:app --reload --port 3000

```powershell```

cd backend

cargo run --releaseOr:

```

```bash

Or (with logging):cd backend

```powershellpip install -r requirements.txt

$env:RUST_LOG="info"; cargo run --releaseuvicorn main:app --reload --port 3000

``````



The server listens on `http://0.0.0.0:3000` by default.The server listens on port 3000 by default.



## Build## Database



To create an optimized binary:SQLite is auto-initialized when the app starts. The database file `votes.db` is created in the `backend/` folder.



```powershell## Development

cargo build --release

```- FastAPI automatically generates API docs at `http://localhost:3000/docs` (Swagger UI)

- Use `--reload` flag to auto-restart on code changes

Binary will be at `target/release/waihus-backend.exe` (Windows) or `target/release/waihus-backend` (Unix).

## Database

SQLite database (`votes.db`) is auto-created on first run in the `backend/` folder.

## Development

Watch for file changes and rebuild:

```powershell
cargo watch -x run
```

(requires `cargo-watch`: `cargo install cargo-watch`)

## Performance

Rust + Actix-web + SQLite is very fast:
- ~1ms per request (depending on hardware)
- Minimal memory footprint
- Safe concurrency with zero data races (enforced by Rust compiler)
