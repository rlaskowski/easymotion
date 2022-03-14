GOBIN := go
BUILDNAME := easymotion

build-docker:
	@-$(MAKE) clean
	docker build -t rlaskowski/easymotion .

build-linux:
	@-$(MAKE) clean
	GOOS=linux go build -o dist/easymotion cmd/easymotion/main.go

build-darwin-arm:
	@-$(MAKE) clean
	GOOS=darwin GOARCH=arm64 go build -o dist/easymotion cmd/easymotion/main.go

build-raspi:
	@-$(MAKE) clean
	GOOS=linux GOARCH=arm go build -o dist/easymotion cmd/easymotion/main.go

get-dependencies:
	@go get -d -v ./...

run:
	go run cmd/easymotion/main.go run

clean:
	@rm -Rf dist 


