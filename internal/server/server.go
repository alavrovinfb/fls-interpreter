package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/alavrovinfb/fls-interpreter/pkg/pb"
	"github.com/alavrovinfb/fls-interpreter/pkg/script"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Interpreter struct {
	logger          *logrus.Logger
	grpcAddr        *string
	grpcPort        *string
	gatewayAddr     *string
	gatewayPort     *string
	gatewayEndpoint *string
	internalAddr    *string
	internalPort    *string
	heathEndpoint   *string
}

func NewInterpreter(logger *logrus.Logger, grpcAddr, grpcPort, gatewayAddr, gatewayPort, gatewayEndpoint, internalAddr,
	internalPort, heathEndpoint *string) *Interpreter {
	return &Interpreter{
		logger:          logger,
		grpcAddr:        grpcAddr,
		grpcPort:        grpcPort,
		gatewayAddr:     gatewayAddr,
		gatewayPort:     gatewayPort,
		gatewayEndpoint: gatewayEndpoint,
		internalAddr:    internalAddr,
		internalPort:    internalPort,
		heathEndpoint:   heathEndpoint,
	}
}

func (is *Interpreter) Run() error {
	doneC := make(chan error)

	go func() { doneC <- is.serveInternal() }()
	go func() { doneC <- is.serveExternal() }()
	go func() { doneC <- is.gatewayServe() }()

	return <-doneC
}

func NewLogger(logLevel string) *logrus.Logger {
	logger := logrus.StandardLogger()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)

	// Set the log level on the default logger based on command line flag
	if level, err := logrus.ParseLevel(logLevel); err != nil {
		logger.Errorf("Invalid %q provided for log level", viper.GetString("logging.level"))
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(level)
	}

	return logger
}

// ServeInternal runs internal endpoints like pprof, metrics, health
func (is *Interpreter) serveInternal() error {
	logger := is.logger
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	mux.Handle(*is.heathEndpoint, http.HandlerFunc(health))

	internalAddr := net.JoinHostPort(*is.internalAddr, *is.internalPort)
	logger.Infof("internal server started on: %s", internalAddr)

	return http.ListenAndServe(internalAddr, mux)
}

// ServeExternal builds and runs the server that listens on ServerAddress and GatewayAddress
func (is *Interpreter) serveExternal() error {
	logger := is.logger
	grpcServer, err := NewGRPCServer(logger)
	if err != nil {
		logger.Fatalln(err)
	}
	grpc_prometheus.Register(grpcServer)
	externalAddr := net.JoinHostPort(*is.grpcAddr, *is.grpcPort)
	grpcL, err := net.Listen("tcp", externalAddr)
	if err != nil {
		logger.Fatalln(err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		logger.Infof("got signal %v, attempting graceful shutdown", <-sigCh)
		grpcServer.GracefulStop()
		logger.Info("server gracefully stopped")

		os.Exit(0)
	}()

	logger.Printf("serving gRPC at %s", externalAddr)

	return grpcServer.Serve(grpcL)
}

func forwardResponseOption(_ context.Context, w http.ResponseWriter, _ protoreflect.ProtoMessage) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	return nil
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (is *Interpreter) gatewayServe() error {
	logger := is.logger
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	gmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: false,
			},
		}))
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	ro := gmux.GetForwardResponseOptions()
	ro = append(ro, forwardResponseOption)
	grpcAddr := net.JoinHostPort(*is.grpcAddr, *is.grpcPort)
	err := pb.RegisterFlsInterpreterHandlerFromEndpoint(ctx, gmux, grpcAddr, opts)
	if err != nil {
		logger.Error(err)
		return err
	}

	prefix := *is.gatewayEndpoint
	mux := http.NewServeMux()
	mux.Handle(prefix, http.StripPrefix(prefix[:len(prefix)-1], gmux))
	gatewayAddr := net.JoinHostPort(*is.gatewayAddr, *is.gatewayPort)
	if err := gmux.HandlePath("POST", "/files", handleBinaryFileUpload); err != nil {
		return err
	}
	logger.Printf("serving http at %s", gatewayAddr)
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(gatewayAddr, mux)
}

// custom handler for grpc gateway, unfortunately grpc-gateway doesn't support files uploading through grpc streams (
func handleBinaryFileUpload(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse form: %s", err.Error()), http.StatusBadRequest)
		return
	}

	f, _, err := r.FormFile("script")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get file 'script': %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer f.Close()

	script.Body.RLock()
	defer script.Body.RUnlock()
	script.Vars.RLock()
	defer script.Vars.RUnlock()
	if err := script.Body.Run(f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outPB, err := pb.OutToPB(script.Body.Out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(outPB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content/type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
