package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	pb "github.com/aagoldingay/time-grpc-go/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// protoc -I pb/ pb/service.proto --go_out=plugins=grpc:pb

const (
	// Port defines port to listen to
	Port = ":50051"

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

// InitiateTimer implements pb.TimeRecord - accepts: *pb.TimeRequest
func (s *server) CompleteTimer(ctx context.Context, in *pb.TimeRequest) (*pb.Confirmation, error) {
	id := in.GetJobID()
	if _, exists := tasks[id]; exists {
		tasks[id].TotalTime += time.Since(tasks[id].StartTime).Hours()
		tasks[id].Status = pb.JobStatus_value["FINISHED"]
		log.Printf("TASK COMPLETED: %d - duration = %.2f hour(s)", id, tasks[id].TotalTime)
		return &pb.Confirmation{JobID: id, JobStatus: pb.JobStatus_FINISHED, Error: pb.Error_OK}, nil
	}
	return &pb.Confirmation{JobID: id, JobStatus: pb.JobStatus_NONE, Error: pb.Error_NOTFOUND}, nil
}

// InitiateTimer implements pb.TimeRecord - accepts: *pb.NewTimeRequest
func (s *server) InitiateTimer(ctx context.Context, in *pb.NewTimeRequest) (*pb.Confirmation, error) {
	id := getNewID()
	tasks[id] = &Task{ID: id, Status: pb.JobStatus_value["NEW"], TotalTime: 0.00}
	log.Printf("NEW TASK: %d", id)
	return &pb.Confirmation{JobID: id, JobStatus: pb.JobStatus_NEW, Error: pb.Error_CREATED}, nil
}

// StartTimer implements pb.TimeRecord - accepts: *pb.TimeRequest
func (s *server) StartTimer(ctx context.Context, in *pb.TimeRequest) (*pb.Confirmation, error) {
	id := in.GetJobID()
	if _, exists := tasks[id]; exists {
		tasks[id].StartTime = time.Now()
		tasks[id].Status = pb.JobStatus_value["STARTED"]
		log.Printf("TASK %d STARTED", id)
		return &pb.Confirmation{JobID: id, JobStatus: pb.JobStatus_STARTED, Error: pb.Error_OK}, nil
	}
	return &pb.Confirmation{JobID: id, JobStatus: pb.JobStatus_NONE, Error: pb.Error_NOTFOUND}, nil
}

// UpdateTimer implements pb.TimeRecord - accepts: *pb.TimeRequest
// Will change the status of a present job to started or paused
func (s *server) UpdateTimer(ctx context.Context, in *pb.TimeRequest) (*pb.Confirmation, error) {
	id := in.GetJobID()
	if _, exists := tasks[id]; exists {
		if tasks[id].Status == pb.JobStatus_value["STARTED"] {
			tasks[id].Status = pb.JobStatus_value["PAUSED"]
			return &pb.Confirmation{JobID: id, JobStatus: pb.JobStatus_PAUSED, Error: pb.Error_OK}, nil
		}
		if tasks[id].Status == pb.JobStatus_value["PAUSED"] {
			tasks[id].Status = pb.JobStatus_value["STARTED"]
			return &pb.Confirmation{JobID: id, JobStatus: pb.JobStatus_STARTED, Error: pb.Error_OK}, nil
		}
		//this return will happen when a task is not started or paused
		return &pb.Confirmation{JobID: id, JobStatus: pb.JobStatus(tasks[id].Status), Error: pb.Error_BADREQUEST}, nil
	}
	return &pb.Confirmation{JobID: id, JobStatus: pb.JobStatus_NONE, Error: pb.Error_NOTFOUND}, nil
}

// Returns new ID using the length of tasks, as int32
func getNewID() int32 {
	return int32(len(tasks) + 1)
}

func main() {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	lis, err := net.Listen("tcp", Port)
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
