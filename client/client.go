package main

import (
	"context"
	"fmt"
	"log"
	pb "time-grpc/pb"

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

	response, err := c.InitiateTimer(context.Background(), &pb.TimeRequest{JobID: 0, Timer: ptypes.TimestampNow(), JobStatus: pb.JobStatus_NEW})
	fmt.Println(response)
}
