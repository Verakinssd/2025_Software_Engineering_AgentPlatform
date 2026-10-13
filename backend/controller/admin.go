package controller

import (
	"backend/exceptions"
	"backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func GetUserPenaltyFor(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("user_id"))
	res, err := service.GetPenaltyFor(uint(userId))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func HandleUserPenalty(c *gin.Context) {
	type UserPenalty struct {
		UsersId     []uint `json:"users_id"`
		PenaltyType string `json:"penalty_type"`
		Duration    string `json:"duration"`
		Reason      string `json:"reason"`
	}
	var userPenalty UserPenalty
	if err := c.ShouldBindJSON(&userPenalty); err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}
	adminId := c.MustGet("userId").(uint)

	duration, _ := strconv.Atoi(userPenalty.Duration)
	res, err2 := service.CreatePenaltyFor(userPenalty.UsersId, adminId, userPenalty.PenaltyType, userPenalty.Reason, time.Duration(duration)*24*time.Hour)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, err2)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Penalty created successfully", "data": res})
}

func UpdateUserPenalty(c *gin.Context) {
	type UserPenalty struct {
		UserId         uint   `json:"user_id" binding:"required"`
		PenaltyId      uint   `json:"penalty_id" binding:"required"`
		NewPenaltyType string `json:"new_penalty_type" binding:"required"`
		NewDuration    string `json:"new_duration" binding:"required"`
		NewReason      string `json:"new_reason"`
	}
	var userPenalty UserPenalty
	if err := c.ShouldBindJSON(&userPenalty); err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}
	adminId := c.MustGet("userId").(uint)

	duration, _ := strconv.Atoi(userPenalty.NewDuration)
	res, err2 := service.UpdatePenaltyFor(userPenalty.UserId, adminId, userPenalty.NewPenaltyType, userPenalty.NewReason, time.Duration(duration)*24*time.Hour)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, err2)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Penalty updated successfully", "data": res})
}

func GetAllUser(c *gin.Context) {
	res, err := service.GetAllUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, res)
}
