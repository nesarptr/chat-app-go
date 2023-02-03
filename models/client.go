package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	UserName     string  `json:"username" validate:"required,min=3,max=32" gorm:"unique"`
	Password string  `json:"password" validate:"required,min=6"`
	Texts    []Text  `json:"texts" validate:"-"`
}