# honghu(鸿鹄)
鸿鹄-推荐引擎

引擎的核心是通过用户信息和物料信息的撮合，给用户一定数量的推荐物料候选集，用户信息和物料信息都是在内存中都是统一格式的tfrecord,这点需要重点注意。

## userctx 目录
    userctx：在推荐请求的最开始就要初始化
    user_ab: 存储用户的命中实验的ab信息

    user_tf： 特征信息（UserFeatures），从用户特征服务器获取

    user_freq：频控规则，从特征服务器里获取最近点击列表，然后map缓存起来，召回的时候会通过CheckFilter间接调用
    
    CheckFreq
    过滤：CheckFilter(source *datasource.Item) bool: 过滤规则（频控、分类过滤），比如某个国家、语言的过滤，在召回的时候会调用
    
    例如国家过滤如何做：
            就是拿到用户特征u_s_country和物料特征d_s_country，然后判断tfrecord格式（int64list、floatlist、stringlist）是否存在交集，函数为source.IntersectionIsNone
## datasource 目录
  db.go：包含物料集和物料的索引集，负责定时下载物料tar.gz文件和定时加载tar.gz中的meta、物料tfrecord、物料indexes索引目录。
     items.go：物料集，每个物料用三个元素:引擎编号InsideIndex，物料主键ID，物料特征ItemTFeature（tfrecord格式）

## strategy 目录
### recall召回
    多路召回，并发执行，并且会去除重复、去调用CheckFilter进行全局性过滤、记录召回信息
### weighed 调权
    可以对不同召回进行调权
### layout 排版
    在没有排序的情况下，通过排版进行不同召回按固定位置进行排版
### rank 排序
     打分排序

### api.proto 冲突问题

    因为本服务提供grpc接口，所以服务本身有个api.proto , 同时，fuku的client也有个api.proto文件，命名冲突，编译时使用 go build -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn" 可解决