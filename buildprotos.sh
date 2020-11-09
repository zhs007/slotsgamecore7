export PATH="$PATH:$(go env GOPATH)/bin"
protoc --proto_path=protos/ --go_out=plugins=grpc:sgc7pb --go_opt=paths=source_relative protos/*.proto
protoc --proto_path=rngprotos/ --go_out=plugins=grpc:dtrngpb --go_opt=paths=source_relative rngprotos/*.proto