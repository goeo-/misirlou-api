package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DB is a simple wrapper around a gorm.DB which has all the methods of the
// models package.
type DB struct {
	db *gorm.DB
}

// CreateDB creates a new instance of the gorm database and sets up the
// migrations
func CreateDB(dsn string) (*DB, error) {
	gdb, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db := &DB{gdb}
	return db, db.runMigrations()
}

func (db *DB) runMigrations() error {
	fmt.Println("Migrating database...")
	start := time.Now()
	if !db.db.HasTable("db_version") {
		db.db.Exec("CREATE TABLE db_version(version INT NOT NULL);")
		db.db.Exec("INSERT INTO db_version(version) VALUES ('-1');")
	}
	if db.db.Error != nil {
		return db.db.Error
	}
	var versions []int
	var res *gorm.DB
	if res = db.db.Table("db_version").Pluck("version", &versions); res.Error != nil {
		return res.Error
	}
	version := versions[0] + 1
	for ; ; version++ {
		data, err := ioutil.ReadFile(fmt.Sprintf("migrations/%d.sql", version))
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println(time.Since(start))
				return nil
			}
			return err
		}
		if res = db.db.Exec(string(data)); res.Error != nil {
			return db.db.Error
		}
		if res = db.db.Exec("UPDATE db_version SET version = ?", version); res.Error != nil {
			return db.db.Error
		}
	}
}
