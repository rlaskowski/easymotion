GOBIN := go
BUILDNAME := easymotion

build-docker:
	@-$(MAKE) clean
	docker build -t rlaskowski/easymotion .

build:
	@-$(MAKE) clean
	${GOBIN} build -o dist/easymotion cmd/easymotion/main.go

build-darwin-arm:
	@-$(MAKE) clean
	GOOS=darwin GOARCH=arm64 ${GOBIN} build -o dist/easymotion cmd/easymotion/main.go

build-raspi:
	@-$(MAKE) clean
	GOOS=linux GOARCH=arm ${GOBIN} build -o dist/easymotion cmd/easymotion/main.go

get-dependencies:
	@-${GOBIN} get -d -v ./...

run:
	${GOBIN} run cmd/easymotion/main.go run

clean:
	@rm -Rf dist data


