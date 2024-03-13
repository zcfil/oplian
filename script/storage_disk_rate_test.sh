#!/bin/bash
 
# 定义要测试的磁盘或分区挂载点路径和设备名称，只测试容量大于 1T 的设备（可根据实际情况修改）
txtfile=txtfile1 
# 根据设备名称循环测试读写速度
for i in `df -h | awk '$2 ~ /T/ && $2+0 >= 1 { print $6}'`
do
sdd=`df -h  |grep  -w  $i |awk   '{print $1}'`
    # 测试随机读写取速率
   # echo "Testing read speed..."
    #dd if=$device/1 of=/dev/null bs=512K count=32k oflag=direct 2>&1 | awk -F, '{print $3}' | awk '{print $1 / 1048576 " MB/s" }'
   DD2=`dd if=/dev/zero of=${devices}/${txtfile}    bs=1024M count=10  oflag=direct    2>&1 | awk -F, '{print  $4}'|grep	 ^$ -v`  
   echo "${sdd} ${i} ${DD2}"|awk '{printf("%-13s%-13s%-6s\n", $1, $2, $3" "$4)}' 
   rm ${devices}/${txtfile}
    # 测试写入速率
    #echo "Testing write speed..."
    #dd if=/dev/zero of=$device/1 bs=512K count=32k oflag=direct 2>&1 | awk -F, '{print $3}' | awk '{print $1 / 1048576 " MB/s" }'
done
