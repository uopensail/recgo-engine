

from pandasql import sqldf

def build_sub_collection(items, subitems):
    ret = []
    for subitem in subitems:
        for i in range(len(items)):

            if subitem['id'] == items[i]['id']['value']:
                # print("xxx",subitem, items[i])
                ret.append(subitem['id'])
                break
    return ret

def gen_subpool(itemdf, items, conditons):
    # 转换为 DataFrame
    ret = {}
    itemdf["d_s_cat2"] = itemdf["d_s_cat2"].apply(lambda x: ",".join(x) if isinstance(x, list) else x)
    for condition in conditons:
        sql = "SELECT * FROM itemdf WHERE " + condition

        result = sqldf(sql, {"itemdf": itemdf})

        subitems = result.to_dict(orient="records")
        ret[condition] = build_sub_collection(items, subitems)
    return ret
