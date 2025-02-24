import sys
import random
import json

class FeatureType:
    Int64Type = 0
    Float32Type = 1
    StringType = 2
    Int64sType = 3
    Float32sType = 4
    StringsType = 5

def gen_record(i):
    scores = [1.3,2.5,4.5,0.909]
    levels = [1,2,3,4,5]
    langs = ['en', 'fr', 'cn']
    ctrys = ['us', 'uk', 'ch','jp']
    cat1s = ["cat11","cat12","cat13","cat14","cat15","cat16"]
    cat2s = ["cat21","cat22","cat23","cat24","cat25","cat26"]
    return {
        "d_s_id":{
            "type":FeatureType.StringType,
            "value":"item_id_"+str(i)
        },
        "d_s_level":{
            "type":FeatureType.Int64Type,
            "value":levels[random.randint(0, 100)%len(levels)]
        },
        "d_s_score":{
            "type":FeatureType.Float32Type,
            "value":scores[random.randint(0, 100)%len(scores)]
        },
        "d_s_language":{
            "type":FeatureType.StringType,
            "value":langs[random.randint(0, 100)%len(langs)]
        },
        "d_s_country":{
            "type":FeatureType.StringType,
            "value":ctrys[random.randint(0, 100)%len(ctrys)]
        },
        "d_s_cat1":{
            "type":FeatureType.StringType,
            "value":cat1s[random.randint(0, 100)%len(cat1s)]
        },
        "d_s_cat2":{
            "type":FeatureType.StringType,
            "value":cat2s[random.randint(0, 100)%len(cat2s)]
        },
        "d_d_ctr":{
            "type":FeatureType.Float32Type,
            "value":random.randint(0,100)/100.0
        }
    }
    

def gen_items(count):
    ret = []
    for i in range(count):
        ret.append(gen_record(i))
    return ret

def write_file(path, items):
    with open(path,'w') as fd:
        i =0
        for item in items:
            line = json.dumps(item)
            fd.write(f"item_id_{i}\t{line}")
            fd.write("\n")
            i+=1

if __name__ == "__main__":
    path = sys.argv[1]
    count = sys.argv[2]
    items = gen_items(int(count))
    write_file(path, items)