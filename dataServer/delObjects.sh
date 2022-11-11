#!/bin/bash
echo 'del objectData start ...'
rm -rf ../objectData/1/temp/* ../objectData/2/temp/* ../objectData/3/temp/* ../objectData/4/temp/* ../objectData/5/temp/* ../objectData/6/temp/* 

rm -rf ../objectData/1/objects/* ../objectData/2/objects/* ../objectData/3/objects/* ../objectData/4/objects/* ../objectData/5/objects/* ../objectData/6/objects/*

curl --location --request POST 'http://qkbyte.orginone.cn:9200/metadata/_delete_by_query' \
--header 'User-Agent: Apifox/1.0.0 (https://www.apifox.cn)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: qkbyte.orginone.cn:9200' \
--header 'Connection: keep-alive' \
--data-raw '{
  "query": {
    "match_all": {}
  }
}'

echo 'del objectData end ...'
