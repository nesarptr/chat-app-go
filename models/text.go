package models

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Text struct {
	gorm.Model
	Body         string `json:"content" validate:"required,min=1"`
	SenderName   string `json:"from" validate:"required,min=3,max=32,nefield=ReceiverName"`
	Sender       Client `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:SenderName;references:UserName" validate:"-"`
	ReceiverName string `json:"to" validate:"required,min=3,max=32,nefield=SenderName"`
	Receiver     Client `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ReceiverName;references:UserName" validate:"-"`
}

func (t *Text) Create(db *gorm.DB) error {
	if db.Create(t).Error != nil {
		return fiber.ErrBadRequest
	}
	return nil
}
