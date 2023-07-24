#!/usr/bin/env bash

# Build services
services=(mystic-forge)

for service in "${services[@]}"
do
  echo "building the ${service}"
	go build -buildvcs=false -o bin/ "./"
  if [ "$?" != 0 ]; then
    exit 1
  fi
done
