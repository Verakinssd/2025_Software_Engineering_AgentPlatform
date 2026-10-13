package controller

import (
	"backend/config"
	"backend/exceptions"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"
)

func contained(ext string, allowedExt *[]string) bool {
	if len(*allowedExt) == 0 {
		return true
	}
	for _, e := range *allowedExt {
		if e == ext {
			return true
		}
	}
	return false
}

func saveFile(c *gin.Context, field string, allowedExt *[]string) (string, *exceptions.Error) {
	file, err := c.FormFile(field)
	if err != nil {
		return "", exceptions.InternalServerError(err)
	}
	ext := filepath.Ext(file.Filename)
	if !contained(ext, allowedExt) {
		return "", exceptions.ErrInvalidFileTypes
	}
	fileName := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), uuid.New().String(), ext)
	savePath := filepath.Join(config.Config.Server.Upload, fileName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return "", exceptions.InternalServerError(err)
	}
	return fileName, nil
}

func GetImage(c *gin.Context) {
	fileName := c.Param("filename")
	filePath := filepath.Join(config.Config.Server.Upload, fileName)
	c.File(filePath)
}

func UploadImage(c *gin.Context) {
	type Req struct {
		Image *multipart.FileHeader `form:"image" binding:"required"`
	}
	var req Req
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	file, err := req.Image.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}
	defer file.Close()

	fileName, e := saveFile(c, "image", &[]string{".jpg", ".png", ".gif", ".webp", ".jpeg"})
	if e != nil {
		c.JSON(http.StatusBadRequest, e)
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": "/api/image/" + fileName})
}
