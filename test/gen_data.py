import sys
import random
import json
import pandas as pd
from pandasql import sqldf
import itertools


class FeatureType:
    Int64Type = 0
    Float32Type = 1
    StringType = 2
    Int64sType = 3
    Float32sType = 4
    StringsType = 5


def item_filed_schema():
    scores = [1.3, 2.5, 4.5, 0.909]
    levels = [1, 2, 3, 4, 5]
    langs = ['en', 'fr', 'cn']
    ctrys = ['us', 'uk', 'ch', 'jp']
    cats = ["cat1", "cat2", "cat3", "cat4", "cat5", "cat6"]
    cat1s = ["cat11", "cat12", "cat13", "cat14", "cat15", "cat16"]
    cat2s = ["cat21", "cat22", "cat23", "cat24", "cat25", "cat26"]
    meta = {
        "id": (FeatureType.StringType, lambda i:   "item_id_"+str(i)),
        "d_s_level": (FeatureType.Int64Type, lambda i:  levels[random.randint(0, 100) % len(levels)]),
        "d_s_score": (FeatureType.Float32Type, lambda i:  scores[random.randint(0, 100) % len(scores)]),
        "d_s_language": (FeatureType.StringType, lambda i:  langs[random.randint(0, 100) % len(langs)]),
        "d_s_country": (FeatureType.StringType, lambda i:  ctrys[random.randint(0, 100) % len(ctrys)]),
        "d_s_cat1": (FeatureType.StringType, lambda i:  cat1s[random.randint(0, 100) % len(cat1s)]),
        "d_s_cat2": (FeatureType.StringsType, lambda i:  random.sample(cat2s, k=2)),
        "d_d_ctr": (FeatureType.Float32Type, lambda i:  random.randint(0, 100)/100.0),
        "d_s_cat": (FeatureType.StringType, lambda i:  cats[random.randint(0, 100) % len(cat2s)]),


    }
    return meta


def write_json_file(path, obj):
    with open(path, 'w') as fd:
        c = json.dumps(obj)
        fd.write(f"{c}")


def write_arrayobj_file(path, items):
    with open(path, 'w') as fd:
        i = 0
        for item in items:
            line = json.dumps(item)
            fd.write(f"{line}\n")
            i += 1


def write_dict_str_file(path, itemdict):
    with open(path, 'w') as fd:
        for k, item in itemdict.items():
            line = item
            fd.write(f"{k}\t{line}\n")


class InvertIndex:
    def __init__(self, itemdf, schema, fileds_list):
        self.itemdf = itemdf
        invert_index = {}
        for field_name in fileds_list:
            filed_invert_index = {}
            if (type(field_name) is str):
                file_name_key = field_name
                ret = invert_index_with_items(
                    itemdf, [field_name], schema)
                for key_t,v in ret.items():
                    key = ','.join(str(x) for x in key_t)
                    filed_invert_index[key] = v
            elif (type(field_name) is list):
                file_name_key = '|'.join(str(x) for x in field_name)
                ret = invert_index_with_items(
                    itemdf, field_name, schema)
                for key_t,v in ret.items():
                    key = ','.join(str(x) for x in key_t)
                    filed_invert_index[key] = v
            else:
                raise "not suppport"
            invert_index[file_name_key] = filed_invert_index
            

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


class Pool:
    def __init__(self, count, dir):
        schema = item_filed_schema()
        self.items = self.gen_items(count, schema)

        meta = self.gen_meta(schema)
        subpool = self.gen_subpool(
            ["d_s_language='en'", "d_s_level=2 and d_d_ctr > 0.5"])
        subpool_filedata = {}
        subpoolid = 1
        for k, v in subpool.items():
            subpool_filedata[subpoolid] = ",".join(str(x) for x in v)
            subpoolid += 1

        write_arrayobj_file(dir + "/pool.txt", self.items)
        write_json_file(dir + "/pool.meta",  meta)
        write_dict_str_file(dir + "/subpool.txt", subpool_filedata)
        InvertIndex(self.itemdf, schema, [
                    "d_s_language", "d_s_level"])
        pass

    @staticmethod
    def build_sub_collection(items, subitems):
        ret = []
        for subitem in subitems:
            for i in range(len(items)):

                if subitem['id'] == items[i]['id']['value']:
                    # print("xxx",subitem, items[i])
                    ret.append(subitem['id'])
                    break
        return ret

    def gen_subpool(self, conditons):
        # 转换为 DataFrame
        ret = {}
        self.itemdf["d_s_cat2"] = self.itemdf["d_s_cat2"].apply(lambda x: ",".join(x) if isinstance(x, list) else x)
        for condition in conditons:
            sql = "SELECT * FROM itemdf WHERE " + condition

            result = sqldf(sql, {"itemdf": self.itemdf})

            subitems = result.to_dict(orient="records")
            ret[condition] = Pool.build_sub_collection(self.items, subitems)
        return ret

    def gen_items(self, count, schema):
        ret = []
        dfdata = []
        for i in range(count):
            record = {}
            item = {}
            for k, v in schema.items():
                item[k] = v[1](i)
                record[k] = {
                    "type": v[0],
                    "value": v[1](i)
                }
            ret.append(record)
            dfdata.append(item)
        df = pd.DataFrame(dfdata)
        self.itemdf = df
        return ret

    def gen_meta(self, schema):
        ret = {}
        for k, v in schema.items():
            ret[k] = v[0]
        return ret


if __name__ == "__main__":
    dir = sys.argv[1]
    count = sys.argv[2]
    pool = Pool(int(count), dir)
