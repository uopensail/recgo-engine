package recalls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/sample"
)

type Request struct {
	Features sample.Features `json:"features"`
	Excludes []string        `json:"excludes"`
	Count    int             `json:"count"`
}

type Response struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data []model.KeyScore `json:"data"`
}

type Modeler struct {
	conf   *model.ModelRecallConfigure
	client *http.Client
}

func NewModeler(conf *model.ModelRecallConfigure) *Modeler {
	client := &http.Client{
		Timeout: 30 * time.Second, // 设置超时时间
		Transport: &http.Transport{
			MaxIdleConns:        100,              // 最大空闲连接数
			MaxIdleConnsPerHost: 10,               // 每个主机的最大空闲连接数
			IdleConnTimeout:     90 * time.Second, // 空闲连接超时时间
		},
	}

	return &Modeler{
		conf:   conf,
		client: client,
	}
}

func (m *Modeler) Do(uCtx *userctx.UserContext) model.Collection {
	reqest := Request{
		Features: uCtx.Features,
		Excludes: uCtx.Filter.Exclude(),
		Count:    m.conf.Count,
	}
	data, err := json.Marshal(reqest)
	if err != nil {
		// fmt.Errorf("JSON序列化失败: %v", err)
		return nil
	}

	// 创建请求
	req, err := http.NewRequest("POST", m.conf.URL, bytes.NewBuffer(data))
	if err != nil {
		// fmt.Errorf("创建请求失败: %v", err)
		return nil
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 发送请求
	response, err := m.client.Do(req)
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

	collection := make([]*model.Entry, 0, len(resp.Data))
	for _, kv := range resp.Data {
		id, feas := uCtx.Items.GetByKey(kv.Key)
		if id < 0 {
			continue
		}

		entry := &model.Entry{
			ID:       id,
			KeyScore: kv,
			Runtime:  *model.NewRuntime(feas),
		}
		entry.AddChan(m.conf.Name, fmt.Sprintf("recall by model: %s", m.conf.URL))
		collection = append(collection, entry)
	}

	return collection
}
