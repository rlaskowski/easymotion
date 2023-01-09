gobin := go
os := linux
arch := amd64
arm := 7
buildname := device


build-docker:
	@-$(MAKE) clean
	docker build --build-arg buildname=$(buildname) -t rlaskowski/easymotion-$(buildname) .

build:
	@-$(MAKE) clean
	@GOOS=$(os) GOARCH=$(arch) ${gobin} build -tags $(buildname) -o dist/$(buildname)/$(buildname) cmd/easymotion/main.go

build-raspi: 
	@-$(MAKE) clean
	@GOOS=$(os) GOARCH=$(arch) GOARM=$(arm) ${gobin} build -tags $(buildname) -o dist/$(buildname)/$(buildname) cmd/easymotion/main.go

get-dependencies:
	@-${gobin} get -d -v ./...

run:
	${gobin} run -tags $(buildname) cmd/easymotion/main.go run

clean:
	@rm -Rf service/grpcservice/proto/opencv
	@rm -Rf dist data

proto:
	protoc --go_out=. --go-grpc_out=. service/grpcservice/proto/opencv.proto


