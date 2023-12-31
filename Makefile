.PHONY:
.SILENT:

generate:
	cd protos/ && protoc -I proto proto/sso/sso.proto --go_out=./gen/go --go_opt=paths=source_relative --go-grpc_out=./gen/go/ --go-grpc_opt=paths=source_relative
build: 
	go build sso/cmd/sso/main.go
run:
	./main