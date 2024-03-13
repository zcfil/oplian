#!/bin/bash
export PATH=$PATH
cat > /root/oplian/script/run_storage.sh <<EOF
#!/bin/bash

WORKER_PORT0=$1
export MINER_API_INFO=$2
MainDisk=$3
export UNSEALED_SERVER=$4

export RUST_BACKTRACE=full
export RUST_LOG=info
export WORKER_PATH=\${MainDisk}/ipfs/data/lotusstorage
export FIL_PROOFS_PARENT_CACHE=\${MainDisk}/filecoin-parents/
export IPFS_GATEWAY="https://proof-parameters.s3.cn-south-1.jdcloud-oss.com/ipfs/"
export FIL_PROOFS_PARAMETER_CACHE=\${MainDisk}/filecoin-proof-parameters
export TRUST_PARAMS_FORCE=1
export FIL_PROOFS_USE_GPU_TREE_BUILDER=1
export FIL_PROOFS_USE_MULTICORE_SDR=1
export FIL_PROOFS_MAXIMIZE_CACHING=1
export FIL_PROOFS_USE_GPU_COLUMN_BUILDER=1


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

    nohup /root/oplian/bin/lotus-storage --worker-repo \${MainDisk}/ipfs/data/lotusstorage/ run --storage=true --listen 0.0.0.0:\${WORKER_PORT0} &
pid=\$!

        sudo prlimit --nofile=1048576 --nproc=unlimited --rtprio=99 --nice=-19 --pid \$!
	wait
EOF

sleep  1

chmod  +x  /root/oplian/script/run_storage.sh || echo "ERROR: No such file or directory" >&2
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

cat >  /etc/supervisor/conf.d/storage.conf  <<EOF
[program:lotus-storage]
command=/root/oplian/script/run_storage.sh
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
stdout_logfile=$3/ipfs/logs/storage.log
EOF

ls /etc/supervisor/conf.d/storage.conf || echo "ERROR: No such file or directory" >&2

