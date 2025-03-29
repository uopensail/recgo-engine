# recgo-engine

# EngineData
## 目录结构
```txt
├── enginedata.123
│   ├── dbmodel.toml // 引擎配置文件
│   └── resources 
│       ├── invert_index // 倒排索引
│       │   ├── i_s_country|i_s_cat
│       │   ├── i_s_language
│       │   └── i_s_level
│       ├── pool.txt //物料池
│       ├── resource.meta.json //物料池meta 声明用户和物料的属性定义
│       ├── scores
│       │   └── ctr //用于排序的分值？
│       └── subpool.txt //物料子集池
├── enginedata.123.success  //目录准备好的标记
├── enginedata.124
```

### 倒排样例
i_s_country|i_s_cat
```txt
0:ch|1:cat1\titem_id_4,item_id_6159
0:ch|1:cat2\titem_id_3,item_id_8197
```
i_s_language
```txt
0|cn\titem_id_0,item_id_8193,item_id_8194,item_id_3,
```

## 物料池
```txt
item_id_0\t{"id": {"type": 2, "value": "item_id_0"}, "i_s_level": {"type": 0, "value": 1}, "i_s_score": {"type": 1, "value": 4.5}, "i_s_language": {"type": 2, "value": "cn"}, "i_s_country": {"type": 2, "value": "us"}, "i_s_cat1": {"type": 2, "value": "cat14"}, "i_s_cat2": {"type": 5, "value": ["cat24", "cat26"]}, "i_d_ctr": {"type": 1, "value": 0.71}, "i_s_cat": {"type": 2, "value": "cat1"}}
```

## 资源meta
resource.meta.json

```json
{
    "field_data_type": {
        "u_s_country": 2,
        "u_s_cat": 5,
        "u_s_language": 2,
        "id": 2,
        "i_s_level": 0,
        "i_s_score": 1,
        "i_s_language": 2,
        "i_s_country": 2,
        "i_s_cat1": 2,
        "i_s_cat2": 5,
        "i_d_ctr": 1,
        "i_s_cat": 2
    }
}
```

## 子物料
通常是用户配置定向过滤条件率先出来的子物料，这里的1是一个id
subpool.txt
```txt
1\titem_id_8,item_id_9,item_id_12,item_id_16,item_id_17,item_id_21,item_id_25,item_id_35,item_id_36,item_id_38,item_id_39,item_id_43,item_id_44,item_id_47,item_id_51,item_id_54,item_id_56,item_id_63,item_id_67,item_id_69,item_id_73,item_id_75,item_id_82,item_id_83,item_id_87,item_id_88,item_id_90,item_id_91,item_id_93,item_id_95,item_id_96,item_id_98,item_id_99,item_id_102,item_id_103,item_id_104,item_id_109,item_id_110,item_id_115,item_id_118,item_id_120,item_id_124,item_id_125,item_id_128,item_id_132,item_id_139,item_id_146,item_id_147,item_id_148,item_id_151,item_id_153,item_id_161,item_id_162,item_id_163,item_id_169,item_id_171,item_id_175,item_id_183,item_id_185,item_id_191,item_id_192,item_id_193,item_id_195,item_id_198,item_id_199,item_id_200,item_id_204,item_id_207,item_id_210,item_id_212,item_id_216,item_id_217,item_id_218
2\titem_id_8,item_id_9,item_id_12,item_id_16,item_id_17,item_id_21,item_id_25,item_id_35,item_id_36,item_id_38,item_id_39,item_id_43,item_id_44,item_id_47,item_id_51,item_id_54,item_id_56,item_id_63,item_id_67,item_id_69,item_id_73,item_id_75,item_id_82,item_id_83,item_id_87,item_id_88,item_id_90,item_id_91,item_id_93,item_id_95,item_id_96,item_id_98,item_id_99,item_id_102,item_id_103,item_id_104,item_id_109,item_id_110,item_id_115,item_id_118,item_id_120,item_id_124,item_id_125,item_id_128,item_id_132,item_id_139,item_id_146,item_id_147,item_id_148,item_id_151,item_id_153,item_id_161,item_id_162,item_id_163,item_id_169,item_id_171,item_id_175,item_id_183,item_id_185,item_id_191,item_id_192,item_id_193,item_id_195,item_id_198,item_id_199,item_id_200,item_id_204,item_id_207,item_id_210,item_id_212,item_id_216,item_id_217,item_id_218
```