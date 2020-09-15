package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost port=5432 user=postgres password=120505Aa dbname=tavern_of_games_db sslmode=disable"
	}

	os.Exit(m.Run())
}
