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
		DBPath: getEnv("DB_PATH", "vigil.db"), // Default to local file if not specified
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

// --- CORS Middleware ---
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow any origin (since we are a local dashboard)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle Preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	config := loadConfig()
	initDB(config.DBPath)
	defer db.Close()

	mux := http.NewServeMux()

	// 1. Health Check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Vigil Server is Online"))
	})

	// 2. Collector Endpoint
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

	// 3. History Endpoint
	mux.HandleFunc("GET /api/history", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, hostname, timestamp FROM reports ORDER BY id DESC LIMIT 50")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var history []map[string]interface{}
		for rows.Next() {
			var id int
			var host, ts string
			rows.Scan(&id, &host, &ts)
			history = append(history, map[string]interface{}{
				"id":        id,
				"hostname":  host,
				"timestamp": ts,
			})
		}
		jsonResponse(w, history)
	})

	fmt.Printf("Vigil Server listening on port %s...\n", config.Port)
	
	// Wrap the mux with CORS middleware
	if err := http.ListenAndServe(":"+config.Port, enableCORS(mux)); err != nil {
		log.Fatal(err)
	}
} 
// <--- This final bracket was likely missing!