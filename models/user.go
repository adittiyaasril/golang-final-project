package models

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
)

type User struct {
	ID              uint64 `gorm:"type:bigint;primaryKey"`
	Username        string `gorm:"unique;not null" valid:"required"`
	Email           string `gorm:"unique;not null" valid:"email,required"`
	Password        string `gorm:"not null" valid:"required"`
	Age             int    `gorm:"not null;check:age > 8" valid:"required"`
	ProfileImageURL string `gorm:"type:text" valid:"url"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Photos       []Photo       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Comments     []Comment     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	SocialMedias []SocialMedia `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type LoginRequest struct {
	Email    string `valid:"email,required"`
	Password string `valid:"required"`
}

type UpdateUserRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Age             int    `json:"age"`
	ProfileImageURL string `json:"profileImageURL"`
}

func (user *User) Validate() error {
	if _, err := govalidator.ValidateStruct(user); err != nil {
		return err
	}
	if len(user.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	return nil
}

func (login *LoginRequest) ValidateLogin() error {
	if _, err := govalidator.ValidateStruct(login); err != nil {
		return err
	}
	if len(login.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	return nil
}

func (updateUserRequest *UpdateUserRequest) Validate() error {
	if updateUserRequest.Username == "" && updateUserRequest.Email == "" && updateUserRequest.Age == 0 && updateUserRequest.ProfileImageURL == "" {
		return errors.New("at least one field must be provided for update")
	}
	return nil
}
