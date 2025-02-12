package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// DB connection global (for simplicity)
var db *sql.DB

// User represents a user in the system
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Budget example struct
type Budget struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

func main() {
	// --- 1. Connect to PostgreSQL ---
	var err error

	// Use environment variables or hard-coded (for demo)
	// e.g. POSTGRES_URI="postgres://user:password@localhost:5432/budgetdb?sslmode=disable"
	dbURI := os.Getenv("POSTGRES_URI")
	if dbURI == "" {
		dbURI = "postgres://admin:admin@localhost:5432/budgetdb?sslmode=disable"
	}

	db, err = sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	// Check if the connection is alive
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}

	log.Println("Connected to PostgreSQL!")

	// --- 2. Initialize Database Tables ---
	err = initDB(db)
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v\n", err)
	}

	// --- 3. Create default admin user if not exists ---
	err = createDefaultAdmin(db)
	if err != nil {
		log.Fatalf("Failed to create default admin: %v\n", err)
	}

	// --- 4. Set up Routes ---
	r := mux.NewRouter()

	r.HandleFunc("/api/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/budgets", budgetsHandler).Methods("GET")

	// Serve static files (frontend) from the "public" folder
	// (We'll place our Node app build files there, or just a static index.html)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))

	// --- 5. Start the server ---
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// initDB creates tables if they do not exist
func initDB(db *sql.DB) error {
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        password VARCHAR(100) NOT NULL
    );
    `

	createBudgetsTable := `
    CREATE TABLE IF NOT EXISTS budgets (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        amount NUMERIC(10,2) NOT NULL,
        description TEXT
    );
    `

	// You can add more tables as needed, e.g. transactions, categories, etc.

	_, err := db.Exec(createUsersTable)
	if err != nil {
		return fmt.Errorf("creating users table: %v", err)
	}

	_, err = db.Exec(createBudgetsTable)
	if err != nil {
		return fmt.Errorf("creating budgets table: %v", err)
	}

	return nil
}

// createDefaultAdmin checks if an admin user exists, if not, creates one
func createDefaultAdmin(db *sql.DB) error {
	// Check if an admin user already exists
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE username=$1`, "admin").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		// In production, you would hash the password
		_, err := db.Exec(`INSERT INTO users (username, password) VALUES ($1, $2)`, "admin", "admin")
		if err != nil {
			return fmt.Errorf("insert default admin: %v", err)
		}
		log.Println("Default admin created: admin / admin")
	} else {
		log.Println("Admin user already exists.")
	}
	return nil
}

// loginHandler - example login endpoint
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds User
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check user in DB
	var dbUser User
	row := db.QueryRow(`SELECT id, username, password FROM users WHERE username=$1`, creds.Username)
	err = row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Compare password (in production: compare hashed)
	if dbUser.Password != creds.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Return success, or you could set a session/cookie
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "login successful"})
}

// budgetsHandler - example budgets endpoint
func budgetsHandler(w http.ResponseWriter, r *http.Request) {
	// For demonstration, let's just return a static array of budgets.
	// In reality, you’d query the `budgets` table.
	mockBudgets := []Budget{
		{ID: 1, Name: "Groceries", Amount: 300.00, Description: "Monthly groceries"},
		{ID: 2, Name: "Rent", Amount: 1200.00, Description: "Monthly rent payment"},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mockBudgets)
}
