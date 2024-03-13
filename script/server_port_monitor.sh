#!/bin/bash
# 检查iperf3是否已经安装
if ! command -v iperf3 &> /dev/null
then
    # 如果未安装，则进行安装
    sudo apt-get update   >/dev/null  2>&1
    sudo apt-get install -y iperf3   >/dev/null  2>&1
else
    # 如果已经安装，则输出当前版本信息
    iperf3 -v   >/dev/null  2>&1
fi
# 定义要启动的端口数组
PORTS=(5201 5202 5203 5204)
# 循环启动每个端口上的iperf3服务器，并将PID写入pidfile
for port in "${PORTS[@]}"; do
    iperf3 -s -p "$port"   >/dev/null  2>&1   &
done
# 停止所有iperf3服务器一段时间后杀死进程
nohup  sleep 1000 &&  kill  -9   $(ps aux | grep  520  |grep      grep -v |awk '{print $2}')  >  /var/log/iperf3.log 2>&1   &
