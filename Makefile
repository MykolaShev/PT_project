.PHONY: run trace parse json

run:
	go run cmd/sipl/main.go run test.sipl

trace:
	go run cmd/sipl/main.go trace test.sipl

parse:
	go run cmd/sipl/main.go parse test.sipl

json:
	go run cmd/sipl/main.go --json run test.sipl