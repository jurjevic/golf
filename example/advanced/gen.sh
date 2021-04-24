#!/usr/bin/env bash

# try './gen.sh green'

COLOR="blue"
if [ $# -eq 0 ]
  then
    echo "No arguments supplied, so we choose the red color..."
else
  COLOR=$1
fi

go install github.com/jurjevic/golf@latest

# golf blue.sh red.sh -i=sh.go
go run ../*.go red.sh ${COLOR}.sh -i=sh.go -- 'var color string = "'$COLOR'"; var debug bool = true'
