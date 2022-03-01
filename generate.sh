protoc greet/greetpb/greet.proto --go-grpc_out=. # generate service stub
protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative greet/greetpb/greet.proto # generate message stub

# combine service and message stubs
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative greet/greetpb/greet.proto
