package main

import (
	"time"

	"github.com/gin-api/controler"
	"github.com/gin-api/handler"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {

	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	time.Local = location

	login := r.Group("/api")
	login.POST("/login", controler.Login)

	authorized := r.Group("/api/product")

	authorized.Use(handler.AuthMiddleware())
	authorized.POST("addProduct", controler.AddProduct)
	authorized.GET("getProductAll", controler.GetProductAll)
	authorized.GET("getProductById/:id", controler.GetProductById)
	authorized.PUT("updateProduct/:id", controler.UpdateProduct)
	authorized.PUT("updateStock/:id", controler.UpdateStock)
	authorized.DELETE("deleteProduct/:id", controler.DeleteStock)
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
