package rest

import (
	"github.com/Nalivayko13/codingTask/gateway/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
)

// @Summary Login
// @Tags AuthUser
// @Description log in user in system with login=qwerty ans pass=secretPass
// @ID AuthUser
// @Accept json
// @Produce json
// @Param input body model.User true "User credentials"
// @Success 200 {object} string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {string} string
// @Router /auth/login [post]
func (h *Handler) AuthUser(c *gin.Context) {
	var user *model.User
	if err := c.BindJSON(&user); err != nil {
		h.logger.Log.Error("Could not binding JSON",
			zap.String("url", "/auth/login"), zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: "Could not binding JSON"})
		log.Printf("Failed to process request: create user: %v", err)
		return
	}

	token, err := h.service.AuthUser(c, user)
	if err != nil {
		h.logger.Log.Error("Internal server error",
			zap.String("url", "/auth/login"), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server cannot process the request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
