#!/bin/bash
##获取磁盘盘符，程序传参：
disknum=$1
a=$2
diskdrive=$(echo $a | sed 's/,/ /g')

##创建md0下的logs目录
mkdir -p /mnt/md0/ipfs/logs/

##判断如存在是否存在块设备md0阵列存在退出
# 获取所有已定义的软件 RAID 块设备名称
RAID_DEVICES=$(cat /proc/mdstat | grep md  |awk   '{print  $4}')
# 判断是否存在 'raid0' 块设备
if [[ "$RAID_DEVICES" == *"raid0"* ]]; then
  echo "raid0 exists."
  if mount | grep /mnt/md0 | grep -q '^/dev/md*'; then
    echo "/md0 挂载在一个 md 设备上"
else
    echo "/md0 未挂载在一个 md 设备上"
    df -h  |grep  'dev/md' |awk '{print "umount -lf " $NF}'  |bash
	  mount  /dev/md*   /mnt/md0/
fi
  exit 1
else

##如果没有RAID进行组阵列
# 检查 /mnt/md0/ 目录是否已经存在，不存在则创建
if [ ! -d "/mnt/md0/" ]; then
  echo "Warning: /mnt/md0/ does not exist, creating..."
  mkdir -p /mnt/md0/
fi

#########

# 卸载md0
umount -lf  /mnt/disk*
sleep 5
umount -lf  /mnt/md0

echo y | mdadm --create /dev/md0  --level=0 --raid-devices=${disknum} ${diskdrive}

if [ $? -ne 0 ]; then
    echo -e "\033[32m failed \033[0m"
else
    echo -e "\033[32m mdo created succeed \033[0m"
fi

# Wait for the array to be created
while [ ! -e /dev/md0 ]; do
   sleep 3
done

#格式化raid0
mkfs.xfs -f /dev/md0

if [ $? -ne 0 ]; then
    echo -e "\033[32m failed \033[0m"
else
    echo -e "\033[32m mdo init succeed \033[0m"
fi

#挂载md0
[ $(mount -l | grep /mnt/md0 | wc -l) -eq 1 ] || ( mkdir -p /mnt/md0 && sudo mount /dev/md0 /mnt/md0  )
diskid=`ls -la /dev/disk/by-uuid/ | grep md0 | awk '{print $9}'`
grep -i md0 /etc/fstab
if [ $? -eq 0 ]; then
    echo -e "\033[32m fstabmdo already exists \033[0m"
else
mdadm -Ds >> /etc/mdadm/mdadm.conf
fi

fi

echo $disknum
echo $diskdrive
