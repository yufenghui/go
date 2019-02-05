package services

import (
	"github.com/jinzhu/gorm"
	"github.com/yufenghui/go/gravitee/config"
	"github.com/yufenghui/go/gravitee/services/health"
	"reflect"
)

var (
	// HealthService ...
	HealthService *health.Service
)

func Init(cfg *config.Config, db *gorm.DB) error {

	if nil == reflect.TypeOf(HealthService) {
		HealthService = health.NewService(db)
	}

	return nil
}

// Close closes any open services
func Close() {

}
