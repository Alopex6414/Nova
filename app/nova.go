package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Nova struct {
}

func New() *Nova {
	return &Nova{}
}

func (nova *Nova) Init() {

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
		// user management
		novaService.POST("/user/:userId", handleCreatUser)
		novaService.DELETE("/user/:userId", handleDeleteUser)
		novaService.PATCH("/user/:userId", handleModifyUser)
		novaService.GET("/user/:userId", handleQueryUser)
	}
	// start http service
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	// listen and server
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
