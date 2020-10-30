export PATH="$PATH:$(go env GOPATH)/bin"
protoc --proto_path=protos/ --go_out=plugins=grpc:sgc7pb --go_opt=paths=source_relative protos/*.proto