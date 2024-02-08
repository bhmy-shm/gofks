package main

//go:generate protoc -I=. --proto_path=. --go_out=./protoc --go-grpc_out=./protoc ./user.proto
