#!/bin/bash
cat > /root/oplian/script/run_lotus.sh<<EOF
#!/bin/bash
set -e
datetime=\$(date +'%Y-%m-%d %H:%M:%S')

export LISTEN_PORT=$1
export DIR=$2
export LISTEN_IP=$3
export LOTUS_ID=$4
export MainDisk=$5

export RUST_BACKTRACE=full
export RUSTFLAGS="-C target-cpu=native -g"
export FFI_BUILD_FROM_SOURCE=1
export RUST_LOG=info
export LOTUS_PATH=\${MainDisk}/ipfs/data/lotus
export IPFS_GATEWAY="https://proof-parameters.s3.cn-south-1.jdcloud-oss.com/ipfs/"
export FIL_PROOFS_PARAMETER_CACHE=\${MainDisk}/filecoin-proof-parameters

unset FIL_PROOFS_MAXIMIZE_CACHING
export SKIP_BASE_EXP_CACHE=1
#export GOLOG_LOG_FMT=json

pid=$(ps -ef | grep 'lotus daemon'| egrep -v "log|grep"  | awk '{print $2}')
if [ x"$pid" = "x" ] ;then
 nohup /root/oplian/bin/lotus daemon --import-snapshot \${DIR}  &
  sudo prlimit --nofile=1048576 --nproc=unlimited --rtprio=99 --nice=-19 --pid \$!
  echo "${datetime}  lotus restarted  successfully!" >> \${MainDisk}/ipfs/logs/restart_press.log
else
  echo "${datetime}  lotus RUNNING  successfully!"
fi
wait

EOF
sleep  1
chmod  +x  /root/oplian/script/run_lotus.sh || echo "ERROR: No such file or directory" >&2
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

cat >  /etc/supervisor/conf.d/lotus.conf  <<EOF
[program:lotus]
command=/root/oplian/script/run_lotus.sh
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
stdout_logfile=$5/ipfs/logs/lotus.log
EOF

ls /etc/supervisor/conf.d/lotus.conf  || echo "ERROR: No such file or directory" >&2
