module=github.com/c-reeder/aoc2022
current=d02

build:
	go build -o ./bin $(module)/cmd/$(current)

# Day-specific rules
# E.g. "make d01"
d%:
	go build -o ./bin $(module)/cmd/$@
