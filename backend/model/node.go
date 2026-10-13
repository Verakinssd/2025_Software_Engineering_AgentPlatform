package model

import (
	dbPkg "backend/db"
	"backend/exceptions"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	NodeTypeStart = "start" // 开始节点
	NodeTypeEnd   = "end"   // 结束节点

	NodeTypeConcat = "concat" // 串联节点
	NodeTypeLLM    = "llm"    // LLM节点
	NodeTypeCode   = "code"   // 代码节点

	NodeTypeCondition = "condition" // 分支节点
	NodeTypeLoop      = "loop"      // 循环节点
	NodeTypeIteration = "iteration" // 迭代节点

	NodeWorkflow      = "workflow"  // 工作流节点
	NodeKnowledgeBase = "knowledge_base" // 知识库节点

)

type ComputedPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Dimensions struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Node struct {
	ComputedPosition ComputedPosition `json:"computedPosition"`
	Data             Data             `json:"data"`
	Dimensions       Dimensions       `json:"dimensions"`
	Dragging         bool             `json:"dragging"`
	Events           interface{}      `json:"events"`
	HandleBounds     interface{}      `json:"handleBounds"`
	Id               string           `json:"id"`
	Initialized      bool             `json:"initialized"`
	IsParent         bool             `json:"isParent"`
	Position         Position         `json:"position"`
	Resizing         bool             `json:"resizing"`
	Selected         bool             `json:"selected"`
	Type             string           `json:"type"`
}

type Data struct {
	Name    string            `json:"name"`
	Inputs  []*NodeDataStream `json:"input"`
	Outputs []*NodeDataStream `json:"output"`

	ConcatPrompt string `json:"concatPrompt"` // Concat节点的拼接规则

	SystemPrompt string `json:"systemPrompt"` // LLM节点的系统提示词
	UserPrompt   string `json:"userPrompt"`   // LLM节点的用户提示词模板

	CodePython     string `json:"codePython"`     // 代码节点的Python代码
	CodeJavaScript string `json:"codeJavaScript"` // 代码节点的JavaScript代码
	CodeType       string `json:"codeType"`       // 代码类型，"Python"或"JavaScript"

	KnowledgeFileId uint   `json:"knowledgeFileId"` // 知识库文件ID
	QueryTemplate   string `json:"queryTemplate"`   // 知识库查询模板
	
	WorkflowId      uint   `json:"workflowId"`      // 工作流节点引用的工作流ID
}

type NodeDataStream struct {
	Value                 any    `json:"value"`
	Key                   string `json:"value_name"`
	Type                  string `json:"value_type"`
	Assigned              bool   `json:"assigned"`
	AssignedNodeId        string `json:"assignedNodeId"`
	AssignedNodeName      string `json:"assignedNodeName"`
	AssignedNodeType      string `json:"assignedNodeType"`
	AssignedAttributeName string `json:"assignedAttributeName"`
	AssignedAttributeType string `json:"assignedAttributeType"`
}

// LLM API请求结构体
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type WebSearchTool struct {
	Enable      bool    `json:"enable"`                // 必填
	SearchQuery *string `json:"search_query,omitempty"`// 选填：给模型一个"首选关键词"
}

type Tool struct {
	Type      string         `json:"type"`                 // 固定写 "web_search"
	WebSearch *WebSearchTool `json:"web_search,omitempty"` // 只用这一种工具就够
}

type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
	Tools       []Tool        `json:"tools,omitempty"`
}


type ChatCompletionChoice struct {
	Index        int              `json:"index"`
	Message      ChatMessage      `json:"message"`
	FinishReason string           `json:"finish_reason"`
	Delta        ChatMessageDelta `json:"delta,omitempty"`
}

type ChatMessageDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
}

func (n *Node) Run() error {
	switch n.Type {
	case NodeTypeStart:
		return n.RunStart()
	case NodeTypeEnd:
		return n.RunEnd()
	case NodeTypeConcat:
		return n.RunConcat()
	case NodeTypeLLM:
		return n.RunLLM()
	case NodeTypeCode:
		return n.RunCodeNode()
	case NodeKnowledgeBase, "knowledge":
		return n.RunKnowledgeBase()
	case NodeWorkflow:
		return n.RunWorkflowNode()
	default:
		return exceptions.ErrNodeTypeInvalid.WithDetail(n.Type)
	}
}

func (n *Node) RunStart() error {
	n.Data.Outputs = n.Data.Inputs
	return nil
}

func (n *Node) RunEnd() error {
	n.Data.Outputs = n.Data.Inputs
	return nil
}

func (n *Node) RunCodeNode() error {
	if len(n.Data.Outputs) == 0 {
		return exceptions.ErrNodeOutputKeyNotFound.WithDetail("output key not found")
	}

	fmt.Println("######   " + n.Data.CodeType)

	switch n.Data.CodeType {
	case "JavaScript":
		if n.Data.CodeJavaScript == "" {
			return exceptions.ErrNodeInputKeyNotFound.WithDetail("JavaScript代码为空")
		}
		// TODO: 实现JavaScript代码执行
		return exceptions.InternalServerError(fmt.Errorf("JavaScript代码执行尚未实现"))
	case "Python", "": // 默认为Python
		if n.Data.CodePython == "" {
			return exceptions.ErrNodeInputKeyNotFound.WithDetail("Python代码为空")
		}
		return n.RunPythonCode()
	default:
		return exceptions.ErrNodeTypeInvalid.WithDetail(fmt.Sprintf("不支持的代码类型: %s", n.Data.CodeType))
	}
}

/* code template
def main(args: Args) -> dict:
    """
    Python代码节点的入口函数，您需要编写此函数的实现。

    参数:
        args: 包含输入参数的对象
            args.params: 包含所有输入值的字典，键名对应输入节点的变量名

    返回:
        dict: 必须返回一个字典，键名对应输出节点的变量名
              如果缺少任何在输出中定义的键名，节点将报错
    """
    # 获取输入参数
    params = args.params

    # 示例：获取名为'input'的输入值
    input_value = params.get('input', '')

    # 处理数据
    processed_value = input_value.upper() if isinstance(input_value, str) else str(input_value)

    # 准备返回结果 - 必须包含所有在输出节点中定义的键名
    result = {
        "key0": processed_value,                # 处理后的输入
        "key1": [1, 2, 3, 4, 5],               # 示例：输出一个列表
        "key2": {"name": "示例", "value": 100}  # 示例：输出一个嵌套对象
    }

    # 返回结果
    return result
*/

func (n *Node) RunPythonCode() error {
	inputParamsBuilder := strings.Builder{}
	inputParamsBuilder.WriteString("class Args:\n")
	inputParamsBuilder.WriteString("    def __init__(self):\n")
	inputParamsBuilder.WriteString("        self.params = {}\n")

	for _, input := range n.Data.Inputs {
		if input.Value != nil {
			var valueStr string
			switch v := input.Value.(type) {
			case string:
				valueStr = fmt.Sprintf("\"%s\"", v)
			default:
				valueStr = fmt.Sprintf("%v", v)
			}
			inputParamsBuilder.WriteString(fmt.Sprintf("        self.params[\"%s\"] = %s\n", input.Key, valueStr))
		}
	}

	pythonCode := fmt.Sprintf(`%s

%s

args = Args()

result = main(args)

import json
print(json.dumps(result))
`, inputParamsBuilder.String(), n.Data.CodePython)

	type CodeExecuteRequest struct {
		Code    string `json:"code"`
		Timeout string `json:"timeout"`
	}

	type CodeExecuteResponse struct {
		Result  string `json:"result"`
		Success bool   `json:"success"`
	}

	requestBody := CodeExecuteRequest{
		Code:    pythonCode,
		Timeout: "30", // 30秒超时
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return exceptions.InternalServerError(err).WithDetail("序列化代码执行请求失败")
	}

	apiURL := "http://localhost:5000/execute"
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return exceptions.InternalServerError(err).WithDetail("创建HTTP请求失败")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return exceptions.InternalServerError(err).WithDetail("发送HTTP请求失败")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return exceptions.InternalServerError(fmt.Errorf("API请求返回状态码: %d, 响应: %s", resp.StatusCode, string(bodyBytes)))
	}

	var execResult CodeExecuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&execResult); err != nil {
		return exceptions.InternalServerError(err).WithDetail("解析API响应失败")
	}

	if !execResult.Success {
		return exceptions.InternalServerError(fmt.Errorf("Python代码执行失败")).WithDetail(execResult.Result)
	}

	outputStr := strings.TrimSpace(execResult.Result)
	if outputStr == "" {
		return exceptions.InternalServerError(fmt.Errorf("Python代码没有输出任何结果"))
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(outputStr), &result); err != nil {
		return exceptions.InternalServerError(err).WithDetail(fmt.Sprintf("解析Python输出结果失败，输出不是有效的JSON格式:\n%s", outputStr))
	}

	missingKeys := []string{}
	for _, output := range n.Data.Outputs {
		if _, ok := result[output.Key]; !ok {
			missingKeys = append(missingKeys, output.Key)
		}
	}

	if len(missingKeys) > 0 {
		return exceptions.InternalServerError(fmt.Errorf("Python输出缺少必需的键")).WithDetail(
			fmt.Sprintf("缺少的键: %s\n实际输出: %s", strings.Join(missingKeys, ", "), outputStr))
	}

	for _, output := range n.Data.Outputs {
		if val, ok := result[output.Key]; ok {
			output.Value = val
		}
	}
	return nil
}

func (n *Node) RunLLM() error {
	if len(n.Data.Outputs) == 0 {
		return exceptions.ErrNodeOutputKeyNotFound.WithDetail("output key not found")
	}

	userPrompt := n.Data.UserPrompt
	for _, input := range n.Data.Inputs {
		key := "{{" + input.Key + "}}"
		userPrompt = strings.ReplaceAll(userPrompt, key, fmt.Sprintf("%v", input.Value))
	}

	if userPrompt == "" {
		return exceptions.ErrNodeInputKeyNotFound.WithDetail("LLM节点需要非空的用户提示词")
	}

	modelName := "glm-4-flash-250414"

	temperature := 0.7

	maxTokens := 1000

	messages := []ChatMessage{
		{Role: "system", Content: n.Data.SystemPrompt},
		{Role: "user", Content: userPrompt},
	}

	tools := []Tool{
		{
			Type: "web_search",
			WebSearch: &WebSearchTool{
				Enable: true,                // 打开联网
				// SearchQuery 也可以留空，让模型自己决定搜什么
			},
		},
	}

	requestBody := ChatCompletionRequest{
		Model:       modelName,
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   maxTokens,
		Stream:      false,
		Tools:       tools,
	}

	// apiKey := os.Getenv("ZHIPU_API_KEY")
	apiKey := "59378aa3c0ec426a8c4f20256009cbd3.pThBZsVUut90jDyD"
	if apiKey == "" {
		return exceptions.InternalServerError(fmt.Errorf("未设置ZHIPU_API_KEY环境变量"))
	}

	apiURL := "https://open.bigmodel.cn/api/paas/v4/chat/completions"

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return exceptions.InternalServerError(err).WithDetail("序列化请求数据失败")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return exceptions.InternalServerError(err).WithDetail("创建HTTP请求失败")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return exceptions.InternalServerError(err).WithDetail("发送HTTP请求失败")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return exceptions.InternalServerError(fmt.Errorf("API请求返回状态码: %d, 响应: %s", resp.StatusCode, string(bodyBytes)))
	}

	var result ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return exceptions.InternalServerError(err).WithDetail("解析API响应失败")
	}

	if len(result.Choices) == 0 {
		return exceptions.InternalServerError(fmt.Errorf("API返回结果中没有选择项"))
	}

	response := result.Choices[0].Message.Content
	n.Data.Outputs[0].Value = response

	fmt.Println(response)
	return nil
}

func (n *Node) RunConcat() error {
	if len(n.Data.Outputs) == 0 {
		return exceptions.ErrNodeOutputKeyNotFound.WithDetail("output key not found")
	}
	concatPrompt := n.Data.ConcatPrompt
	for _, input := range n.Data.Inputs {
		key := "{{" + input.Key + "}}"
		concatPrompt = strings.ReplaceAll(concatPrompt, key, fmt.Sprintf("%v", input.Value))
	}
	n.Data.Outputs[0].Value = concatPrompt
	return nil
}

func (n *Node) RunKnowledgeBase() error {
	if len(n.Data.Outputs) == 0 {
		return exceptions.ErrNodeOutputKeyNotFound.WithDetail("output key not found")
	}

	queryTemplate := n.Data.QueryTemplate
	if queryTemplate == "" {
		return exceptions.ErrNodeInputKeyNotFound.WithDetail("知识库节点需要提供查询模板(QueryTemplate)")
	}

	finalQuery := queryTemplate
	for _, input := range n.Data.Inputs {
		placeholder := "{{" + input.Key + "}}"
		if strings.Contains(finalQuery, placeholder) {
			finalQuery = strings.ReplaceAll(finalQuery, placeholder, fmt.Sprintf("%v", input.Value))
		}
	}

	if n.Data.KnowledgeFileId == 0 {
		return exceptions.ErrNodeInputKeyNotFound.WithDetail("知识库节点未配置文件ID")
	}

	type KnowledgeFile struct {
		ID          uint
		FileName    string
		Content     string
		UserId      uint
		Description string
	}

	var file KnowledgeFile
	result := dbPkg.DB.Table("knowledge_files").Where("id = ?", n.Data.KnowledgeFileId).First(&file)
	if result.Error != nil {
		return exceptions.InternalServerError(result.Error).WithDetail(fmt.Sprintf("找不到ID为%d的知识库文件", n.Data.KnowledgeFileId))
	}

	ragResponse := retrieveFromKnowledgeBase(file.Content, finalQuery)

	n.Data.Outputs[0].Value = ragResponse

	return nil
}

func retrieveFromKnowledgeBase(content string, query string) string {
	form := url.Values{}
	form.Add("query", query)
	form.Add("doc", content)

	apiURL := "http://114.116.204.215:5002/api/query/"
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Sprintf("创建RAG请求失败: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("调用RAG API失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("RAG API返回错误状态码: %d", resp.StatusCode)
	}

	var result struct {
		Answer string `json:"answer"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Sprintf("解析RAG API响应失败: %v", err)
	}

	return result.Answer
}

func (n *Node) RunWorkflowNode() error {
	if n.Data.WorkflowId == 0 {
		return exceptions.InternalServerError(fmt.Errorf("未指定工作流ID"))
	}

	// 获取当前工作流上下文，防止工作流递归调用自身
	// 从数据库查询当前节点所属的工作流
	var currentWorkflow Workflow
	result := dbPkg.DB.Where("content LIKE ?", "%"+n.Id+"%").First(&currentWorkflow)
	if result.Error == nil {
		// 检测循环依赖 - 从当前工作流开始检测，看是否会回到自己
		visited := make(map[uint]bool)
		if err := n.detectWorkflowCycle(currentWorkflow.ID, visited); err != nil {
			return err
		}
	}

	var referencedWorkflow Workflow
	if err := dbPkg.DB.First(&referencedWorkflow, n.Data.WorkflowId).Error; err != nil {
		return exceptions.InternalServerError(err).WithDetail(fmt.Sprintf("获取工作流(ID:%d)失败", n.Data.WorkflowId))
	}

	if err := referencedWorkflow.ParseContent(); err != nil {
		return exceptions.InternalServerError(err).WithDetail("解析引用工作流内容失败")
	}

	var startNode, endNode *Node
	for _, node := range referencedWorkflow.Nodes {
		if node.Type == NodeTypeStart {
			startNode = node
		}
		if node.Type == NodeTypeEnd {
			endNode = node
		}
	}

	if startNode == nil {
		return exceptions.InternalServerError(fmt.Errorf("引用的工作流中找不到开始节点"))
	}
	if endNode == nil {
		return exceptions.InternalServerError(fmt.Errorf("引用的工作流中找不到结束节点"))
	}

	var missingInputs []string
	var extraInputs []string
	
	currentInputMap := make(map[string]*NodeDataStream)
	for _, input := range n.Data.Inputs {
		currentInputMap[input.Key] = input
	}
	
	for _, input := range startNode.Data.Inputs {
		if input.Key != "" {
			if _, exists := currentInputMap[input.Key]; !exists {
				missingInputs = append(missingInputs, input.Key)
			}
		}
	}
	
	for _, input := range n.Data.Inputs {
		found := false
		for _, startInput := range startNode.Data.Inputs {
			if startInput.Key == input.Key {
				found = true
				break
			}
		}
		if !found {
			extraInputs = append(extraInputs, input.Key)
		}
	}
	
	var missingOutputs []string
	var extraOutputs []string
	
	currentOutputMap := make(map[string]*NodeDataStream)
	for _, output := range n.Data.Outputs {
		currentOutputMap[output.Key] = output
	}
	
	for _, output := range endNode.Data.Outputs {
		if output.Key != "" {
			if _, exists := currentOutputMap[output.Key]; !exists {
				missingOutputs = append(missingOutputs, output.Key)
			}
		}
	}
	
	for _, output := range n.Data.Outputs {
		found := false
		for _, endOutput := range endNode.Data.Outputs {
			if endOutput.Key == output.Key {
				found = true
				break
			}
		}
		if !found {
			extraOutputs = append(extraOutputs, output.Key)
		}
	}
	
	var errorMsgs []string
	if len(missingInputs) > 0 {
		errorMsgs = append(errorMsgs, fmt.Sprintf("缺少必要的输入变量: %s", strings.Join(missingInputs, ", ")))
	}
	if len(extraInputs) > 0 {
		errorMsgs = append(errorMsgs, fmt.Sprintf("存在多余的输入变量: %s", strings.Join(extraInputs, ", ")))
	}
	if len(missingOutputs) > 0 {
		errorMsgs = append(errorMsgs, fmt.Sprintf("缺少必要的输出变量: %s", strings.Join(missingOutputs, ", ")))
	}
	if len(extraOutputs) > 0 {
		errorMsgs = append(errorMsgs, fmt.Sprintf("存在多余的输出变量: %s", strings.Join(extraOutputs, ", ")))
	}
	
	if len(errorMsgs) > 0 {
		return exceptions.InternalServerError(fmt.Errorf("工作流节点与引用工作流接口不匹配")).WithDetail(strings.Join(errorMsgs, "\n"))
	}
	
	inputs := make(map[string]interface{})
	for _, input := range n.Data.Inputs {
		if input.Value != nil {
			inputs[input.Key] = input.Value
		}
	}
	
	workflowCopy := &Workflow{
		Model:         referencedWorkflow.Model,
		UserId:        referencedWorkflow.UserId,
		BelongsTo:     referencedWorkflow.BelongsTo,
		Name:          referencedWorkflow.Name,
		Description:   referencedWorkflow.Description,
		Cover:         referencedWorkflow.Cover,
		Clicks:        referencedWorkflow.Clicks,
		Likes:         referencedWorkflow.Likes,
		LikeStatus:    referencedWorkflow.LikeStatus,
		Collects:      referencedWorkflow.Collects,
		CollectStatus: referencedWorkflow.CollectStatus,
		Status:        referencedWorkflow.Status,
		Content:       referencedWorkflow.Content,
	}
	
	if err := workflowCopy.ParseContent(); err != nil {
		return exceptions.InternalServerError(err).WithDetail("解析工作流副本内容失败")
	}
	
	for _, node := range workflowCopy.Nodes {
		if node.Type == NodeTypeStart {
			for _, input := range node.Data.Inputs {
				if inputValue, exists := inputs[input.Key]; exists {
					input.Value = inputValue
				}
			}
			node.Data.Outputs = node.Data.Inputs
			break
		}
	}
	
	nodesMap := make(map[string]*Node)
	var endNodeInCopy *Node
	
	for _, node := range workflowCopy.Nodes {
		nodesMap[node.Id] = node
		if node.Type == NodeTypeEnd {
			endNodeInCopy = node
			node.Data.Inputs = node.Data.Outputs
		}
	}
	
	runningOrder, err := workflowCopy.GetRunningOrder(endNodeInCopy)
	if err != nil {
		return exceptions.InternalServerError(err)
	}
	
	for _, node := range runningOrder {
		for _, input := range node.Data.Inputs {
			if input.Assigned {
				var succeed bool = false
				for _, output := range nodesMap[input.AssignedNodeId].Data.Outputs {
					if output.Key == input.AssignedAttributeName {
						input.Value = output.Value
						succeed = true
						break
					}
				}
				if !succeed {
					return exceptions.InternalServerError(fmt.Errorf("节点 %v 输入 %v 未找到对应输出", node.Id, input.Key))
				}
			}
		}
		
		if err := node.Run(); err != nil {
			return exceptions.InternalServerError(err).WithDetail(fmt.Sprintf("引用工作流节点 %v 运行失败", node.Id))
		}
	}
	
	for i, output := range n.Data.Outputs {
		for _, endOutput := range endNodeInCopy.Data.Outputs {
			if endOutput.Key == output.Key {
				n.Data.Outputs[i].Value = endOutput.Value
				break
			}
		}
	}
	
	return nil
}

func (n *Node) detectWorkflowCycle(currentWorkflowID uint, visited map[uint]bool) error {
	// 如果已经访问过当前工作流，说明存在循环
	if visited[currentWorkflowID] {
		return exceptions.InternalServerError(fmt.Errorf("检测到工作流循环引用")).WithDetail(
			fmt.Sprintf("工作流调用链中存在循环：工作流ID %d", currentWorkflowID))
	}
	
	// 标记当前工作流为已访问
	visited[currentWorkflowID] = true
	defer func() {
		delete(visited, currentWorkflowID) // 回溯时清除标记
	}()
	
	// 获取当前工作流
	var currentWorkflow Workflow
	if err := dbPkg.DB.First(&currentWorkflow, currentWorkflowID).Error; err != nil {
		return exceptions.InternalServerError(err).WithDetail(fmt.Sprintf("获取工作流(ID:%d)失败", currentWorkflowID))
	}
	
	// 解析工作流内容
	if err := currentWorkflow.ParseContent(); err != nil {
		return exceptions.InternalServerError(err).WithDetail(fmt.Sprintf("解析工作流(ID:%d)内容失败", currentWorkflowID))
	}
	
	// 检查当前工作流中的所有工作流节点
	for _, node := range currentWorkflow.Nodes {
		if node.Type == NodeWorkflow && node.Data.WorkflowId != 0 {
			// 递归检查引用的工作流
			if err := n.detectWorkflowCycle(node.Data.WorkflowId, visited); err != nil {
				return err
			}
		}
	}
	
	return nil
}
