package repository

import (
	"backend/exceptions"
	"backend/model"
	"errors"

	"github.com/jinzhu/gorm"
)

// 移除硬编码的回复数量限制
// const SecondaryCommentLimit = 3

// ===================================
//          private methods
// ===================================

func checkLikeStatus(commentId, userId uint) bool {
	var record model.WorkflowCommentLikeRecord
	err := db.Where("comment_id = ? AND user_id = ?", commentId, userId).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		panic(err) // FIXME
	}
	return true
}

func fetchCommentLikesCount(commentIds []uint) (map[uint]uint, error) {
	var records []struct {
		CommentID uint
		Likes     uint
	}
	err := db.Model(&model.WorkflowCommentLikeRecord{}).
		Select("comment_id, COUNT(*) as likes").
		Where("comment_id IN (?)", commentIds).
		Group("comment_id").
		Scan(&records).Error
	if err != nil {
		return nil, err
	}

	idMap := make(map[uint]uint)
	for _, record := range records {
		idMap[record.CommentID] = record.Likes
	}
	return idMap, nil
}

// ===================================
//           public methods
// ===================================

// GetCommentsForWorkflowOrderBy returns all comments for a workflow, ordered by time or likes
func GetCommentsForWorkflowOrderBy(order, workflowId, page, pageSize uint) ([]*model.WorkflowComment, *exceptions.Error) {
	var comments []*model.WorkflowComment
	// order == 1: by time desc
	// order == 2: by likes desc
	if order != 1 && order != 2 {
		panic("invalid order")
	}
	var err error
	// find all comments for workflow
	if order == 1 {
		err = db.Where("workflow_id = ?", workflowId).
			Order("created_at DESC").
			Offset((page - 1) * pageSize).
			Limit(pageSize).
			Preload("Commenter").
			Find(&comments).Error
	} else {
		err = db.Where("workflow_id = ?", workflowId).
			Order("likes DESC").
			Offset((page - 1) * pageSize).
			Limit(pageSize).
			Preload("Commenter").
			Find(&comments).Error
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrCommentNotFound
		}
		return nil, exceptions.InternalServerError(err)
	}

	var commentIds []uint
	for _, comment := range comments {
		commentIds = append(commentIds, comment.ID)
	}

	// fetch comment likes' count
	newMap, err := fetchCommentLikesCount(commentIds)

	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	// find all secondary comments for each comment
	// all secondary comments are sorted by time asc (just like bilibili)
	for _, comment := range comments {
		var secondary []*model.WorkflowComment
		err2 := db.Where("parent_id = ?", comment.ID).
			Order("created_at").
			// 移除Limit(SecondaryCommentLimit)，获取所有回复
			Preload("Commenter").
			Preload("ReplyToUser").
			Find(&secondary).Error
		if err2 != nil && !errors.Is(err2, gorm.ErrRecordNotFound) {
			return nil, exceptions.InternalServerError(err2)
		}

		var subCommentIds []uint
		for _, subComment := range secondary {
			subCommentIds = append(subCommentIds, subComment.ID)
		}
		// fetch sub-comment likes' count
		subMap, err3 := fetchCommentLikesCount(subCommentIds)
		if err3 != nil {
			return nil, exceptions.InternalServerError(err3)
		}

		// update sub-comment likes' count
		for _, subComment := range secondary {
			subComment.Likes = subMap[subComment.ID]
		}

		comment.Likes = newMap[comment.ID]
		comment.Children = secondary
	}
	return comments, nil
}

// CreateWorkflowComment replies to a workflow directly
func CreateWorkflowComment(workflowId, userId uint, content string) (*model.WorkflowComment, *exceptions.Error) {
	comment := model.WorkflowComment{
		CommentBase: model.CommentBase{
			CommenterID: userId,
			Content:     content,
		},
		WorkflowID: &workflowId,
		ParentID:   nil,
		ReplyTo:    nil,
	}
	err := db.Create(&comment).Error
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	return &comment, nil
}

// CreateSecondaryComment replies to a secondary comment or a workflow comment
func CreateSecondaryComment(parentId, userId, replyTo uint, content string) (*model.WorkflowComment, *exceptions.Error) {
	comment := model.WorkflowComment{
		CommentBase: model.CommentBase{
			CommenterID: userId,
			Content:     content,
		},
		WorkflowID: nil,
		ParentID:   &parentId,
		ReplyTo:    &replyTo,
	}
	err := db.Create(&comment).Error
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	return &comment, nil
}

func DeleteComment(commentId, userId uint) *exceptions.Error {
	var comment model.WorkflowComment
	err := db.Where("id = ?", commentId).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.ErrCommentNotFound
		}
		return exceptions.InternalServerError(err)
	}

	if comment.CommenterID != userId {
		return exceptions.ErrForbidden
	}

	if comment.WorkflowID != nil {
		return deleteWorkflowComment(comment)
	}

	if comment.ParentID != nil {
		return deleteSecondaryComment(comment)
	}

	panic("unknown comment type, exiting...")
}

func deleteWorkflowComment(comment model.WorkflowComment) *exceptions.Error {
	// find all secondary comments for the comment
	var secondary []*model.WorkflowComment
	err2 := db.Where("parent_id = ?", comment.ID).Find(&secondary).Error
	if err2 != nil && !errors.Is(err2, gorm.ErrRecordNotFound) {
		return exceptions.InternalServerError(err2)
	}
	// delete all secondary comments
	for _, child := range secondary {
		err3 := db.Delete(child).Error
		if err3 != nil {
			return exceptions.InternalServerError(err3)
		}
	}
	// delete the workflow comment
	err4 := db.Delete(comment).Error
	if err4 != nil {
		return exceptions.InternalServerError(err4)
	}
	return nil
}

func deleteSecondaryComment(comment model.WorkflowComment) *exceptions.Error {
	// delete the secondary comment
	err2 := db.Delete(comment).Error
	if err2 != nil {
		return exceptions.InternalServerError(err2)
	}
	return nil
}

func DeleteCorrespondingCommentsAssociatedWithWorkflow(workflowId uint) *exceptions.Error {
	var comments []*model.WorkflowComment
	err := db.Where("workflow_id = ?", workflowId).Find(&comments).Error
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	var commentIds []uint
	for _, comment := range comments {
		commentIds = append(commentIds, comment.ID)
	}

	// delete all secondary comments
	err2 := db.Where("parent_id IN (?)", commentIds).Delete(&model.WorkflowComment{}).Error
	if err2 != nil {
		return exceptions.InternalServerError(err2)
	}

	// delete all workflow comments
	err3 := db.Where("id IN (?)", commentIds).Delete(&model.WorkflowComment{}).Error
	if err3 != nil {
		return exceptions.InternalServerError(err3)
	}

	return nil
}

// GetLikeRecordsForUser returns all like records for a user on a list of comments
func GetLikeRecordsForUser(userId uint, commentIds []uint) ([]*model.WorkflowCommentLikeRecord, *exceptions.Error) {
	var records []*model.WorkflowCommentLikeRecord
	if err := db.Where("user_id = ? AND comment_id IN (?)", userId, commentIds).
		Find(&records).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, exceptions.InternalServerError(err)
	}
	return records, nil
}

func LikeWorkflowComment(commentId, userId uint) (bool, *exceptions.Error) {
	// first return value represents whether the like is successful or not
	// false -> already liked, so canceled / just failed
	// true  -> like successful
	if checkLikeStatus(commentId, userId) {
		db.Where("user_id = ? AND comment_id = ?", userId, commentId).Delete(&model.WorkflowCommentLikeRecord{})
		return false, nil
	}
	record := model.WorkflowCommentLikeRecord{
		LikeRecordBase: model.LikeRecordBase{
			CommentID: commentId,
			UserID:    userId,
		},
	}
	err := db.Create(&record).Error
	if err != nil {
		return false, exceptions.InternalServerError(err)
	}
	return true, nil
}
