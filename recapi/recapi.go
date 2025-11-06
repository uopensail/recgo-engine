package recapi

import (
	"github.com/uopensail/ulib/sample"
)

type Request struct {
	TraceId  string                  `json:"trace_id,omitempty"`  // 请求id
	UserId   string                  `json:"user_id"`             // 用户id
	RelateId string                  `json:"relate_id,omitempty"` // 相关推荐：详情页的商品id
	Count    int32                   `son:"count,omitempty"`      // 请求的数量
	Context  *sample.MutableFeatures `json:"context,omitempty"`   // 上下文特征
	Featues  *sample.MutableFeatures `json:"featues,omitempty"`   // 外部特征
}

type ItemInfo struct {
	Item     string   `json:"item,omitempty"`
	Channels []string `json:"channels,omitempty"`
	Reasons  []string `json:"reasons,omitempty"`
}

type Response struct {
	TraceId  string      `json:"trace_id,omitempty"` // 请求id
	UserId   string      `json:"user_id"`            // 用户ID
	Pipeline string      `json:"pipeline"`           // 策略管道
	Items    []*ItemInfo `json:"items"`
}
