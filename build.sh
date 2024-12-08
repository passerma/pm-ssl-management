#!/bin/bash

export GOOS=linux
export GOARCH=amd64

name=$(cat app.yml | grep -w name | head -n 1 | awk '{print $2}')

rm -rf $name

go build -x -o ./dist/$name
