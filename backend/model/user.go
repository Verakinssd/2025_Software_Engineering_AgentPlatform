package model

import (
	"backend/db"
	"backend/exceptions"
	"backend/util"
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username             string `gorm:"type:varchar(255);not null;unique"`
	Password             string `gorm:"type:varchar(255);not null"`
	Role                 string `gorm:"type:varchar(255);not null"`
	Email                string `gorm:"type:varchar(255);not null;unique"`
	Avatar               string `gorm:"type:varchar(255);not null"`
	Profile              string `gorm:"type:text;not null"`
	FavouriteWorkflowCnt uint   `gorm:"-"`
	PublishedWorkflowCnt uint   `gorm:"-"`
	FollowersCnt         uint   `gorm:"-"`
}

func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID                   uint   `json:"id"`
		Username             string `json:"username"`
		Role                 string `json:"role"`
		Email                string `json:"email"`
		Profile              string `json:"profile"`
		Avatar               string `json:"avatar"`
		FavouriteWorkflowCnt uint   `json:"favourite_workflow_cnt"`
		PublishedWorkflowCnt uint   `json:"published_workflow_cnt"`
		FollowersCnt         uint   `json:"followers_cnt"`
	}{
		ID:                   u.ID,
		Username:             u.Username,
		Role:                 u.Role,
		Email:                u.Email,
		Profile:              u.Profile,
		Avatar:               u.Avatar,
		FavouriteWorkflowCnt: u.FavouriteWorkflowCnt,
		PublishedWorkflowCnt: u.PublishedWorkflowCnt,
		FollowersCnt:         u.FollowersCnt,
	})
}

// UserUpdateProfileOption is option parameter for updating user profile
type UserUpdateProfileOption func(*User) *exceptions.Error

func WithUsername(username string) UserUpdateProfileOption {
	return func(u *User) *exceptions.Error {
		var user User
		if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.InternalServerError(err)
		} else if user != (User{}) {
			return exceptions.ErrUserNameAlreadyExists
		}
		u.Username = username
		return nil
	}
}

func WithEmail(email string) UserUpdateProfileOption {
	return func(u *User) *exceptions.Error {
		var user User
		if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.InternalServerError(err)
		} else if user != (User{}) {
			return exceptions.ErrUserEmailAlreadyExists
		}
		u.Email = email
		return nil
	}
}

func WithPassword(oldPassword, newPassword string) UserUpdateProfileOption {
	return func(u *User) *exceptions.Error {
		res := util.CheckPassword(u.Password, oldPassword)
		if !res {
			return exceptions.ErrPasswordIncorrect
		}
		if oldPassword == newPassword {
			return exceptions.ErrPasswordSame
		}
		hashedPassword, err := util.HashPassword(newPassword)
		if err != nil {
			return exceptions.InternalServerError(err)
		}
		u.Password = hashedPassword
		return nil
	}
}

func WithProfile(profile string) UserUpdateProfileOption {
	return func(u *User) *exceptions.Error {
		u.Profile = profile
		return nil
	}
}

func init() {
	db.DB.AutoMigrate(User{})
}
