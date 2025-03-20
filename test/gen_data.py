import sys
import random
import json
import pandas as pd
from build_invert_index import InvertIndex
from model import FeatureType
from build_subpool import gen_subpool

 

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

class Pool:
    @staticmethod
    def write_itempool_file(path, items):
        with open(path, 'w') as fd:
            i = 0
            for item in items:
                id = item["id"]["value"]
                line = json.dumps(item)
                fd.write(f"{id}\t{line}\n")
                i += 1
    def __init__(self, count, dir):
        schema = item_filed_schema()
        self.items = self.gen_items(count, schema)

        meta = self.gen_meta(schema)
        subpool = gen_subpool(self.itemdf, self.items,
            ["d_s_language='en'", "d_s_level=2 and d_d_ctr > 0.5"])
        subpool_filedata = {}
        subpoolid = 1
        for k, v in subpool.items():
            subpool_filedata[subpoolid] = ",".join(str(x) for x in v)
            subpoolid += 1

        Pool.write_itempool_file(dir + "/pool.txt", self.items)
        write_json_file(dir + "/resource.meta.json",  meta)
        write_dict_str_file(dir + "/subpool.txt", subpool_filedata)
        InvertIndex(dir,self.itemdf, schema, [
                    "d_s_language", "d_s_level",["d_s_country","d_s_cat"]])
        pass

   
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
