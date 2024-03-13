#!/bin/bash
pids=$(ps -ef | grep "bench" | grep -v grep | awk '{print $2}')
if [ -n "$pids" ]; then
echo "Found processes: $pids"
kill -9 $pids
echo "Killed processes: $pids"
fi
pids=$(ps -ef | grep "iperf3" | grep -v grep | awk '{print $2}')
if [ -n "$pids" ]; then
echo "Found processes: $pids"
kill -9 $pids
echo "Killed processes: $pids"
fi
pids=$(ps -ef | grep "kill_" | grep -v grep | awk '{print $2}')
if [ -n "$pids" ]; then
echo "Found processes: $pids"
kill -9 $pids
echo "Killed processes: $pids"
fi