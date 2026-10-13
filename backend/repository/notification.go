package repository

import (
	"backend/exceptions"
	"backend/model"
)

func CreateBroadcastNotification(userId uint, title string, content string) *exceptions.Error {
	fromUser, e := GetUserById(userId)
	if e != nil {
		return e
	}
	allUsers, e2 := GetAllUsers()
	if e2 != nil {
		return e2
	}
	transaction := db.Begin()
	if transaction.Error != nil {
		return exceptions.InternalServerError(transaction.Error)
	}
	for _, user := range allUsers {
		n := &model.Notification{
			Title:    title,
			Content:  content,
			Type:     "broadcast",
			FromUser: fromUser,
			ToUser:   user,
		}
		err := db.Create(n).Error
		if err != nil {
			transaction.Rollback()
			return exceptions.InternalServerError(err)
		}
	}
	transaction.Commit()
	return nil
}

func CreateNotification(fromUserId uint, toUserId uint, title string, content string) *exceptions.Error {
	fromUser, e := GetUserById(fromUserId)
	if e != nil {
		return e
	}
	toUser, e2 := GetUserById(toUserId)
	if e2 != nil {
		return e2
	}
	n := &model.Notification{
		Title:    title,
		Content:  content,
		Type:     "personal",
		FromUser: fromUser,
		ToUser:   toUser,
	}
	err := db.Create(n).Error
	if err != nil {
		return exceptions.InternalServerError(err)
	}
	return nil
}

func GetUsersNotifications(userId uint, page uint, pageSize uint) ([]*model.Notification, *exceptions.Error) {
	var notifications []*model.Notification
	err := db.Preload("FromUser").
		Where("to_user_id = ?", userId).
		Order("created_at desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&notifications).Error
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	return notifications, nil
}

func GetNotificationById(id uint) (*model.Notification, *exceptions.Error) {
	var notification model.Notification
	err := db.Preload("FromUser").
		Preload("ToUser").
		First(&notification, id).Error
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	return &notification, nil
}

func ReadAllNotifications(userId uint) *exceptions.Error {
	err := db.Model(&model.Notification{}).
		Where("to_user_id = ?", userId).
		Update("is_read", true).Error
	if err != nil {
		return exceptions.InternalServerError(err)
	}
	return nil
}
