.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build:
	go build -o jamql main.go

.PHONY: run
run:
	ENV=dev go run main.go

.PHONY: test
test:
	go test -count=1 -v ./...

.PHONY: race
race:
	go test -race -count=1 ./...

.PHONY: cover
cover:
	go test -coverprofile=c.out -coverpkg=./... -count=1 ./...
	go tool cover -html=c.out

.PHONY: format
format:
	go fmt ./...

.PHONY: clean
clean:
	rm -fr jamql c.out dist/
