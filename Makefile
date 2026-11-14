.PHONY: run trace parse ast dot text html web test cover

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

cst:
	go run cmd/sipl/main.go cst test.sipl

web: html

html:
	go run cmd/sipl/main.go web test.sipl
ifeq ($(OS),Windows_NT)
	@cmd.exe /c start "" tree.html
else
	@(command -v xdg-open >/dev/null 2>&1 && xdg-open tree.html) || \
	 (command -v open >/dev/null 2>&1 && open tree.html) || \
	 echo "Generated tree.html - open it manually in your browser."
endif

cover:
	go test -coverpkg=./pkg/ast,./pkg/parser,./pkg/interpreter -coverprofile=coverage.out ./...
	go tool cover -func coverage.out