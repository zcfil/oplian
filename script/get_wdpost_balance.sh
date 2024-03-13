#!/bin/bash
miner=`/usr/local/sbin/lotus-miner info --hide-sectors-info | grep Miner | grep sectors | awk '{print $2}'`
/usr/local/sbin/lotus-miner  actor control list  > /tmp/lotus-miner.txt
balnce=`/usr/local/sbin/lotus-miner  actor control list | grep post | awk '{print $5}'| sed 's/\x1b\[[^\x1b]*m//g'|awk  -F. '{print $1}'`
echo "post_balnce:${balnce}"