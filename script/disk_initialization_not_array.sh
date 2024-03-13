#!/bin/bash
##创建md0下的logs目录
mkdir -p /mnt/md0/ipfs/logs/

first_disk_uuid=$(lsblk  -f  |grep   /$ -v | awk '{print $3}'|grep  - | head -n 1)
first_disk_path="/dev/disk/by-uuid/$first_disk_uuid"
mkdir -p /mnt/md0
mount "$first_disk_path" /mnt/md0

for disk_uuid in $(lsblk  -f  |grep   /$ -v | awk '{print $3}'|grep  - | tail -n +2); do
    disk_path="/dev/disk/by-uuid/$disk_uuid"
    mount_path="/mnt/disk$(lsblk  -f  |grep   /$ -v | awk '{print $3}'|grep  -| grep -n "$disk_uuid" | cut -d: -f1)"
    mkdir -p "$mount_path"
    mount "$disk_path" "$mount_path"
done