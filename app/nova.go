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
	rc    *RedisCache
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
	// create redis data cache
	if nova.conf.Configure.Cache.CacheType == "redis" {
		logger.Info("Create Nova redis data cache...")
		nova.rc, err = NewRedisCache()
		if err != nil {
			logger.Fatalf("Failed to create redis data cache: %s\n", err)
			fmt.Printf("Failed to create redis data cache: %s\n", err)
			os.Exit(4)
		}
		logger.Info("Successfully create redis data cache.")
	}
	// create database
	logger.Info("Create Nova database...")
	nova.db, err = NewDB("file:nova.db?cache=shared")
	if err != nil {
		logger.Fatalf("Failed to create database: %s\n", err)
		fmt.Printf("Failed to create database: %s\n", err)
		os.Exit(5)
	}
	logger.Info("Successfully create database.")
	// create tables
	logger.Info("Create tables...")
	err = nova.db.CreateTables()
	if err != nil {
		logger.Fatalf("Failed to create tables: %s\n", err)
		fmt.Printf("Failed to create tables: %s\n", err)
		os.Exit(6)
	}
	logger.Info("Successfully create tables.")
	// query users from database
	logger.Info("Query users from database...")
	users, err := nova.db.QueryUsers()
	if err != nil {
		logger.Fatalf("Failed to query users from database: %s\n", err)
		fmt.Printf("Failed to query users from database: %s\n", err)
		os.Exit(7)
	}
	for _, user := range users {
		nova.cache.userCache.userSet = append(nova.cache.userCache.userSet, *user)
	}
	logger.Info("Successfully query users from database.")
	// query single-choice questions from database
	logger.Info("Query single-choice questions from database...")
	singleChoiceQuestions, err := nova.db.QueryQuestionsSingleChoice()
	if err != nil {
		logger.Fatalf("Failed to query single-choice questions from database: %s\n", err)
		fmt.Printf("Failed to query single-choice questions from database: %s\n", err)
		os.Exit(8)
	}
	for _, question := range singleChoiceQuestions {
		nova.cache.questionsCache.singleChoiceCache.singleChoiceSet = append(nova.cache.questionsCache.singleChoiceCache.singleChoiceSet, *question)
	}
	logger.Info("Successfully query single-choice questions from database.")
	// query multiple-choice questions from database
	logger.Info("Query multiple-choice questions from database...")
	multipleChoiceQuestions, err := nova.db.QueryQuestionsMultipleChoice()
	if err != nil {
		logger.Fatalf("Failed to query multiple-choice questions from database: %s\n", err)
		fmt.Printf("Failed to query multiple-choice questions from database: %s\n", err)
		os.Exit(9)
	}
	for _, question := range multipleChoiceQuestions {
		nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet = append(nova.cache.questionsCache.multipleChoiceCache.multipleChoiceSet, *question)
	}
	logger.Info("Successfully query multiple-choice questions from database.")
	// query judgement questions from database
	logger.Info("Query judgement questions from database...")
	judgementQuestions, err := nova.db.QueryQuestionsJudgement()
	if err != nil {
		logger.Fatalf("Failed to query judgement questions from database: %s\n", err)
		fmt.Printf("Failed to query judgement questions from database: %s\n", err)
		os.Exit(10)
	}
	for _, question := range judgementQuestions {
		nova.cache.questionsCache.judgementCache.judgementSet = append(nova.cache.questionsCache.judgementCache.judgementSet, *question)
	}
	logger.Info("Successfully query judgement questions from database.")
	// query essay questions from database
	logger.Info("Query essay questions from database...")
	essayQuestions, err := nova.db.QueryQuestionsEssay()
	if err != nil {
		logger.Fatalf("Failed to query essay questions from database: %s\n", err)
		fmt.Printf("Failed to query essay questions from database: %s\n", err)
		os.Exit(11)
	}
	for _, question := range essayQuestions {
		nova.cache.questionsCache.essayCache.essaySet = append(nova.cache.questionsCache.essayCache.essaySet, *question)
	}
	logger.Info("Successfully query essay questions from database.")
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
		novaService.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "hello Nova\n") })
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
		/* question management */
		// questionId related
		novaService.POST("/question/Id", nova.HandleCreateQuestionId)
		// question related
		novaService.POST("/question/single-choice/:Id", nova.HandleCreateQuestionSingleChoice)
		novaService.PUT("/question/single-choice/:Id", nova.HandleUpdateQuestionSingleChoice)
		novaService.DELETE("/question/single-choice/:Id", nova.HandleDeleteQuestionSingleChoice)
		novaService.PATCH("/question/single-choice/:Id", nova.HandleModifyQuestionSingleChoice)
		novaService.GET("/question/single-choice/:Id", nova.HandleQueryQuestionSingleChoice)
		novaService.POST("/question/multiple-choice/:Id", nova.HandleCreateQuestionMultipleChoice)
		novaService.PUT("/question/multiple-choice/:Id", nova.HandleUpdateQuestionMultipleChoice)
		novaService.DELETE("/question/multiple-choice/:Id", nova.HandleDeleteQuestionMultipleChoice)
		novaService.PATCH("/question/multiple-choice/:Id", nova.HandleModifyQuestionMultipleChoice)
		novaService.GET("/question/multiple-choice/:Id", nova.HandleQueryQuestionMultipleChoice)
		novaService.POST("/question/judgement/:Id", nova.HandleCreateQuestionJudgement)
		novaService.PUT("/question/judgement/:Id", nova.HandleUpdateQuestionJudgement)
		novaService.DELETE("/question/judgement/:Id", nova.HandleDeleteQuestionJudgement)
		novaService.PATCH("/question/judgement/:Id", nova.HandleModifyQuestionJudgement)
		novaService.GET("/question/judgement/:Id", nova.HandleQueryQuestionJudgement)
		novaService.POST("/question/essay/:Id", nova.HandleCreateQuestionEssay)
		novaService.PUT("/question/essay/:Id", nova.HandleUpdateQuestionEssay)
		novaService.DELETE("/question/essay/:Id", nova.HandleDeleteQuestionEssay)
		novaService.PATCH("/question/essay/:Id", nova.HandleModifyQuestionEssay)
		novaService.GET("/question/essay/:Id", nova.HandleQueryQuestionEssay)
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
				logger.Panicf("Failed to read CA file: %s\n", err)
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
		// creat timeout context
		logger.Info("Create timeout context...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// stop & clean up resources
		logger.Info("Stop server & Clean Up resources...")
		if err := nova.Stop(); err != nil {
			logger.Fatalf("Server forced to shutdown: %v\n", err)
		}
		logger.Info("Server stop & clean up successfully")
		// graceful Shutting down server
		logger.Info("Shutting down server...")
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
		// creat timeout context
		logger.Info("Create timeout context...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// stop & clean up resources
		logger.Info("Stop server & Clean Up resources...")
		if err := nova.Stop(); err != nil {
			logger.Fatalf("Server forced to shutdown: %v\n", err)
		}
		logger.Info("Server stop & clean up successfully")
		// graceful Shutting down server
		logger.Info("Shutting down server...")
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Server forced to shutdown: %v\n", err)
		}
		logger.Info("Server exiting")
	}
}

func (nova *Nova) Stop() error {
	// stop redis cache (configure redis)
	if nova.conf.Configure.Cache.CacheType == "redis" {
		// stop redis cache
		if err := nova.rc.Close(); err != nil {
			return err
		}
	}
	// stop sqlite3 database
	if err := nova.db.Close(); err != nil {
		return err
	}
	return nil
}
