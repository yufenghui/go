package cmd

import (
	"github.com/jinzhu/gorm"
	"github.com/yufenghui/go/gravitee/config"
	"github.com/yufenghui/go/gravitee/database"
)

func initConfig(configFile string) (*config.Config, *gorm.DB, error) {

	// config
	cfg := config.NewConfig(configFile)

	// database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		return nil, nil, err
	}

	return cfg, db, nil
}
