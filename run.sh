#!/bin/env zsh
go build .
for i in {1..18}
do
    echo "Running cubes for n = $i"
    time ./cubes -n $i -f "results/cubes.txt"
    echo ""
done
