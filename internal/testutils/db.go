package testutils

import (
	"github.com/jmoiron/sqlx"
	"go-labs/internal/repository/db"
	"log"
	"sync"
)

var (
	testDB *sqlx.DB
	dbOnce sync.Once
)

func DB() *sqlx.DB {
	dbOnce.Do(func() {
		const driver = "sqlite3"
		var err error

		testDB, err = sqlx.Connect(driver, "file::memory:?cache=shared&mode=rwc")
		if err != nil {
			log.Panicf("connect: %v", err)
		}

		if err = db.Migrate(testDB, driver); err != nil {
			log.Panicf("db migrate: %v\n", err)
		}
	})

	return testDB
}
