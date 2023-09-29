package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-api/repository"
	"github.com/gin-api/service"
	"github.com/gin-gonic/gin"
)

type productHandler struct {
	productHand service.ProductService
}

func NewProductHandler(productHand service.ProductService) productHandler {
	return productHandler{productHand: productHand}
}

func (h productHandler) GetProducts(c *gin.Context) {
	name := c.Query("name")
	category := c.Query("category")
	productCode := c.Query("productCode")

	queryKey := []service.QueryNameCategoryCode{
		{Field: "name", Key: name},
		{Field: "category", Key: category},
		{Field: "productCode", Key: productCode},
	}

	priceMin := c.Query("priceMin")
	priceMax := c.Query("priceMax")

	availableMin := c.Query("availableMin")
	availableMax := c.Query("availableMax")

	reservedMin := c.Query("reservedMin")
	reservedMax := c.Query("reservedMax")

	soldMin := c.Query("soldMin")
	soldMax := c.Query("soldMax")

	queryStock := []service.QueryStock{
		{Field: "price", Min: priceMin, Max: priceMax},
		{Field: "stock.available", Min: availableMin, Max: availableMax},
		{Field: "stock.reserved", Min: reservedMin, Max: reservedMax},
		{Field: "stock.sold", Min: soldMin, Max: soldMax},
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	products, err := h.productHand.GetProductsService(page, pageSize, queryStock, queryKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h productHandler) GetProductById(c *gin.Context) {
	id := c.Param("id")
	product, err := h.productHand.GetProductService(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h productHandler) CreatedProduct(c *gin.Context) {
	product := &repository.Product{}

	if err := c.Bind(product); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	message, err := h.productHand.CreateProductService(*product)

	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, map[string]interface{}{"message": *message})

}

func (h productHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	product := &repository.ProductUpdate{}

	if err := c.Bind(product); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	product.UpdatedAt = time.Now()
	productRes, err := h.productHand.UpdateProductService(id, *product)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, *productRes)
}

func (h productHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := h.productHand.DeleteProductService(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h productHandler) UpdateStockProduct(c *gin.Context) {
	id := c.Param("id")
	stock := &repository.Stock{}

	if err := c.Bind(stock); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	product, err := h.productHand.UpdateStockProductService(id, *stock)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, *product)
}
