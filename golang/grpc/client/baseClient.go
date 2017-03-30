package client

import (
	"time"

	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	log "github.com/Sirupsen/logrus"

	"../../pool"
	proto "../../error"
	"fmt"
)

type  Constructor func(*grpc.ClientConn) interface{}

type BaseClient struct {
	ConfigPath string
	ConfigName string
	ClientName string
	Construct Constructor

	GrpcClient interface{}
	Pool *pool.Pool
}

// not connect the server when newClient
// in the withGrpcClient like Gorm connect database
// read endpoint from configPath
// multiconfig pf13/viper toml switch dev and product
func NewGrpcClient(configPath, clientName string, construct Constructor ) *BaseClient {
	return &BaseClient{
		ConfigPath: configPath,
		ClientName: clientName,
		Construct: construct,
		ConfigName: clientName+"_conf",
	}
}

// use func as parameter, in this func(WithGrpcClient) monitor some indicators
// user panic failure fast
func (cli *BaseClient)WithGrpcClient(useClient func(clientInterface interface{}, ctx context.Context, opts ...grpc.CallOption) (interface{}, error)) (interface{},  proto.ServerError)  {
	serverPath := "localhost:50051"
	if cli.Pool == nil {
		p, err := pool.New(func() (*grpc.ClientConn, error) {
			conn, err := grpc.Dial(
				serverPath,
				grpc.WithInsecure(),
				grpc.WithBlock(),
				//grpc.WithTimeout(time.Duration(1000)),
			)
			if err != nil {
				fmt.Println("dial err")
				panic(err)
			}
			return conn, nil
		}, 1, 1, 0)
		if err != nil {
			log.Errorf("The pool returned an error: %s", err.Error())
			panic(err.Error())
		}
		cli.Pool = p
	}

	conn, err := cli.Pool.Get(context.Background())
	if err != nil {
		log.Errorf("The pool.Get returned an error: %s", err.Error())
		panic(err.Error())
	}
	if conn.ClientConn == nil {
		panic("pool get nil")
	}

	cli.GrpcClient = cli.Construct(conn.ClientConn)
	//timeOut := clientConfig.SocketTimeout()
	socketTimeout := time.Duration(100) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), socketTimeout*time.Millisecond)
	defer cancel()

	var header, trailer metadata.MD

	r, err := useClient(cli.GrpcClient, ctx, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Errorf("useClient returned an error: %s", err.Error())
	}

	return r, nil
}

func (cli *BaseClient)CloseConnection() {
	log.Debugf("%s.Client --- Close Connection", cli.ClientName)
	if cli.Pool != nil {
		cli.Pool.Close()
	}
}