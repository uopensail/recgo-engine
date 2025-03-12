DATA_DIR=/tmp/sunmao/resources/
ITEM_COUNT=10000
mkdir -pv ${DATA_DIR}
python3 gen_filter.py "redis://127.0.0.1:6379/1" ${ITEM_COUNT}
python3 gen_pool.py ${DATA_DIR}/pool.json ${ITEM_COUNT}
python3 gen_cat_hotest.py ${DATA_DIR}/cat_hotest.txt ${ITEM_COUNT}
python3 gen_w2v.py ${DATA_DIR}/w2v.txt ${ITEM_COUNT}