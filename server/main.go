package main

import (
	"flag"
	"net"
	"log"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "app/todo"
)

var (
	port = flag.Int("port", 50051, "The port on which the gRPC server will run.")
)

type server struct {
	todos []*pb.OneTodo
}

// CREATE one todo
func (s *server) CreateTodo(ctx context.Context, request *pb.CreateTodoRequest) (*pb.OneTodo, error) {
	s.todos = append(s.todos, request.Message)
	return &pb.OneTodo{Id: request.Message.Id, Task: request.Message.Task}, nil
}

// READ one todo
func (s *server) GetTodo(ctx context.Context, request *pb.GetTodoRequest) (*pb.OneTodo, error) {
	var oneTodo *pb.OneTodo
	for _, todo := range s.todos {
		if todo.Id == request.Id {
			oneTodo = todo
		}
	}
	if oneTodo == nil {
		return &pd.OneTodo{}, nil
	}
	return &pb.OneTodo{Id: oneTodo.Id, Task: oneTodo.Task}, nil
}

// READ all todos
func (s *server) GetTodos(_ *empty.Empty, stream pb.Todo_GetTodosServer) error {
	for _, todo := range s.todos {
		if err := stream.Send(todo); err != nil {
			return err
		}
	}
	return nil
}

// UPDATE one todo


// DELETE one todo
func (s *server) DeleteTodo(ctx context.Context, request *pb.DeleteTodoRequest) (*empty.Empty, error) {
	y := s.todos[:0]
	for _, todo := range s.todos {
    		if todo.Id != request.Id {
			y = append(y, todo)
		}
	}
	s.todos = y
	return &empty.Empty{}, nil
}

// DELETE all todos
func (s *server) DeleteTodos(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	s.todos = s.todos[:0]
	return &empty.Empty{}, nil
}

// GET health
func (s *server) GetHealth(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func main() {
	flag.Parse() 								// parse declared command line flags
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))	// start a server listening on localhost + port
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer() 							// create an instance of a gRPC server
	pb.RegisterTodoServer(s, &server{}) 					// register our service with the gRPC server
	reflection.Register(s)							// provide information about publically-accessible gRPC services on the server
	if err := s.Serve(lis); err != nil {					// accept incomming connections on the listener; log errors otherwise
		log.Fatalf("failed to serve: %s", err)
	}
}
