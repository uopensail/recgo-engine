import sys
import random
import json

 

def gen_record(item_id,max_count):
    return ("%s\t%s"  %("item_id_"+str(item_id),",".join(["item_id_"+str(i) for i in random.sample(range(0, max_count),50)])))
    

def gen_items(max_count):
    ret = [gen_record(1599,max_count),gen_record(4408,max_count)]
    index_list = random.sample(range(0, max_count),int(max_count/2))
    for i in index_list:
        ret.append(gen_record(i,max_count))
    return ret

def write_file(path, items):
    with open(path,'w') as fd:
        for line in items:
            fd.write(line)
            fd.write("\n")

if __name__ == "__main__":
    dir = sys.argv[1]
    count = sys.argv[2]
    items = gen_items(int(count))
    write_file(dir +"/w2v.txt", items)