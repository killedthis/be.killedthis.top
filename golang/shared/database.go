package shared

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

const queryServiceProviders = "SELECT DISTINCT serviceprovider FROM killedthis"
const queryKilledShowsByProvider = "SELECT * FROM killedthis WHERE serviceprovider = ? ORDER BY date"

type Database struct {
	Db *sql.DB
}

type KilledShow struct {
	Id              int        `db:"index"`
	Title           string     `db:"title"`
	ServiceProvider string     `db:"serviceprovider"`
	Brand           *string    `db:"brand"`
	Date            *time.Time `db:"data"`
	DateAdded       *time.Time `db:"dateadded"`
	Reason          *string    `db:"reason"`
	TmdbId          *int       `db:"tmdbid"`
}

func OpenDatabase(dbCfg *DatabaseConfig) *Database {
	cfg := mysql.Config{
		User:      dbCfg.Username,
		Passwd:    dbCfg.Password,
		Net:       "tcp",
		Addr:      dbCfg.Hostname,
		DBName:    dbCfg.Schema,
		ParseTime: true,
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
		killedShow, err := scanShow(rows)
		if err != nil {
			log.Println("failed to parse row: ", err)
			return nil
		}

		queriedShows = append(queriedShows, killedShow)
	}

	return queriedShows
}

func scanShow(rows *sql.Rows) (show KilledShow, err error) {
	err = rows.Scan(
		&show.Id,
		&show.Title,
		&show.ServiceProvider,
		&show.Brand,
		&show.Date,
		&show.DateAdded,
		&show.Reason,
		&show.TmdbId)

	return show, err
}

func (show KilledShow) Year() int {
	return show.Date.Year()
}

func (show KilledShow) Month() string {
	i := int(show.Date.Month())
	if i < 10 {
		return fmt.Sprintf("0%d", i)
	} else {
		return fmt.Sprintf("%d", i)
	}
}

func (show KilledShow) TmdbPoster() string {
	if show.TmdbId != nil {
		return fmt.Sprintf("%d.jpg", *show.TmdbId)
	} else {
		// Return not available poster
		return "not_available.jpg"
	}
}
