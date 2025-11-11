const express = require('express');
const cors = require('cors');
const { db } = require('./db');
const path = require('path');

const app = express();
app.use(cors());
app.use(express.json());

// POST /api/vote
app.post('/api/vote', (req, res) => {
  const { characterId, name, image_url } = req.body;
  if (!characterId) return res.status(400).json({ error: 'characterId is required' });
  const stmt = db.prepare('INSERT INTO votes (character_id, name, image_url) VALUES (?, ?, ?)');
  const info = stmt.run(String(characterId), name || null, image_url || null);
  res.status(201).json({ id: info.lastInsertRowid });
});

// GET /api/report - aggregated counts per character
app.get('/api/report', (req, res) => {
  const rows = db.prepare(`
    SELECT character_id as characterId, name, image_url as imageUrl, COUNT(*) as votes
    FROM votes
    GROUP BY character_id
    ORDER BY votes DESC
  `).all();
  res.json(rows);
});

// GET /api/votes - raw votes
app.get('/api/votes', (req, res) => {
  const rows = db.prepare('SELECT * FROM votes ORDER BY created_at DESC').all();
  res.json(rows);
});

// GET /api/export - CSV of aggregated report
app.get('/api/export', (req, res) => {
  const rows = db.prepare(`
    SELECT character_id as characterId, name, COUNT(*) as votes
    FROM votes
    GROUP BY character_id
    ORDER BY votes DESC
  `).all();

  let csv = 'characterId,name,votes\n';
  for (const r of rows) {
    // simple CSV escaping
    const name = r.name ? '"' + String(r.name).replace(/"/g, '""') + '"' : '';
    csv += `${r.characterId},${name},${r.votes}\n`;
  }

  res.setHeader('Content-Type', 'text/csv');
  res.setHeader('Content-Disposition', 'attachment; filename="report.csv"');
  res.send(csv);
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => console.log(`Backend running on port ${PORT}`));
