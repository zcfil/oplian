#!/bin/bash

MainDisk=$1
startNo=$2
endNo=$3

export RUST_BACKTRACE=full
export RUST_LOG=info
export FIL_PROOFS_PARENT_CACHE=\${MainDisk}/filecoin-parents/
export IPFS_GATEWAY="https://proof-parameters.s3.cn-south-1.jdcloud-oss.com/ipfs/"
export FIL_PROOFS_PARAMETER_CACHE=\${MainDisk}/filecoin-proof-parameters
export TRUST_PARAMS_FORCE=1
export FIL_PROOFS_USE_GPU_TREE_BUILDER=1
export FIL_PROOFS_USE_MULTICORE_SDR=1
export P1_CORES=1
export FIL_PROOFS_MAXIMIZE_CACHING=1
export FIL_PROOFS_USE_GPU_COLUMN_BUILDER=1
export RUST_LOG=info

export HOST_WORKER_COUNT=1
numactl -C ${startNo}-${endNo} nohup /root/oplian/bin/worker-p2 > ${MainDisk}/ipfs/logs/worker-p2.log 2>&1 &