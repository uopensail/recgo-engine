package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/recgo-engine/resources"
	"github.com/uopensail/recgo-engine/services"
	"github.com/uopensail/recgo-engine/strategy"
	"github.com/uopensail/ulib/prome"
	"github.com/uopensail/ulib/zlog"
	"go.uber.org/zap"
)

var gitCommitInfo = "" // injected at build time
var buildTime = ""     // injected at build time

// PingPongHandler returns a basic "PONG" response for health checking.
func PingPongHandler(c *gin.Context) {
	pStat := prome.NewStat("PingPongHandler")
	defer pStat.End()
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "PONG"})
}

// GitHashHandler returns the git commit info and build time of the running build.
func GitHashHandler(c *gin.Context) {
	pStat := prome.NewStat("GitHashHandler")
	defer pStat.End()
	c.JSON(http.StatusOK, gin.H{
		"git_info":   gitCommitInfo,
		"build_time": buildTime,
	})
}

// run initializes configuration, resources, strategy, and starts the HTTP server.
func run(configPath, logDir string) *services.Services {
	// Load application config
	if err := config.AppConfigInstance.Init(configPath); err != nil {
		panic(err)
	}

	// Initialize resource manager and recommendation strategy
	resources.ResourceManagerInstance = resources.NewResourceManager(config.AppConfigInstance)
	strategy.StrategyInstance = strategy.NewStrategy(config.AppConfigInstance)

	// Initialize logger
	zlog.InitLogger(config.AppConfigInstance.ProjectName, config.AppConfigInstance.Debug, logDir)

	// Create services
	srv := services.NewServices()

	// Start HTTP server
	httpSrv := newHTTPServer(srv.RegisterGinRouter)
	go func() {
		if err := httpSrv.Start(); err != nil {
			zlog.LOG.Fatal("HTTP server error", zap.Error(err))
		}
	}()

	return srv
}

// newHTTPServer sets up the Gin engine and routes, and wraps it as a simple HTTP server.
func newHTTPServer(registerFunc func(*gin.Engine)) *httpServer {
	engine := gin.New()
	engine.Use(gin.Recovery())

	// Health check endpoints
	engine.GET("/ping", PingPongHandler)
	engine.GET("/git_hash", GitHashHandler)

	// Business routes
	registerFunc(engine)

	return &httpServer{
		engine: engine,
		addr:   fmt.Sprintf(":%d", config.AppConfigInstance.ServerConfig.HttpServerConfig.HTTPPort),
	}
}

// httpServer wraps Gin engine for easy start/stop.
type httpServer struct {
	engine *gin.Engine
	addr   string
}

// Start runs the HTTPS server.
func (s *httpServer) Start() error {
	return http.ListenAndServe(s.addr, s.engine)
}

// runPProf starts pprof if a port is configured.
func runPProf(port int) {
	if port > 0 {
		go func() { _ = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil) }()
	}
}

// runProme starts Prometheus metrics exporter.
func runProme(projectName string, port int) *prome.Exporter {
	exporter := prome.NewExporter(projectName)
	go func() {
		if err := exporter.Start(port); err != nil {
			panic(err)
		}
	}()
	return exporter
}

// main is the entry point of the application.
func main() {
	// Read CLI flags
	configPath := flag.String("config", "conf/config.toml", "Configuration file path")
	logDir := flag.String("log", "./logs", "Log directory")
	flag.Parse()

	if *configPath == "" {
		panic("Configuration file path cannot be empty")
	}

	// Initialize and start application
	appService := run(*configPath, *logDir)

	if config.AppConfigInstance.ProjectName == "" {
		panic("config.ProjectName is empty")
	}

	// Start debug and monitoring services
	runPProf(config.AppConfigInstance.PProfPort)
	promeExport := runProme(config.AppConfigInstance.ProjectName, config.AppConfigInstance.PromePort)

	// Wait for termination signal
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "app running...")
	<-signalCh

	// Graceful shutdown
	appService.Close()
	promeExport.Close()
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "app exited")
}
