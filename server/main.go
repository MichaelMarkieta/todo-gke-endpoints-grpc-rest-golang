package "main"

import (
	"flag"
	"net"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "app/todo"
)

var (
	port = flag.Int("port", 50051, "The port on which the gRPC server will run.")
)

type server struct {
	todos []*pb.Todo
}

// CREATE one todo
func (s *server) CreateTodo(ctx context.Context, todo *pb.Todo) (*pb.Todo, error) {
	s.todos, err = append(s.todos, todo)
	if err != nil {
		return err
	}
	return &pb.Todo{Id: todo.Id, Task: todo.Task}, nil
}

// READ one todo
func (s *server) GetTodo(request *pb.GetTodoRequest) (*pb.Todo, error) {
	var oneTodo *pb.Todo
	for _, todo := range s.todos {
		if todo.Id == request.Id {
			oneTodo = todo
		}
	}
	if oneTodo == nil {
		return nil
	}
	return &pb.Todo{Id: oneTodo.Id, Task: oneTodo.Task}, nil
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
func (s *server) DeleteTodo(ctx context.Context, request *pb.DeleteTodoRequest) (_ *empty.Empty, error) {
	y := s.todos[:0]
	for _, todo := range s.todos {
    		if todo.Id != request.Id {
			y = append(y, todo)
		}
	}
	s.todos = y
	return nil
}

// DELETE all todos
func (s *server) DeleteTodos(ctx context.Context, _ *empty.Empty) (_ *empty.Empty, error) {
	s.todos = []*pb.Todo
	return nil
}

func main() {
	flag.Parse() 								// parse declared command line flags
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))	// start a server listening on localhost + port
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer() 							// create an instance of a gRPC server
	pb.RegisterTodoServer(s, %server{}) 					// register our service with the gRPC server
	reflection.Register(s)							// provide information about publically-accessible gRPC services on the server
	if err := s.Serve(lis); err != nil {					// accept incomming connections on the listener; log errors otherwise
		log.Fatalf("failed to serve: %s", err)
	}
}
