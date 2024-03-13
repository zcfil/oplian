#!/bin/bash
MainDisk=$1
export FFI_USE_CUDA=1
export FFI_BUILD_FROM_SOURCE=1
export LD_LIBRARY_PATH=/usr/local/cuda-11.2/lib64
export CUDA_HOME=/usr/local/cuda
export PATH=$PATH:/usr/local/go/bin:/usr/local/cuda-11.2/bin

export FIL_PROOFS_PARAMETER_CACHE="${MainDisk}/filecoin-proof-parameters"

export FIL_PROOFS_MAXIMIZE_CACHING=1



export FIL_PROOFS_MAXIMIZE_CACHING=1
export FIL_PROOFS_USE_GPU_COLUMN_BUILDER=1
export FIL_PROOFS_USE_GPU_TREE_BUILDER=1
export FIL_PROOFS_USE_MULTICORE_SDR=1
export BELLMAN_CPU_UTILIZATION=0
export C2_COUNT=2
# GPU
export CUDA_VISIBLE_DEVICES=0
\rm  -r /tmp/gpu003
mkdir /tmp/gpu003
export TMPDIR="/tmp/gpu003"
export TRUST_PARAMS=1
#numactl -C 16-31 -m 0 nohup ./lotus-bench prove 32G-c2-input.json > bench3.log 2>&1 &
nohup ./lotus-bench prove 32G-c2-input.json > bench3.log 2>&1 &
