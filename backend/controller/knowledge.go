package controller

import (
	"backend/config"
	"backend/exceptions"
	"backend/model"
	"backend/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadKnowledgeFile(c *gin.Context) {
	userId, _ := c.Get("userId")

	type Req struct {
		File        *multipart.FileHeader `form:"file" binding:"required"`
		Description string                `form:"description"`
	}
	var req Req
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	ext := filepath.Ext(req.File.Filename)
	// 支持.txt、.md、.docx、.pdf文件
	supportedExtensions := map[string]bool{
		".txt":  true,
		".md":   true,
		".docx": true,
		".pdf":  true,
		".zip":  true,
		".pptx": true,
		".xls":  true,
		".xlsx": true,
	}
	
	if !supportedExtensions[strings.ToLower(ext)] {
		c.JSON(http.StatusBadRequest, exceptions.ErrInvalidFileTypes.WithDetail("支持文件格式：txt, md, docx, pdf, zip, pptx, xls, xlsx"))
		return
	}

	// 保存上传的文件
	knowledgeDir := filepath.Join(config.Config.Server.Upload, "knowledge")
	if err := os.MkdirAll(knowledgeDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError(err))
		return
	}

	fileName := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), uuid.New().String(), ext)
	filePath := filepath.Join(knowledgeDir, fileName)
	if err := c.SaveUploadedFile(req.File, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError(err))
		return
	}

	// 解析不同类型文件的内容
	var content string
	fileType := strings.TrimPrefix(strings.ToLower(ext), ".")

	parsedContent, err := parseDocPdfFile(filePath, req.File.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError(fmt.Errorf("解析docx文件失败: %w", err)))
		return
	}
	content = parsedContent

	knowledgeFile := &model.KnowledgeFile{
		FileName:    req.File.Filename,
		FilePath:    fileName,
		FileSize:    req.File.Size,
		FileType:    fileType,
		UserId:      userId.(uint),
		Content:     content,
		Description: req.Description,
	}

	if err := repository.CreateKnowledgeFile(knowledgeFile); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_id":   knowledgeFile.ID,
		"file_name": knowledgeFile.FileName,
		"file_size": knowledgeFile.FileSize,
	})
}

func parseDocPdfFile(filePath, fileName string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", fmt.Errorf("创建表单字段失败: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("写入文件内容失败: %w", err)
	}
	
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("关闭writer失败: %w", err)
	}

	parseAPI := "http://114.116.204.215:5001/api/parse/"
	req, err := http.NewRequest("POST", parseAPI, body)
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	var result struct {
		Filename string `json:"filename"`
		Size     int    `json:"size"`
		Text     string `json:"text"`
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应内容失败: %w", err)
	}

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", fmt.Errorf("解析JSON响应失败: %w", err)
	}

	return result.Text, nil
}

func GetUserKnowledgeFiles(c *gin.Context) {
	userId, _ := c.Get("userId")

	files, err := repository.GetKnowledgeFilesByUserId(userId.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

func GetKnowledgeFile(c *gin.Context) {
	userId, _ := c.Get("userId")
	fileIdStr := c.Param("file_id")
	
	fileId, err := strconv.Atoi(fileIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}
	
	file, errRepo := repository.GetKnowledgeFileById(uint(fileId))
	
	if errRepo != nil {
		c.JSON(http.StatusInternalServerError, errRepo)
		return
	}
	
	if file == nil {
		c.JSON(http.StatusNotFound, exceptions.InternalServerError(fmt.Errorf("找不到ID为%d的知识库文件", fileId)))
		return
	}

	if file.UserId != userId.(uint) {
		c.JSON(http.StatusForbidden, exceptions.ErrPermissionDenied)
		return
	}


	c.JSON(http.StatusOK, file)
}

func DownloadKnowledgeFile(c *gin.Context) {
	userId, _ := c.Get("userId")
	fileIdStr := c.Param("file_id")
	
	fileId, err := strconv.Atoi(fileIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}
	
	file, errRepo := repository.GetKnowledgeFileById(uint(fileId))
	if errRepo != nil {
		c.JSON(http.StatusInternalServerError, errRepo)
		return
	}
	
	if file == nil {
		c.JSON(http.StatusNotFound, exceptions.InternalServerError(fmt.Errorf("找不到ID为%d的知识库文件", fileId)))
		return
	}

	if file.UserId != userId.(uint) {
		c.JSON(http.StatusForbidden, exceptions.ErrPermissionDenied)
		return
	}

	knowledgeDir := filepath.Join(config.Config.Server.Upload, "knowledge")
	filePath := filepath.Join(knowledgeDir, file.FilePath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, exceptions.InternalServerError(fmt.Errorf("文件不存在: %s", filePath)))
		return
	}

	fileName := file.FileName
	fileNameEscaped := url.QueryEscape(fileName)
	fileNameUTF8 := url.PathEscape(fileName)
	disposition := fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, fileNameEscaped, fileNameUTF8)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", disposition)
	c.Header("Content-Type", getContentType(file.FileType))
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	
	c.File(filePath)
}

func getContentType(fileType string) string {
	switch fileType {
	case "txt":
		return "text/plain"
	case "md":
		return "text/markdown"
	case "docx", "doc":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case "pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

func DeleteKnowledgeFile(c *gin.Context) {
	userId, _ := c.Get("userId")
	fileIdStr := c.Param("file_id")
	
	fileId, err := strconv.Atoi(fileIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}
	
	// 首先验证文件是否存在且属于当前用户
	file, errRepo := repository.GetKnowledgeFileById(uint(fileId))
	if errRepo != nil {
		c.JSON(http.StatusInternalServerError, errRepo)
		return
	}
	
	if file == nil {
		c.JSON(http.StatusNotFound, exceptions.InternalServerError(fmt.Errorf("找不到ID为%d的知识库文件", fileId)))
		return
	}

	if file.UserId != userId.(uint) {
		c.JSON(http.StatusForbidden, exceptions.ErrPermissionDenied)
		return
	}

	// 删除物理文件
	knowledgeDir := filepath.Join(config.Config.Server.Upload, "knowledge")
	filePath := filepath.Join(knowledgeDir, file.FilePath)
	
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			c.JSON(http.StatusInternalServerError, exceptions.InternalServerError(fmt.Errorf("删除物理文件失败: %w", err)))
			return
		}
	}

	// 删除数据库记录
	if err := repository.DeleteKnowledgeFile(uint(fileId)); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文件删除成功",
		"file_id": fileId,
	})
} 