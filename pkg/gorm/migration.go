package gorm

import (
	"gorm.io/gorm"

	"github.com/cewuandy/go-restful-dns/internal/repository/db/models"
)

func AutoMigrate(db *gorm.DB) error {
	var err error
	err = db.AutoMigrate(&models.Record{})
	if err != nil {
		return err
	}
	return nil
}
