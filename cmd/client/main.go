package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	proto "challenge/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	addr = flag.String("addr", "localhost:8088", "the address to connect to")
)

type ContextKey string

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewChallengeServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	mData := metadata.New(map[string]string{"key": "key_value"})
	ctxVal := metadata.NewIncomingContext(ctx, mData)

	md, ok := metadata.FromIncomingContext(ctxVal)
	if ok {
		fmt.Println(md.Get("key")[0])
	} else {
		fmt.Println("key not find")
	}

	defer cancel()
	r, err := c.ReadMetadata(ctxVal, &proto.Placeholder{Data: "Data"})

	c.MakeShortLink(ctxVal, &proto.Link{Data: "https://qwertyuiogrtjrgtkldfksgdsgpfw.com"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	res, err := c.StartTimer(ctxVal, &proto.Timer{Name: "name", Seconds: 20, Frequency: 1})

	var i = 0
	if err == nil {
		for {
			resp, err := res.Recv()
			println("resp:", resp, "err: ", err)
			if err == io.EOF {
				fmt.Println("err == io.EOF i = ", i)
				break

			}
			if err != nil {
				log.Fatalf("cannot receive %v", err)
			}
			fmt.Printf("%+v\n", resp)

			i++
		}

	} else {
		println(err)
	}

	log.Printf("Greeting: %s", r.GetData())
}
