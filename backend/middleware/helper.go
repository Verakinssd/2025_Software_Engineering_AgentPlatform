package middleware

import (
	"backend/exceptions"
	"backend/repository"
	"errors"
	"time"
)

type Operation string

const (
	Comment         Operation = "Comment"
	PublishWorkflow Operation = "PublishWorkflow"
	// ... other operations
)

var operationMap = map[Operation]func(userId uint) *exceptions.Error{
	Comment:         authenticateComment,
	PublishWorkflow: authenticatePublishWorkflow,
	// ... other operations
}

func AuthenticateLogin(userId uint) *exceptions.Error {
	record, err := repository.GetPenaltyByUserId(userId)
	if err != nil {
		if errors.Is(err, exceptions.ErrPenaltyNotFound) {
			// user is not banned for Login
			return nil
		}
		return err
	}
	if record.Type == "temporary_ban" || record.Type == "permanent_ban" && record.ExpireAt.After(time.Now()) {
		return exceptions.ErrForbidden
	}
	return nil
}

func authenticateComment(userId uint) *exceptions.Error {
	record, err := repository.GetPenaltyByUserId(userId)
	if err != nil {
		if errors.Is(err, exceptions.ErrPenaltyNotFound) {
			// user is not banned for Comment
			return nil
		}
		return err
	}
	if record.Type == "restrict_comment" && record.ExpireAt.After(time.Now()) {
		return exceptions.ErrForbidden
	}
	return nil
}

func authenticatePublishWorkflow(userId uint) *exceptions.Error {
	record, err := repository.GetPenaltyByUserId(userId)
	if err != nil {
		if errors.Is(err, exceptions.ErrPenaltyNotFound) {
			// user is not banned for publishing workflow
			return nil
		}
		return err
	}
	if record.Type == "restrict_agent_publish" && record.ExpireAt.After(time.Now()) {
		return exceptions.ErrForbidden
	}
	return nil
}
