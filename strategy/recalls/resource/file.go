package resource

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/model/dbmodel/table"

	"github.com/uopensail/ulib/datastruct"
	"github.com/uopensail/ulib/finder"
	"github.com/uopensail/ulib/pool"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/utils"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

type fileIndeces struct {
	updateTime int64
	indexes    map[string][]int64 //索引集合cache
}

type FileResource struct {
	cfg table.RecallResourceMeta
	*fileIndeces
}

func NewFileResource(envCfg config.EnvConfig, cfg table.RecallResourceMeta, pl *pool.Pool) *FileResource {
	cfg.ParseFileSource()
	fs := &FileResource{
		cfg: cfg,
		fileIndeces: &fileIndeces{
			indexes: make(map[string][]int64),
		},
	}
	indeces := load(envCfg, cfg.FileResourceConfig, pl)
	fs.fileIndeces = indeces
	return fs
}

func getFile(envCfg config.EnvConfig, location string) string {
	if strings.HasPrefix(location, "oss://") || strings.HasPrefix(location, "s3://") {
		baseName := filepath.Base(location)
		ts := time.Now().Unix()
		localPath := filepath.Join(envCfg.WorkDir, fmt.Sprintf("%s-%d", baseName, ts))
		myFinder := finder.GetFinder(&envCfg.Finder)
		myFinder.Download(location, localPath)
		return localPath
	} else {
		return location
	}

}

func load(envCfg config.EnvConfig, meta table.FileResourceConfig, pl *pool.Pool) *fileIndeces {
	//load (key, []string).txt
	stat := prome.NewStat("FileResource.load")
	defer stat.End()
	targetFilePath := getFile(envCfg, meta.Location)
	file, err := os.Open(targetFilePath)
	if err != nil {
		zlog.LOG.Error("failed to open file", zap.Error(err))
		stat.MarkErr()
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	indeces := fileIndeces{
		indexes: make(map[string][]int64),
	}
	for scanner.Scan() {
		line := scanner.Text()
		vvs := strings.Split(line, "\t")
		if len(vvs) < 2 {
			continue
		}

		key := vvs[0]
		vv := utils.String2Int64List(vvs[1], ",")
		indeces.indexes[key] = vv
	}
	stat.SetCounter(len(indeces.indexes))
	return &indeces
}

func (res *FileResource) Get(keys []string, pl *pool.Pool) [][]datastruct.Tuple[int, float32] {
	ret := make([][]datastruct.Tuple[int, float32], len(keys))
	for i := 0; i < len(keys); i++ {
		v := res.indexes[keys[i]]

		ret[i] = make([]datastruct.Tuple[int, float32], len(v))
		for j := 0; j < len(v); j++ {
			ret[i][j] = datastruct.Tuple[int, float32]{
				First:  int(v[j]),
				Second: 1.0,
			}
		}
	}
	return ret
}

func (res *FileResource) CheckResourceUpdate(envCfg config.EnvConfig, poolUpdate bool) bool {
	if poolUpdate {
		return true
	}
	myFinder := finder.GetFinder(&envCfg.Finder)
	nUpdateTime := myFinder.GetUpdateTime(res.cfg.Location)
	if res.updateTime < nUpdateTime {
		return true
	}
	return false
}

func (res *FileResource) Meta() *table.RecallResourceMeta {
	return &res.cfg
}

func (res *FileResource) Close() {

}
