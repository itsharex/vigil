package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// Config holds our environment variables
type Config struct {
	Port   string
	DBPath string
}

func loadConfig() Config {
	return Config{
		Port:   getEnv("PORT", "8090"),
		DBPath: getEnv("DB_PATH", "vigil.db"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func initDB(path string) {
	var err error
	db, err = sql.Open("sqlite", path)
	if err != nil {
		log.Fatalf("❌ Failed to open database at %s: %v", path, err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS reports (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		hostname TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		data JSON
	);`

	if _, err := db.Exec(query); err != nil {
		log.Fatalf("❌ Failed to create table: %v", err)
	}
	fmt.Printf("✅ Database connected: %s\n", path)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	config := loadConfig()
	initDB(config.DBPath)
	defer db.Close()

	mux := http.NewServeMux()

	// --- API Endpoints ---

	// 1. Health Check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Vigil Server is Online"))
	})

	// 2. Collector Endpoint (Agents send data here)
	mux.HandleFunc("POST /api/report", func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		hostname := fmt.Sprintf("%v", payload["hostname"])
		jsonData, _ := json.Marshal(payload)

		_, err := db.Exec("INSERT INTO reports (hostname, data) VALUES (?, ?)", hostname, string(jsonData))
		if err != nil {
			log.Printf("❌ DB Write Error: %v", err)
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		fmt.Printf("💾 Report saved: %s | %s\n", hostname, time.Now().Format("15:04:05"))
		w.Write([]byte("Ack"))
	})

	// 3. History Endpoint (UI reads data here)
	mux.HandleFunc("GET /api/history", func(w http.ResponseWriter, r *http.Request) {
		// Get the latest report for each unique hostname
		query := `
		SELECT hostname, timestamp, data 
		FROM reports 
		WHERE id IN (
			SELECT MAX(id) 
			FROM reports 
			GROUP BY hostname
		)
		ORDER BY timestamp DESC`

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var history []map[string]interface{}
		for rows.Next() {
			var host, ts string
			var dataRaw []byte
			rows.Scan(&host, &ts, &dataRaw)
			
			// Parse the stored JSON to send back as proper object
			var dataMap map[string]interface{}
			json.Unmarshal(dataRaw, &dataMap)

			history = append(history, map[string]interface{}{
				"hostname":  host,
				"timestamp": ts,
				"details":   dataMap,
			})
		}
		jsonResponse(w, history)
	})

	// --- STATIC FILE SERVER (The New UI) ---
	// This serves everything in the "web" folder at the root URL "/"
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)

	fmt.Printf("Vigil Server listening on port %s...\n", config.Port)
	if err := http.ListenAndServe(":"+config.Port, mux); err != nil {
		log.Fatal(err)
	}
}