#!/bin/bash
while true; do
pids=$(ps -ef | grep "bench" | grep -v grep | awk '{print $2}')
if [ -n "$pids" ]; then
echo "Found processes: $pids"
kill -9 $pids
echo "Killed processes: $pids"
pids=$(ps -ef | grep "kill_net" | grep -v grep | awk '{print $2}')
if [ -n "$pids" ]; then
echo "Found processes: $pids"
kill -9 $pids
echo "Killed processes: $pids"
fi
break
else
sleep 3
fi
done