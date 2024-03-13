#!/bin/bash
#判断文件夹是否存在 -d
LOTUS_MINER_PATH="/mnt/md0/ipfs/data/lotusminer"
LOTUS_PATH="/mnt/md0/ipfs/data/lotus"
if [[ ! -d "$LOTUS_MINER_PATH" ]] && [[ ! -d "$LOTUS_PATH" ]]; then
 echo false
else
 echo true
fi