#!/bin/bash
MainDisk=$1
time=$(cat ${Minesize}/ipfs/log*/miner.log | grep "总耗时" | awk -F : '{print $5}' | awk -F "." '{print $1}' | tail -n 1)
if [[ ${time} -ge 30 ]]; then
echo true
else
echo false
fi