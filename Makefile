BINARY_NAME=wordle
VERSION=0.1.0

build:
	GOARCH=amd64 GOOS=darwin go build -o ./target/${BINARY_NAME}-${VERSION}-darwin-amd64 . 
	GOARCH=amd64 GOOS=linux go build -o ./target/${BINARY_NAME}-${VERSION}-linux-amd64 .
	GOARCH=amd64 GOOS=windows go build -o ./target/${BINARY_NAME}-${VERSION}-windows-amd64 .
	GOARCH=arm64 GOOS=darwin go build -o ./target/${BINARY_NAME}-${VERSION}-darwin-arm64 . 
	GOARCH=arm64 GOOS=linux go build -o ./target/${BINARY_NAME}-${VERSION}-linux-arm64 .

run: build
	./target/${BINARY_NAME}

clean:
	go clean
	rm -rf ./target
