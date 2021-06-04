package main

import (
	"context"
	"fmt"
	"github.com/adamhicks/grpc-examples/grpc/streams"
	"google.golang.org/grpc"
	"math/rand"
)

func runClientUntilError(ctx context.Context) error {
	conn, err := grpc.DialContext(ctx, "127.0.0.1:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	cli := streams.NewTestClient(conn)

	stream, err := cli.Stream(ctx, &streams.StreamRequest{Nonce: rand.Int63()})
	if err != nil {
		panic(err)
	}

	for {
		upd, err := stream.Recv()
		if err != nil {
			return err
		}
		fmt.Println(upd.Counter)
	}
}
