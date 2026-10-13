package controller

import (
	dbPkg "backend/db"
	"backend/exceptions"
	"backend/model"
	"backend/repository"
	"backend/service"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateWorkflow(c *gin.Context) {
	type Req struct {
		Name        string `form:"name" binding:"required"`
		Description string `form:"description" binding:"required"`
	}
	var req Req
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	userId, _ := c.Get("userId")
	workflow, e := service.CreateWorkflow(userId.(uint), req.Name, req.Description)
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}
	c.JSON(http.StatusOK, gin.H{"workflow_id": workflow.ID})
}

func GetWorkflowList(c *gin.Context) {
	userId, _ := c.Get("userId")
	workflows, e := service.GetWorkflowList(userId.(uint))
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}
	c.JSON(http.StatusOK, gin.H{"workflows": workflows})
}

func GetAllPublicWorkflows(c *gin.Context) {
	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err).
			WithDetail("Invalid page number, should be a positive integer greater than 0, got "+pageStr))
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	var workflows []*model.Workflow
	var e *exceptions.Error = nil

	searchKeyword := c.Query("key")
	if searchKeyword == "" {
		workflows, e = service.GetAllPublicWorkflows(uint(page), uint(pageSize))
		if e != nil {
			c.JSON(http.StatusInternalServerError, e)
			return
		}
	} else {
		workflows, e = service.GetPublicWorkflowsWithSearchKeyword(uint(page), uint(pageSize), searchKeyword)
		if e != nil {
			c.JSON(http.StatusInternalServerError, e)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"workflows": workflows})
}

func GetPopularWorkflows(c *gin.Context) {
	pageSizeStr := c.Query("pageSize")
	if pageSizeStr == "" {
		pageSizeStr = "10"
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}
	workflows, e := service.GetPopularWorkflows(uint(pageSize))
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}
	c.JSON(http.StatusOK, gin.H{"workflows": workflows})
}

func GetWorkflowById(c *gin.Context) {
	userId, _ := c.Get("userId")
	workflowIdStr := c.Param("workflow_id")
	workflowId, err := strconv.Atoi(workflowIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}
	w, e := service.GetWorkflowById(userId.(uint), uint(workflowId))
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}
	c.JSON(http.StatusOK, w)
}

func UpdateWorkflow(c *gin.Context) {
	userId, _ := c.Get("userId")
	workflowIdStr := c.Param("workflow_id")
	workflowId, err := strconv.Atoi(workflowIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}
	type Req struct {
		Name        string                `form:"name"`
		Description string                `form:"description"`
		Cover       *multipart.FileHeader `form:"cover"`
		Status      string                `form:"status"`
	}
	var req Req
	err = c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	w, e := service.GetWorkflowById(userId.(uint), uint(workflowId))
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	w.Name = req.Name
	w.Description = req.Description
	for _, s := range model.WorkflowStatuses {
		if s == req.Status {
			w.Status = req.Status
		}
	}
	if req.Cover != nil {
		filename, e := saveFile(c, "cover", &[]string{".jpg", ".jpeg", ".png", ".gif", ".webp"})
		if e != nil {
			c.JSON(http.StatusInternalServerError, e)
			return
		}
		w.Cover = "/api/image/" + filename
	}
	w, e = repository.UpdateWorkflow(w)
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	c.JSON(http.StatusOK, gin.H{"workflow": w})
}

func ForkWorkflow(c *gin.Context) {
	userId, _ := c.Get("userId")
	workflowIdStr := c.Param("workflow_id")
	workflowId, err := strconv.Atoi(workflowIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	w, e := service.ForkWorkflow(userId.(uint), uint(workflowId))
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	c.JSON(http.StatusOK, gin.H{"workflow": w})
}

func DeleteWorkflow(c *gin.Context) {
	userId, _ := c.Get("userId")
	workflowIdStr := c.Param("workflow_id")
	workflowId, err := strconv.Atoi(workflowIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}
	e := service.DeleteWorkflow(userId.(uint), uint(workflowId))
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func UpdateWorkflowContent(c *gin.Context) {
	userId, _ := c.Get("userId")
	workflowIdStr := c.Param("workflow_id")
	workflowId, err := strconv.Atoi(workflowIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	type Req struct {
		Content interface{} `json:"content"`
	}

	var req Req
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	w, e := service.GetWorkflowById(userId.(uint), uint(workflowId))
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	if w.UserId != userId.(uint) {
		c.JSON(http.StatusForbidden, exceptions.ErrPermissionDenied.WithDetail("你没有权限修改此工作流"))
		return
	}

	s, err := json.Marshal(req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}
	w.Content = string(s)
	w, e = repository.UpdateWorkflow(w)
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}
	c.JSON(http.StatusOK, gin.H{"workflow": w})
}

func RunWorkflow(c *gin.Context) {
	userId, _ := c.Get("userId")
	workflowIdStr := c.Param("workflow_id")
	workflowId, err := strconv.Atoi(workflowIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}
	w, e := service.GetWorkflowById(userId.(uint), uint(workflowId))
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	nodeId := c.Query("node_id")

	if e := service.RunWorkflow(w, nodeId); e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	c.JSON(http.StatusOK, gin.H{"workflow": w})
}

func FavoriteWorkflow(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	workflowIdStr := c.Param("workflow_id")
	workflowId, err := strconv.Atoi(workflowIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	res, err2 := service.FavoriteWorkflow(userId, uint(workflowId))
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, err2)
	}

	// Return:
	// 	true if this favorite is added,
	// 	false if it is removed
	c.JSON(http.StatusOK, gin.H{"is_favorite": res})
}

func GetFavoriteWorkflowList(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	workflows, err := service.GetFavoriteWorkflowList(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	// TODO: add pagination
	c.JSON(http.StatusOK, gin.H{"workflows": workflows})
}

func LikeWorkflow(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	workflowIdStr := c.Param("workflow_id")
	workflowId, err := strconv.Atoi(workflowIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	res, err2 := service.LikeWorkflow(userId, uint(workflowId))
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, err2)
		return
	}

	// Return:
	// 	true if this like is added,
	// 	false if it is removed
	c.JSON(http.StatusOK, gin.H{"is_liked": res})
}

func RunWorkflowWithInputs(c *gin.Context) {
	userId, _ := c.Get("userId")
	workflowIdStr := c.Param("workflow_id")
	workflowId, err := strconv.Atoi(workflowIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}
	
	w, e := service.GetWorkflowById(userId.(uint), uint(workflowId))
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	nodeId := c.Query("node_id")
	
	var inputs map[string]interface{}
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err).WithDetail("解析输入参数失败"))
		return
	}
	
	result, e := service.RunWorkflowWithInputs(w, nodeId, inputs)
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	c.JSON(http.StatusOK, gin.H{"workflow": result})
}

func GetAvailableWorkflows(c *gin.Context) {
	userId, _ := c.Get("userId")
	
	var userWorkflows []model.Workflow
	if err := dbPkg.DB.Where("user_id = ?", userId).Find(&userWorkflows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError(err))
		return
	}
	
	var publicWorkflows []model.Workflow
	if err := dbPkg.DB.Where("status = ? AND user_id != ?", model.StatusPublic, userId).Find(&publicWorkflows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError(err))
		return
	}
	
	allWorkflows := append(userWorkflows, publicWorkflows...)
	
	type WorkflowInfo struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		UserId      uint   `json:"userId"`
		Status      string `json:"status"`
	}
	
	result := make([]WorkflowInfo, 0, len(allWorkflows))
	for _, w := range allWorkflows {
		result = append(result, WorkflowInfo{
			ID:          w.ID,
			Name:        w.Name,
			Description: w.Description,
			UserId:      w.UserId,
			Status:      w.Status,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{"workflows": result})
}

func GetPublicWorkflows(c *gin.Context) {
	userId, _ := c.Get("userId")
	
	var publicWorkflows []model.Workflow
	if err := dbPkg.DB.Where("status = ? AND user_id != ?", model.StatusPublic, userId).Find(&publicWorkflows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError(err))
		return
	}
	
	type WorkflowInfo struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		UserId      uint   `json:"userId"`
		Status      string `json:"status"`
	}
	
	result := make([]WorkflowInfo, 0, len(publicWorkflows))
	for _, w := range publicWorkflows {
		result = append(result, WorkflowInfo{
			ID:          w.ID,
			Name:        w.Name,
			Description: w.Description,
			UserId:      w.UserId,
			Status:      w.Status,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{"workflows": result})
}

func GetPersonalWorkflows(c *gin.Context) {
	userId, _ := c.Get("userId")
	
	var userWorkflows []model.Workflow
	if err := dbPkg.DB.Where("user_id = ?", userId).Find(&userWorkflows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError(err))
		return
	}
	
	type WorkflowInfo struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		UserId      uint   `json:"userId"`
		Status      string `json:"status"`
	}
	
	result := make([]WorkflowInfo, 0, len(userWorkflows))
	for _, w := range userWorkflows {
		result = append(result, WorkflowInfo{
			ID:          w.ID,
			Name:        w.Name,
			Description: w.Description,
			UserId:      w.UserId,
			Status:      w.Status,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{"workflows": result})
}

func GetWorkflowIOStructure(c *gin.Context) {
	userId, _ := c.Get("userId")
	workflowIdStr := c.Param("workflow_id")
	workflowId, err := strconv.Atoi(workflowIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.InvalidRequest(err))
		return
	}

	// 获取工作流
	w, e := service.GetWorkflowById(userId.(uint), uint(workflowId))
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	// 解析工作流内容
	if err := w.ParseContent(); err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError(err).WithDetail("解析工作流内容失败"))
		return
	}

	// 查找开始和结束节点
	var startNode, endNode *model.Node
	for _, node := range w.Nodes {
		if node.Type == model.NodeTypeStart {
			startNode = node
		}
		if node.Type == model.NodeTypeEnd {
			endNode = node
		}
	}

	if startNode == nil || endNode == nil {
		c.JSON(http.StatusInternalServerError, exceptions.InternalServerError(fmt.Errorf("工作流缺少开始或结束节点")))
		return
	}

	// 返回开始节点的输入和结束节点的输出
	c.JSON(http.StatusOK, gin.H{
		"inputs": startNode.Data.Inputs,
		"outputs": endNode.Data.Outputs,
	})
}
