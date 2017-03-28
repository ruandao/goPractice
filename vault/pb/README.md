
安装依赖库：
    go get -u github.com/golang/protobuf/{proto,protoc-gen-go}


生成 .go 文件：
    protoc vault.proto --go_out=plugins=grpc:.