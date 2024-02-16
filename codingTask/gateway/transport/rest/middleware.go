package rest

import (
	"errors"
	"github.com/Nalivayko13/codingTask/gateway/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var UrlAuthServiceValidate string

const (
	authHeader = "Authorization"
)

func (h *Handler) ParseAuthHeader(ctx *gin.Context) {

	header := ctx.GetHeader(authHeader)
	if header == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: "empty auth header"})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: "invalid header"})
		return
	}

	if len(headerParts[1]) == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: "token is empty"})
		return
	}

	err := h.validateToken(header)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}

}

func (h *Handler) CorsMiddleware(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "*")
	ctx.Header("Access-Control-Allow-Headers", "*")
	ctx.Header("Content-Type", "application/json")

	if ctx.Request.Method != "OPTIONS" {
		ctx.Next()
	} else {
		ctx.AbortWithStatus(http.StatusOK)
	}
}

func (h *Handler) validateToken(tokenStr string) error {
	header := map[string]string{
		authHeader: tokenStr,
	}

	resp, err := utils.HttpGetCallWithHeader(UrlAuthServiceValidate, header)
	if err != nil {

	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("invalid token")
	}

	return nil
}
