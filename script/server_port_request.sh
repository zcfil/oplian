#!/bin/bash
SERVER_IP=$1
SERVER_PORT=$2
TME=$3
LOG_FILE="/tmp/iperf.log"
iperf3  -c  ${SERVER_IP} -p  ${SERVER_PORT}  -t   ${TME}   >  ${LOG_FILE}
cat  ${LOG_FILE} |grep      MBytes | awk '{ sum += $7 } END { print  sum/NR " Gbits/sec" }'
rm  ${LOG_FILE}

