#!/bin/bash
mount -l | grep 10.0. | awk '{print $3}' |grep -v md0| awk -F "/" '{print $3}' | uniq > /tmp/ip.txt
ip=`cat /tmp/ip.txt`
touch /tmp/true
for i in $ip
do
    ping -c 3 -i 0.01 -W 3 $i &> /dev/null
    if [ $? -eq 0 ]
    then
        echo "$i:true" >> /tmp/true
    else
        echo "$i:false"
    fi
done
    file="/tmp/true"
    size=$(wc -c < "$file")
    if [[ $size -eq 0 ]];then
#        echo "未挂载"
        echo "true"
    else
    cat /tmp/true | awk -F: '{print $2}' | uniq | xargs echo
    fi
> /tmp/true