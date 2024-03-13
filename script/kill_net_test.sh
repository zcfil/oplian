#!/bin/bash
while true; do
pids=$(ps -ef | grep "iperf3" | grep -v grep | awk '{print $2}')
if [ -n "$pids" ]; then
echo "Found processes: $pids"
kill -9 $pids
echo "Killed processes: $pids"
break
else
sleep 3
fi
done