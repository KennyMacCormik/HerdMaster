package gorm

import (
	"fmt"
	"github.com/KennyMacCormik/HerdMaster/internal/config"
	"github.com/KennyMacCormik/HerdMaster/internal/storage/gorm/dict_tables"
	"github.com/KennyMacCormik/HerdMaster/internal/storage/gorm/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqLite struct {
	db *gorm.DB
}

func (s *SqLite) Close() error {
	db, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB connection: %w", err)
	}
	return db.Close()
}

func New(conf config.Config) (*SqLite, error) {
	sql := SqLite{}

	db, err := gorm.Open(sqlite.Open(conf.DB.ConnString), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}
	db.Exec("PRAGMA foreign_keys = ON;")

	if conf.DB.AutoMigrate {
		err = db.AutoMigrate(&models.Owner{}, &models.ShepherdDog{})
		if err != nil {
			return nil, err
		}
	}

	if conf.DB.AutoFillDict {
		dict_tables.FillDictCountries(db)
		dict_tables.FillDictStates(db)
		dict_tables.FillDictCities(db)
	}

	sql.db = db

	return &sql, nil
}
