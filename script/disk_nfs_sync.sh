#!/bin/bash
source   /etc/profile
MD5Context1=$(md5sum  /etc/exports  |awk  '{print $1}')
output=$(lsblk -f | grep - | grep /$ -v|grep  /boot -v  | awk '{print $NF}' | grep ^/ | sort -nr | uniq)
for i in $output; do
  if [ -f "/etc/exports" ]; then
    if grep -wq "^$i" /etc/exports; then
      echo "$i 已存在于 /etc/exports 文件中，忽略操作"
    else
      echo "$i *(rw,no_root_squash,no_subtree_check,async)" >> /etc/exports
      echo "已将结果 $i 追加到 /etc/exports 文件"
    fi
  else
    echo "/etc/exports 文件不存在"
  fi
done
MD5Context2=$(md5sum  /etc/exports  |awk  '{print $1}')
  if [ "$MD5Context1" != "$MD5Context2" ]; then
      echo "文件内容发生了变化,重启服务端"
      systemctl restart nfs-kernel-server.service
    else
      echo  "未操作"
  fi

