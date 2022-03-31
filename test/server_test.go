package test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"challenge/pkg/proto"
	. "challenge/pkg/proto"
	. "challenge/pkg/server"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:8088", "the address to connect to")
)

func startServer(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	grpcServer.RegisterService(&ChallengeService_ServiceDesc, TestRPCServer{})

	grpcServer.Serve(listener)
}

func startClient() (*grpc.ClientConn, ChallengeServiceClient) {

	conn, _ := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	return conn, proto.NewChallengeServiceClient(conn)

}

func TestMetaData(t *testing.T) {

	go startServer(8088)

	time.Sleep(time.Second)
	conn, client := startClient()
	defer conn.Close()
	mdata, err := client.ReadMetadata(context.Background(), &proto.Placeholder{Data: "Data"})

	x := map[string]string{}
	json.Unmarshal([]byte(mdata.Data), &x)

	println("metadata from local server: ", mdata.Data, err)

	for key, val := range x {
		fmt.Println("key: ", key, "val: ", val)
		if !((key == ":authority") || (key == "content-type") || (key == "user-agent")) {
			t.Error("no metadata")
		}
	}
	fmt.Println("json: ", x)
	log.Default().Println("TestMetaData done")

}

func TestShortLink(t *testing.T) {

	go startServer(8089)

	time.Sleep(time.Second)
	conn, client := startClient()
	defer conn.Close()
	os.Setenv("BITLY_OAUTH_TOKEN", "448151789bf0264b0596dac054cdc900c10a1b40")

	link, _ := client.MakeShortLink(context.Background(), &proto.Link{Data: "https://qwertyuiogrtjrgtkldfksgdsgpfw.com"})
	println(link.Data)

	x := map[string]interface{}{}
	err := json.Unmarshal([]byte(link.Data), &x)

	fmt.Println("\n\r\n\rerr: ", err, "\n\r", "json: ", (x["id"]))

	if x["id"] != "bit.ly/3wLOGyH" {
		t.Error("reference and read links dont match")
	}
	log.Default().Println("TestShortLink done")

}
