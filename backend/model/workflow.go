package model

import (
	"backend/db"
	"backend/exceptions"
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
)

const (
	StatusPrivate = "private"
	StatusTeam    = "team"
	StatusPublic  = "public"
	StatusOffline = "offline"
)

var WorkflowStatuses = []string{StatusPrivate, StatusTeam, StatusPublic, StatusOffline}

type Workflow struct {
	gorm.Model
	UserId        uint    `gorm:"not null"`
	BelongsTo     User    `gorm:"foreignKey:UserId"`
	Name          string  `gorm:"varchar(255);not null"`
	Description   string  `gorm:"type:text;not null"`
	Cover         string  `gorm:"varchar(255);"`
	Clicks        uint    `gorm:"not null;default:0"`
	Likes         uint    `gorm:"-"`
	LikeStatus    bool    `gorm:"-"`
	Collects      uint    `gorm:"-"`
	CollectStatus bool    `gorm:"-"`
	Status        string  `gorm:"varchar(16);not null;default:'private'"`
	Content       string  `gorm:"type:longtext;not null"`
	Nodes         []*Node `gorm:"-"`
	Edges         []*Edge `gorm:"-"`
}

func (w *Workflow) ParseContent() error {
	type Content struct {
		Nodes []*Node `json:"nodes"`
		Edges []*Edge `json:"edges"`
	}
	var c Content
	err := json.Unmarshal([]byte(w.Content), &c)
	if err != nil {
		return err
	}
	w.Nodes = c.Nodes
	w.Edges = c.Edges
	return nil
}

func (w *Workflow) SaveToContent() error {
	type Content struct {
		Nodes []*Node `json:"nodes"`
		Edges []*Edge `json:"edges"`
	}
	c := Content{
		Nodes: w.Nodes,
		Edges: w.Edges,
	}
	content, err := json.Marshal(c)
	if err != nil {
		return err
	}
	w.Content = string(content)
	return nil
}

func (w *Workflow) GetRunningOrder(endNode *Node) ([]*Node, error) {
	nodesMap := make(map[string]*Node)
	for _, node := range w.Nodes {
		nodesMap[node.Id] = node
	}

	order := make([]*Node, 0)
	visited := make(map[string]bool)
	var dfs func(*Node) error
	dfs = func(now *Node) error {
		fmt.Println(fmt.Sprintf("visit %v", now.Id))
		if visited[now.Id] {
			if now.Id == endNode.Id {
				return exceptions.InternalServerError(fmt.Errorf("cycle detected"))
			}
			return nil
		}
		visited[now.Id] = true

		prev := make(map[string]*Node)

		for _, input := range now.Data.Inputs {
			if input.Assigned {
				if nodesMap[input.AssignedNodeId] == nil {
					return exceptions.InternalServerError(fmt.Errorf("node %v (%v) not found", input.AssignedNodeId, input.AssignedNodeName))
				}
				prev[input.AssignedNodeId] = nodesMap[input.AssignedNodeId]
			}
		}

		for _, prevNode := range prev {
			fmt.Println(fmt.Sprintf("%v need %v", now.Id, prevNode.Id))
			if err := dfs(prevNode); err != nil {
				return err
			}
		}

		order = append(order, now)
		return nil
	}

	return order, dfs(endNode)
}

type FavoriteWorkflowRecord struct {
	gorm.Model
	WorkflowId            uint      `gorm:"not null;index"`
	UserId                uint      `gorm:"not null;index"`
	CorrespondingWorkflow *Workflow `gorm:"foreignKey:WorkflowId"`
	BelongsTo             *User     `gorm:"foreignKey:UserId"`
}

type LikeWorkflowRecord struct {
	gorm.Model
	WorkflowId            uint      `gorm:"not null;index"`
	UserId                uint      `gorm:"not null;index"`
	CorrespondingWorkflow *Workflow `gorm:"foreignKey:WorkflowId"`
	BelongsTo             *User     `gorm:"foreignKey:UserId"`
}

func (w *Workflow) MarshalJSON() ([]byte, error) {
	var content interface{}
	e := json.Unmarshal([]byte(w.Content), &content)
	if e != nil {
		return nil, e
	}
	return json.Marshal(struct {
		ID            uint        `json:"id"`
		CreatedAt     string      `json:"createdAt"`
		UpdatedAt     string      `json:"updatedAt"`
		BelongsTo     User        `json:"belongsTo"`
		Name          string      `json:"name"`
		Description   string      `json:"description"`
		Cover         string      `json:"cover"`
		Clicks        uint        `json:"clicks"`
		Likes         uint        `json:"likes"`
		LikeStatus    bool        `json:"likeStatus"`
		Collects      uint        `json:"collects"`
		CollectStatus bool        `json:"collectStatus"`
		Status        string      `json:"status"`
		Content       interface{} `json:"content"`
	}{
		ID:            w.ID,
		CreatedAt:     w.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     w.UpdatedAt.Format("2006-01-02 15:04:05"),
		BelongsTo:     w.BelongsTo,
		Name:          w.Name,
		Description:   w.Description,
		Cover:         w.Cover,
		Clicks:        w.Clicks,
		Likes:         w.Likes,
		LikeStatus:    w.LikeStatus,
		Collects:      w.Collects,
		CollectStatus: w.CollectStatus,
		Status:        w.Status,
		Content:       content,
	})
}

func (w *FavoriteWorkflowRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID                    uint      `json:"id"`
		CreatedAt             string    `json:"createdAt"`
		WorkflowId            uint      `json:"workflowId"`
		UserId                uint      `json:"userId"`
		CorrespondingWorkflow *Workflow `json:"correspondingWorkflow"`
	}{
		ID:                    w.ID,
		CreatedAt:             w.CreatedAt.Format("2006-01-02 15:04:05"),
		WorkflowId:            w.WorkflowId,
		UserId:                w.UserId,
		CorrespondingWorkflow: w.CorrespondingWorkflow,
	})
}

func (w *Workflow) BeforeCreate(scope *gorm.Scope) error {
	// 确保Clicks字段有默认值
	if w.Clicks == 0 {
		scope.SetColumn("Clicks", 0)
	}
	return nil
}

func init() {
	db.DB.AutoMigrate(&Workflow{}, &FavoriteWorkflowRecord{}, &LikeWorkflowRecord{})
}
