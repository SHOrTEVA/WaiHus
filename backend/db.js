const path = require('path');
const Database = require('better-sqlite3');

const dbPath = path.join(__dirname, 'votes.db');
const db = new Database(dbPath);

// Initialize schema
db.prepare(`
  CREATE TABLE IF NOT EXISTS votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id TEXT NOT NULL,
    name TEXT,
    image_url TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
  );
`).run();

module.exports = {
  db
};
