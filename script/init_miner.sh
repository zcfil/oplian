#!/bin/bash
##钱包地址
Walletas=$1
Minesize=$2
MainDisk=$3
Dir='/tmp/MineIn.log'

#########################测试初始旷工位置##############################
export PATH=$PATH
datetime=`date +'%Y-%m-%d %H:%M:%S'`
export RUST_BACKTRACE=full
export RUSTFLAGS="-C target-cpu=native -g"
export FFI_BUILD_FROM_SOURCE=1
export RUST_LOG=info
export LOTUS_PATH="${MainDisk}/ipfs/data/lotus"
export LOTUS_MINER_PATH="${MainDisk}/ipfs/data/lotusminer"
export LOTUS_BACKUP_BASE_PATH="${MainDisk}/test/"
export IPFS_GATEWAY="https://proof-parameters.s3.cn-south-1.jdcloud-oss.com/ipfs/"
export FIL_PROOFS_PARAMETER_CACHE="${MainDisk}/filecoin-proof-parameters"

export FIL_PROOFS_MAXIMIZE_CACHING=1
export TRUST_PARAMS_FORCE=1
export SKIP_BASE_EXP_CACHE=1
export GOLOG_LOG_FMT=json
export FIL_PROOFS_USE_MULTICORE_SDR=0
export TRUST_PARAMS=1
export path="${MainDisk}/test"

export FULLNODE_API_INFO=$4
export LISTEN_IP=$5
export LISTEN_PORT=$6

#########################测试初始旷工位置##############################
#判断lotus 状态
#timeout 120   /root/oplian/bin/lotus  sync wait > /dev/null
#if [[ $? -ne 0 ]];then
#echo 'lotus_failed'
#exit
#   else
#echo "ok" > /dev/null
#fi
#wait
/root/oplian/bin/lotus-miner init --owner=${Walletas}  --sector-size=${Minesize}GiB     > ${Dir}     2>&1   #初始化一个新的矿工，lotus上可以查看
wait
sleep  1
MineIt=`cat  ${Dir} |grep   'Created new miner'  |awk  '{print $NF}'|egrep -o  [a-Z].*[0-9]`
#输出旷工ID
echo  ${Walletas}
echo  ${Minesize}
echo  ${MineIt}
