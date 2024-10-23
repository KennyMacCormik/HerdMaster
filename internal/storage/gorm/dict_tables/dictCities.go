package dict_tables

import (
	"github.com/KennyMacCormik/HerdMaster/internal/storage/gorm/models"
	"gorm.io/gorm"
)

var dictCities = []models.DictCity{
	{
		Name:    "Los Angeles",
		StateID: 1,
	},
	{
		Name:    "San Francisco",
		StateID: 1,
	},
	{
		Name:    "Houston",
		StateID: 2,
	},
	{
		Name:    "Dallas",
		StateID: 2,
	},
	{
		Name:    "New York City",
		StateID: 3,
	},
	{
		Name:    "Buffalo",
		StateID: 3,
	},
	{
		Name:    "Toronto",
		StateID: 4,
	},
	{
		Name:    "Ottawa",
		StateID: 4,
	},
	{
		Name:    "Vancouver",
		StateID: 5,
	},
	{
		Name:    "Victoria",
		StateID: 5,
	},
	{
		Name:    "Montreal",
		StateID: 6,
	},
	{
		Name:    "Quebec City",
		StateID: 6,
	},
}

func FillDictCities(db *gorm.DB) {
	for _, state := range dictCities {
		var result models.DictCity
		db.First(&result, "name = ?", state.Name)

		if result.Name == "" {
			db.Create(&state)
		} else {
			db.Model(&models.DictCity{}).Where(
				"id = ?", result.ID).Updates(
				models.DictCity{Name: state.Name, StateID: state.StateID})
		}
	}
}
