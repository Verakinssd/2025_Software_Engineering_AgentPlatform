package service

//// 创建节点
//func CreateNode(userId uint, workflowId uint, name string, x float32, y float32, nodeType string, inputs []*nodes.NodeInput, outputs []*nodes.NodeOutput, config interface{}) (*nodes.Node, *exceptions.Error) {
//	workflow, err := GetWorkflowById(userId, workflowId)
//	if err != nil {
//		return nil, err
//	}
//
//	node, err := repository.CreateNode(workflow, name, x, y, nodeType, inputs, outputs, config)
//	if err != nil {
//		return nil, err
//	}
//
//	return node, nil
//}
//
//// 修改RunNode函数参数和返回值类型
//func RunNode(userId uint, nodeId uint, inputs *[]map[string]interface{}) (*[]map[string]interface{}, []uint, *exceptions.Error) {
//	node, e1 := repository.GetNodeById(userId, nodeId)
//	if e1 != nil {
//		return nil, nil, e1
//	}
//	outputs, nextNodeIds, e2 := runners.RunNode(node, inputs)
//	return outputs, nextNodeIds, e2
//}
//
//// 更新节点信息
//func UpdateNode(userId uint, nodeId uint, name string, x float32, y float32, inputs []*nodes.NodeInput, outputs []*nodes.NodeOutput, config interface{}) (*nodes.Node, *exceptions.Error) {
//	_, err := repository.GetNodeById(userId, nodeId)
//	if err != nil {
//		return nil, err
//	}
//
//	updatedNode, err := repository.UpdateNode(nodeId, name, x, y, inputs, outputs, config)
//	if err != nil {
//		return nil, err
//	}
//
//	return updatedNode, nil
//}
//
//// 删除节点
//func DeleteNode(userId uint, nodeId uint) *exceptions.Error {
//	_, err := repository.GetNodeById(userId, nodeId)
//	if err != nil {
//		return err
//	}
//
//	return repository.DeleteNode(nodeId)
//}
//
//// 执行整个工作流
//func RunWorkflow(userId uint, workflowId uint, initialInputs *[]map[string]interface{}) (*[]map[string]interface{}, *exceptions.Error) {
//	workflow, err := GetWorkflowById(userId, workflowId)
//	if err != nil {
//		return nil, err
//	}
//
//	beginNodeCount := 0
//	for _, node := range workflow.Nodes {
//		if node.Type == nodes.NodeTypeBegin {
//			beginNodeCount++
//			if beginNodeCount > 1 {
//				return nil, exceptions.ErrNodeTypeInvalid.WithDetail("工作流中存在多个开始节点")
//			}
//		}
//	}
//
//	var startNode *nodes.Node
//	for _, node := range workflow.Nodes {
//		if node.Type == nodes.NodeTypeBegin {
//			startNode = node
//			break
//		}
//	}
//
//	if startNode == nil {
//		return nil, exceptions.ErrNodeTypeInvalid.WithDetail("找不到开始节点")
//	}
//
//	endNodeCount := 0
//	for _, node := range workflow.Nodes {
//		if node.Type == nodes.NodeTypeEnd {
//			endNodeCount++
//			if endNodeCount > 1 {
//				return nil, exceptions.ErrNodeTypeInvalid.WithDetail("工作流中存在多个结束节点")
//			}
//		}
//	}
//
//	endNodeFound := false
//	for _, node := range workflow.Nodes {
//		if node.Type == nodes.NodeTypeEnd {
//			endNodeFound = true
//			break
//		}
//	}
//
//	if !endNodeFound {
//		return nil, exceptions.ErrNodeTypeInvalid.WithDetail("找不到结束节点")
//	}
//
//	return executeWorkflowFromNode(userId, startNode, initialInputs)
//}
//
//// 从特定节点执行工作流
//func executeWorkflowFromNode(userId uint, startNode *nodes.Node, inputs *[]map[string]interface{}) (*[]map[string]interface{}, *exceptions.Error) {
//	visitedNodes := make(map[uint]bool)
//
//	nodeQueue := []*nodes.Node{startNode}
//	currentInputs := inputs
//
//	for len(nodeQueue) > 0 {
//		currentNode := nodeQueue[0]
//		nodeQueue = nodeQueue[1:]
//
//		if visitedNodes[currentNode.ID] {
//			continue
//		}
//
//		visitedNodes[currentNode.ID] = true
//
//		if currentNode.Type == nodes.NodeTypeEnd {
//			// 运行结束节点，获取输出结果
//			outputs, _, err := runners.RunNode(currentNode, currentInputs)
//			if err != nil {
//				return nil, err
//			}
//			return outputs, nil
//		}
//
//		outputs, nextNodeIds, err := runners.RunNode(currentNode, currentInputs)
//		if err != nil {
//			return nil, err
//		}
//
//		currentInputs = outputs
//
//		if len(nextNodeIds) > 0 {
//			for _, nextNodeId := range nextNodeIds {
//				nextNode, err := repository.GetNodeById(userId, nextNodeId)
//				if err != nil {
//					return nil, err
//				}
//				nodeQueue = append(nodeQueue, nextNode)
//			}
//		}
//	}
//
//	return currentInputs, nil
//}
