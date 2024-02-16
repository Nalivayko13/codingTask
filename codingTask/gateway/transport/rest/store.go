package rest

import (
	"github.com/Nalivayko13/codingTask/gateway/model"
	"github.com/Nalivayko13/codingTask/gateway/transport/rest/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// @Summary CreateStore
// @Security ApiKeyAuth
// @Description create new store
// @Accept json
// @Produce json
// @Param input body model.Store true "Store"
// @Success 201
// @Failure 400 {object} ErrorResponse
// @Failure 500 {string} string
// @Router /store/ [post]
func (h *Handler) Create(ctx *gin.Context) {
	var store *model.Store
	if err := ctx.BindJSON(&store); err != nil {
		h.logger.Log.Error("Could not binding JSON", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			ErrorResponse{Message: "Could not binding JSON"})
		return
	}
	if err := h.service.CreateStore(ctx.Request.Context(), store); err != nil {
		h.logger.Log.Error("Internal server error", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "ok"})
}

// @Summary CreateVersion
// @Security ApiKeyAuth
// @Description create new store version
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param input body model.Version true "Version"
// @Success 201
// @Failure 400 {object} ErrorResponse
// @Failure 500 {string} string
// @Router /store/{id}/version [post]
func (h *Handler) CreateVersion(ctx *gin.Context) {
	storeID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Log.Error("Incorrect storeID input", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"message": err.Error()})
		return
	}
	var version *model.Version
	if err := ctx.BindJSON(&version); err != nil {
		h.logger.Log.Error("Could not binding JSON", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			ErrorResponse{Message: "Could not binding JSON"})
		return
	}
	if storeID != 0 {
		version.StoreID = storeID
	}
	if err := h.service.CreateVersion(ctx.Request.Context(), version); err != nil {
		h.logger.Log.Error("Internal server error", zap.Int("id", storeID), zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// @Summary DeleteStore
// @Security ApiKeyAuth
// @Description delete store by setting flag is_deleted=true
// @Param id path int true "id"
// @Produce json
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500 {string} string
// @Router /store/{id}  [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	storeID := ctx.Param("id")
	if storeID == "" {
		h.logger.Log.Error("incorrect creator input: empty storeID field")
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"message": "incorrect storeID"})
		return
	}

	if err := h.service.DeleteStore(ctx, storeID); err != nil {
		h.logger.Log.Error("Internal server error", zap.String("id", storeID), zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// @Summary DeleteVersion
// @Security ApiKeyAuth
// @Description delete store version by setting flag is_deleted=true
// @Param id path int true "id"
// @Param version_id path int true "version_id"
// @Param creator path int true "creator"
// @Produce json
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500 {string} string
// @Router /store/{id}/verion/{version_id}/{creator}  [delete]
func (h *Handler) DeleteByVersion(ctx *gin.Context) {
	storeID := ctx.Param("id")
	if storeID == "" {
		h.logger.Log.Error("incorrect creator input: empty storeID field")
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"message": "incorrect storeID"})
		return
	}
	versionID := ctx.Param("version_id")
	if storeID == "" {
		h.logger.Log.Error("incorrect creator input: empty versionID field")
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"message": "incorrect versionID"})
		return
	}
	creator := ctx.Param("creator")
	if creator == "" {
		h.logger.Log.Error("incorrect creator input: empty creator field")
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"message": "input creator"})
		return
	}

	if err := h.service.DeleteVersion(ctx, creator, storeID, versionID); err != nil {
		h.logger.Log.Error("Internal server error", zap.String("id", storeID), zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// @Summary GetStore
// @Security ApiKeyAuth
// @Description get store by id
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.Store
// @Failure 400 {object} ErrorResponse
// @Failure 500 {string} string
// @Failure default {string} string
// @Router /store/{id} [get]
func (h *Handler) GetStore(ctx *gin.Context) {
	storeID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Log.Error("Incorrect storeID input", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"message": err.Error()})
		return
	}

	store, err := h.service.GetStore(ctx.Request.Context(), storeID)
	if err != nil {
		h.logger.Log.Error("internal server error", zap.Int("id", storeID), zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, store)
}

// @Summary StoreHistory
// @Security ApiKeyAuth
// @Description get store history by id
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500 {string} string
// @Failure default {string} string
// @Router /store/{id}/history [get]
func (h *Handler) StoreHistory(ctx *gin.Context) {
	storeID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Log.Error("incorrect storeID input", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"message": "incorrect storeID input"})
		return
	}

	versions, store, err := h.service.GetHistory(ctx.Request.Context(), storeID)
	if err != nil {
		h.logger.Log.Error("internal server error", zap.Int("id", storeID), zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK,
		gin.H{
			"versions": versions,
			"store":    store,
		})
}

// @Summary StoreVersion
// @Security ApiKeyAuth
// @Description get store version by id
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param version_id path int true "version_id"
// @Success 200 {object} response.VersionResp
// @Failure 400 {object} ErrorResponse
// @Failure 500 {string} string
// @Failure default {string} string
// @Router /store/{id}/version/{version_id} [get]
func (h *Handler) StoreVersion(ctx *gin.Context) {
	storeID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Log.Error("incorrect storeID input", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"message": "incorrect storeID input"})
		return
	}
	versionID, err := strconv.Atoi(ctx.Param("version_id"))
	if err != nil {
		h.logger.Log.Error("incorrect versionID input", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"message": "incorrect versionID input"})
		return
	}

	version, store, err := h.service.GetVersion(ctx.Request.Context(), storeID, versionID)
	if err != nil {
		h.logger.Log.Error("Internal server error", zap.Int("id", storeID), zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.CreateVersionResp(version, store))
}
