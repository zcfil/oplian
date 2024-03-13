#!/bin/bash
#节点号
MainDisk=$1
miner=`/usr/local/sbin/lotus-miner info --hide-sectors-info | grep Miner | grep sectors | awk '{print $2}'`
news=$(cat ${Minesize}/ipfs/log*/*miner.log | grep "failed: exit 16" | grep "Submitting window post" | awk '{print $8}' | wc -l)
messages=$(cat ${Minesize}/ipfs/log*/*miner.log | grep "failed: exit 16" | grep "Submitting window post" | awk '{print $8}')
if [[ ${news} -ge 1 ]]; then
echo false
else
echo true
fi