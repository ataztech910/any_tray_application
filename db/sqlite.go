package sqlite

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (db *sql.DB, err error) {
	os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, err := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	return sqliteDatabase, nil
}

func CreateUsersTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE sessions (
		"_id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"sessionKey" TEXT
	  );` // SQL Statement for Create Table

	log.Println("Create sessions table...")
	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("sessions table created")
}

func InstertJwt(db *sql.DB, key string) {
	log.Println("Inserting student record ...")
	insertSessionsSQL := `INSERT INTO sessions(sessionKey) VALUES (?)`
	statement, err := db.Prepare(insertSessionsSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(key)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
