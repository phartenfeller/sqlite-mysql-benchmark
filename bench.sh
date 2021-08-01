#!/bin/bash
bench () {
  for i in {1..5}
  do
    echo "Running autocannon ($1) $i"
    autocannon 51.116.188.78:8098/api/randomRead --warmup [ -c 2 -d 10 ] --connections $1 --workers 6 -d 60 --json >> ./autocannon.json
    sleep 15s
  done
}

TS=$(date '+%Y-%m-%d_%H:%M:%S')
echo "Date is $TS"

mkdir "./results"/$TS
cd "./results"/$TS

bench 5
bench 10
bench 25
bench 50
bench 75
bench 100
bench 200
bench 400
#bench 500
#bench 1000
