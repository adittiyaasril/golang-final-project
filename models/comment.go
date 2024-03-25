package models

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Comment struct {
	ID        uint64 `gorm:"type:bigint;primaryKey"`
	UserID    uint64 `gorm:"not null"`
	PhotoID   uint64 `gorm:"not null"`
	Message   string `gorm:"not null" valid:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (comment *Comment) Validate() error {
	if _, err := govalidator.ValidateStruct(comment); err != nil {
		return err
	}
	return nil
}
