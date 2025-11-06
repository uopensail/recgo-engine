package rank

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/recgo-engine/userctx"
	"github.com/uopensail/ulib/sample"
)

type Request struct {
	Features   sample.Features `json:"features"`
	Collection []string        `json:"collection"`
}

type Response struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data []model.KeyScore `json:"data"`
}

type Modeler struct {
	conf   *model.ModelBasedRankConfigure
	client *http.Client
}

func NewModeler(conf *model.ModelBasedRankConfigure) *Modeler {
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

func (m *Modeler) Do(uCtx *userctx.UserContext, collection model.Collection) model.Collection {
	arr := make([]string, 0, len(collection))
	for _, entry := range collection {
		arr = append(arr, entry.Key)
	}
	reqest := Request{
		Features:   uCtx.Features,
		Collection: arr,
	}
	data, err := json.Marshal(reqest)
	if err != nil {
		// fmt.Errorf("JSON序列化失败: %v", err)
		return collection
	}

	// 创建请求
	req, err := http.NewRequest("POST", m.conf.URL, bytes.NewBuffer(data))
	if err != nil {
		// fmt.Errorf("创建请求失败: %v", err)
		return collection
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 发送请求
	response, err := m.client.Do(req)
	if err != nil {
		// fmt.Errorf("发送请求失败: %v", err)
		return collection
	}

	// 读取响应体
	body, err := io.ReadAll(response.Body)
	if err != nil {
		// fmt.Errorf("读取响应失败: %v", err)
		return collection
	}

	// 解析JSON响应
	resp := Response{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return collection
	}

	dict := make(map[string]float32, len(resp.Data))
	for _, kv := range resp.Data {
		dict[kv.Key] = kv.Score
	}

	for _, entry := range collection {
		if score, ok := dict[entry.Key]; ok {
			entry.Score = score
		} else {
			entry.Score = 0
		}
	}
	sort.Stable(collection)
	return collection
}
