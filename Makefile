GOBIN := go
BUILDNAME := easymotion

build:
	@-${MAKE} clean
	go build -o dist/ cmd/easymotion.go

build-docker:
	docker build -t rlaskowski/easymotion:latest .

run:
	go run cmd/easymotion.go

clean:
	@rm -Rf dist 
