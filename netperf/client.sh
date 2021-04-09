#!/usr/bin/env bash
set -o errexit -o nounset -o pipefail
IFS=$'\n\t\v'
cd `dirname "${BASH_SOURCE[0]:-$0}"`

remote=$1
filename="${2:-noname}"
testamount=${3:-1000}

for requestsize in 1 500 1000 1500; do
  for testname in "TCP_RR" "UDP_RR" "TCP_STREAM" "UDP_STREAM"; do
    file="$filename-$testname-${requestsize}b-x$testamount.txt"
    echo "performing test $testname $testamount times with request size $requestsize bytes."
    if [[ -f "$file" ]]; then
      echo "SKIP: $file"
      continue
    fi
    echo "writing to $file"
    rm -f $file
    for (( i=0; i<$testamount; i++ )); do
      P=0
      if [ $i -eq 0 ]; then
        P=1
      fi
      netperf -H $remote -I99,5 -j -c -l -1000 -t TCP_RR -P $P -- -D -r $requestsize -O THROUGHPUT_UNITS,THROUGHPUT,MEAN_LATENCY,MIN_LATENCY,MAX_LATENCY,P50_LATENCY,P90_LATENCY,P99_LATENCY,STDDEV_LATENCY,LOCAL_CPU_UTIL,SOCKET_TYPE,PROTOCOL,REQUEST_SIZE,LOCAL_CPU_UTIL,LOCAL_SD,CONFIDENCE_LEVEL >> $file
    done
  done
done