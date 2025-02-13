package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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

// Budget represents a budgeting record (including period and user_id)
type Budget struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Period      string  `json:"period"`  // e.g., daily, weekly, monthly, yearly, or custom
	UserID      int     `json:"user_id"` // foreign key to users table
}

// Charge represents a charge record
type Charge struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"` // quick description
	Amount     float64 `json:"amount"`
	ChargeType string  `json:"charge_type"` // "subscription" or "one-time"
	Periodical string  `json:"periodical"`  // daily, weekly, monthly, yearly, or custom (if subscription)
	UserID     int     `json:"user_id"`     // foreign key to users table
	CreatedAt  string  `json:"created_at"`  // timestamp string for example purposes
}

// BankAccount represents a bank account record
type BankAccount struct {
	ID        int    `json:"id"`
	Nickname  string `json:"nickname"`   // friendly name for the account
	Bank      string `json:"bank"`       // bank name
	API       string `json:"api"`        // API key/endpoint to auto-fetch charges
	UserID    int    `json:"user_id"`    // foreign key to users table
	CreatedAt string `json:"created_at"` // timestamp string for example purposes
}

func main() {
	// --- 1. Connect to PostgreSQL ---
	var err error

	// Use environment variables or hard-coded (for demo)
	// Example: POSTGRES_URI="postgres://user:password@localhost:5432/budgetdb?sslmode=disable"
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

	// budgetsHandler: GET to list budgets, POST to add a new budget
	r.HandleFunc("/api/budgets", budgetsHandler).Methods("GET", "POST")

	// budgetHandler: PUT to update and DELETE to remove a budget by id
	r.HandleFunc("/api/budgets/{id}", budgetHandler).Methods("PUT", "DELETE")

	r.HandleFunc("/api/login", loginHandler).Methods("POST")
	// Additional endpoints for charges and bank accounts can be added similarly

	// Serve static files (frontend) from the "public" folder
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))

	// --- 5. Start the server ---
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// initDB creates tables if they do not exist
func initDB(db *sql.DB) error {
	// Create users table
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        password VARCHAR(100) NOT NULL
    );
    `

	// Create budgets table with new columns: period and user_id
	createBudgetsTable := `
    CREATE TABLE IF NOT EXISTS budgets (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        amount NUMERIC(10,2) NOT NULL,
        description TEXT,
        period VARCHAR(20), -- daily, weekly, monthly, yearly, or custom
        user_id INTEGER NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
    `

	// Create charges table
	createChargesTable := `
    CREATE TABLE IF NOT EXISTS charges (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        amount NUMERIC(10,2) NOT NULL,
        charge_type VARCHAR(50) NOT NULL, -- "subscription" or "one-time"
        periodical VARCHAR(20),           -- daily, weekly, monthly, yearly, or custom
        user_id INTEGER NOT NULL,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
    `

	// Create bank_accounts table
	createBankAccountsTable := `
    CREATE TABLE IF NOT EXISTS bank_accounts (
        id SERIAL PRIMARY KEY,
        nickname VARCHAR(100) NOT NULL,
        bank VARCHAR(100) NOT NULL,
        api TEXT,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
    `

	// Execute table creation queries; IF NOT EXISTS makes them no-ops if the table exists
	_, err := db.Exec(createUsersTable)
	if err != nil {
		return fmt.Errorf("creating users table: %v", err)
	}

	_, err = db.Exec(createBudgetsTable)
	if err != nil {
		return fmt.Errorf("creating budgets table: %v", err)
	}

	_, err = db.Exec(createChargesTable)
	if err != nil {
		return fmt.Errorf("creating charges table: %v", err)
	}

	_, err = db.Exec(createBankAccountsTable)
	if err != nil {
		return fmt.Errorf("creating bank_accounts table: %v", err)
	}

	return nil
}

// createDefaultAdmin checks if an admin user exists, and only creates one if not
func createDefaultAdmin(db *sql.DB) error {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE username=$1`, "admin").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		// In production, the password should be hashed instead of plain text
		_, err := db.Exec(`INSERT INTO users (username, password) VALUES ($1, $2)`, "admin", "admin")
		if err != nil {
			return fmt.Errorf("insert default admin: %v", err)
		}
		log.Println("Default admin created: admin / admin")
	} else {
		log.Println("Admin user already exists, skipping creation.")
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

	// Compare password (in production, compare hashed passwords)
	if dbUser.Password != creds.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Return success (in a real app, you might set a session or JWT)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "login successful"})
}

// budgetsHandler handles GET and POST requests for budgets.
func budgetsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Retrieve all budgets
		rows, err := db.Query("SELECT id, name, amount, description, period, user_id FROM budgets")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying budgets: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var budgets []Budget
		for rows.Next() {
			var b Budget
			err := rows.Scan(&b.ID, &b.Name, &b.Amount, &b.Description, &b.Period, &b.UserID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning budget: %v", err), http.StatusInternalServerError)
				return
			}
			budgets = append(budgets, b)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(budgets)

	case "POST":
		// Add a new budget
		var b Budget
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Insert the new budget record and return the new id
		err = db.QueryRow(
			"INSERT INTO budgets (name, amount, description, period, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			b.Name, b.Amount, b.Description, b.Period, b.UserID,
		).Scan(&b.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error inserting budget: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(b)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// budgetHandler handles PUT (update) and DELETE requests for a specific budget.
func budgetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the id parameter to an integer.
	budgetID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "PUT":
		// Update an existing budget
		var b Budget
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		result, err := db.Exec(
			"UPDATE budgets SET name=$1, amount=$2, description=$3, period=$4, user_id=$5 WHERE id=$6",
			b.Name, b.Amount, b.Description, b.Period, b.UserID, budgetID,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating budget: %v", err), http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching update result: %v", err), http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			http.Error(w, "Budget not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Budget updated successfully"})

	case "DELETE":
		// Delete a budget
		result, err := db.Exec("DELETE FROM budgets WHERE id=$1", budgetID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error deleting budget: %v", err), http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching delete result: %v", err), http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			http.Error(w, "Budget not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Budget deleted successfully"})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
