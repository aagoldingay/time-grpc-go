package main

import (
	"context"
	"testing"

	pb "github.com/aagoldingay/time-grpc-go/pb"
	"github.com/golang/protobuf/ptypes"
)

func Test_InitiateTimer(t *testing.T) {
	s := server{}

	tests := []struct {
		id, status int32
	}{
		{
			id:     0,
			status: 0,
		},
		{
			id:     0,
			status: 0,
		},
	}

	for _, tst := range tests {
		req := &pb.TimeRequest{JobID: tst.id, Timer: ptypes.TimestampNow(), JobStatus: pb.JobStatus_NEW}
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
