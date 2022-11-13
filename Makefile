gobin := go
os := linux
arch := amd64
buildname := device


build-docker:
	@-$(MAKE) clean
	docker build --build-arg buildname=$(buildname) -t rlaskowski/easymotion-$(buildname) .

build:
	@-$(MAKE) clean
	@GOOS=$(os) GOARCH=$(arch) ${gobin} build -tags $(buildname) -o dist/$(buildname)/$(buildname) cmd/easymotion/main.go

get-dependencies:
	@-${gobin} get -d -v ./...

run:
	${gobin} run -tags $(buildname) cmd/easymotion/main.go run

clean:
	@rm -Rf dist data


