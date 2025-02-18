package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	// Import statements
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Global variables
var (
	jwtSecret []byte
	db        *sql.DB
)

// --------------------------
//        Data Models
// --------------------------

// User: plaintext username, bcrypt-hashed password
type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Permissions string `json:"permissions"`
}

// Budget: belongs to a user
type Budget struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
	Period   string  `json:"period"`
	UserID   int     `json:"user_id"`
}

// Charge: belongs to a user
type Charge struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
	Category   string  `json:"category"`
	Periodical string  `json:"periodical"`
	UserID     int     `json:"user_id"`
	CreatedAt  string  `json:"created_at"`
}

// Share: user_id shares something with user_share_id
type Share struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	UserShareID int    `json:"user_share_id"`
	Access      string `json:"access"`
}

// --------------------------
//      CORS Middleware
// --------------------------

// corsMiddleware adds the necessary headers to allow cross-origin requests.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow any origin (or restrict to a specific one)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow specific methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Allow specific headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// --------------------------
//  Initialization
// --------------------------

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}
	jwtSecret = []byte(secret)
}

func main() {
	var err error

	// Connect to PostgreSQL
	dbURI := os.Getenv("POSTGRES_URI")
	if dbURI == "" {
		// Default (change as needed)
		dbURI = "postgres://admin:admin@localhost:5432/budgetdb?sslmode=disable"
	}

	db, err = sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	log.Println("Connected to PostgreSQL!")

	// Create tables if needed
	if err := initDB(db); err != nil {
		log.Fatalf("Failed to initialize DB: %v\n", err)
	}

	// Create default admin user if not exists
	if err := createDefaultAdmin(db); err != nil {
		log.Fatalf("Failed to create default admin: %v\n", err)
	}

	// Create default admin user if not exists
	if err := createDefaultUsers(db); err != nil {
		log.Fatalf("Failed to create default Users: %v\n", err)
	}

	// --------------------------
	//         ROUTES
	// --------------------------
	r := mux.NewRouter()

	// Users (admin-only)
	r.HandleFunc("/api/users", createUserHandler).Methods("POST")
	r.HandleFunc("/api/users/{id}", updateUserHandler).Methods("PUT")
	r.HandleFunc("/api/users/{id}", deleteUserHandler).Methods("DELETE")
	r.HandleFunc("/api/users", getUsersHandler).Methods("GET")

	// Login
	r.HandleFunc("/api/login", loginHandler).Methods("POST")

	// Budgets
	r.HandleFunc("/api/budgets", getBudgetsHandler).Methods("GET")
	r.HandleFunc("/api/budgets", createBudgetHandler).Methods("POST")
	r.HandleFunc("/api/budgets/{id}", updateBudgetHandler).Methods("PUT")
	r.HandleFunc("/api/budgets/{id}", deleteBudgetHandler).Methods("DELETE")

	// Charges
	r.HandleFunc("/api/charges", getChargesHandler).Methods("GET")
	r.HandleFunc("/api/charges", createChargeHandler).Methods("POST")
	r.HandleFunc("/api/charges/{id}", updateChargeHandler).Methods("PUT")
	r.HandleFunc("/api/charges/{id}", deleteChargeHandler).Methods("DELETE")

	// Shares
	r.HandleFunc("/api/shares", getSharesHandler).Methods("GET")
	r.HandleFunc("/api/shares", createShareHandler).Methods("POST")
	r.HandleFunc("/api/shares/{id}", deleteShareHandler).Methods("DELETE")

	// Serve static files (optional front-end)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))

	// Wrap the router with the CORS middleware
	handler := corsMiddleware(r)

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// --------------------------
//   Database Initialization
// --------------------------

func initDB(db *sql.DB) error {
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(255) UNIQUE NOT NULL, 
        password VARCHAR(255) NOT NULL,         
        permissions VARCHAR(50) NOT NULL DEFAULT 'user'
    );
    `
	createBudgetsTable := `
    CREATE TABLE IF NOT EXISTS budgets (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        amount NUMERIC(10,2) NOT NULL,
        category TEXT,
        period VARCHAR(20),
        user_id INTEGER NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
    `
	createChargesTable := `
    CREATE TABLE IF NOT EXISTS charges (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        amount NUMERIC(10,2) NOT NULL,
        category TEXT NOT NULL, 
        periodical VARCHAR(20),
        user_id INTEGER NOT NULL,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
    `
	createSharesTable := `
    CREATE TABLE IF NOT EXISTS shares (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        user_share_id INTEGER NOT NULL,
        access TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY (user_share_id) REFERENCES users(id) ON DELETE CASCADE
    );
    `

	if _, err := db.Exec(createUsersTable); err != nil {
		return fmt.Errorf("creating users table: %v", err)
	}
	if _, err := db.Exec(createBudgetsTable); err != nil {
		return fmt.Errorf("creating budgets table: %v", err)
	}
	if _, err := db.Exec(createChargesTable); err != nil {
		return fmt.Errorf("creating charges table: %v", err)
	}
	if _, err := db.Exec(createSharesTable); err != nil {
		return fmt.Errorf("creating shares table: %v", err)
	}

	return nil
}

func createDefaultAdmin(db *sql.DB) error {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE permissions='admin'`).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		hashedPass, err := hashPassword("admin")
		if err != nil {
			return fmt.Errorf("hash admin password: %v", err)
		}
		_, err = db.Exec(`
            INSERT INTO users (username, password, permissions)
            VALUES ($1, $2, $3)
        `, "admin", hashedPass, "admin")
		if err != nil {
			return fmt.Errorf("insert default admin: %v", err)
		}
		log.Println("Default admin created: username=admin / password=admin (permissions=admin)")
	} else {
		log.Println("Admin user already exists, skipping creation.")
	}
	return nil
}

func createDefaultUsers(db *sql.DB) error {
	// Check if user "alice" exists
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE username=$1`, "alice").Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking for alice: %v", err)
	}
	if count == 0 {
		hashedPass, err := hashPassword("alice") // Replace with desired password
		if err != nil {
			return fmt.Errorf("error hashing alice's password: %v", err)
		}
		_, err = db.Exec(`
            INSERT INTO users (username, password, permissions)
            VALUES ($1, $2, $3)
        `, "alice", hashedPass, "user")
		if err != nil {
			return fmt.Errorf("error inserting alice: %v", err)
		}
		log.Println("Default user created: username=alice / password=alice (permissions=user)")
	} else {
		log.Println("User 'alice' already exists, skipping creation.")
	}

	// Check if user "bob" exists
	err = db.QueryRow(`SELECT COUNT(*) FROM users WHERE username=$1`, "bob").Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking for bob: %v", err)
	}
	if count == 0 {
		hashedPass, err := hashPassword("bob") // Replace with desired password
		if err != nil {
			return fmt.Errorf("error hashing bob's password: %v", err)
		}
		_, err = db.Exec(`
            INSERT INTO users (username, password, permissions)
            VALUES ($1, $2, $3)
        `, "bob", hashedPass, "user")
		if err != nil {
			return fmt.Errorf("error inserting bob: %v", err)
		}
		log.Println("Default user created: username=bob / password=bob (permissions=user)")
	} else {
		log.Println("User 'bob' already exists, skipping creation.")
	}

	return nil
}

// --------------------------
//    Password + JWT
// --------------------------

func hashPassword(plainPass string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt error: %v", err)
	}
	return string(hashedBytes), nil
}

func checkPasswordHash(plainPass, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plainPass))
	return err == nil
}

func generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// getUserIDFromToken parses the JWT from the "Authorization: Bearer <token>" header
// and returns the user_id claim. Returns an error if invalid or missing.
func getUserIDFromToken(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("no auth header")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return 0, fmt.Errorf("invalid auth header format")
	}
	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return 0, fmt.Errorf("token parse error: %v", err)
	}
	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid claims")
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("no user_id in token")
	}
	return int(userIDFloat), nil
}

// isAdmin checks if the requesting user is an admin.
func isAdmin(r *http.Request) bool {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		return false
	}
	var permissions string
	err = db.QueryRow(`SELECT permissions FROM users WHERE id=$1`, userID).Scan(&permissions)
	if err != nil {
		return false
	}
	return permissions == "admin"
}

// GET /api/users => return all users with ID, username, and permissions (admin-only)
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Error(w, "Forbidden - Admins only", http.StatusForbidden)
		return
	}

	rows, err := db.Query(`SELECT id, username, permissions FROM users`)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching users: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Define a struct for the view (excluding the password)
	type UserView struct {
		ID          int    `json:"id"`
		Username    string `json:"username"`
		Permissions string `json:"permissions"`
	}

	var users []UserView
	for rows.Next() {
		var u UserView
		if err := rows.Scan(&u.ID, &u.Username, &u.Permissions); err != nil {
			http.Error(w, fmt.Sprintf("Error scanning user: %v", err), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// --------------------------
//        User Handlers
// --------------------------

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var dbUser User
	err := db.QueryRow(`
        SELECT id, username, password, permissions
        FROM users
        WHERE username=$1
    `, creds.Username).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password, &dbUser.Permissions)
	if err != nil {
		http.Error(w, "User not found or DB error", http.StatusUnauthorized)
		return
	}

	if !checkPasswordHash(creds.Password, dbUser.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	tokenString, err := generateJWT(dbUser.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "login successful",
		"token":       tokenString,
		"permissions": dbUser.Permissions,
	})
}

// Admin-only: create user
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Error(w, "Forbidden - Admins only", http.StatusForbidden)
		return
	}

	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPass, err := hashPassword(newUser.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	err = db.QueryRow(`
        INSERT INTO users (username, password, permissions)
        VALUES ($1, $2, $3)
        RETURNING id
    `, newUser.Username, hashedPass, newUser.Permissions).Scan(&newUser.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// Admin-only: update user
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Error(w, "Forbidden - Admins only", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPass, err := hashPassword(updatedUser.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	result, err := db.Exec(`
        UPDATE users
        SET username=$1,
            password=$2,
            permissions=$3
        WHERE id=$4
    `, updatedUser.Username, hashedPass, updatedUser.Permissions, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating user: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking rows affected: %v", err), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

// Admin-only: delete user
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		http.Error(w, "Forbidden - Admins only", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM users WHERE id=$1", userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting user: %v", err), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking rows affected: %v", err), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

// --------------------------
//   Budgets Handlers
// --------------------------

// GET /api/budgets => returns budgets belonging to the JWT user
func getBudgetsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	rows, err := db.Query(`
		SELECT id, name, amount, category, period, user_id
		FROM budgets
		WHERE user_id=$1
	`, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying budgets: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var budgets []Budget
	for rows.Next() {
		var b Budget
		if err := rows.Scan(&b.ID, &b.Name, &b.Amount, &b.Category, &b.Period, &b.UserID); err != nil {
			http.Error(w, fmt.Sprintf("Error scanning budget: %v", err), http.StatusInternalServerError)
			return
		}
		budgets = append(budgets, b)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(budgets)
}

// POST /api/budgets => create a new budget for the JWT user
func createBudgetHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var b Budget
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Override the user_id from the token
	b.UserID = userID

	err = db.QueryRow(`
		INSERT INTO budgets (name, amount, category, period, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id

	`, b.Name, b.Amount, b.Category, b.Period, b.UserID).Scan(&b.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting budget: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(b)
}

// PUT /api/budgets/{id} => update a budget that belongs to the JWT user
func updateBudgetHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	budgetIDStr := vars["id"]
	budgetID, err := strconv.Atoi(budgetIDStr)
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}

	var b Budget
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Only update if user_id matches the JWT user
	result, err := db.Exec(`
		UPDATE budgets
		SET name=$1, amount=$2, category=$3, period=$4
		WHERE id=$5 AND user_id=$6

	`,
		b.Name, b.Amount, b.Category, b.Period, budgetID, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating budget: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Budget not found or not owned by user", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Budget updated successfully"})
}

// DELETE /api/budgets/{id} => delete a budget that belongs to the JWT user
func deleteBudgetHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	budgetIDStr := vars["id"]
	budgetID, err := strconv.Atoi(budgetIDStr)
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		DELETE FROM budgets
		WHERE id=$1 AND user_id=$2
	`, budgetID, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting budget: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Budget not found or not owned by user", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Budget deleted successfully"})
}

// --------------------------
//    Charges Handlers
// --------------------------

// GET /api/charges => get all charges for the JWT user
func getChargesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.Query(`
		SELECT id, name, amount, category, periodical, user_id, created_at
		FROM charges
		WHERE user_id=$1

    `, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying charges: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var charges []Charge
	for rows.Next() {
		var c Charge
		if err := rows.Scan(&c.ID, &c.Name, &c.Amount, &c.Category, &c.Periodical, &c.UserID, &c.CreatedAt); err != nil {
			http.Error(w, fmt.Sprintf("Error scanning charge: %v", err), http.StatusInternalServerError)
			return
		}
		charges = append(charges, c)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(charges)
}

// POST /api/charges => create a new charge for the JWT user
func createChargeHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var c Charge
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Force user_id to the JWT user
	c.UserID = userID

	err = db.QueryRow(`
		INSERT INTO charges (name, amount, category, periodical, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
    `, c.Name, c.Amount, c.Category, c.Periodical, c.UserID).Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting charge: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

// PUT /api/charges/{id} => update a charge that belongs to the JWT user
func updateChargeHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	chargeIDStr := vars["id"]
	chargeID, err := strconv.Atoi(chargeIDStr)
	if err != nil {
		http.Error(w, "Invalid charge ID", http.StatusBadRequest)
		return
	}

	var c Charge
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Only update if charge belongs to user
	result, err := db.Exec(`
		UPDATE charges
		SET name=$1, amount=$2, category=$3, periodical=$4
		WHERE id=$5 AND user_id=$6
    `, c.Name, c.Amount, c.Category, c.Periodical, chargeID, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating charge: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Charge not found or not owned by user", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Charge updated successfully"})
}

// DELETE /api/charges/{id} => delete a charge that belongs to the JWT user
func deleteChargeHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	chargeIDStr := vars["id"]
	chargeID, err := strconv.Atoi(chargeIDStr)
	if err != nil {
		http.Error(w, "Invalid charge ID", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
        DELETE FROM charges
        WHERE id=$1 AND user_id=$2
    `, chargeID, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting charge: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Charge not found or not owned by user", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Charge deleted successfully"})
}

// --------------------------
//       Shares Handlers
// --------------------------

// GET /api/shares => return any shares where the JWT user is user_id OR user_share_id
func getSharesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.Query(`
        SELECT id, user_id, user_share_id, access
        FROM shares
        WHERE user_id=$1 OR user_share_id=$1
    `, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching shares: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var shares []Share
	for rows.Next() {
		var s Share
		if err := rows.Scan(&s.ID, &s.UserID, &s.UserShareID, &s.Access); err != nil {
			http.Error(w, fmt.Sprintf("Error scanning share: %v", err), http.StatusInternalServerError)
			return
		}
		shares = append(shares, s)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(shares)
}

// POST /api/shares => create a share
// Request body might look like: { "shareUsername": "bob", "access": "read-only" }
func createShareHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get the user_id from JWT
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Parse request
	var requestBody struct {
		ShareUsername string `json:"shareUsername"`
		Access        string `json:"access"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 3. Look up the user_share_id by the given username
	var userShareID int
	err = db.QueryRow(`
        SELECT id FROM users WHERE username=$1
    `, requestBody.ShareUsername).Scan(&userShareID)
	if err != nil {
		http.Error(w, "No user found with that username", http.StatusNotFound)
		return
	}

	// 4. Insert into 'shares' table
	var newShare Share
	newShare.UserID = userID
	newShare.UserShareID = userShareID
	newShare.Access = requestBody.Access

	err = db.QueryRow(`
        INSERT INTO shares (user_id, user_share_id, access)
        VALUES ($1, $2, $3)
        RETURNING id
    `, newShare.UserID, newShare.UserShareID, newShare.Access).Scan(&newShare.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating share: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newShare)
}

// DELETE /api/shares/{id} => delete a share if the JWT user is either user_id or user_share_id
func deleteShareHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	shareIDStr := vars["id"]
	shareID, err := strconv.Atoi(shareIDStr)
	if err != nil {
		http.Error(w, "Invalid share ID", http.StatusBadRequest)
		return
	}

	// We only allow delete if the current user is user_id or user_share_id
	result, err := db.Exec(`
        DELETE FROM shares
        WHERE id=$1
          AND (user_id=$2 OR user_share_id=$2)
    `, shareID, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting share: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Share not found or you are not allowed to delete it", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Share deleted successfully"})
}
