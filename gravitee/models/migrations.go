package models

import (
	"fmt"
	"github.com/yufenghui/go/gravitee/log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/spring2go/gravitee/util/migrations"
)

var (
	list = []migrations.MigrationStage{
		{
			Name:     "initial",
			Function: migrate0001,
		},
	}
)

func init() {
	if err := os.Chdir("../"); err != nil {
		log.ERROR.Fatal(err)
	}
}

// MigrateAll executes all migrations
func MigrateAll(db *gorm.DB) error {
	return migrations.Migrate(db, list)
}

func migrate0001(db *gorm.DB, name string) error {
	//-------------
	// OAUTH models
	//-------------

	// Create tables
	if err := db.CreateTable(new(OauthClient)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_clients table: %s", err)
	}

	return nil
}
