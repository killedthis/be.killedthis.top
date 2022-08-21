package builder

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

const queryServiceProviders = "SELECT DISTINCT serviceprovider FROM killedthis"
const queryKilledShowsByProvider = "SELECT * FROM killedthis WHERE serviceprovider = ? ORDER BY date"

type Database struct {
	Db *sql.DB
}

type KilledShow struct {
	Id              int64   `db:"index"`
	Title           string  `db:"title"`
	ServiceProvider string  `db:"serviceprovider"`
	Brand           *string `db:"brand"`
	Date            *string `db:"data"`
	DateAdded       *string `db:"dateadded"`
	Reason          *string `db:"reason"`
	TmdbId          *int64  `db:"tmdbid"`
}

func OpenDatabase() *Database {
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_HOST"),
		DBName: os.Getenv("DB_NAME"),
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	return &Database{
		Db: db,
	}
}

func (m *Database) GetServiceProviders() []string {
	rows, err := m.Db.Query(queryServiceProviders)
	if err != nil {
		log.Println("failed to query database: ", err)
		return nil
	}

	serviceProviders := make([]string, 0)
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Println("failed to read name fields: ", err)
			return nil
		}

		serviceProviders = append(serviceProviders, name)
	}

	return serviceProviders
}

func (m *Database) GetShowsByProvider(provider string) []KilledShow {
	if provider == "" {
		log.Println("invalid serviceprovider")
		return nil
	}

	rows, err := m.Db.Query(queryKilledShowsByProvider, provider)
	if err != nil {
		log.Println("failed to query shows: ", err)
		return nil
	}

	queriedShows := make([]KilledShow, 0)
	for rows.Next() {
		var killedShow KilledShow
		err := rows.Scan(
			&killedShow.Id,
			&killedShow.Title,
			&killedShow.ServiceProvider,
			&killedShow.Brand,
			&killedShow.Date,
			&killedShow.DateAdded,
			&killedShow.Reason,
			&killedShow.TmdbId)
		if err != nil {
			log.Println("failed to parse row: ", err)
			return nil
		}

		queriedShows = append(queriedShows, killedShow)
	}

	return queriedShows
}
