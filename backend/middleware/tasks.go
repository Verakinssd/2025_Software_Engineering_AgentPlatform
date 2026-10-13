package middleware

import (
	"backend/repository"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"time"
)

const RedisAddr = "127.0.0.1:6379"

const (
	TypePenaltyExpiration       = "penalty:expiration"
	TypeRecordsCasCadeDeletion  = "records:cascade_deletion"
	TypeCommentsCasCadeDeletion = "comments:cascade_deletion"
)

var (
	RedisClient = asynq.RedisClientOpt{Addr: RedisAddr}
	AsynqClient = asynq.NewClient(RedisClient)
)

var taskHandlers = map[string]func(context.Context, *asynq.Task) error{
	TypePenaltyExpiration:       HandlePenaltyExpirationTask,
	TypeRecordsCasCadeDeletion:  HandleRecordsCasCadeDeletionTask,
	TypeCommentsCasCadeDeletion: HandleCommentsCasCadeDeletionTask,
}

func RegisterTaskHandlers(mux *asynq.ServeMux) {
	for taskType, handler := range taskHandlers {
		mux.HandleFunc(taskType, handler)
	}
}

type PenaltyExpirationPayload struct {
	UserId uint
}

type RecordsCasCadeDeletionPayload struct {
	WorkflowId uint
}

type CommentsCasCadeDeletionPayload struct {
	WorkflowId uint
}

//---------------------------------------------------------------
// 						  Tasks Creators.
//---------------------------------------------------------------

func NewPenaltyExpirationTask(userId uint, duration time.Duration) (*asynq.Task, error) {
	payload, err := json.Marshal(PenaltyExpirationPayload{UserId: userId})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypePenaltyExpiration, payload, asynq.MaxRetry(3), asynq.ProcessIn(duration)), nil
}

func NewRecordsCasCadeDeletionTask(workflowId uint) (*asynq.Task, error) {
	payload, err := json.Marshal(RecordsCasCadeDeletionPayload{WorkflowId: workflowId})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeRecordsCasCadeDeletion, payload), nil
}

func NewCommentsCasCadeDeletionTask(workflowId uint) (*asynq.Task, error) {
	payload, err := json.Marshal(CommentsCasCadeDeletionPayload{WorkflowId: workflowId})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeCommentsCasCadeDeletion, payload), nil
}

//---------------------------------------------------------------
// 						  Tasks Handlers.
//---------------------------------------------------------------

func HandlePenaltyExpirationTask(ctx context.Context, t *asynq.Task) error {
	var payload PenaltyExpirationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	// Perform penalty expiration task
	if err := repository.DeletePenaltyRecord(payload.UserId); err != nil {
		return fmt.Errorf("repository.DeletePenaltyRecord failed: %v: %w", err, asynq.SkipRetry)
	}
	return nil
}

func HandleRecordsCasCadeDeletionTask(ctx context.Context, t *asynq.Task) error {
	var payload RecordsCasCadeDeletionPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	if err := repository.DeleteCorrespondingRecordsAssociatedWithWorkflow(payload.WorkflowId); err != nil {
		return fmt.Errorf("repository.DeleteCorrespondingRecordsAssociatedWithWorkflow failed: %v: %w", err, asynq.SkipRetry)
	}
	return nil
}

func HandleCommentsCasCadeDeletionTask(ctx context.Context, t *asynq.Task) error {
	var payload CommentsCasCadeDeletionPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	if err := repository.DeleteCorrespondingCommentsAssociatedWithWorkflow(payload.WorkflowId); err != nil {
		return fmt.Errorf("repository.DeleteCorrespondingCommentsAssociatedWithWorkflow failed: %v: %w", err, asynq.SkipRetry)
	}
	return nil
}
