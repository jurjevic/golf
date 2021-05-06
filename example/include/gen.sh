#!/usr/bin/env bash

# try './gen.sh green'

COLOR="blue"
if [ $# -eq 0 ]
  then
    echo "No arguments supplied, so we choose the blue color..."
else
  COLOR=$1
fi

go install github.com/jurjevic/golf@latest

# golf blue.sh red.sh
go run ../../*.go red.sh ${COLOR}.sh -- 'var color string = "'$COLOR'"; var debug bool = true'
