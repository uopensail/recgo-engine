from os import path
import os
import random


def gem_item_scores(name,items_feature):
    item_score = {}
    for item in items_feature:
        item_id = item['id']['value']
        score = random.randint(0, 100000)/100000.0
        item_score[item_id] = score
    return item_score

class ItemScores:
    def __init__(self, dir,items_feature):
        path.join(dir, 'item_scores')
        self.dir = dir
        item_scores = {}
        item_scores["ctr"] = gem_item_scores("ctr",items_feature)
        ItemScores.write_file(dir, item_scores)   
    @staticmethod
    def write_file(workdir, item_scores):
        invert_index_dir = workdir+"/scores"
        os.makedirs(invert_index_dir, exist_ok=True)
        for score_name, score_dict in  item_scores.items():
            fullpath = invert_index_dir + "/" +score_name
            with open(fullpath, 'w') as fd:
                for item_id, score in score_dict.items():
                    fd.write(f"{item_id}\t{score}\n")
        
    