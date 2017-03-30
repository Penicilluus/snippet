package client

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	proto "../../error"
	pb "../../protocol/gen-go/helloworld"
)

type Client struct {
	*BaseClient
}

var gloalClient *Client
// NewClientWithConfigPath Create a new Client instance
// with the specific config path
func NewClientWithConfigPath(configPath string) *Client {
	// use The singleton pattern, the Client should be global
	// all client use the single client and should not connect server frequently.
	if gloalClient == nil {
		gloalClient = &Client{NewGrpcClient(configPath, "hello", func(conn *grpc.ClientConn) interface{} { return pb.NewGreeterClient(conn)})}
	}
	return gloalClient
}

// GetUserByUserName ...
func (ac *Client) SayHello(userName string) (*pb.HelloReply, proto.ServerError) {

	r, err := ac.WithGrpcClient(func(client interface{}, ctx context.Context, opts ...grpc.CallOption) (interface{}, error) {
		return client.(pb.GreeterClient).SayHello(ctx, &pb.HelloRequest{Name: userName, Num: "1"}, opts...)
	})

	if err != nil {
		log.Errorln("failed to GetUserByUserName:", userName, err)
		return nil, err
	}
	return r.(*pb.HelloReply), nil
}