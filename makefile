SHELL := cmd.exe
.SHELLFLAGS := /Q /C

PROTOC := C:/tools/protoc/bin/protoc.exe

r:
	go run .

gengrpc:
	"$(PROTOC)" -I proto -I third_party\googleapis --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative proto\v1\*.proto
