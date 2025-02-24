import sys
import random
import json

def gen_record(i):
    langs = ['en', 'fr', 'cn']
    ctrys = ['us', 'uk', 'ch','jp']
    cats = [1,2,3,4,5,6,7]
    return {
        "item_id":{
            "type":2,
            "value":"item_id_"+str(i)
        },
        "lang":{
            "type":2,
            "value":langs[random.randint(0, 100)%len(langs)]
        },
        "coutry":{
            "type":2,
            "value":ctrys[random.randint(0, 100)%len(ctrys)]
        },
        "cat":{
            "type":0,
            "value":cats[random.randint(0, 100)%len(cats)]
        }
    }
    

def gen_items(count):
    ret = []
    for i in range(count):
        ret.append(gen_record(i))
    return ret

def write_file(path, items):
    with open(path,'w') as fd:
        for item in items:
            line = json.dumps(item)
            fd.write(line)
            fd.write("\n")

if __name__ == "__main__":
    path = sys.argv[1]
    count = sys.argv[2]
    items = gen_items(int(count))
    write_file(path, items)