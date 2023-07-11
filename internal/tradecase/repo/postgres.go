package repo

import (
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
}

func New(dsn string, conf *gorm.Config) {
}
