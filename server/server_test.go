package main

import (
	"context"
	"testing"
	"time"

	pb "github.com/aagoldingay/time-grpc-go/pb"
	"github.com/golang/protobuf/ptypes"
)

func Test_InitiateTimer(t *testing.T) {
	s := server{}
	stamp, _ := time.Parse(TimeFormat, "2018-06-10T19:30:45.000Z")
	tests := []struct {
		time time.Time
	}{
		{
			time: time.Now(),
		},
		{
			time: stamp,
		},
	}

	for _, tst := range tests {
		ts, _ := ptypes.TimestampProto(tst.time)
		req := &pb.NewTimeRequest{Timer: ts}
		resp, err := s.InitiateTimer(context.Background(), req)
		if err != nil {
			t.Errorf("InitiateTimer(%v) got unexpected error", req)
		}
		if resp.Error != pb.Error_CREATED {
			t.Errorf("InitaiteTimer(%v)=%v, wanted %v", req, resp.Error, pb.Error_CREATED)
		}
		if resp.JobStatus != pb.JobStatus_NEW {
			t.Errorf("InitiateTimer(%v)=%v, wanted %v", req, resp.JobStatus, pb.JobStatus_NEW)
		}
	}
}

func Test_CompleteTimerSuccessful(t *testing.T) {
	s := server{}
	tasks[1] = &Task{ID: 1, StartTime: time.Now(), TotalTime: 0.00}

	stamp, _ := time.Parse(TimeFormat, "2018-06-10T19:30:45.000Z")
	tasks[2] = &Task{ID: 2, StartTime: stamp, TotalTime: 0.10}
	tests := []struct {
		id   int32
		time time.Time
	}{
		{
			id:   1,
			time: time.Now(),
		},
		{
			id:   2,
			time: stamp,
		},
	}

	for _, tst := range tests {
		ts, _ := ptypes.TimestampProto(tst.time)
		req := &pb.CompleteRequest{JobID: tst.id, Timer: ts}
		resp, err := s.CompleteTimer(context.Background(), req)
		if err != nil {
			t.Errorf("InitiateTimer(%v) got unexpected error", req)
		}
		if resp.Error != pb.Error_OK {
			t.Errorf("InitaiteTimer(%v)=%v, wanted %v", req, resp.Error, pb.Error_OK)
		}
		if resp.JobStatus != pb.JobStatus_FINISHED {
			t.Errorf("InitiateTimer(%v)=%v, wanted %v", req, resp.JobStatus, pb.JobStatus_FINISHED)
		}
	}
}

func Test_CompleteTimerUnSuccessful(t *testing.T) {
	s := server{}
	tasks[1] = &Task{ID: 1, StartTime: time.Now(), TotalTime: 0.00}

	stamp, _ := time.Parse(TimeFormat, "2018-06-10T19:30:45.000Z")
	tasks[2] = &Task{ID: 2, StartTime: stamp, TotalTime: 0.10}
	tests := []struct {
		id   int32
		time time.Time
	}{
		{
			id:   0,
			time: time.Now(),
		},
		{
			id:   3,
			time: stamp,
		},
	}

	for _, tst := range tests {
		ts, _ := ptypes.TimestampProto(tst.time)
		req := &pb.CompleteRequest{JobID: tst.id, Timer: ts}
		resp, err := s.CompleteTimer(context.Background(), req)
		if err != nil {
			t.Errorf("InitiateTimer(%v) got unexpected error", req)
		}
		if resp.Error != pb.Error_NOTFOUND {
			t.Errorf("InitaiteTimer(%v)=%v, wanted %v", req, resp.Error, pb.Error_NOTFOUND)
		}
		if resp.JobStatus != pb.JobStatus_NEW {
			t.Errorf("InitiateTimer(%v)=%v, wanted %v", req, resp.JobStatus, pb.JobStatus_NEW)
		}
	}
}
