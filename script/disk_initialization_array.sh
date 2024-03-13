#!/bin/bash
##获取磁盘盘符，程序传参：
disknum=$1
a=$2
diskdrive=$(echo $a | sed 's/,/ /g')

##创建md0下的logs目录
mkdir -p /mnt/md0/ipfs/logs/

##判断如存在md0则直接恢复raid
md0mum=`lsblk -l | grep md0 |wc -l`
if [ $md0mum -ge 1 ]; then

    echo -e "\033[32m Raid0 has been created without execution \033[0m"
    [ $(mount -l | grep /mnt/md0 | wc -l) -eq 1 ] || ( mkdir -p /mnt/md0 && sudo mount /dev/md0 /mnt/md0  )
#恢复raid
#mdadm -A  /dev/md0
mount  /dev/md0   /mnt/md0/

else

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
sleep 3
wait

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