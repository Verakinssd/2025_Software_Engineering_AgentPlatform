package service

import (
	_db "backend/db"
	"backend/exceptions"
	"backend/model"
	"backend/repository"
	"backend/util"
	"errors"
)

var db = _db.DB

func UserRegister(username, password, email, role string) (*model.User, *exceptions.Error) {
	user, err := repository.GetUserBy("username", username)
	user2, _ := repository.GetUserBy("email", email)
	if err != nil && !errors.Is(err, exceptions.ErrUserNotFound) {
		return nil, err
	}
	if user != nil {
		return nil, exceptions.ErrUserNameAlreadyExists
	}
	if user2 != nil {
		return nil, exceptions.ErrUserEmailAlreadyExists
	}
	return repository.CreateUser(username, password, email, role)
}

func UserLogin(username string, password string) (*model.User, *exceptions.Error) {
	user, err := repository.GetUserBy("username", username)
	if err != nil {
		return nil, err
	}
	if util.CheckPassword(user.Password, password) {
		return user, nil
	} else {
		return nil, exceptions.ErrPasswordIncorrect
	}
}

func UserProfile(UserId uint) (*model.User, *exceptions.Error) {
	return repository.GetUserById(UserId)
}

func UserUpdateProfile(UserId uint, opts ...model.UserUpdateProfileOption) (*model.User, *exceptions.Error) {
	user, err := repository.GetUserById(UserId)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		// opt are callbacks
		if err = opt(user); err != nil {
			return nil, err
		}
	}
	user, err2 := repository.SaveUser(user)
	if err2 != nil {
		return nil, err2
	}
	return user, nil
}
