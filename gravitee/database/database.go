package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yufenghui/go/gravitee/config"
	"time"
)

func init() {
	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}
}

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	databaseType := cfg.Database.Type

	if databaseType == "mysql" {

		args := fmt.Sprintf(
			"%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DatabaseName,
		)

		db, err := gorm.Open(databaseType, args)
		if err != nil {
			return nil, err
		}

		db.DB().SetMaxIdleConns(cfg.Database.MaxIdleConns)
		db.DB().SetMaxOpenConns(cfg.Database.MaxOpenConns)

		db.LogMode(cfg.IsDevelop)

		return db, nil
	}

	return nil, fmt.Errorf("Database type %s not supported", databaseType)
}
