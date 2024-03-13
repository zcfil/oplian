#!/bin/bash
export PATH=$PATH
cat > /root/oplian/script/run_worker.sh <<EOF
#!/bin/bash
WORKER_PORT0=$1
StartNo=$5
EndNo=$6
export MINER_API_INFO=$2
export MainDisk=$3
export UNSEALED_SERVER=$4

export RUST_BACKTRACE=full
export RUST_LOG=info
export WORKER_PATH=\${MainDisk}/ipfs/data/lotusworker
export FIL_PROOFS_PARENT_CACHE=\${MainDisk}/filecoin-parents/
export IPFS_GATEWAY="https://proof-parameters.s3.cn-south-1.jdcloud-oss.com/ipfs/"
export FIL_PROOFS_PARAMETER_CACHE=\${MainDisk}/filecoin-proof-parameters
export TRUST_PARAMS_FORCE=1
export FIL_PROOFS_USE_GPU_TREE_BUILDER=1
export FIL_PROOFS_USE_MULTICORE_SDR=1
export P1_CORES=1
export FIL_PROOFS_MAXIMIZE_CACHING=1
export FIL_PROOFS_USE_GPU_COLUMN_BUILDER=1
export MINER32_S32G_PATH=\${MainDisk}/lotusworker/s-32
export MINER64_S64G_PATH=\${MainDisk}/lotusworker/s-32
export MINER32_TREE_D_PATH=\${MainDisk}/lotusworker/s-64
export MINER64_TREE_D_PATH=\${MainDisk}/lotusworker/s-64

export RUST_LOG=info

export HOST_WORKER_COUNT=1
unset TMPDIR
unset LOTUS_MONITOR
function handle_term {
  #kill -TERM \$pid
  kill -9 \$pid
  wait
  exit 0
}
trap 'handle_term' TERM

numactl -C ${StartNo}-${EndNo} nohup /root/oplian/bin/lotus-worker --worker-repo \${MainDisk}/ipfs/data/lotusworker/ run --precommit1=true --precommit2=true --commit=false --listen 0.0.0.0:\${WORKER_PORT0} &

pid=\$!

        sudo prlimit --nofile=1048576 --nproc=unlimited --rtprio=99 --nice=-19 --pid \$!
	wait
EOF

sleep  1

chmod  +x  /root/oplian/script/run_worker.sh
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

cat >  /etc/supervisor/conf.d/worker.conf  <<EOF
[program:lotus-worker]
command=/root/oplian/script/run_worker.sh
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
stdout_logfile=$3/ipfs/logs/worker.log
EOF
