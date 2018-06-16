package main

import (
	"context"
	"testing"
	"time"

	pb "github.com/aagoldingay/time-grpc-go/pb"
)

func Test_InitiateTimer(t *testing.T) {
	s := server{}

	for i := 0; i < 2; i++ {
		req := &pb.NewTimeRequest{New: true}
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
	if _, exists := tasks[1]; !exists {
		t.Errorf("Object 1 do not exist in task map")
	}
	if _, exists := tasks[2]; !exists {
		t.Errorf("Object 2 do not exist in task map")
	}
}

func Test_CompleteTimerSuccessful(t *testing.T) {
	s := server{}
	tasks[1] = &Task{ID: 1, StartTime: time.Now(), TotalTime: 0.00}

	stamp, _ := time.Parse(TimeFormat, "2018-06-10T19:30:45.000Z")
	tasks[2] = &Task{ID: 2, StartTime: stamp, TotalTime: 0.10}

	for i := int32(1); i < 3; i++ {
		req := &pb.TimeRequest{JobID: i}
		resp, err := s.CompleteTimer(context.Background(), req)
		if err != nil {
			t.Errorf("CompleteTimer(%v) got unexpected error", req)
		}
		if resp.Error != pb.Error_OK {
			t.Errorf("CompleteTimer(%v)=%v, wanted %v", req, resp.Error, pb.Error_OK)
		}
		if resp.JobStatus != pb.JobStatus_FINISHED {
			t.Errorf("CompleteTimer(%v)=%v, wanted %v", req, resp.JobStatus, pb.JobStatus_FINISHED)
		}
		val, _ := tasks[i]
		if val.Status != 4 {
			t.Errorf("Task %d status is not finished", i)
		}
	}
}

func Test_CompleteTimerUnsuccessful(t *testing.T) {
	s := server{}

	for i := int32(3); i < 5; i++ {
		req := &pb.TimeRequest{JobID: i}
		resp, err := s.CompleteTimer(context.Background(), req)
		if err != nil {
			t.Errorf("CompleteTimer(%v) got unexpected error", req)
		}
		if resp.Error != pb.Error_NOTFOUND {
			t.Errorf("CompleteTimer(%v)=%v, wanted %v", req, resp.Error, pb.Error_NOTFOUND)
		}
		if resp.JobStatus != pb.JobStatus_NONE {
			t.Errorf("CompleteTimer(%v)=%v, wanted %v", req, resp.JobStatus, pb.JobStatus_NONE)
		}
	}
}
func Test_StartTimerSuccessful(t *testing.T) {
	s := server{}
	tasks[1] = &Task{ID: 1, TotalTime: 0.00}

	req := &pb.TimeRequest{JobID: 1}
	resp, err := s.StartTimer(context.Background(), req)
	if err != nil {
		t.Errorf("StartTimer(%v) got unexpected error", req)
	}
	if resp.Error != pb.Error_OK {
		t.Errorf("StartTimer(%v)=%v, wanted %v", req, resp.Error, pb.Error_OK)
	}
	if resp.JobStatus != pb.JobStatus_STARTED {
		t.Errorf("StartTimer(%v)=%v, wanted %v", req, resp.JobStatus, pb.JobStatus_STARTED)
	}
}

func Test_StartTimerUnsuccessful(t *testing.T) {
	s := server{}

	req := &pb.TimeRequest{JobID: 5}
	resp, err := s.StartTimer(context.Background(), req)
	if err != nil {
		t.Errorf("StartTimer(%v) got unexpected error", req)
	}
	if resp.Error != pb.Error_NOTFOUND {
		t.Errorf("StartTimer(%v)=%v, wanted %v", req, resp.Error, pb.Error_NOTFOUND)
	}
	if resp.JobStatus != pb.JobStatus_NONE {
		t.Errorf("StartTimer(%v)=%v, wanted %v", req, resp.JobStatus, pb.JobStatus_NONE)
	}
}

func Test_UpdateTimerToStarted(t *testing.T) {
	s := server{}
	tasks[1] = &Task{ID: 1, Status: int32(pb.JobStatus_PAUSED), StartTime: time.Now(), TotalTime: 0.00}

	req := &pb.TimeRequest{JobID: 1}
	resp, err := s.UpdateTimer(context.Background(), req)
	if err != nil {
		t.Errorf("UpdateTimer(%v) got unexpected error", req)
	}
	if resp.Error != pb.Error_OK {
		t.Errorf("UpdateTimer(%v)=%v, wanted %v", req, resp.Error, pb.Error_OK)
	}
	if resp.JobStatus != pb.JobStatus_STARTED {
		t.Errorf("UpdateTimer(%v)=%v, wanted %v", req, resp.JobStatus, pb.JobStatus_STARTED)
	}
	val, _ := tasks[1]
	if val.Status != int32(pb.JobStatus_STARTED) {
		t.Errorf("Task %d status is not started", 1)
	}
}

func Test_UpdateTimerToPaused(t *testing.T) {
	s := server{}
	tasks[1] = &Task{ID: 1, Status: int32(pb.JobStatus_STARTED), StartTime: time.Now(), TotalTime: 0.00}

	req := &pb.TimeRequest{JobID: 1}
	resp, err := s.UpdateTimer(context.Background(), req)
	if err != nil {
		t.Errorf("UpdateTimer(%v) got unexpected error", req)
	}
	if resp.Error != pb.Error_OK {
		t.Errorf("UpdateTimer(%v)=%v, wanted %v", req, resp.Error, pb.Error_OK)
	}
	if resp.JobStatus != pb.JobStatus_PAUSED {
		t.Errorf("UpdateTimer(%v)=%v, wanted %v", req, resp.JobStatus, pb.JobStatus_PAUSED)
	}
	val, _ := tasks[1]
	if val.Status != int32(pb.JobStatus_PAUSED) {
		t.Errorf("Task %d status is not paused", 1)
	}
}

func Test_UpdateTimerInvalidStatus(t *testing.T) {
	s := server{}
	tasks[1] = &Task{ID: 1, Status: int32(pb.JobStatus_FINISHED), StartTime: time.Now(), TotalTime: 0.00}

	req := &pb.TimeRequest{JobID: 1}
	resp, err := s.UpdateTimer(context.Background(), req)
	if err != nil {
		t.Errorf("UpdateTimer(%v) got unexpected error", req)
	}
	if resp.Error != pb.Error_BADREQUEST {
		t.Errorf("UpdateTimer(%v)=%v, wanted %v", req, resp.Error, pb.Error_BADREQUEST)
	}
	if resp.JobStatus != pb.JobStatus_FINISHED {
		t.Errorf("UpdateTimer(%v)=%v, wanted %v", req, resp.JobStatus, pb.JobStatus_FINISHED)
	}
	val, _ := tasks[1]
	if val.Status != int32(pb.JobStatus_FINISHED) {
		t.Errorf("Task %d status updated", 1)
	}
}

func Test_UpdateTimerUnsuccessful(t *testing.T) {
	s := server{}
	req := &pb.TimeRequest{JobID: 6}
	resp, err := s.UpdateTimer(context.Background(), req)
	if err != nil {
		t.Errorf("UpdateTimer(%v) got unexpected error", req)
	}
	if resp.Error != pb.Error_NOTFOUND {
		t.Errorf("UpdateTimer(%v)=%v, wanted %v", req, resp.Error, pb.Error_NOTFOUND)
	}
	if resp.JobStatus != pb.JobStatus_NONE {
		t.Errorf("UpdateTimer(%v)=%v, wanted %v", req, resp.JobStatus, pb.JobStatus_NONE)
	}
}
