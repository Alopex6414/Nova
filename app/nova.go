package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	. "nova/configure"
	"nova/logger"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Nova struct {
	conf  *Config
	cache *Cache
	db    *DB
}

func New() *Nova {
	return &Nova{
		conf:  NewConfig("./configure/nova_configure.yaml"),
		cache: NewCache(),
	}
}

func (nova *Nova) Init() {
	// init logger
	cfg := logger.Config{
		Level:      logger.DebugLevel,
		Filename:   "./logs/nova.log",
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
		Console:    false,
	}
	if err := logger.Init(cfg); err != nil {
		fmt.Printf("Failed to init logger: %s\n", err)
		os.Exit(2)
	}
	// load configure
	logger.Info("Loading configuration...")
	err := nova.conf.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load config: %s\n", err)
		fmt.Printf("Failed to load config: %s\n", err)
		os.Exit(3)
	}
	logger.Info("Successfully loaded configuration.")
	logger.Debug("YAML configuration:", nova.conf)
	// create database
	logger.Info("Create Nova database...")
	nova.db, err = NewDB("file:nova.db?cache=shared")
	if err != nil {
		logger.Fatalf("Failed to create database: %s\n", err)
		fmt.Printf("Failed to create database: %s\n", err)
		os.Exit(4)
	}
	logger.Info("Successfully create database.")
	// create tables
	logger.Info("Create tables...")
	err = nova.db.CreateTables()
	if err != nil {
		logger.Fatalf("Failed to create tables: %s\n", err)
		fmt.Printf("Failed to create tables: %s\n", err)
		os.Exit(5)
	}
	logger.Info("Successfully create tables.")
}

func (nova *Nova) Start() {
	// apply default Gin service
	router := gin.Default()
	// apply Gin logger & recovery middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// create router group for nova
	novaService := router.Group("nova/v1")
	{
		novaService.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "hello Gin\n") })
		/* user management */
		// userId related
		novaService.POST("/user/userId", nova.HandleCreateUserId)
		novaService.GET("/user/userId", nova.HandleQueryUserId)
		// user related
		novaService.POST("/user/:userId", nova.HandleCreateUser)
		novaService.PUT("/user/:userId", nova.HandleUpdateUser)
		novaService.DELETE("/user/:userId", nova.HandleDeleteUser)
		novaService.PATCH("/user/:userId", nova.HandleModifyUser)
		novaService.GET("/user/:userId", nova.HandleQueryUser)
		// user login related
		novaService.POST("/user/login/:userId", nova.HandleCreateUserLogin)
	}
	// enable tls settings
	var tlsConfig *tls.Config
	tlsSettings := nova.conf.Configure.TLS
	if tlsSettings.TLSType != "non-tls" {
		var minVersion uint16
		// tls version
		switch tlsSettings.TLSMinVersion {
		case "1.2":
			minVersion = tls.VersionTLS12
		case "1.3":
			minVersion = tls.VersionTLS13
		default:
			minVersion = tls.VersionTLS13
		}
		// one-way tls
		tlsConfig = &tls.Config{
			MinVersion: minVersion,
			CurvePreferences: []tls.CurveID{
				tls.X25519,
				tls.CurveP256,
			},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
		}
		// mutual tls
		if tlsSettings.TLSType == "mutual-tls" {
			// read CA certificate
			caFile := nova.conf.Configure.TLS.CAFile
			caCert, err := os.ReadFile(caFile)
			if err != nil {
				panic(err)
			}
			// create CA certificate pool
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			// specific CA certificate pool
			tlsConfig.ClientCAs = caCertPool
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
		// start https service
		port := nova.conf.Configure.Port
		certFile := tlsSettings.CertFile
		keyFile := tlsSettings.KeyFile
		server := &http.Server{
			Addr:      ":" + strconv.Itoa(port),
			Handler:   router,
			TLSConfig: tlsConfig,
		}
		// listen and server
		go func() {
			if err := server.ListenAndServeTLS(certFile, keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Fatalf("Failed to start server: %s\n", err)
			}
		}()
		// waiting for close server signal
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Info("Shutting down server...")
		// creat timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// graceful Shutting down server
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Server forced to shutdown: %v\n", err)
		}
		logger.Info("Server exiting")
	} else {
		// start http service
		port := nova.conf.Configure.Port
		server := &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: router,
		}
		// listen and server
		go func() {
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Fatalf("Failed to start server: %s\n", err)
			}
		}()
		// waiting for close server signal
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Info("Shutting down server...")
		// creat timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// graceful Shutting down server
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Server forced to shutdown: %v\n", err)
		}
		logger.Info("Server exiting")
	}
}
