package models

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Text struct {
	gorm.Model
	Body     string `json:"body" validate:"required,min=1"`
	From     string `json:"from" validate:"required,min=3,max=32,nefield=To"`
	Sender   Client `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:From;references:UserName" validate:"-"`
	To       string `json:"to" validate:"required,min=3,max=32,nefield=From"`
	Receiver Client `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:To;references:UserName" validate:"-"`
}

func (t *Text) Create(db *gorm.DB) error {
	if db.Create(t).Error != nil {
		return fiber.ErrBadRequest
	}
	return nil
}
