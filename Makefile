GOBIN := go
OS := linux
ARCH := amd64
BUILDNAME := easymotion

build-docker:
	@-$(MAKE) clean
	docker build -t rlaskowski/easymotion .

build:
	@-$(MAKE) clean
	@GOOS=$(OS) GOARCH=$(ARCH) ${GOBIN} build -o dist/device cmd/device/main.go

build-darwin-arm:
	@-$(MAKE) clean
	GOOS=darwin GOARCH=arm64 ${GOBIN} build -o dist/device cmd/device/main.go

build-raspi:
	@-$(MAKE) clean
	GOOS=linux GOARCH=arm ${GOBIN} build -o dist/device cmd/device/main.go

get-dependencies:
	@-${GOBIN} get -d -v ./...

run:
	${GOBIN} run cmd/device/main.go run

clean:
	@rm -Rf dist data


