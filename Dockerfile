FROM golang:alpine

RUN apk update && apk add git protobuf

# Install Golang gRPC and Protobuf packages
RUN go get -u google.golang.org/grpc
RUN go get -u github.com/golang/protobuf/proto
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN export PATH=$PATH:$GOPATH/bin

# Copy source into app directory
COPY . /go/src/app
WORKDIR /go/src/app

# Retrieve the Google configurations for mapping RPC to HTTP REST
RUN git clone https://github.com/googleapis/googleapis.git

# Compile protocol buffers
RUN protoc \
        --proto_path=.:googleapis \
        --go_out=plugins=grpc:. \
        todo/todo.proto

# Install the Golang gRPC server
RUN go get -v app/server
RUN go install app/server

ENTRYPOINT ["go/bin/server"]
