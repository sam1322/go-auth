package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	// _ "github.com/lib/pq" // Import the PostgreSQL driver (not being used as it is in maintenance mode is not actively developed)
)

type Service interface {
	Health() map[string]string
	// QueryBoxes() map[string]string
	QueryBoxes() map[string]interface{}
	GetDB() *sql.DB
}

type service struct {
	db *sql.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() Service {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)

	if err != nil {
		log.Fatal(err)
	}
	s := &service{db: db}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	// this will be printed in the terminal, confirming the connection to the database
	fmt.Println("The Database is connected")
	return s
}

func (s *service) GetDB() *sql.DB {
	return s.db
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) QueryBoxes() map[string]interface{} {
	// func (s *service) QueryBoxes() map[string]string {
	type user struct {
		Boxid     int       `json:"box_id"`
		Username  string    `json:"name"`
		Address   string    `json:"address"`
		CreatedAt time.Time `json:"created_at"`
		// UpdatedAt time.Time `json:"updated_at"`
	}

	query := `SELECT box_id, name, address, created_at  FROM boxes`

	rows, err := s.db.Query(query)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer rows.Close()

	fmt.Println("Rows", rows)

	// err = rows.Err()

	// CheckError(err)

	var users []user
	for rows.Next() {
		var u user
		// err := rows.Scan(&u.Boxid, &u.Username, &u.Address)
		err := rows.Scan(&u.Boxid, &u.Username, &u.Address, &u.CreatedAt)
		if err != nil {
			log.Fatalf(err.Error())
		}
		users = append(users, u)
	}

	// fmt.Printf("%#v", users)
	// for _, user := range users {
	// fmt.Printf("%#v\n", user)
	empJSON, err := json.MarshalIndent(&users, "", " ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("Marshal function output %s\n", string(empJSON))

	return map[string]interface{}{
		"message": "It's healthy",
		"users":   users,
	}
}
