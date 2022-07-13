# service stubs
protoc \
    --go-grpc_out=paths=source_relative:. \
    greet/greetpb/greet.proto

# combine service and message stubs
protoc \
    --go_out=paths=source_relative:. \
    --go-grpc_out=paths=source_relative:. \
    greet/greetpb/greet.proto
