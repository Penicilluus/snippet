package client

import (
	"log"
	"os"

	pb "../../protocol/gen-go/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"fmt"
	//"github.com/grpc-ecosystem/go-grpc-prometheus"
	"testing"
)

const (
	address     = "localhost:50000"
	defaultName = "world"
)

func start() {
	context.Background()
	// Set up a connection to the server.
	conn, err := grpc.Dial(address,
		grpc.WithInsecure(),
		//grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		//grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
	)
	fmt.Println(conn)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name, Num: "2"})
	r1, err := c.SayHelloAgain(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println(r1.Message)
	log.Printf("Greeting: %s", r.Message)
}

func Test_SayHello(t *testing.T) {
	client := NewClientWithConfigPath("../account_client.conf")
	// if use connection pool, need not close connection
	//defer client.CloseConnection()

	resp, err := client.SayHello("hello")

	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}
