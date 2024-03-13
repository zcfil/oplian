#!/bin/bash

#节点worker钱包
worker=$1
full=FULLNODE_API_INFO=$2
miner=MINER_API_INFO=$3
#获取 lotus token
export $full
# 获取miner  token
export $miner
# 获取boostd 专用miner token
export APISEALER=$miner
#获取boostd 专用miner token
export APISECTORINDEX=$miner
export LISTEN_IP=$4
export LISTEN_PORT=$5
export MainDisk=$6
mkdir -p  $MainDisk/ipfs/data/boost
#boost地址
export BOOST_PATH=$MainDisk/ipfs/data/boost

# 初始化boostd
/root/oplian/bin/boostd  --vv init \
       --api-sealer=$APISEALER \
       --api-sector-index=$APISECTORINDEX \
       --wallet-publish-storage-deals=$worker\
       --wallet-deal-collateral=$worker\
       --max-staging-deals-bytes=50000000000


wait
cat >  /root/oplian/script/run_boost.sh << eof
#!/usr/bin/env bash
set -e
datetime=\$(date +'%Y-%m-%d %H:%M:%S')
export MainDisk=$6
export RUST_LOG=info
export IP=\$(hostname -I | awk '{print \$NF}')
export LOTUS_PATH=\${MainDisk}/ipfs/data/lotus
export LOTUS_MINER_PATH=\${MainDisk}/ipfs/data/lotusminer
export IPFS_GATEWAY="https://proof-parameters.s3.cn-south-1.jdcloud-oss.com/ipfs/"
export FIL_PROOFS_PARAMETER_CACHE=\${MainDisk}/filecoin-proof-parameters
export FIL_PROOFS_PARENT_CACHE=\${MainDisk}/filecoin-parents/
#机械硬盘打开下面两个
export REMOVE="true"
export SSD="true"
export REMOVE="true"
export LISTEN_IP=$4
export LISTEN_PORT=$5

sysctl -w net.core.rmem_max=2500000
pid=\$(ps -aux | grep 'boostd --vv run'|grep -v grep | awk '{print \$2}')
if [ x"\$pid" = "x" ]; then
        export $full
        export $miner
        export APISEALER=$miner
        export APISECTORINDEX=$miner
        export BOOST_PATH=\${MainDisk}/ipfs/data/boost
        export BOOST_BIN="/root/oplian/bin/boostd"
                  nohup \${BOOST_BIN} --vv run >> \${MainDisk}/ipfs/logs/boost.log &
                              sudo prlimit --nofile=1048576 --nproc=unlimited --rtprio=99 --nice=-19 --pid \$!
                                  else
                                                        echo "\${datetime}  lotus RUNNING  successfully!"
fi
wait
eof
chmod +x  /root/oplian/script/run_boost.sh

cat >  /etc/supervisor/conf.d/boost.conf  <<EOF
[program:boost]
command=/root/oplian/script/run_boost.sh
user=root
autostart=true
autorestart=true
stopwaitsecs=60
startretries=100
stopasgroup=true
killasgroup=true
stdout_logfile_maxbytes = 10000000MB
stdout_logfile_backups=25
redirect_stderr=true
stdout_logfile=$6/ipfs/logs/boost.log
EOF

wait
supervisorctl   update