@echo off
echo Generating gRPC code...
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/user.proto

echo Generating GraphQL code...
go run github.com/99designs/gqlgen generate

echo Done!
pause 