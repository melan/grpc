package api

//go:generate protoc --go_out=. --go-grpc_out=. -I proto api.v1.proto
//go:generate protoc --go_out=. --go-grpc_out=. -I proto api.v2.proto
