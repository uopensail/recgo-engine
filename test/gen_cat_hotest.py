import sys
import random
import json

 

def gen_record(key,max_count):
    return ("%s\t%s"  %(key,",".join([str(i) for i in random.sample(range(0, max_count),50)])))
    

def gen_items(max_count):
    ret = []
    field = "d_s_cat1"
    cat1s = ["cat11","cat12","cat13","cat14","cat15","cat16"]
    for cat1 in cat1s:
        
        ret.append(gen_record(field+"_"+cat1,max_count))
    return ret

def write_file(path, items):
    with open(path,'w') as fd:
        for line in items:
            fd.write(line)
            fd.write("\n")

if __name__ == "__main__":
    path = sys.argv[1]
    count = sys.argv[2]
    
    items = gen_items(int(count))
    write_file(path, items)