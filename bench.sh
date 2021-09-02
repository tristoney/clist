#!/bin/bash

go test -v -run=NOTEST -cpu=12 -bench=. -benchtime=1000x -count=20 > bench.txt
benchstat bench.txt