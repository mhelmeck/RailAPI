package dbutils

import (
	"database/sql"
	"log"
)

func Initialize(dbDriver *sql.DB) {
	statement, _ := dbDriver.Prepare(train)
	statement.Exec()

	statement, _ = dbDriver.Prepare(station)
	statement.Exec()

	statement, _ = dbDriver.Prepare(schedule)
	statement.Exec()

	log.Println("All tables created/initialized successfully!")
}
