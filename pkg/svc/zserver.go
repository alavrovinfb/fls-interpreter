package svc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/alavrovinfb/fls-interpreter/pkg/pb"
)

const (
	// version is the current version of the service
	version = "0.0.1"
)

// Default implementation of the FlsInterpreter server interface
type server struct{}

// GetVersion returns the current version of the service
func (server) GetVersion(context.Context, *empty.Empty) (*pb.VersionResponse, error) {
	return &pb.VersionResponse{Version: version}, nil
}

// implement script execute endpoint
func (server) ScriptExecute(context.Context, *pb.ScriptRequest) (*pb.ScriptResponse, error) {
	return nil, nil
}

// NewBasicServer returns an instance of the default server interface
func NewBasicServer() (pb.FlsInterpreterServer, error) {
	return &server{}, nil
}
