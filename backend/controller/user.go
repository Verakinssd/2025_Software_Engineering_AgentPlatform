package controller

import (
	"backend/config"
	_db "backend/db"
	"backend/exceptions"
	"backend/middleware"
	"backend/model"
	"backend/repository"
	"backend/service"
	"errors"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
)

var db = _db.DB

func UserRegister(c *gin.Context) {
	type RegisterReq struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
		Email    string `form:"email"    binding:"required"`
	}
	var req RegisterReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	user, err2 := service.UserRegister(req.Username, req.Password, req.Email, "user")
	if err2 != nil {
		c.JSON(http.StatusBadRequest, err2)
		return
	}

	c.JSON(http.StatusOK, user)
}

func UserLogin(c *gin.Context) {
	type LoginReq struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}
	var req LoginReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	user, err2 := service.UserLogin(req.Username, req.Password)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, err2)
		return
	}
	// Check if user is banned for login
	if err3 := middleware.AuthenticateLogin(user.ID); err3 != nil {
		c.JSON(http.StatusForbidden, err3)
		return
	}

	token, _ := middleware.GenerateToken(user.ID, user.Role)
	c.SetCookie(config.Config.Jwt.Token.Name, token, int(config.Config.Jwt.Token.Expiration), "/", "", false, true)
	c.JSON(http.StatusOK, user)
}

func UserProfile(c *gin.Context) {
	userId, _ := c.Get("userId")
	user, err := service.UserProfile(userId.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func UserUpdateProfile(c *gin.Context) {
	type UpdateProfileReq struct {
		NewUsername string `form:"new_username"`
		OldPassword string `form:"old_password"`
		NewPassword string `form:"new_password"`
		NewEmail    string `form:"new_email"`
		Profile     string `form:"profile"`
	}
	var req UpdateProfileReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	var opts []model.UserUpdateProfileOption
	if _, exists := c.Request.Form["new_username"]; exists {
		opts = append(opts, model.WithUsername(req.NewUsername))
	}
	if _, exists := c.Request.Form["old_password"]; exists {
		// If old_password is provided, new_password must be provided as well
		if _, newExist := c.Request.Form["new_password"]; !newExist {
			c.JSON(
				http.StatusBadRequest,
				exceptions.InvalidRequest(errors.New("old_password provided, but new_password is missing")),
			)
			return
		}
		opts = append(opts, model.WithPassword(req.OldPassword, req.NewPassword))
	}
	if _, exists := c.Request.Form["new_email"]; exists {
		opts = append(opts, model.WithEmail(req.NewEmail))
	}
	if _, exists := c.Request.Form["profile"]; exists {
		opts = append(opts, model.WithProfile(req.Profile))
	}
	user, err2 := service.UserUpdateProfile(c.MustGet("userId").(uint), opts...)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, err2)
		return
	}
	c.JSON(http.StatusOK, user)
}

func UserLogout(c *gin.Context) {
	c.SetCookie(config.Config.Jwt.Token.Name, "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successfully"})
}

func UserUpdateAvatar(c *gin.Context) {
	type AvatarReq struct {
		Avatar *multipart.FileHeader `form:"avatar" binding:"required"`
	}
	var req AvatarReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	file, err := req.Avatar.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}
	defer file.Close()

	userId, _ := c.Get("userId")
	user, e := repository.GetUserById(userId.(uint))
	if e != nil {
		c.JSON(http.StatusBadRequest, e)
		return
	}

	fileName, e := saveFile(c, "avatar", &[]string{".jpg", ".png", ".gif", ".webp", ".jpeg"})
	if e != nil {
		c.JSON(http.StatusBadRequest, e)
		return
	}

	user.Avatar = "/api/image/" + fileName
	user, e = repository.SaveUser(user)
	if e != nil {
		c.JSON(http.StatusBadRequest, e)
		return
	}

	c.JSON(http.StatusOK, user)
}
