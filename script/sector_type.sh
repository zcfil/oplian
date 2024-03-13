#!/bin/bash
MainDisk=$2
source ${MainDisk}/ipfs/conf/.lotusprofile 2> /dev/null
sector=$1
deals=$(lotus-miner sectors status $sector 2> /dev/null |  grep 'Deals:\s*\[[0-9]*\]' | awk '{print $2}' | tr -d '[]')
if [[ $deals -gt 0 ]]; then
  echo "DC"
else
  echo "CC"
fi