package repository

import (
	_db "backend/db"
	"backend/exceptions"
	"backend/model"
	"backend/util"
	"errors"

	"github.com/jinzhu/gorm"
)

var db = _db.DB

func prepareSocialDataForUserProfile(user *model.User) *exceptions.Error {
	favouriteWorkflowCount, err := getUserFavouriteWorkflowCount(user.ID)
	if err != nil {
		return err
	}
	user.FavouriteWorkflowCnt = favouriteWorkflowCount

	publishedWorkflowCount, err2 := getUserPublishedWorkflowCount(user.ID)
	if err2 != nil {
		return err2
	}
	user.PublishedWorkflowCnt = publishedWorkflowCount

	followersCount, err3 := getUserFollowersCount(user.ID)
	if err3 != nil {
		return err3
	}
	user.FollowersCnt = followersCount

	// other data...
	return nil
}

func GetUserBy(field, value string) (*model.User, *exceptions.Error) {
	var user model.User
	err := db.Where(field+" = ?", value).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrUserNotFound
		} else {
			return nil, exceptions.InternalServerError(err)
		}
	} else {
		err2 := prepareSocialDataForUserProfile(&user)
		if err2 != nil {
			return nil, err2
		}

		return &user, nil
	}
}

func GetUserById(id uint) (*model.User, *exceptions.Error) {
	var user model.User
	err := db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrUserNotFound
		} else {
			return nil, exceptions.InternalServerError(err)
		}
	} else {
		err2 := prepareSocialDataForUserProfile(&user)
		if err2 != nil {
			return nil, err2
		}

		return &user, nil
	}
}

func CreateUser(username, password, email, role string) (*model.User, *exceptions.Error) {
	var user model.User
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return &user, exceptions.InternalServerError(err)
	}
	user = model.User{
		Username: username,
		Password: hashedPassword,
		Role:     role,
		Email:    email,
	}
	err = db.Create(&user).Error
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	} else {
		return &user, nil
	}
}

func GetAllUsers() ([]*model.User, *exceptions.Error) {
	var users []*model.User
	err := db.Find(&users).Error
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	} else {
		return users, nil
	}
}

func SaveUser(user *model.User) (*model.User, *exceptions.Error) {
	err := db.Save(user).Error
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	} else {
		return user, nil
	}
}
