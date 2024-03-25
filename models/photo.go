package models

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
)

type Photo struct {
	ID        uint64 `gorm:"type:bigint;primaryKey"`
	Title     string `gorm:"not null" valid:"required"`
	Caption   string
	PhotoURL  string `gorm:"not null" valid:"url,required"`
	UserID    uint64 `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Comments []Comment `gorm:"foreignKey:PhotoID;constraint:OnDelete:CASCADE"`
}

func (photo *Photo) Validate() error {
	if _, err := govalidator.ValidateStruct(photo); err != nil {
		return err
	}
	return nil
}

func (photo *Photo) ValidateUpdate(updatedPhoto Photo) error {
	if updatedPhoto.Title == "" && updatedPhoto.Caption == "" && updatedPhoto.PhotoURL == "" {
		return errors.New("no data provided for update")
	}

	if updatedPhoto.PhotoURL != "" {
		if !govalidator.IsURL(updatedPhoto.PhotoURL) {
			return errors.New("invalid photo url format")
		}
	}

	return nil
}
