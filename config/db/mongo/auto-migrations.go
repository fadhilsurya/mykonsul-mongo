package db

import (
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrationData() {
	direction := flag.String("direction", "up", "Migration direction (up/down)")
	flag.Parse()

	m, err := migrate.New(
		"file:./migrations",
		"mongodb://localhost:27017/mykonsul")
	if err != nil {
		log.Fatal(err)
	}

	if *direction == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migration applied successfully (up).")
	} else if *direction == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migration rolled back successfully (down).")
	} else {
		log.Fatalf("Invalid direction: %s. Use 'up' or 'down'.", *direction)
	}
}
