package models

import "gorm.io/gorm"

type Owner struct {
	gorm.Model
	Name      string
	Email     string
	Phone     string
	AddressID uint `gorm:"index"`
	Address   Address
}

type Address struct {
	gorm.Model
	Street     string
	PostalCode string
	CityID     uint `gorm:"index"`
	City       DictCity
}

type DictCountry struct {
	ID   uint   `gorm:"autoIncrement"`
	Name string `gorm:"unique;not null;uniqueIndex:idx_name_code"`
	Code string `gorm:"unique;not null;uniqueIndex:idx_name_code;check:length(Code) < 4"`
}

type DictState struct {
	ID        uint   `gorm:"autoIncrement"`
	Name      string `gorm:"not null;uniqueIndex:idx_name_code_country_id"`
	Code      string `gorm:"not null;uniqueIndex:idx_name_code_country_id;check:length(Code) < 6"`
	CountryID uint   `gorm:"not null;uniqueIndex:idx_name_code_country_id"`
	Country   DictCountry
}

type DictCity struct {
	ID      uint   `gorm:"autoIncrement"`
	Name    string `gorm:"not null;uniqueIndex:idx_name_state_id"`
	StateID uint   `gorm:"not null;uniqueIndex:idx_name_state_id"`
	State   DictState
}
