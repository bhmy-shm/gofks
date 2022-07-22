package gofks

//go:generate protoc315 -I=. --proto_path=E:\Go\gofks\transport\example --go_out=. .\transport\proto\hello.proto
//go:generate protoc315 -I=. --proto_path=E:\Go\gofks\transport\example --go-grpc_out=. .\transport\proto\hello.proto
