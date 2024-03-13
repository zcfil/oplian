#!/bin/bash
 
# 文件路径和名称
file=$1
 
# 块大小和数量
block_size="1M"
block_count=$((32*1024))

# 测试参数
test_time=60     # 测试时长（秒）

# 初始化测试文件，以确保测试时正确反映磁盘性能
#echo "Initializing test file..."
dd if=/dev/zero of=$file bs=$block_size count=$block_count  > /dev/null 2>&1  && sync  > /dev/null 2>&1 

# 测试读取速率
#echo "Testing read speed..."
dd if=$file of=/dev/null bs=$block_size count=$block_count iflag=direct 2>&1 |awk -F, '{print $4}' 

# 测试写入速率
#echo "Testing write speed..."
dd if=/dev/zero of=$file bs=$block_size count=$block_count oflag=direct 2>&1 |awk -F, '{print $4}' 

rm  ${file}
#echo "Done."
