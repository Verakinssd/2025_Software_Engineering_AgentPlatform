package service

import (
	"backend/exceptions"
	"backend/model"
	"backend/repository"
)

// ========= private methods =========
func getCommentLikes(userId uint, commentIds []uint) (map[uint]bool, *exceptions.Error) {
	records, err := repository.GetLikeRecordsForUser(userId, commentIds)
	if err != nil {
		return nil, err
	}
	newMap := make(map[uint]bool)
	for _, record := range records {
		newMap[record.CommentID] = true
	}
	return newMap, nil
}

// ========= public methods =========

// GetAllCommentsOrderByTime returns all comments for a workflow ordered by time
func GetAllCommentsOrderByTime(workflowId, userId, page, pageSize uint) ([]*model.WorkflowCommentPublic, *exceptions.Error) {
	comments, err := repository.GetCommentsForWorkflowOrderBy(1, workflowId, page, pageSize)
	if err != nil {
		return nil, err
	}

	// In order to avoid N + 1 problem, we should fetch all likes for all comments at once
	commentIds := make([]uint, 0)
	for _, comment := range comments {
		commentIds = append(commentIds, comment.ID)
	}

	recordMap, err2 := getCommentLikes(userId, commentIds)
	if err2 != nil {
		return nil, err2
	}

	var result []*model.WorkflowCommentPublic
	for _, comment := range comments {
		result = append(result, &model.WorkflowCommentPublic{
			WorkflowComment: *comment,
			LikeStatus:      recordMap[comment.ID],
		})
	}
	return result, nil
}

// GetAllCommentsOrderByHot returns all comments for a workflow ordered by hot
func GetAllCommentsOrderByHot(workflowId, userId, page, pageSize uint) ([]*model.WorkflowCommentPublic, *exceptions.Error) {
	comments, err := repository.GetCommentsForWorkflowOrderBy(2, workflowId, page, pageSize)
	if err != nil {
		return nil, err
	}

	commentIds := make([]uint, 0)
	for _, comment := range comments {
		commentIds = append(commentIds, comment.ID)
	}

	recordMap, err2 := getCommentLikes(userId, commentIds)
	if err2 != nil {
		return nil, err2
	}

	var result []*model.WorkflowCommentPublic
	for _, comment := range comments {
		result = append(result, &model.WorkflowCommentPublic{
			WorkflowComment: *comment,
			LikeStatus:      recordMap[comment.ID],
		})
	}
	return result, nil
}

func CreateWorkflowComment(workflowId, userId uint, content string) (*model.WorkflowCommentPublic, *exceptions.Error) {
	comment, err := repository.CreateWorkflowComment(workflowId, userId, content)
	if err != nil {
		return nil, err
	}
	return &model.WorkflowCommentPublic{
		WorkflowComment: *comment,
		LikeStatus:      false,
	}, nil
}

func CreateSecondaryComment(parentId, userId, replyTo uint, content string) (*model.WorkflowCommentPublic, *exceptions.Error) {
	comment, err := repository.CreateSecondaryComment(parentId, userId, replyTo, content)
	if err != nil {
		return nil, err
	}
	return &model.WorkflowCommentPublic{
		WorkflowComment: *comment,
		LikeStatus:      false,
	}, nil
}

func DeleteComment(commentId, userId uint) *exceptions.Error {
	return repository.DeleteComment(commentId, userId)
}

func LikeComment(commentId, userId uint) (bool, *exceptions.Error) {
	return repository.LikeWorkflowComment(commentId, userId)
}
