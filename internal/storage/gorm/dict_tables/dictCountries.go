package dict_tables

import (
	"github.com/KennyMacCormik/HerdMaster/internal/storage/gorm/models"
	"gorm.io/gorm"
)

var dictCountries = []models.DictCountry{
	{
		ID:   1,
		Name: "United States",
		Code: "USA",
	},
	{
		ID:   2,
		Name: "Canada",
		Code: "CAN",
	},
}

func isEmptyCountry(m models.DictCountry) bool {
	if m.Code == "" && m.Name == "" {
		return true
	}
	return false
}

func FillDictCountries(db *gorm.DB) {
	for _, country := range dictCountries {
		var result models.DictCountry
		db.First(&result, "code = ?", country.Code)
		if isEmptyCountry(result) {
			db.First(&result, "name = ?", country.Name)
		}

		if isEmptyCountry(result) {
			db.Create(&country)
		} else {
			db.Model(&models.DictCountry{}).Where(
				"id = ?", result.ID).Updates(
				models.DictCountry{Name: country.Name, Code: country.Code})
		}
	}
}
