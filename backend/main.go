package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() error {
	var err error
	// open using modernc.org/sqlite driver
	db, err = sql.Open("sqlite", "file:votes.db?_busy_timeout=5000")
	if err != nil {
		return err
	}
	// create table if not exists
	stmt := `
	CREATE TABLE IF NOT EXISTS votes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		character_id TEXT NOT NULL,
		name TEXT,
		image_url TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = db.Exec(stmt)
	return err
}

type VoteReq struct {
	CharacterId string `json:"characterId"`
	Name        string `json:"name"`
	ImageURL    string `json:"image_url"`
}

func allowCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var vr VoteReq
	if err := json.NewDecoder(r.Body).Decode(&vr); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if vr.CharacterId == "" {
		http.Error(w, "characterId is required", http.StatusBadRequest)
		return
	}
	res, err := db.Exec("INSERT INTO votes (character_id, name, image_url) VALUES (?, ?, ?)", vr.CharacterId, vr.Name, vr.ImageURL)
	if err != nil {
		log.Println("insert error:", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	rows, err := db.Query(`
		SELECT character_id as characterId, name, image_url as imageUrl, COUNT(*) as votes
		FROM votes
		GROUP BY character_id
		ORDER BY votes DESC
	`)
	if err != nil {
		log.Println("query error:", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	out := []map[string]interface{}{}
	for rows.Next() {
		var characterId string
		var name sql.NullString
		var imageUrl sql.NullString
		var votes int
		if err := rows.Scan(&characterId, &name, &imageUrl, &votes); err != nil {
			log.Println("scan error:", err)
			continue
		}
		m := map[string]interface{}{
			"characterId": characterId,
			"name":        nil,
			"imageUrl":    nil,
			"votes":       votes,
		}
		if name.Valid {
			m["name"] = name.String
		}
		if imageUrl.Valid {
			m["imageUrl"] = imageUrl.String
		}
		out = append(out, m)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func votesHandler(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	rows, err := db.Query("SELECT id, character_id, name, image_url, created_at FROM votes ORDER BY created_at DESC")
	if err != nil {
		log.Println("query error:", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	out := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var characterId string
		var name sql.NullString
		var imageUrl sql.NullString
		var createdAt string
		if err := rows.Scan(&id, &characterId, &name, &imageUrl, &createdAt); err != nil {
			log.Println("scan error:", err)
			continue
		}
		m := map[string]interface{}{"id": id, "characterId": characterId, "name": nil, "imageUrl": nil, "createdAt": createdAt}
		if name.Valid {
			m["name"] = name.String
		}
		if imageUrl.Valid {
			m["imageUrl"] = imageUrl.String
		}
		out = append(out, m)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func exportCSVHandler(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	rows, err := db.Query(`
		SELECT character_id as characterId, name, COUNT(*) as votes
		FROM votes
		GROUP BY character_id
		ORDER BY votes DESC
	`)
	if err != nil {
		log.Println("query error:", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=report.csv")
	writer := csv.NewWriter(w)
	writer.Write([]string{"characterId", "name", "votes"})
	for rows.Next() {
		var characterId string
		var name sql.NullString
		var votes int
		if err := rows.Scan(&characterId, &name, &votes); err != nil {
			log.Println("scan error:", err)
			continue
		}
		n := ""
		if name.Valid {
			n = name.String
		}
		writer.Write([]string{characterId, n, strconv.Itoa(votes)})
	}
	writer.Flush()
}

func main() {
	log.Println("Starting Go backend...")
	if err := initDB(); err != nil {
		log.Fatal("DB init failed:", err)
	}
	defer db.Close()

	http.HandleFunc("/api/vote", voteHandler)
	http.HandleFunc("/api/report", reportHandler)
	http.HandleFunc("/api/votes", votesHandler)
	http.HandleFunc("/api/export", exportCSVHandler)

	// simple health
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		allowCORS(w)
		w.Write([]byte(fmt.Sprintf("WaiHus Go backend - %s", time.Now().Format(time.RFC3339))))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	addr := ":" + port
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
