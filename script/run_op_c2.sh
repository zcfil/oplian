#!/bin/bash

export GPU_NO=$1
export PORT=$2
export MainDisk=$3
export FFI_USE_CUDA=1
export FFI_BUILD_FROM_SOURCE=1
export FIL_PROOFS_MAXIMIZE_CACHING=1
export FIL_PROOFS_USE_GPU_COLUMN_BUILDER=1
export FIL_PROOFS_USE_GPU_TREE_BUILDER=1
export FIL_PROOFS_USE_MULTICORE_SDR=1
export BELLMAN_CPU_UTILIZATION=0
export FIL_PROOFS_USE_MULTICORE_SDR=0
export RUST_LOG=info
export FFI_USE_CUDA_SUPRASEAL=1
export FIL_PROOFS_PARAMETER_CACHE=${MainDisk}/filecoin-proof-parameters
export FIL_PROOFS_PARENT_CACHE=${MainDisk}/filecoin-parents
# GPU
export PORT_C2=$2
export CUDA_VISIBLE_DEVICES=${GPU_NO}
export TMPDIR=/tmp/gpu${GPU_NO}
export TRUST_PARAMS=1

nohup /root/oplian/oplian-sectors-c2 run-c2 ${PORT} > ${MainDisk}/ipfs/logs/oplian-sectors-c2-${PORT}.log 2>&1 &
