GOBIN := go
BUILDNAME := easymotion

build:
	@-${MAKE} clean
	go build -o dist/ cmd/easymotion.go

build-raspi:
	@-$(MAKE) clean
	GOOS=linux GOARCH=arm go build -o dist/ cmd/easymotion.go

build-docker:
	docker build -t rlaskowski/easymotion:latest .

get-dependencies:
	@go get -d -v ./...

run:
	go run cmd/easymotion.go

clean:
	@rm -Rf dist 


