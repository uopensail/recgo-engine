package userctx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/recapi"
	"github.com/uopensail/recgo-engine/resources"
	"github.com/uopensail/ulib/sample"
)

type Request struct {
	UserId string `json:"user_id"`
}

type Response struct {
	Code int                     `json:"code"`
	Msg  string                  `json:"msg"`
	Data *sample.MutableFeatures `json:"data"`
}

type UserContext struct {
	context.Context
	Request  *recapi.Request
	Items    *model.Items
	Filter   model.IFliter
	Features *sample.MutableFeatures
	Related  *sample.ImmutableFeatures
}

func NewUserContext(ctx context.Context, req *recapi.Request) *UserContext {
	items := resources.ResourceManagerInstance.GetItems()
	var related *sample.ImmutableFeatures
	if len(req.RelateId) > 0 {
		id, feas := items.GetByKey(req.RelateId)
		if id >= 0 {
			related = feas
		}
	}

	uCtx := UserContext{
		Context:  ctx,
		Request:  req,
		Items:    items,
		Filter:   nil,
		Features: nil,
		Related:  related,
	}

	var features *sample.MutableFeatures
	features = uCtx.fetchFeatures()
	if features == nil {
		features = sample.NewMutableFeatures()
	}

	merge := func(key string, feature sample.Feature) error {
		features.Set(key, feature)
		return nil
	}
	uCtx.Request.Featues.ForEach(merge)
	uCtx.Request.Context.ForEach(merge)
	uCtx.Features = features

	return &uCtx
}

func (uCtx *UserContext) fetchFeatures() *sample.MutableFeatures {
	client := &http.Client{
		Timeout: 30 * time.Second, // 设置超时时间
		Transport: &http.Transport{
			MaxIdleConns:        100,              // 最大空闲连接数
			MaxIdleConnsPerHost: 10,               // 每个主机的最大空闲连接数
			IdleConnTimeout:     90 * time.Second, // 空闲连接超时时间
		},
	}

	reqest := Request{
		UserId: uCtx.Request.UserId,
	}
	data, err := json.Marshal(reqest)
	if err != nil {
		// fmt.Errorf("JSON序列化失败: %v", err)
		return nil
	}

	// 创建请求
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(data))
	if err != nil {
		// fmt.Errorf("创建请求失败: %v", err)
		return nil
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 发送请求
	response, err := client.Do(req)
	if err != nil {
		// fmt.Errorf("发送请求失败: %v", err)
		return nil
	}

	// 读取响应体
	body, err := io.ReadAll(response.Body)
	if err != nil {
		// fmt.Errorf("读取响应失败: %v", err)
		return nil
	}

	resp := Response{}
	// 解析JSON响应
	if err := json.Unmarshal(body, &resp); err != nil {
		// fmt.Errorf("JSON反序列化失败: %v", err)
		return nil
	}
	return resp.Data
}
