package main

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	//
	defaultMode = "cli"
	// Server
	defaultServerAddress = "0.0.0.0"
	defaultServerPort    = "9090"

	// Gateway
	defaultGatewayEnable      = false
	defaultGatewayAddress     = "0.0.0.0"
	defaultGatewayPort        = "8080"
	defaultGatewayURL         = "/fls-interpreter/v1/"
	defaultGatewaySwaggerFile = "pkg/pb/service.swagger.json"

	// Health
	defaultInternalEnable  = true
	defaultInternalAddress = "0.0.0.0"
	defaultInternalPort    = "8081"
	defaultInternalHealth  = "/healthz"

	// Heartbeat
	defaultKeepaliveTime    = 10
	defaultKeepaliveTimeout = 20

	// Logging
	defaultLoggingLevel = "debug"
)

var (
	defaultFiles = []string{}
	// define flag overrides
	flagMode        = pflag.String("mode", defaultMode, "FLS mode cli or server default cli")
	flagScriptFiles = pflag.StringSlice("script.files", defaultFiles, "comma separated list of files")

	flagServerAddress = pflag.String("server.address", defaultServerAddress, "adress of gRPC server")
	flagServerPort    = pflag.String("server.port", defaultServerPort, "port of gRPC server")

	flagGatewayAddress     = pflag.String("gateway.address", defaultGatewayAddress, "address of gateway server")
	flagGatewayPort        = pflag.String("gateway.port", defaultGatewayPort, "port of gateway server")
	flagGatewayURL         = pflag.String("gateway.endpoint", defaultGatewayURL, "endpoint of gateway server")
	flagGatewaySwaggerFile = pflag.String("gateway.swaggerFile", defaultGatewaySwaggerFile, "directory of swagger.json file")

	flagInternalAddress = pflag.String("internal.address", defaultInternalAddress, "address of internal http server")
	flagInternalPort    = pflag.String("internal.port", defaultInternalPort, "port of internal http server")
	flagInternalHealth  = pflag.String("internal.health", defaultInternalHealth, "endpoint for health checks")

	flagKeepaliveTime    = pflag.Int("config.keepalive.time", defaultKeepaliveTime, "default value, in seconds, of the keepalive time")
	flagKeepaliveTimeout = pflag.Int("config.keepalive.timeout", defaultKeepaliveTimeout, "default value, in seconds, of the keepalive timeout")

	flagLoggingLevel = pflag.String("logging.level", defaultLoggingLevel, "log level of application")
)

func init() {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
