package pubsub

import "encoding/json"

type Opts struct {
	AppKey    string      `json:"appKey"`
	AppSecret string      `json:"appSecret"`
	AgentId   json.Number `json:"agentId"`
	DeptId    json.Number `json:"deptId"`
}
