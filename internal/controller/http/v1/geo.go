package v1

import (
	geoEntity "domashka-backend/internal/entity/geo"
	"domashka-backend/internal/utils/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type GeoHandler struct {
	geoUsecase geoUsecase
}

// AddClientAddress godoc
// @Summary      Удаление товара из корзины
// @Description  Удаление товара из корзины
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        request   body      RemoveCartItemRequest  true  "body"
// @Router       /cart/remove_item [post]
func (h *GeoHandler) AddClientAddress(ctx *gin.Context) {
	var address geoEntity.Address
	clientID := ctx.Param("client_id")
	clientIDInt, err := strconv.Atoi(clientID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid client_id"})
		return
	}

	if err := ctx.BindJSON(&address); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	if !validation.IsAddressInRussia(address) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Address is not in Russia"})
		return
	}

	if err := h.geoUsecase.AddClientAddress(ctx, clientIDInt, address); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *GeoHandler) UpdateClientAddress(ctx *gin.Context) {
	var address geoEntity.Address
	clientID := ctx.Param("client_id")
	addressID := ctx.Param("address_id")
	clientIDInt, err := strconv.Atoi(clientID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid client_id"})
		return
	}

	addressIDInt, err := strconv.Atoi(addressID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid address_id"})
		return
	}

	if err := ctx.BindJSON(&address); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	if !validation.IsAddressInRussia(address) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Address is not in Russia"})
		return
	}

	if err := h.geoUsecase.UpdateClientAddress(ctx, clientIDInt, addressIDInt, address); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *GeoHandler) GetClientAddresses(ctx *gin.Context) {
	clientID := ctx.Param("client_id")
	clientIDInt, err := strconv.Atoi(clientID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid client_id"})
		return
	}

	addresses, err := h.geoUsecase.GetClientAddresses(ctx, clientIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": addresses})
}

func (h *GeoHandler) GetChefAddress(ctx *gin.Context) {
	chefID := ctx.Param("chef_id")
	chefIDInt, err := strconv.Atoi(chefID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid chef_id"})
		return
	}

	address, err := h.geoUsecase.GetChefAddress(ctx, chefIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": address})
}

func (h *GeoHandler) UpdateOrPostChefAddress(ctx *gin.Context) {
	var address geoEntity.Address
	chefID := ctx.Param("chef_id")
	chefIDInt, err := strconv.Atoi(chefID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid chef_id"})
		return
	}

	if err := ctx.BindJSON(&address); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	if !validation.IsAddressInRussia(address) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Address is not in Russia"})
		return
	}

	if _, err := h.geoUsecase.GetChefAddress(ctx, chefIDInt); err != nil {
		if err := h.geoUsecase.AddChefAddress(ctx, chefIDInt, address); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
	} else {
		if err := h.geoUsecase.UpdateChefAddress(ctx, chefIDInt, address); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *GeoHandler) GetChefsAddrByRange(ctx *gin.Context) {
	clientAddressID := ctx.Param("client_address_id")
	radius := ctx.Param("radius")
	clientAddressIDInt, err := strconv.Atoi(clientAddressID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid client_address_id"})
		return
	}

	radiusFloat, err := strconv.ParseFloat(radius, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid radius"})
		return
	}

	addresses, err := h.geoUsecase.FindChefsNearAddress(ctx, clientAddressIDInt, radiusFloat)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": addresses})
}

func (h *GeoHandler) GetClientsAddrByRange(ctx *gin.Context) {
	chefAddressID := ctx.Param("chef_id")
	radius := ctx.Query("radius")
	chefIDInt, err := strconv.Atoi(chefAddressID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid chef_address_id"})
		return
	}

	radiusFloat, err := strconv.ParseFloat(radius, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid radius"})
		return
	}

	addresses, err := h.geoUsecase.FindClientsNearAddress(ctx, chefIDInt, radiusFloat)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": addresses})
}

func newGeoHandler(router *gin.RouterGroup, g geoUsecase) {

	h := GeoHandler{geoUsecase: g}

	geoGroup := router.Group("/geo")
	{
		geoGroup.POST("/clients/:client_id/addresses", h.AddClientAddress)
		geoGroup.PUT("/clients/:client_id/addresses/:address_id", h.UpdateClientAddress)
		geoGroup.GET("/clients/:client_id/addresses", h.GetClientAddresses)
		geoGroup.GET("/chefs/:chef_id/address", h.GetChefAddress)
		geoGroup.POST("/chefs/:chef_id/address", h.UpdateOrPostChefAddress)
		geoGroup.GET("/search/chefs-near-client-address/:client_address_id", h.GetChefsAddrByRange)
		geoGroup.GET("/search/clients-near-chef/:chef_id", h.GetClientsAddrByRange)
	}
}
