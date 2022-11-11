# dos2unix start.sh

ps -ef | grep main

openssl dgst -sha256 -binary 3.jpg | base64

mkdir -p 1/objects/ 2/objects/ 3/objects/ 4/objects/ 5/objects/ 6/objects/

rm -rf 1/temp/* 2/temp/* 3/temp/* 4/temp/* 5/temp/* 6/temp/* 

rm -rf 1/objects/* 2/objects/* 3/objects/* 4/objects/* 5/objects/* 6/objects/*

curl -v 127.0.0.1:6060/OSS/objects/test5 -XPUT -d "this object will be separate to 4+2 shards" -H "Digest: SHA-256=MBMxWHrPMsuOBaVYHkwScZQRyTRMQyiKp2oelpLZza8="

curl  127.0.0.1:6060/OSS/objects/test5?version=1

curl  127.0.0.1:6060/OSS/locate/MBMxWHrPMsuOBaVYHkwScZQRyTRMQyiKp2oelpLZza8=


