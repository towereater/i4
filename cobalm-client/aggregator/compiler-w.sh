#!/bin/bash
set -e 

echo "Compiling..."
#env GOOS=windows GOARCH=amd64 go build -o ./bin/aggregator-win64 .
env GOOS=windows GOARCH=386 go build -o ./bin/aggregator-win32 .

echo "Done."