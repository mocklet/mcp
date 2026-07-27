.PHONY: build test format clean

build:
	go build -o mocklet-mcp .

test:
	go test ./...

format:
	gofmt -w .

clean:
	rm -f mocklet-mcp
