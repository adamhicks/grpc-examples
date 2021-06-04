package main

import (
	"context"
	"fmt"
	"github.com/adamhicks/grpc-examples/grpc/streams"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc
	srv    *grpc.Server
	wg     sync.WaitGroup
}

func (s *Server) Run(ctx context.Context) error {
	s.srv = grpc.NewServer()
	s.ctx, s.cancel = context.WithCancel(ctx)

	streams.RegisterTestServer(s.srv, s)

	addr := ":8080"
	sock, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}

	fmt.Println("listening on", addr)

	s.wg.Add(1)
	go func() {
		err := s.srv.Serve(sock)
		if err != nil {
			panic(err)
		}
		s.wg.Done()
	}()

	return nil
}

func (s *Server) Stop() {
	fmt.Println("stopping server")
	s.cancel()
	s.srv.GracefulStop()
	s.wg.Wait()
}

func (s *Server) Stream(request *streams.StreamRequest, server streams.Test_StreamServer) error {
	tick := time.NewTicker(time.Millisecond * 100)
	defer tick.Stop()

	i := request.Nonce

	for {
		select {
		case <-tick.C:
			err := server.Send(&streams.Update{Counter: i})
			if err != nil {
				return err
			}
			i++
		case <-server.Context().Done():
			return nil
		case <-s.ctx.Done():
			return nil
		}
	}
}
