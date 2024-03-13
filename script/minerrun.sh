#!/bin/bash
export PATH=$PATH
cat > /root/oplian/script/run_miner.sh <<EOF
#!/bin/bash
set -e

export PLEDGE_MINER=$1
export wnpost=$2
export wdpost=$3
export PARTITIONS=$4
export FULLNODE_API_INFO=$5
export LISTEN_IP=$6
export LISTEN_PORT=$7
export ACTOR=$8
export MainDisk=$9
export LOCAL_PROVER=$10
export FIL_PROOFS_USE_GPU_TREE_BUILDER=1
export FIL_PROOFS_USE_MULTICORE_SDR=1
export FIL_PROOFS_MAXIMIZE_CACHING=1
export FIL_PROOFS_USE_GPU_COLUMN_BUILDER=1

datetime=\$(date +'%Y-%m-%d %H:%M:%S')
export RUST_BACKTRACE=full
export RUSTFLAGS="-C target-cpu=native -g"
export FFI_BUILD_FROM_SOURCE=1
export RUST_LOG=info
export LOTUS_PATH=\${MainDisk}/ipfs/data/lotus
export LOTUS_MINER_PATH=\${MainDisk}/ipfs/data/lotusminer
export IPFS_GATEWAY="https://proof-parameters.s3.cn-south-1.jdcloud-oss.com/ipfs/"
export FIL_PROOFS_PARAMETER_CACHE="\${MainDisk}/filecoin-proof-parameters"

export SKIP_BASE_EXP_CACHE=1
export TRUST_PARAMS_FORCE=1
export CHECK_UNLOCAK="true"
export REMOVE_FORCE="true"
export FINFORCE="true"

#爆快延长
export PROPAGATION_DELAY_SECS=25
IP=`hostname -I | awk '{print \$1}'`
# /usr/bin/bash  /root/oplian/script/mount_hdd.sh
nvidia-smi
if [ $? -ne 0 ]; then
  echo "ERROR: no GPU detected \$IP"
  exit
fi

nohup /root/oplian/bin/lotus-miner  run &
sudo prlimit --nofile=1048576 --nproc=unlimited --rtprio=99 --nice=-19 --pid \$!
echo "\${datetime}  lotus-miner restarted  successfully!" >> $9/ipfs/logs/restart_press.log
wait
EOF

sleep  1

chmod  +x  /root/oplian/script/run_miner.sh || echo "ERROR: No such file or directory" >&2
cat >  /etc/supervisor/supervisord.conf <<'EOF'
[unix_http_server]
file=/var/run/supervisor.sock   ; (the path to the socket file)
chmod=0700                       ; sockef file mode (default 0700)

[supervisord]
logfile=/var/log/supervisor/supervisord.log ; (main log file;default $CWD/supervisord.log)
pidfile=/var/run/supervisord.pid ; (supervisord pidfile;default supervisord.pid)
childlogdir=/var/log/supervisor            ; ('AUTO' child log dir, default $TEMP)
[rpcinterface:supervisor]
supervisor.rpcinterface_factory = supervisor.rpcinterface:make_main_rpcinterface
[supervisorctl]
serverurl=unix:///var/run/supervisor.sock ; use a unix:// URL  for a unix socket
[include]
files = /etc/supervisor/conf.d/*.conf
EOF

sleep  1

cat >  /etc/supervisor/conf.d/miner.conf  <<EOF
[program:lotus-miner]
command=/root/oplian/script/run_miner.sh
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
stdout_logfile=$9/ipfs/logs/miner.log
EOF

ls  /etc/supervisor/conf.d/miner.conf || echo "ERROR: No such file or directory" >&2
