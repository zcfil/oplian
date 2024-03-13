#!/bin/bash
Dir=$1
pid=
MainDisk=$2
cd  ${Dir}
while ! pgrep -x "lotus-bench" > /dev/null; do
bash  ${Dir}run_bench1.sh ${MainDisk} &
    pid=$!
    wait $pid 2> /dev/null
    sleep 10
done

while pgrep -x "lotus-bench" > /dev/null; do
    wait $(pgrep -f "lotus-bench") 2> /dev/null
    sleep 10
done
bash ${Dir}run_bench2.sh ${MainDisk} &
pid=$!
wait $pid 2> /dev/null
sleep 10

while pgrep -x "lotus-bench" > /dev/null; do
    wait $(pgrep -f "lotus-bench") 2> /dev/null
    sleep 10
done
bash ${Dir}run_bench3.sh ${MainDisk} &
pid=$!
wait $pid 2> /dev/null
sleep 10
while pgrep -x "lotus-bench" > /dev/null; do
    wait $(pgrep -f "lotus-bench") 2> /dev/null
    sleep 10
done
cat bench*  |grep   'commit phase' |awk '{ sum += $5 } END { print "time: " sum/NR }' > cc.txt

DD=`cat bench*  |grep   'commit phase'  |awk '{ sum += $5 } END { print  sum/NR }'`
echo  ${DD}

