FROM grpc/go:1.0

RUN apt-get update && apt-get install git -y

# Install Golang gRPC and Protobuf packages
#RUN go get -u github.com/golang/protobuf/proto
#RUN export PATH=$PATH:$GOPATH/bin

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
