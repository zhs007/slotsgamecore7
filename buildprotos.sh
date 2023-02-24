export PATH="$PATH:$(go env GOPATH)/bin"
protoc --proto_path=protos/ --go_out=./sgc7pb --go_opt=paths=source_relative protos/*.proto
protoc --proto_path=protos/ --go-grpc_out=./sgc7pb --go-grpc_opt=paths=source_relative protos/*.proto
protoc --proto_path=rngprotos/ --go_out=./dtrngpb --go_opt=paths=source_relative rngprotos/*.proto
protoc --proto_path=rngprotos/ --go-grpc_out=./dtrngpb --go-grpc_opt=paths=source_relative rngprotos/*.proto