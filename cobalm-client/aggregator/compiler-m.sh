#!/bin/bash
set -e 

echo "Compiling..."
go build -o ./bin/aggregator .

echo "Done."