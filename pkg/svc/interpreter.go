package svc

import (
	"context"
	"net/http"

	"github.com/alavrovinfb/fls-interpreter/pkg/pb"
	"github.com/alavrovinfb/fls-interpreter/pkg/script"
	"github.com/golang/protobuf/ptypes/empty"
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

// ScriptExecute executes FLS script then returns results
func (server) ScriptExecute(_ context.Context, req *pb.ScriptRequest) (*pb.ScriptResponse, error) {
	script.Body.RLock()
	defer script.Body.RUnlock()
	script.Vars.RLock()
	defer script.Vars.RUnlock()
	if err := script.Parse(req.Script.AsMap(), script.Vars, script.Body); err != nil {
		return nil, err
	}
	script.Body.RestOut()
	if err := script.Body.Execute(script.InitFunc, nil); err != nil {
		return nil, err
	}

	res, err := pb.OutToPB(script.Body.Out)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetVersion returns the current version of the service
func (server) Reset(context.Context, *empty.Empty) (*pb.ResetResponse, error) {
	script.Body.Reset()
	script.Vars.Reset()

	return &pb.ResetResponse{
		Message: "reset done",
		Status:  http.StatusOK,
	}, nil
}

// NewBasicServer returns an instance of the default server interface
func NewBasicServer() (pb.FlsInterpreterServer, error) {
	return &server{}, nil
}
