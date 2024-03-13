#!/bin/bash
#lotus高度
lotus=$(/usr/local/sbin/lotus  sync  status  |grep 'Height diff' |uniq |awk  '{print $NF}'|sort -nr |head -n1)
if [[ ${lotus} -ge 2 ]];then
echo false
else
  echo true
fi