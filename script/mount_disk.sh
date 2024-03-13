ip=$1
disk=$2
DISK_ALL=$(echo  $disk|awk -F',' '{for(i=1; i<=NF; i++) print $i}')
for  i  in  ${DISK_ALL}
do
    { [ $(mount -l | grep -w  "/mnt/$ip$(echo ${i}|sed 's#/mnt##g')" | wc -l) -ge 1 ] || ( if [ -d /mnt/$ip$(echo ${i}|sed 's#/mnt##g') ]; then  mount -t nfs -o noatime $ip:${i} /mnt/$ip$(echo ${i}|sed 's#/mnt##g'); else  mkdir -pv  /mnt/$ip$(echo ${i}|sed 's#/mnt##g')  &&  mount -t nfs -o noatime $ip:${i}  /mnt/$ip$(echo ${i}|sed 's#/mnt##g') ; fi ) }
done