#!/bin/bash
export LC_ALL=en_US.UTF-8
mounted_disks=()
counter=1
{
    lsblk -f | grep raid -v | egrep '/$|/boot|LVM' -v | awk '{print $3}'
    lsblk -f | grep raid -v | egrep '/$|/boot|LVM' -v | awk '{print $4}'
} | grep - | sort -nr | uniq |
while read -r disk_uuid; do
    if [[ ! " ${mounted_disks[@]} " =~ " $disk_uuid " ]]; then
        disk_path="/dev/disk/by-uuid/$disk_uuid"
        disk_path_dev=$(ls  -l  ${disk_path}   |  awk  -F '/' '{print $NF}')  
	    disk_path_2=$(lsblk -f |grep  -w  ${disk_path_dev}  |awk   '{print $NF}'|sort  -nr|uniq)  
        mount_path="/mnt/disk$counter"
        if [[ $(mount -l | mount -l | egrep  ^/dev|grep   -w "$disk_path_2" | wc -l) -eq 0 ]]; then
          while [[ $(mount -l |egrep  ^/dev|grep   -w "$mount_path" | wc -l) -ne 0 ]]; do
                ((counter++))
                mount_path="/mnt/disk$counter"
            done
            mkdir -p "$mount_path"
            mount "$disk_path" "$mount_path"
            mounted_disks+=("$disk_uuid")
        fi

    fi
done
