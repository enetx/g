package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	. "github.com/enetx/g"
	// _ "modernc.org/sqlite"
)

type TestDB struct {
	ID        int64             `db:"id"`
	Name      Option[String]    `db:"name"`
	Code      Option[Float]     `db:"code"`
	UpdatedAt Option[time.Time] `db:"updated_at"`
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE testdb (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			code REAL,
			updated_at DATETIME
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	c := TestDB{
		Name:      Some(String("Test")),
		Code:      Some(Float(42.4)),
		UpdatedAt: None[time.Time](), // NULL
	}

	_, err = db.Exec(
		`INSERT INTO testdb (name, code, updated_at) VALUES (?, ?, ?)`,
		c.Name, c.Code, c.UpdatedAt,
	)
	if err != nil {
		log.Fatal("Insert error:", err)
	}

	var testdb TestDB
	row := db.QueryRow(`SELECT id, name, code, updated_at FROM testdb WHERE code = ?`, 42.4)

	if err := row.Scan(
		&testdb.ID,
		&testdb.Name,
		&testdb.Code,
		&testdb.UpdatedAt,
	); err != nil {
		log.Fatal("Scan error:", err)
	}

	fmt.Println("ID:", testdb.ID)

	if testdb.Name.IsSome() {
		fmt.Println("Name:", testdb.Name.Unwrap())
	}

	if testdb.Code.IsSome() {
		fmt.Println("Code:", testdb.Code.Unwrap())
	}

	if testdb.UpdatedAt.IsNone() {
		fmt.Println("UpdatedAt is NULL")
	}

	testdb.UpdatedAt = Some(time.Now())
	fmt.Println("Set UpdatedAt:", testdb.UpdatedAt.Unwrap())
}
