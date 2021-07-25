#!/bin/bash
go build && \
source ./local.env && \
export $(cut -d= -f1 ./local.env) && \
./dev.hartenfeller.sqlite-mysql-benchmark
