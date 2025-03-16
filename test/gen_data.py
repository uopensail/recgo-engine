import sys
import random
import json
import pandas as pd
from pandasql import sqldf

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
    cat1s = ["cat11", "cat12", "cat13", "cat14", "cat15", "cat16"]
    cat2s = ["cat21", "cat22", "cat23", "cat24", "cat25", "cat26"]
    meta =  {
        "id": (FeatureType.StringType,lambda i:   "item_id_"+str(i)),
        "d_s_level": (FeatureType.Int64Type, lambda i:  levels[random.randint(0, 100) % len(levels)]),
        "d_s_score": (FeatureType.Float32Type, lambda i:  scores[random.randint(0, 100) % len(scores)]),
        "d_s_language": (FeatureType.StringType, lambda i:  langs[random.randint(0, 100) % len(langs)]),
        "d_s_country": (FeatureType.StringType, lambda i:  ctrys[random.randint(0, 100) % len(ctrys)]),
        "d_s_cat1": (FeatureType.StringType, lambda i:  cat1s[random.randint(0, 100) % len(cat1s)]),
        "d_s_cat2": (FeatureType.StringType, lambda i:  cat2s[random.randint(0, 100) % len(cat2s)]),
        "d_d_ctr": (FeatureType.Float32Type, lambda i:  random.randint(0, 100)/100.0)

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
        for k,item in itemdict.items():
            line = item
            fd.write(f"{k}\t{line}\n")

class Pool:
    def __init__(self, count, dir):
        schema = item_filed_schema()
        self.items = self.gen_items(count, schema)
      
        meta= self.gen_meta(schema)
        subpool = self.gen_subpool(["d_s_language='en'","d_s_level=2 and d_d_ctr > 0.5"])
        subpool_filedata = {}
        subpoolid =0
        for k,v in subpool.items():
            subpool_filedata[subpoolid] = ",".join(str(x) for x in v)
            subpoolid +=1

        write_arrayobj_file(dir + "/pool.txt", self.items)
        write_json_file(dir + "/pool.meta",  meta)
        write_dict_str_file(dir + "/subpool.txt", subpool_filedata)
        pass
        

    @staticmethod
    def build_sub_collection(items, subitems):
        ret = []
        for subitem in subitems:
            for i in range(len(items)):
                
                if subitem['id'] == items[i]['id']['value']:
                    #print("xxx",subitem, items[i])
                    ret.append(subitem['id'])
                    break
        return ret
    
    def gen_subpool(self,conditons):
          # 转换为 DataFrame
        ret = {}
        for condition in conditons:
            sql = "SELECT id FROM itemdf WHERE " + condition
           
            result = sqldf(sql, {"itemdf": self.itemdf})
            
            subitems = result.to_dict(orient="records")
            ret[condition] = Pool.build_sub_collection(self.items, subitems)
        return ret
 
    def gen_items(self, count, schema):
        ret = []
        dfdata = []
        for i in range(count):
            record = {}
            item ={}
            for k, v in schema.items():
                item[k] =v[1](i)
                record[k] = {
                    "type": v[0],
                    "value": v[1](i)
                }
            ret.append(record)
            dfdata.append(item)  
        df = pd.DataFrame(dfdata)
        print(df)
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