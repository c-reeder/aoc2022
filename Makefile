module=github.com/c-reeder/aoc2022
current=d01

d01:
	go run $(module)/cmd/d01

run:
	go run $(module)/cmd/$(current)
