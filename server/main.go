package "main"

import (
	"flag"
	"net"
	"log"

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

// READ one todo

// READ all todos
func (s, *server) GetTodos(_ *empty.Empty, stream pb.Todo_GetTodosServer) error {
	for _, todo := range s.todos {
		if err := stream.Send(todo); err != nil {
			return err
		}
	}
	return nil
}

// UPDATE one todo

// DELETE one todo

// DELETE all todos

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
