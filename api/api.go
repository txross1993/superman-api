package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/txross1993/superman-api/models"
	"github.com/txross1993/superman-api/superman"
)

// Config holds the api configuration for the bind host and port and the
// superman service
type Config struct {
	Host     string
	Port     string
	Superman *superman.Service
}

// API configures the superman api
type API struct {
	Config
	router *gin.Engine
}

// NewAPI configures a new instance of the superman api
func NewAPI(cfg Config) *API {
	router := gin.Default()
	api := &API{
		Config: cfg,
		router: router,
	}
	api.SetupRoutes()
	return api
}

// Run starts the API server
func (api *API) Run() error {
	return api.router.Run()

}

// SetupRoutes declares the routes and handlers for the API server
func (api *API) SetupRoutes() {
	v1 := api.router.Group("/v1")
	{
		v1.POST("/", api.AnalyzeLoginEvent)
	}
}

// AnalyzeLoginEvent binds the request to the expected format and hands
// the request to the Superman service for analysis
func (api *API) AnalyzeLoginEvent(c *gin.Context) {
	var event models.UserIPAccessEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := api.Superman.AnalyzeEvent(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, resp)
	return
}
