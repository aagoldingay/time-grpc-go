package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/aagoldingay/time-grpc-go/pb"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051" //TODO
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTimeRecordClient(conn)

	response, err := c.InitiateTimer(context.Background(), &pb.NewTimeRequest{Timer: ptypes.TimestampNow()})
	fmt.Println(response)
}
