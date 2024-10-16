package models

import "gorm.io/gorm"

type Record struct {
	gorm.Model
	Name   string
	RrType uint16
	Class  uint16
	Record string
}
