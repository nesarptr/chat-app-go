package models

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	UserName   string `json:"username" validate:"required,min=3,max=32" gorm:"unique"`
	Password   string `json:"-" validate:"required,min=6"`
	LatestText string `json:"latestMessage" validate:"-"`
}

func (c *Client) Create(db *gorm.DB) error {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	c.Password = string(hashedPw)
	if db.Create(c).Error != nil {
		return fiber.ErrBadRequest
	}
	return nil
}
