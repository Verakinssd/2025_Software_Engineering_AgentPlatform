package controller

import (
	"backend/exceptions"
	"backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func BroadcastNotification(c *gin.Context) {
	type Req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	var req Req
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	userId, _ := c.Get("userId")
	e := service.BroadcastNotification(userId.(uint), req.Title, req.Content)
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification sent successfully"})
}

func SendNotification(c *gin.Context) {
	type Req struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content" binding:"required"`
		ToUserId uint   `json:"to_user_id" binding:"required"`
	}
	var req Req
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	userId, _ := c.Get("userId")
	e := service.SendNotification(userId.(uint), req.ToUserId, req.Title, req.Content)
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification sent successfully"})
}

func GetNotifications(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("per_page", "10")
	page, e1 := strconv.Atoi(pageStr)
	pageSize, e2 := strconv.Atoi(pageSizeStr)
	if e1 != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(e1))
		return
	}
	if e2 != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(e2))
		return
	}

	userId, _ := c.Get("userId")
	notifications, e := service.GetUsersNotifications(userId.(uint), uint(page), uint(pageSize))
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}
	type ApiNotification struct {
		Id           uint      `json:"id"`
		Title        string    `json:"title"`
		Content      string    `json:"content"`
		Date         time.Time `json:"date"`
		IsRead       bool      `json:"is_read"`
		FromUserName string    `json:"from_user_name"`
	}
	var apiNotifications []ApiNotification
	for _, notification := range notifications {
		apiNotifications = append(apiNotifications, ApiNotification{
			Id:           notification.ID,
			Title:        notification.Title,
			Content:      notification.Content,
			Date:         notification.CreatedAt,
			IsRead:       notification.IsRead,
			FromUserName: notification.FromUser.Username,
		})
	}
	type Res struct {
		Notifications []ApiNotification `json:"notifications"`
	}
	c.JSON(http.StatusOK, Res{Notifications: apiNotifications})
}

func ReadNotification(c *gin.Context) {
	idStr := c.Param("notification_id")
	id, e := strconv.Atoi(idStr)
	if e != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(e))
		return
	}

	userId, _ := c.Get("userId")
	err := service.ReadNotification(userId.(uint), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification read successfully"})
}

func ReadAllNotifications(c *gin.Context) {
	userId, _ := c.Get("userId")
	err := service.ReadAllNotifications(userId.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All notifications read successfully"})
}
