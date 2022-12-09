module=github.com/c-reeder/aoc2022
current=d08

build:
	go build -o ./bin $(module)/cmd/$(current)

test:
	go test $(module)/cmd/$(current)

# Day-specific rules
# E.g. "make d01"
d%:
	go build -o ./bin $(module)/cmd/$@
