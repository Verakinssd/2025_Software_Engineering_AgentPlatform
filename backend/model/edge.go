package model

type Edge struct {
	Data         interface{} `json:"data"`
	Events       interface{} `json:"events"`
	Id           string      `json:"id"`
	Label        string      `json:"label"`
	MarkerEnd    string      `json:"markerEnd"`
	Source       string      `json:"source"`
	SourceHandle any         `json:"sourceHandle"`
	SourceNode   Node        `json:"sourceNode"`
	SourceX      float64     `json:"sourceX"`
	SourceY      float64     `json:"sourceY"`
	Target       string      `json:"target"`
	TargetHandle any         `json:"targetHandle"`
	TargetNode   Node        `json:"targetNode"`
	TargetX      float64     `json:"targetX"`
	TargetY      float64     `json:"targetY"`
	Type         string      `json:"type"`
}
