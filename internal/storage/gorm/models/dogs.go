package models

import (
	"gorm.io/gorm"
	"time"
)

type ShepherdDog struct {
	gorm.Model
	name            string
	BreedID         uint `gorm:"index"`
	Breed           DictBreed
	Age             int
	GenderID        uint `gorm:"index"`
	Gender          DictGender
	WeightG         int
	HeightMm        int
	CoatID          uint `gorm:"index"`
	Coat            DictCoat
	MicrochipNumber int
	OwnerID         uint `gorm:"index"`
	Owner           Owner
	Vaccinated      bool
	DateOfBirth     time.Time
}

type DictCoat struct {
	ID   uint `gorm:"primary_key"`
	name string
}

type DictBreed struct {
	ID     uint `gorm:"primary_key"`
	name   string
	origin string
}

type DictGender struct {
	ID   uint `gorm:"primary_key"`
	name string
}
