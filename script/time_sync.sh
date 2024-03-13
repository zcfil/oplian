#!/bin/bash
if ! command -v ntpdate &> /dev/null
then
    echo "ntpdate is not installed. Installing now..."
    sudo apt-get update
    sudo apt-get install ntpdate
fi

local_time=$(date +%s)
# 获取本地时间与网络时间的时间戳
sudo ntpdate -q pool.ntp.org > /dev/null  # 查询 NTP 服务器，将系统时间设置为网络时间
ntp_time=$(date +%s)  # 获取系统时间的时间戳，即网络时间的时间戳


# 判断时间是否同步
time_diff=$(( $ntp_time - $local_time ))
if [[ $time_diff -gt 2 || $time_diff -lt -2 ]]; then
    echo false
else
    echo true
fi