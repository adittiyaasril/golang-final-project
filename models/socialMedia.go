package models

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type SocialMedia struct {
	ID             uint64 `gorm:"type:bigint;primaryKey"`
	Name           string `gorm:"not null" valid:"required"`
	SocialMediaURL string `gorm:"not null" valid:"url,required"`
	UserID         uint64 `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (socialMedia *SocialMedia) Validate() error {
	if _, err := govalidator.ValidateStruct(socialMedia); err != nil {
		return err
	}
	return nil
}
