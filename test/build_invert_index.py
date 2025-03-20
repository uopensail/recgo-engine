import pandas as pd

import itertools
import os
from model import FeatureType

def invert_index_with_items(dataframe: pd.DataFrame, 
                           field_name_list: list, 
                           meta: dict) -> dict:
    """
    返回笛卡尔积组合对应的完整数据项
    
    返回结构:
    {
        (val1, val2): [
            {"field1": val1, "field2": val2, ...},  # 完整数据项
            ...
        ]
    }
    """
    
    # 构建基础倒排索引（存储原始索引）
    field_indices = {}
    for field in field_name_list:
        field_type = meta[field][0]
        
        if field_type == FeatureType.StringsType:
            exploded = dataframe[[field]].explode(field)
            field_index = exploded.groupby(field).apply(lambda x: set(x.index))
        elif field_type == FeatureType.FloatType or field_type == FeatureType.FloatsType:
            raise "not suppport float type build invert_index"
        else:
            field_index = dataframe.groupby(field).apply(lambda x: set(x.index))
        
        field_indices[field] = field_index.to_dict()

    # 生成笛卡尔积组合
    value_combinations = itertools.product(
        *[field_indices[field].keys() for field in field_name_list]
    )

    # 构建最终结果字典
    result = {}
    for combination in value_combinations:
        indices = None
        for i, field in enumerate(field_name_list):
            current = field_indices[field].get(combination[i], set())
            indices = current if indices is None else indices & current
            if not indices:
                break
        
        if indices:
            # 获取完整数据项并转换为字典列表
            items = dataframe.loc[list(indices)].to_dict('records')
            # 过滤多值字段的无效组合
            filtered = [
                item for item in items
                if all(
                    (combination[i] in (item[field] if isinstance(item[field], list) else [item[field]]))
                    for i, field in enumerate(field_name_list)
                )
            ]
            if filtered:
                result[combination] = filtered
    
    return result


class InvertIndex:
    def __init__(self, workdir, itemdf, schema, fileds_list):
        self.itemdf = itemdf
        invert_index = {}
        for field_name in fileds_list:
            filed_invert_index = {}
            if (type(field_name) is str):
                file_name_key = field_name
                ret = invert_index_with_items(
                    itemdf, [field_name], schema)
                for key_t,v in ret.items():
                    filed_invert_index[key_t] = v
            elif (type(field_name) is list):
                file_name_key = '|'.join(str(x) for x in field_name)
                ret = invert_index_with_items(
                    itemdf, field_name, schema)
                for key_t,v in ret.items():

                    filed_invert_index[key_t] = v
            else:
                raise "not suppport"
            invert_index[file_name_key] = filed_invert_index
        InvertIndex.write_file(workdir, invert_index)

    @staticmethod
    def write_file(workdir, invert_indexes):
        invert_index_dir = workdir+"/invert_index"
        os.makedirs(invert_index_dir, exist_ok=True)
        for file_name_key, list_v in  invert_indexes.items():
            field_name = file_name_key.split("|")
            fullpath = invert_index_dir + "/" +file_name_key
            with open(fullpath, 'w') as fd:
                for k, vv in list_v.items():
                    line_value = ",".join(str(x['id']) for x in vv)
                    if len(field_name) > 1:
                        line_key = ""
                        for i in range(len(field_name)):
                            if (i != 0):
                                line_key+="|"
                            line_key += str(i)
                            line_key += ":"
                            line_key += k[i]

                    else:
                        line_key = f"0|{k[0]}"
                    fd.write(f"{line_key}\t{line_value}\n")



