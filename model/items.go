package model

import (
	"bufio"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/bytedance/sonic"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/sample"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type Items struct {
	arena     *sample.Arena
	dict      map[string]int
	array     []*sample.ImmutableFeatures
	filePath  string
	updatTime int64
}

func NewItems(filepath string) (Resource, error) {
	stat := prome.NewStat("NewItems")
	defer stat.End()
	file, err := os.Open(filepath)
	if err != nil {
		zlog.LOG.Error("failed to open file", zap.Error(err))
		stat.MarkErr()
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	items := &Items{
		arena:    sample.NewArena(),
		array:    make([]*sample.ImmutableFeatures, 0, 4096),
		dict:     make(map[string]int, 4096),
		filePath: filepath,
	}

	index := 0
	for scanner.Scan() {
		line := scanner.Text()
		ss := strings.Split(line, "\t")
		if len(ss) != 2 {
			zlog.LOG.Warn("ingore line", zap.String("line", line))
			continue
		}

		feas := sample.NewImmutableFeatures(items.arena)
		err = sonic.Unmarshal(unsafe.Slice(unsafe.StringData(ss[1]), len(ss[1])), feas)
		if err != nil {
			zlog.LOG.Error("unmarshal immutableFeatures error", zap.String("data", line))
			continue
		}
		items.array = append(items.array, feas)
		items.dict[ss[0]] = index
		index++
	}

	if err := scanner.Err(); err != nil {
		zlog.LOG.Error("error while scanning file", zap.Error(err))
		stat.MarkErr()
		return nil, err
	}

	stat.SetCounter(index)
	items.updatTime = time.Now().Unix()
	return items, nil
}

func (items *Items) GetByKey(key string) (int, *sample.ImmutableFeatures) {
	if id, ok := items.dict[key]; ok {
		return id, items.array[id]
	}
	return -1, nil
}

func (items *Items) GetByID(id int) *sample.ImmutableFeatures {
	if 0 <= id && id < len(items.array) {
		return items.array[id]
	}
	return nil
}

func (items *Items) GetUpdateTime() (int64, error) {
	return items.updatTime, nil
}

func (items *Items) GetURL() (string, error) {
	return items.filePath, nil
}
