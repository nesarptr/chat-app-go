package models

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	UserName string `json:"username" validate:"required,min=3,max=32" gorm:"unique"`
	Password string `json:"password" validate:"required,min=6"`
	Texts    []Text `json:"texts" validate:"-" gorm:"ForeignKey:SenderID;AssociationForeignKey:ID"`
}

func (c *Client) Create(db *gorm.DB) error {
	if db.Create(c).Error != nil {
		return fiber.ErrBadRequest
	}
	return nil
}
