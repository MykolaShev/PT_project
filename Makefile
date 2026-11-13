.PHONY: run trace parse json test

run:
	go run cmd/sipl/main.go run test.sipl

trace:
	go run cmd/sipl/main.go trace test.sipl

parse:
	go run cmd/sipl/main.go parse test.sipl

json:
	go run cmd/sipl/main.go --json run test.sipl

test:
	go test -v ./...

ast:
	go run cmd/sipl/main.go ast test.sipl

dot:
	go run cmd/sipl/main.go dot test.sipl

text:
	go run cmd/sipl/main.go text test.sipl