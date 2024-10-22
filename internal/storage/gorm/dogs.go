package gorm

import (
	"gorm.io/gorm"
	"time"
)

type ShepherdDogs struct {
	gorm.Model
	name            string
	BreedID         uint
	Breed           DictBreed
	Age             int
	GenderID        uint
	Gender          DictGender
	WeightG         int
	HeightMm        int
	CoatID          uint
	Coat            DictCoat
	MicrochipNumber int
	OwnerID         uint
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
