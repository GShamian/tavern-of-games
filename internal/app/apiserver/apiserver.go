package apiserver

import (
	"database/sql"
	"net/http"

	"github.com/GShamian/tavern-of-games/internal/app/store/sqlstore"
)

// Start function. Starts server.
func Start(config *Config) error {
	//Getting a pointer to our db and getting an access to it.
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	//Creating Store instance with our db. Check store.go documentation.
	store := sqlstore.New(db)
	//Creating server instance with our store. Check server.go documentation.
	srv := newServer(store)
	//Starting srv server with address from config
	return http.ListenAndServe(config.BindAddr, srv)
}

// new DB functions. Constructor for DB. Importing a db url to get an access to db.
// As a result we get a pointer on our target db.
func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
