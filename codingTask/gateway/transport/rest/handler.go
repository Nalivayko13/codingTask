package rest

import (
	_ "github.com/Nalivayko13/codingTask/gateway/docs"
	"github.com/Nalivayko13/codingTask/gateway/logging"
	"github.com/Nalivayko13/codingTask/gateway/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service service.GatewayService
	logger  logging.Logger
}

func NewHandler(service service.GatewayService, logger logging.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	auth.POST("/login", h.AuthUser)

	store := router.Group("store")
	store.Use(h.ParseAuthHeader)
	store.GET("/:id", h.GetStore)
	store.GET("/:id/history", h.StoreHistory)
	store.GET("/:id/version/:version_id", h.StoreVersion)
	store.POST("/", h.Create)
	store.POST("/:id/version", h.CreateVersion)
	store.DELETE("/:id/", h.Delete)
	store.DELETE("/:id/version/:version_id/:creator", h.DeleteByVersion)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}

type ErrorResponse struct {
	Message       string `json:"message"`
	ResponseError string `json:"response_error"`
	Status        string `json:"status"`
}
