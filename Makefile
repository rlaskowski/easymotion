gobin := go
os := linux
arch := amd64
buildname := device


build-docker:
	@-$(MAKE) clean
	docker build --build-arg buildname=$(buildname) -t rlaskowski/easymotion-$(buildname) .

build:
	@-$(MAKE) clean
	@GOOS=$(os) GOARCH=$(arch) ${gobin} build -o dist/$(buildname)/$(buildname) cmd/$(buildname)/main.go

get-dependencies:
	@-${gobin} get -d -v ./...

run:
	${gobin} run cmd/device/main.go run

clean:
	@rm -Rf dist data


