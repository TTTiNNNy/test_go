package test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
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

func startClient(port int) (*grpc.ClientConn, ChallengeServiceClient) {

	conn, _ := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	return conn, proto.NewChallengeServiceClient(conn)

}

func TestMetaData(t *testing.T) {

	go startServer(8088)

	time.Sleep(time.Second)
	conn, client := startClient(8088)
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
	conn, client := startClient(8089)
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

func TestStartTimer(t *testing.T) {
	log.Default().Println("start timer with the same name test")
	go startServer(8090)
	time.Sleep(time.Second)

	conn1, client1 := startClient(8090)
	defer conn1.Close()
	conn2, client2 := startClient(8090)
	defer conn2.Close()

	_, err1 := client1.StartTimer(context.Background(), &proto.Timer{Name: "timer", Seconds: 8, Frequency: 2})
	time.Sleep(time.Second * 4)
	_, err2 := client2.StartTimer(context.Background(), &proto.Timer{Name: "timer", Seconds: int64(rand.Intn(100)), Frequency: int64(rand.Intn(100))})

	if err1 == nil && err2 == nil {
		var i = 0

		var last_key string
		for key, _ := range TimerPull {
			i++
			last_key = key
		}
		if last_key != "timer" && i != 1 {
			fmt.Println("map:", TimerPull)
			t.Error("Create more than one timer.")
		}
	} else {
		println("err1: ", err1.Error(), "\n\rerr2: ", err1.Error())
		t.Error()

	}

}

func startTimer(port int, name string, frequency int, duration int) {

	log.Default().Printf("connecting to : %d port with name %s, frequency %d and duration %d", port, name, frequency, duration)

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)

	}
	defer conn.Close()
	c := proto.NewChallengeServiceClient(conn)

	res, err := c.StartTimer(context.Background(), &proto.Timer{Name: name, Seconds: int64(duration), Frequency: int64(frequency)})

	if err == nil {
		for {
			resp, err := res.Recv()
			println("resp:", resp, "err: ", err)
			if err == io.EOF {
				fmt.Println("Done ")
				break

			}
			if err != nil {
				log.Fatalf("cannot receive %v", err)
			}
			fmt.Printf("%+v\n", resp)

		}

	} else {
		println(err)
	}
}
