import sys
import random
import redis
import time

from urllib.parse import urlparse

def write_filter_data(url,user_id,items): 
    redis_client = redis.StrictRedis.from_url(url, decode_responses=True)

    for item in items:
        redis_client.rpush(user_id, item)
    redis_client.expire(user_id, 30*86400)

def gen_items(count, max_id):
    items = []
    scenes = ["foru","relate","cat11"]
    actions = ["imp","click","imp_v","play"]
    for i in range(count):
        item_id = "item_id_"+str(random.randint(0,max_id))
        ts = time.time()- random.randint(0, 7*86400)
        scene = scenes[random.randint(0,len(scenes)-1)]
        action = actions[random.randint(0,len(actions)-1)]
        items.append("%s|%d|%s|%s" %(item_id, ts, scene, action))

    return items

if __name__ == "__main__":
    url = sys.argv[1]
    count = sys.argv[2]
    items = gen_items(200,int(count))
    write_filter_data(url,"test_user_id",items)