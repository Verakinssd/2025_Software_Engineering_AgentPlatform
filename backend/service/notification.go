package service

import (
	"backend/exceptions"
	"backend/model"
	"backend/repository"
)

func BroadcastNotification(userId uint, title string, content string) *exceptions.Error {
	return repository.CreateBroadcastNotification(userId, title, content)
}

func SendNotification(fromUserId uint, toUserId uint, title string, content string) *exceptions.Error {
	return repository.CreateNotification(fromUserId, toUserId, title, content)
}

func GetUsersNotifications(userId uint, page uint, pageSize uint) ([]*model.Notification, *exceptions.Error) {
	return repository.GetUsersNotifications(userId, page, pageSize)
}

func ReadNotification(userId uint, notificationId uint) *exceptions.Error {
	notification, e := repository.GetNotificationById(notificationId)
	if e != nil {
		return e
	}
	if notification.ToUserId != userId {
		return exceptions.ErrPermissionDenied.WithDetail("you are not allowed to read this notification")
	}
	notification.IsRead = true
	db.Save(notification)
	return nil
}

func ReadAllNotifications(userId uint) *exceptions.Error {
	return repository.ReadAllNotifications(userId)
}
