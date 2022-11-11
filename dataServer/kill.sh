#!/bin/bash
echo "kill exe start "
int=1
while(( $int<=6 ))
do
    echo $int
    go=`ps -ef | grep 'go run main.go -f ./conf/'$int'.json' | xargs | cut -b 5-12`
    echo $go
    kill $go 

    gomain=`ps -ef | grep '/exe/main -f ./conf/'$int'.json' | xargs | cut -b 5-12`
    echo $gomain
    kill $gomain 

    let "int++"
done
echo "kill exe end "
