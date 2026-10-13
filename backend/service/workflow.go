package service

import (
	"backend/config"
	"backend/exceptions"
	"backend/middleware"
	"backend/model"
	"backend/repository"
	"fmt"
)

func CreateWorkflow(userId uint, name string, description string) (*model.Workflow, *exceptions.Error) {
	w, e := repository.CreateWorkflow(userId, name, description)
	if e != nil {
		return nil, exceptions.InternalServerError(e)
	}
	return w, nil
}

func DeleteWorkflow(userId uint, workflowId uint) *exceptions.Error {
	if config.Config.Server.Environment == "server" {
		recordsDeletionTask, err := middleware.NewRecordsCasCadeDeletionTask(workflowId)
		if err != nil {
			return exceptions.InternalServerError(err)
		}
		middleware.AsynqClient.Enqueue(recordsDeletionTask)

		commentDeletionTask, err2 := middleware.NewCommentsCasCadeDeletionTask(workflowId)
		if err2 != nil {
			return exceptions.InternalServerError(err2)
		}
		middleware.AsynqClient.Enqueue(commentDeletionTask)
	}
	return repository.DeleteWorkflow(userId, workflowId)
}

func GetWorkflowList(userId uint) ([]*model.Workflow, *exceptions.Error) {
	w, e := repository.GetWorkflowList(userId)
	if e != nil {
		return nil, exceptions.InternalServerError(e)
	}
	return w, nil
}

func GetAllPublicWorkflows(page, pageSize uint) ([]*model.Workflow, *exceptions.Error) {
	w, e := repository.GetAllPublicWorkflows(page, pageSize)
	if e != nil {
		return nil, exceptions.InternalServerError(e)
	}
	return w, nil
}

func GetPublicWorkflowsWithSearchKeyword(page, pageSize uint, keyword string) ([]*model.Workflow, *exceptions.Error) {
	w, e := repository.GetPublicWorkflowsWithSearchKeywords(page, pageSize, keyword)
	if e != nil {
		return nil, exceptions.InternalServerError(e)
	}
	return w, nil
}

func GetPopularWorkflows(pageSize uint) ([]*model.Workflow, *exceptions.Error) {
	w, e := repository.GetPopularWorkflows(pageSize)
	if e != nil {
		return nil, exceptions.InternalServerError(e)
	}
	return w, nil
}

func GetWorkflowById(userId uint, workflowId uint) (*model.Workflow, *exceptions.Error) {
	u, e1 := repository.GetUserById(userId)
	if e1 != nil {
		return nil, exceptions.InternalServerError(e1)
	}
	w, e2 := repository.GetWorkflow(userId, workflowId)
	if e2 != nil {
		return nil, exceptions.InternalServerError(e2)
	}
	if u.ID != w.UserId {
		// TODO: Team permission control
		if w.Status != model.StatusPublic {
			return nil, exceptions.ErrPermissionDenied
		}
	}
	return w, nil
}

func RunWorkflow(workflow *model.Workflow, nodeId string) *exceptions.Error {

	var endNode *model.Node
	nodesMap := make(map[string]*model.Node)
	for _, node := range workflow.Nodes {
		nodesMap[node.Id] = node

		if nodeId == "" && node.Type == model.NodeTypeEnd {
			endNode = node
		} else if node.Id == nodeId {
			endNode = node
		}
		// 将开始和结束节点的输入输出变成一样的
		if node.Type == model.NodeTypeStart {
			node.Data.Outputs = node.Data.Inputs
		}
		if node.Type == model.NodeTypeEnd {
			node.Data.Inputs = node.Data.Outputs
		}
	}
	if endNode == nil {
		return exceptions.InvalidRequest(fmt.Errorf("节点 %v 不存在", nodeId))
	}

	runningOrder, err := workflow.GetRunningOrder(endNode)
	if err != nil {
		return exceptions.InternalServerError(err)
	}

	for _, node := range runningOrder {
		fmt.Printf("Running node %v\n", node.Id)
		// 给每个节点的输入赋值
		for _, input := range node.Data.Inputs {
			if input.Assigned {
				var succeed bool = false
				fmt.Printf("node %v input %v is assigned with node %v key %v\n", node.Id, input.Key, input.AssignedNodeId, input.AssignedAttributeName)
				for _, output := range nodesMap[input.AssignedNodeId].Data.Outputs {
					if output.Key == input.AssignedAttributeName {
						input.Value = output.Value
						succeed = true
						fmt.Printf("Assigning key %v with value %v to %v\n", input.AssignedAttributeName, output.Value, input.Key)
						break
					}
				}
				if !succeed {
					return exceptions.InternalServerError(fmt.Errorf("节点 %v 输入 %v 未找到对应输出", node.Id, input.Key)).WithDetail(fmt.Sprintf("节点 %v 输入 %v 未找到对应输出", node.Id, input.Key))
				}
			}
		}
		if err := node.Run(); err != nil {
			return exceptions.InternalServerError(err).WithDetail(fmt.Sprintf("节点 %v 运行失败", node.Id))
		}
	}

	if err := workflow.SaveToContent(); err != nil {
		return exceptions.InternalServerError(err).WithDetail("保存工作流内容失败")
	}
	return nil

}

func ForkWorkflow(userId uint, workflowId uint) (*model.Workflow, *exceptions.Error) {
	w, e := repository.GetWorkflow(userId, workflowId)
	if e != nil {
		return nil, e
	}
	newW, e := repository.CreateWorkflow(userId, "Fork of "+w.Name, "This is a fork of "+w.Name+".\n\n"+w.Description)
	if e != nil {
		return nil, e
	}
	newW.Content = w.Content
	newW.Cover = w.Cover
	db.Save(&newW)
	return newW, nil
}

func FavoriteWorkflow(userId uint, workflowId uint) (bool, *exceptions.Error) {
	return repository.FavoriteWorkflow(userId, workflowId)
}

func GetFavoriteWorkflowList(userId uint) ([]*model.FavoriteWorkflowRecord, *exceptions.Error) {
	return repository.GetFavoriteWorkflowList(userId)
}

func LikeWorkflow(userId uint, workflowId uint) (bool, *exceptions.Error) {
	return repository.LikeWorkflow(userId, workflowId)
}

func RunWorkflowWithInputs(workflow *model.Workflow, nodeId string, inputs map[string]interface{}) (*model.Workflow, *exceptions.Error) {
	// 创建工作流的深拷贝
	workflowCopy := &model.Workflow{
		Model:       workflow.Model,
		UserId:      workflow.UserId,
		BelongsTo:   workflow.BelongsTo,
		Name:        workflow.Name,
		Description: workflow.Description,
		Cover:       workflow.Cover,
		Clicks:      workflow.Clicks,
		Likes:       workflow.Likes,
		LikeStatus:  workflow.LikeStatus,
		Collects:    workflow.Collects,
		CollectStatus: workflow.CollectStatus,
		Status:      workflow.Status,
		Content:     workflow.Content,
		Nodes:       nil, // 稍后重新解析
		Edges:       nil, // 稍后重新解析
	}
	
	// 解析内容到Nodes和Edges
	if err := workflowCopy.ParseContent(); err != nil {
		return nil, exceptions.InternalServerError(err).WithDetail("解析工作流内容失败")
	}
	
	// 处理输入 - 将自定义输入应用到开始节点
	for _, node := range workflowCopy.Nodes {
		if node.Type == model.NodeTypeStart {
			// 首先应用自定义输入到开始节点
			if node.Data.Inputs != nil {
				for _, input := range node.Data.Inputs {
					if input.Key != "" && inputs[input.Key] != nil {
						input.Value = inputs[input.Key]
					}
				}
			}
			
			// 也可以添加新的输入
			for key, value := range inputs {
				// 检查是否已有此输入
				found := false
				for _, input := range node.Data.Inputs {
					if input.Key == key {
						found = true
						break
					}
				}
				
				// 如果没有找到，添加新输入
				if !found {
					valueType := "string"
					switch value.(type) {
					case int, int32, int64:
						valueType = "integer"
					case float32, float64:
						valueType = "number"
					case bool:
						valueType = "boolean"
					}
					
					node.Data.Inputs = append(node.Data.Inputs, &model.NodeDataStream{
						Key:   key,
						Value: value,
						Type:  valueType,
					})
				}
			}
			break
		}
	}
	
	// ===== 以下是与原RunWorkflow函数完全相同的运行逻辑 =====
	var endNode *model.Node
	nodesMap := make(map[string]*model.Node)
	for _, node := range workflowCopy.Nodes {
		nodesMap[node.Id] = node

		if nodeId == "" && node.Type == model.NodeTypeEnd {
			endNode = node
		} else if node.Id == nodeId {
			endNode = node
		}
		// 将开始和结束节点的输入输出变成一样的
		if node.Type == model.NodeTypeStart {
			node.Data.Outputs = node.Data.Inputs
		}
		if node.Type == model.NodeTypeEnd {
			node.Data.Inputs = node.Data.Outputs
		}
	}
	if endNode == nil {
		return nil, exceptions.InvalidRequest(fmt.Errorf("节点 %v 不存在", nodeId))
	}

	runningOrder, err := workflowCopy.GetRunningOrder(endNode)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	for _, node := range runningOrder {
		fmt.Printf("Running node %v\n", node.Id)
		// 给每个节点的输入赋值
		for _, input := range node.Data.Inputs {
			if input.Assigned {
				var succeed bool = false
				fmt.Printf("node %v input %v is assigned with node %v key %v\n", node.Id, input.Key, input.AssignedNodeId, input.AssignedAttributeName)
				for _, output := range nodesMap[input.AssignedNodeId].Data.Outputs {
					if output.Key == input.AssignedAttributeName {
						input.Value = output.Value
						succeed = true
						fmt.Printf("Assigning key %v with value %v to %v\n", input.AssignedAttributeName, output.Value, input.Key)
						break
					}
				}
				if !succeed {
					return nil, exceptions.InternalServerError(fmt.Errorf("节点 %v 输入 %v 未找到对应输出", node.Id, input.Key)).WithDetail(fmt.Sprintf("节点 %v 输入 %v 未找到对应输出", node.Id, input.Key))
				}
			}
		}
		if err := node.Run(); err != nil {
			return nil, exceptions.InternalServerError(err).WithDetail(fmt.Sprintf("节点 %v 运行失败", node.Id))
		}
	}

	if err := workflowCopy.SaveToContent(); err != nil {
		return nil, exceptions.InternalServerError(err).WithDetail("保存工作流内容失败")
	}
	
	return workflowCopy, nil
}
