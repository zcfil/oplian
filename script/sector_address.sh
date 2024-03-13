#!/bin/bash
MainDisk=$3
source ${MainDisk}/ipfs/conf/.lotusprofile 2> /dev/null
MIid=$1
ida=$2

MIid1=`echo ${MIid}  |sed  's/f/s-t/'`

for i in  ${ida}
do
pwd=$(lotus-miner storage find ${i} 2> /dev/null |grep -A 3 Cache|grep Local|awk '{print $2}'|awk -F "(" '{print $2}'|awk -F ")" '{print $1}')
##echo ------------$i-------------------
echo    ${pwd}/cache/${MIid1}-${i}
done
for i in ${ida}
do
pwd=$(lotus-miner storage find ${i} 2> /dev/null |grep -A 3 Sealed|grep Local|awk '{print $2}'|awk -F "(" '{print $2}'|awk -F ")" '{print $1}')
#echo ------------$i-------------------
echo    ${pwd}/sealed/${MIid1}-${i}
done
