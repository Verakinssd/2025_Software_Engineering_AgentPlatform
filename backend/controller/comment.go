package controller

import (
	"backend/exceptions"
	"backend/model"
	"backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllComments(c *gin.Context) {
	// order:
	// 		 1 => order by time DESC,
	// 		 2 => order by hot DESC
	// notice that sub-comments are always ordering by time ASC (like bilibili)
	orderStr := c.Query("order")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")

	order, err := strconv.Atoi(orderStr)
	if err != nil || (order != 1 && order != 2) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order, should be either 1 or 2"})
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pageSize"})
		return
	}

	workflowId, err := strconv.Atoi(c.Param("workflow_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}
	uid, exist := c.Get("userId")
	var userId uint
	if !exist {
		userId = uint(0)
	} else {
		userId = uid.(uint)
	}

	var comments []*model.WorkflowCommentPublic
	var err2 *exceptions.Error

	if order == 1 {
		// order by time DESC
		comments, err2 = service.GetAllCommentsOrderByTime(uint(workflowId), userId, uint(page), uint(pageSize))
	} else {
		// order by hot DESC
		comments, err2 = service.GetAllCommentsOrderByHot(uint(workflowId), userId, uint(page), uint(pageSize))
	}

	if err2 != nil {
		c.JSON(http.StatusBadRequest, err2)
		return
	}

	c.JSON(http.StatusOK, comments)
}

func CreateComment(c *gin.Context) {
	type CreateCommentRequest struct {
		// CommentType: 1 => workflow comment, 2 => secondary-comment
		CommentType int    `form:"comment_type" binding:"required"`
		Content     string `form:"content" binding:"required"`
		WorkflowId  *uint  `form:"workflow_id"`
		ParentId    *uint  `form:"parent_id"`
		ReplyTo     *uint  `form:"reply_to"`
	}

	var req CreateCommentRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if comment_type is valid
	if req.CommentType != 1 && req.CommentType != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment_type"})
		return
	}

	if req.CommentType == 1 && (req.WorkflowId == nil || req.ParentId != nil || req.ReplyTo != nil) {
		// Workflow comment must have workflow_id but not parent_id and reply_to
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment_type or missing required fields"})
		return
	} else if req.CommentType == 2 && (req.WorkflowId != nil || req.ParentId == nil || req.ReplyTo == nil) {
		// Secondary comment must have parent_id and reply_to but not workflow_id
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment_type or missing required fields"})
		return
	}

	var comment *model.WorkflowCommentPublic
	var err2 *exceptions.Error

	userId := c.MustGet("userId").(uint)

	if req.CommentType == 1 {
		// Workflow comment
		comment, err2 = service.CreateWorkflowComment(*req.WorkflowId, userId, req.Content)
	} else {
		// Secondary comment
		comment, err2 = service.CreateSecondaryComment(*req.ParentId, userId, *req.ReplyTo, req.Content)
	}

	if err2 != nil {
		c.JSON(http.StatusBadRequest, err2)
		return
	}

	c.JSON(http.StatusOK, comment)
}

func LikeComment(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	commentId, err := strconv.Atoi(c.Param("comment_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid commentId"})
		return
	}

	res, err2 := service.LikeComment(uint(commentId), userId)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, err2)
		return
	}

	// `res` represents like result:
	// true  => like success
	// false => like cancel
	c.JSON(http.StatusOK, gin.H{"like_result": res})
}

func DeleteComment(c *gin.Context) {
	commentId, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment_id"})
		return
	}

	userId := c.MustGet("userId").(uint)

	err2 := service.DeleteComment(uint(commentId), userId)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, err2)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
}
