package init

import (
	"github.com/KennyMacCormik/HerdMaster/internal/config"
	"github.com/KennyMacCormik/HerdMaster/internal/storage"
	"github.com/KennyMacCormik/HerdMaster/internal/storage/gorm"
)

func StorageDB(conf config.Config) (storage.DB, error) {
	db, err := gorm.New(conf)
	if err != nil {
		return nil, err
	}
	return db, nil
}
