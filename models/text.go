package models

import "gorm.io/gorm"

type Text struct {
	gorm.Model
	Body string `json:"body" validate:"required,min=1"`
	SenderID      uint    `json:"sender" validate:"required"`
	Sender      Client    `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:SenderID;references:ID" validate:"-"`
	ReceiverID      uint    `json:"receiver" validate:"required"`
	Receiver      Client    `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ReceiverID;references:ID" validate:"-"`
}