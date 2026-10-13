package repository

import (
	"backend/exceptions"
	"backend/model"
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

func GetPenaltyByUserId(userId uint) (*model.UserPenaltyRecord, *exceptions.Error) {
	var record model.UserPenaltyRecord
	err := db.First(&record, "user_id = ?", userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrPenaltyNotFound
		}
		return nil, exceptions.InternalServerError(err)
	}
	return &record, nil
}

func UpdatePenaltyForUser(adminId, userId uint, penaltyType, reason string, duration time.Duration) (*model.UserPenaltyRecord, *exceptions.Error) {
	var existingRecord model.UserPenaltyRecord
	err := db.First(&existingRecord, "user_id = ?", userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrPenaltyNotFound
		}
		return nil, exceptions.InternalServerError(err)
	}
	// update record by deleting old record and creating a new one
	db.Delete(&existingRecord)
	newRecord, err2 := CreatePenaltyForUser(adminId, userId, penaltyType, reason, duration)
	if err2 != nil {
		if errors.Is(err2, exceptions.ErrPermissionDenied) {
			return nil, exceptions.ErrPermissionDenied
		}
		return nil, err2
	}
	return newRecord, nil
}

func CreatePenaltyForUser(adminId, userId uint, penaltyType, reason string, duration time.Duration) (*model.UserPenaltyRecord, *exceptions.Error) {
	admin, err := GetUserById(adminId)
	if err != nil {
		if errors.Is(err, exceptions.ErrUserNotFound) {
			return nil, exceptions.ErrPermissionDenied
		} else {
			return nil, exceptions.InternalServerError(err)
		}
	}
	if admin.Role != "admin" {
		return nil, exceptions.ErrPermissionDenied
	}

	user, err := GetUserById(userId)
	if err != nil {
		if errors.Is(err, exceptions.ErrUserNotFound) {
			return nil, exceptions.ErrUserNotFound
		} else {
			return nil, exceptions.InternalServerError(err)
		}
	}

	var existingRecord model.UserPenaltyRecord
	err2 := db.First(&existingRecord, "user_id = ?", userId).Error
	if err2 != nil && !errors.Is(err2, gorm.ErrRecordNotFound) {
		return nil, exceptions.InternalServerError(err2)
	}

	if existingRecord.ID != 0 {
		db.Delete(&existingRecord)
	}

	newRecord := model.UserPenaltyRecord{
		AdminId:  adminId,
		Admin:    admin,
		UserId:   userId,
		User:     user,
		Type:     penaltyType,
		Reason:   reason,
		ExpireAt: time.Now().Add(duration),
	}
	err3 := db.Create(&newRecord).Error
	if err3 != nil {
		return nil, exceptions.InternalServerError(err3)
	}
	return &newRecord, nil
}

func DeletePenaltyRecord(userId uint) *exceptions.Error {
	err := db.Delete(&model.UserPenaltyRecord{}, "user_id = ?", userId).Error
	if err != nil {
		return exceptions.InternalServerError(err)
	}
	return nil
}

func GetAllUser() ([]*model.User, *exceptions.Error) {
	var users []*model.User
	err := db.Find(&users).Error
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	return users, nil
}
