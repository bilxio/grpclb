/*
* Copyright 2017 bilxio.
*
* @File: client.go
* @Author: Bill Xiong
* @Date:   2017-10-17 12:25:55
* @Last Modified by:   Bill Xiong
* @Last Modified time: 2017-10-17 12:28:00
*/

package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	grpclb "github.com/bilxio/grpclb/naming/etcdv3"
	pb "github.com/bilxio/grpclb/example/helloworld"
)

var (
	serv = flag.String("service", "hello_service", "service name")
	reg = flag.String("reg", "http://127.0.0.1:2379", "register etcd address")
)

func main() {
	flag.Parse()
	r := grpclb.NewResolver(*serv)
	b := grpc.RoundRobin(r)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, *reg, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		client := pb.NewGreeterClient(conn)
		resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
		if err == nil {
			fmt.Printf("%v: Reply is %s\n", t, resp.Message)
		}
	}
}
