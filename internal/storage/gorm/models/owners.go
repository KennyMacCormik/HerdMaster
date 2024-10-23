package models

import "gorm.io/gorm"

type Owner struct {
	gorm.Model
	Name      string  `gorm:"not null"`
	Email     string  `gorm:"not null"`
	Phone     string  `gorm:"not null"`
	AddressID uint    `gorm:"not null;index"`
	Address   Address `gorm:"constraint:OnUpdate:CASCADE"`
}

type Address struct {
	gorm.Model
	Street     string   `gorm:"not null;uniqueIndex:idx_street_postal_code_city_id"`
	PostalCode string   `gorm:"not null;uniqueIndex:idx_street_postal_code_city_id"`
	CityID     uint     `gorm:"not null;uniqueIndex:idx_street_postal_code_city_id"`
	City       DictCity `gorm:"constraint:OnUpdate:CASCADE"`
}

type DictCountry struct {
	ID   uint   `gorm:"autoIncrement"`
	Name string `gorm:"unique;not null;uniqueIndex:idx_name_code"`
	Code string `gorm:"unique;not null;uniqueIndex:idx_name_code;check:length(Code) < 4"`
}

type DictState struct {
	ID        uint        `gorm:"autoIncrement"`
	Name      string      `gorm:"not null;uniqueIndex:idx_name_code_country_id"`
	Code      string      `gorm:"not null;uniqueIndex:idx_name_code_country_id;check:length(Code) < 6"`
	CountryID uint        `gorm:"not null;uniqueIndex:idx_name_code_country_id"`
	Country   DictCountry `gorm:"constraint:OnUpdate:CASCADE"`
}

type DictCity struct {
	ID      uint      `gorm:"autoIncrement"`
	Name    string    `gorm:"not null;uniqueIndex:idx_name_state_id"`
	StateID uint      `gorm:"not null;uniqueIndex:idx_name_state_id"`
	State   DictState `gorm:"constraint:OnUpdate:CASCADE"`
}
