package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
	pb "time-grpc/pb"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	// PORT defines port to listen to
	PORT = ":50051"

	// TimeFormat defines format for Time
	TimeFormat = "2006-01-02T15:04:05.000Z"
)

// implements pb.TimeRecord
type server struct{}

// Task is a localised struct stored in tasks map
type Task struct {
	ID, Status int32
	StartTime  time.Time
	TotalTime  float64
}

// data structure - **use Task.ID as identifier**
var tasks = make(map[int32]*Task)

// InitiateTimer implements pb.TimeRecord
func (s *server) InitiateTimer(ctx context.Context, in *pb.TimeRequest) (*pb.Confirmation, error) {
	if in.GetJobID() == 0 {
		id := getNewID()
		status := pb.JobStatus_value[in.GetJobStatus().String()]
		t, err := ptypes.Timestamp(in.GetTimer())
		if err != nil {
			log.Printf("date provided is invalid [ID:%b; TIME:%v]", id, t)
		}
		t = t.Add(time.Hour * 1)
		tasks[id] = &Task{id, status, t, 0.00}
		log.Printf("NEW TASK: %b - start time = %v", id, t)
		return &pb.Confirmation{JobID: id, JobStatus: pb.JobStatus_NEW, Error: pb.Error_CREATED}, nil
	}
	return &pb.Confirmation{JobID: 0, JobStatus: pb.JobStatus_NONE, Error: pb.Error_BADREQUEST}, nil
}

// InitiateTimer implements pb.TimeRecord
func (s *server) CompleteTimer(ctx context.Context, in *pb.TimeRequest) (*pb.Confirmation, error) {
	id := in.GetJobID()
	if _, exists := tasks[id]; exists {
		t, err := time.Parse(TimeFormat, in.GetTimer().String())
		if err != nil {
			log.Printf("date provided is invalid [ID:%b; TIME:%v]", id, t)
		}
		dur := t.Sub(tasks[id].StartTime).Hours()
		tasks[id].TotalTime += dur
		return &pb.Confirmation{JobID: id, JobStatus: pb.JobStatus_FINISHED, Error: pb.Error_OK}, nil
	}
	return &pb.Confirmation{JobID: id, JobStatus: in.GetJobStatus(), Error: pb.Error_NOTFOUND}, nil
}

// Returns new ID using the length of tasks, as int32
func getNewID() int32 {
	return int32(len(tasks) + 1)
}

func main() {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTimeRecordServer(s, &server{})

	go func() {
		// Register reflection service on gRPC server.
		reflection.Register(s)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	<-stop
	fmt.Printf("shutting down service")

	// TODO: save to file etc.
	s.Stop()
}
