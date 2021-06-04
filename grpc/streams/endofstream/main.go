package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"
)

func main() {
	srv := &Server{}

	err := srv.Run(context.Background())
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(time.Second * 5)
		srv.Stop()
	}()

	err = runClientUntilError(context.Background())
	fmt.Println("error:", err)
	if errors.Is(err, io.EOF) {
		fmt.Println("gotcha")
	}
}
