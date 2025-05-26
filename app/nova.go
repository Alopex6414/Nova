package app

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/gin-gonic/gin"
	"net/http"
	. "nova/configure"
	"os"
	"strconv"
)

type Nova struct {
	conf  *Config
	cache *Cache
}

func New() *Nova {
	return &Nova{
		conf:  NewConfig("./configure/nova_configure.yaml"),
		cache: NewCache(),
	}
}

func (nova *Nova) Init() {
	err := nova.conf.LoadConfig()
	if err != nil {
		panic(err)
	}
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
		err := server.ListenAndServeTLS(certFile, keyFile)
		if err != nil {
			panic(err)
		}
	} else {
		// start http service
		port := nova.conf.Configure.Port
		server := &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: router,
		}
		// listen and server
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}
}
