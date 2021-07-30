#!/bin/bash
bench () {
  for i in {1..2}
  do
    echo "Running autocannon ($1) $i"
    autocannon http://localhost:8098/api/avgLaptimes --warmup [ -c 2 -d 10 ] --connections $1 --workers 6 --json >> ./autocannon.json
  done
}

TS=$(date '+%Y-%m-%d_%H:%M:%S')
echo "Date is $TS"

mkdir "./results"/$TS
cd "./results"/$TS

bench 5
bench 20
#bench 50
