package dict_tables

import (
	"github.com/KennyMacCormik/HerdMaster/internal/storage/gorm/models"
	"gorm.io/gorm"
)

var dictStates = []models.DictState{
	{
		Name:      "California",
		Code:      "CA",
		CountryID: 1,
	},
	{
		Name:      "Texas",
		Code:      "TX",
		CountryID: 1,
	},
	{
		Name:      "New York",
		Code:      "NY",
		CountryID: 1,
	},
	{
		Name:      "Ontario",
		Code:      "ON",
		CountryID: 2,
	},
	{
		Name:      "British Columbia",
		Code:      "BC",
		CountryID: 2,
	},
	{
		Name:      "Quebec",
		Code:      "QC",
		CountryID: 2,
	},
}

func isEmptyState(m models.DictState) bool {
	if m.Code == "" && m.Name == "" {
		return true
	}
	return false
}

func FillDictStates(db *gorm.DB) {
	for _, state := range dictStates {
		var result models.DictState
		db.First(&result, "code = ?", state.Code)
		if isEmptyState(result) {
			db.First(&result, "name = ?", state.Name)
		}

		if isEmptyState(result) {
			db.Create(&state)
		} else {
			db.Model(&models.DictState{}).Where(
				"id = ?", result.ID).Updates(
				models.DictState{Name: state.Name, Code: state.Code, CountryID: state.CountryID})
		}
	}
}
