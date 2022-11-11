#!/bin/bash
echo 'run start ...'
int=1
while (( $int<=6 ))
do
    nohup go run main.go -f ./conf/$int.json >logs/$int.log 2>&1 & 
    let "int++"

done
echo 'run end ...'
