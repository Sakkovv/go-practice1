package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	ID      int     `db:"id"`
	Name    string  `db:"name"`
	Email   string  `db:"email"`
	Balance float64 `db:"balance"`
}

func main() {
	connStr := "host=localhost port=5430 user=user password=password dbname=mydatabase sslmode=disable"

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	fmt.Println("✅ Successfully connected to Postgres!")

	db.MustExec("TRUNCATE users RESTART IDENTITY")
	fmt.Println("\n--- ('users' table truncated for test) ---")

	fmt.Println("Inserting users...")
	user1 := User{Name: "Alice", Email: "alice@example.com", Balance: 1000.0}
	user2 := User{Name: "Bob", Email: "bob@example.com", Balance: 500.0}

	if err := InsertUser(db, user1); err != nil {
		log.Printf("Failed to insert Alice: %v", err)
	}
	if err := InsertUser(db, user2); err != nil {
		log.Printf("Failed to insert Bob: %v", err)
	}

	printAllUsers(db)

	fmt.Println("\nFinding user with ID=1...")
	alice, err := GetUserByID(db, 1)
	if err != nil {
		log.Printf("Failed to find user 1: %v", err)
	} else {
		fmt.Printf("Found: %s (Email: %s)\n", alice.Name, alice.Email)
	}

	fmt.Println("\nTransferring 150.50 from Alice (ID 1) to Bob (ID 2)...")
	err = TransferBalance(db, 1, 2, 150.50)
	if err != nil {
		log.Printf("Transfer error: %v", err)
	} else {
		fmt.Println("Transfer successful!")
	}
	printAllUsers(db)

	fmt.Println("\nAttempting to transfer 10000.0 from Bob (ID 2) to Alice (ID 1)...")
	err = TransferBalance(db, 2, 1, 10000.0)
	if err != nil {
		fmt.Printf("Transfer failed (as expected): %v\n", err)
	} else {
		fmt.Println("ERROR: Transfer succeeded but should not have!")
	}
	printAllUsers(db)
}

func InsertUser(db *sqlx.DB, user User) error {
	query := `INSERT INTO users (name, email, balance) VALUES (:name, :email, :balance)`
	_, err := db.NamedExec(query, user)
	return err
}

func GetAllUsers(db *sqlx.DB) ([]User, error) {
	var users []User
	err := db.Select(&users, "SELECT * FROM users ORDER BY id ASC")
	return users, err
}

func GetUserByID(db *sqlx.DB, id int) (User, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	return user, err
}

func TransferBalance(db *sqlx.DB, fromID int, toID int, amount float64) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var senderBalance float64
	err = tx.Get(&senderBalance, "SELECT balance FROM users WHERE id = $1 FOR UPDATE", fromID)
	if err != nil {
		return fmt.Errorf("failed to get sender balance: %w", err)
	}

	if senderBalance < amount {
		return fmt.Errorf("insufficient funds for user %d", fromID)
	}

	_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE id = $2", amount, fromID)
	if err != nil {
		return fmt.Errorf("failed to deduct funds: %w", err)
	}

	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", amount, toID)
	if err != nil {
		return fmt.Errorf("failed to add funds: %w", err)
	}

	return tx.Commit()
}

func printAllUsers(db *sqlx.DB) {
	users, err := GetAllUsers(db)
	if err != nil {
		log.Printf("Error getting users: %v", err)
		return
	}

	fmt.Println("--- Current Users in DB ---")
	if len(users) == 0 {
		fmt.Println("(Empty)")
	}
	for _, u := range users {
		fmt.Printf("  [ID: %d] %s - %s (Balance: %.2f)\n", u.ID, u.Name, u.Email, u.Balance)
	}
	fmt.Println("---------------------------------")
}
