package init

import (
	"flag"
	"lianxi/bff/dasic/config"
	__ "lianxi/srv/proto/goods"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	initDB()
}
func initDB() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	config.GoodsClient = __.NewGoodsClient(conn)
}
