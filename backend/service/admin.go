package service

import (
	"backend/config"
	"backend/exceptions"
	"backend/middleware"
	"backend/model"
	"backend/repository"
	"errors"
	"time"
)

func CreatePenaltyFor(users []uint, adminId uint, penaltyType, reason string, duration time.Duration) ([]*model.UserPenaltyRecord, *exceptions.Error) {
	var records []*model.UserPenaltyRecord
	if penaltyType != "warning" && penaltyType != "restrict_comment" && penaltyType != "restrict_agent_publish" && penaltyType != "temporary_ban" && penaltyType != "permanent_ban" {
		return records, exceptions.InvalidRequest(
			errors.New(
				"invalid penalty type, should be one of " +
					"['warning', 'restrict_comment', 'restrict_agent_publish', 'temporary_ban' or 'permanent_ban']",
			),
		)
	}
	for _, userId := range users {
		var newDuration = duration
		if penaltyType == "permanent_ban" {
			newDuration = 9999 * 24 * time.Hour
		}

		res, err := repository.CreatePenaltyForUser(adminId, userId, penaltyType, reason, newDuration)
		if err != nil {
			return records, err
		}

		// create a task to expire it
		if config.Config.Server.Environment == "server" {
			task, err2 := middleware.NewPenaltyExpirationTask(userId, newDuration)
			if err2 != nil {
				return records, exceptions.InternalServerError(err2)
			}
			middleware.AsynqClient.Enqueue(task)
		}

		records = append(records, res)
	}
	return records, nil
}

func UpdatePenaltyFor(userId, adminId uint, penaltyType, reason string, duration time.Duration) (*model.UserPenaltyRecord, *exceptions.Error) {
	if config.Config.Server.Environment == "server" {
		task, err := middleware.NewPenaltyExpirationTask(userId, duration)
		if err != nil {
			return nil, exceptions.InternalServerError(err)
		}
		middleware.AsynqClient.Enqueue(task)
	}
	return repository.UpdatePenaltyForUser(adminId, userId, penaltyType, reason, duration)
}

func GetPenaltyFor(userId uint) (*model.UserPenaltyRecord, *exceptions.Error) {
	return repository.GetPenaltyByUserId(userId)
}

func GetAllUser() ([]*model.User, *exceptions.Error) {
	return repository.GetAllUser()
}
