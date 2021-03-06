/*
* Copyright 2017 bilxio.
*
* @File: server.go
* @Author: Bill Xiong
* @Date:   2017-10-17 12:24:40
* @Last Modified by:   Bill Xiong
* @Last Modified time: 2017-10-17 12:27:39
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	grpclb "github.com/bilxio/grpclb/naming/etcdv3"
	pb "github.com/bilxio/grpclb/example/helloworld"
)

var (
	serv = flag.String("service", "hello_service", "service name")
	port = flag.Int("port", 50001, "listening port")
	reg = flag.String("reg", "http://127.0.0.1:2379", "register etcd address")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		panic(err)
	}

	err = grpclb.Register(*serv, "127.0.0.1", *port, *reg, time.Second*10, 15)
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		log.Printf("receive signal '%v'", s)
		grpclb.UnRegister()
		os.Exit(1)
	}()

	log.Printf("starting hello service at %d", *port)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	s.Serve(lis)
}

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Printf("%v: Receive is %s\n", time.Now(), in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
