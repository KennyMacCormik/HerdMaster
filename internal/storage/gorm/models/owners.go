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
	ID   uint `gorm:"primary_key"`
	Name string
	Code string
}

type DictState struct {
	ID        uint `gorm:"primary_key"`
	Name      string
	CountryID uint `gorm:"index"`
	Country   DictCountry
}

type DictCity struct {
	ID      uint `gorm:"primary_key"`
	Name    string
	StateID uint `gorm:"index"`
	State   DictState
}
