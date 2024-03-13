#!/bin/bash

MainDisk=$1
startNo=$2
endNo=$3
proof=$4
cache=$5
unsealed=$6
sealed=$7
number=$8
miner=$9
ticket=${10}
pieces=${11}
producers=${12}
export FIL_PROOFS_MULTICORE_SDR_PRODUCERS=${producers}
export RUST_BACKTRACE=full
export RUST_LOG=info
export FIL_PROOFS_PARENT_CACHE=${MainDisk}/filecoin-parents/
export IPFS_GATEWAY="https://proof-parameters.s3.cn-south-1.jdcloud-oss.com/ipfs/"
export FIL_PROOFS_PARAMETER_CACHE=${MainDisk}/filecoin-proof-parameters
export TRUST_PARAMS_FORCE=1
export FIL_PROOFS_USE_GPU_TREE_BUILDER=1
export FIL_PROOFS_USE_MULTICORE_SDR=1
export FIL_PROOFS_MAXIMIZE_CACHING=1
export FIL_PROOFS_USE_GPU_COLUMN_BUILDER=1
export RUST_LOG=info

numactl -C ${startNo}-${endNo} nohup /root/oplian/bin/worker-p1 run --proof-type ${proof} --cache ${cache} --unsealed ${unsealed} --sealed ${sealed} --number ${number} --miner ${miner} --ticket ${ticket} --pieces ${pieces} >> ${MainDisk}/ipfs/logs/${endNo}worker-p1.log 2>&1 &