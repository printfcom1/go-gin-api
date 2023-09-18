package main

import (
	"time"

	ctrl "github.com/gin-api/src"
	handler "github.com/gin-api/util"
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
	login.POST("/login", ctrl.Login)

	authorized := r.Group("/api/product")

	authorized.Use(handler.AuthMiddleware())
	authorized.POST("addProduct", ctrl.AddProduct)
	authorized.GET("getProductAll", ctrl.GetProductAll)
	authorized.GET("getProductById/:id", ctrl.GetProductById)
	authorized.PUT("updateProduct/:id", ctrl.UpdateProduct)
	authorized.PUT("updateStock/:id", ctrl.UpdateStock)
	authorized.DELETE("deleteProduct/:id", ctrl.DeleteStock)
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
