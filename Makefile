GOBIN := go
BUILDNAME := easymotion

build:
	@-$(MAKE) clean
	GOOS=linux GOARCH=arm go build -o dist/ cmd/easymotion/main.go

get-dependencies:
	@go get -d -v ./...

run:
	go run cmd/easymotion/main.go run

clean:
	@rm -Rf dist 


