package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-api/handler"
	"github.com/gin-api/repository"
	"github.com/gin-api/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupRouter() *gin.Engine {

	db := initMongoDB()
	userRepository := repository.NewUserRepositoryDB(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	productRepository := repository.NewProductRepositoryDB(db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	login := r.Group("/api/user")

	login.POST("/login", userHandler.Login)

	auth := gin.BasicAuth(getDataAdim())
	admin := r.Group("/api/admin", auth)
	admin.POST("/register", userHandler.RegisterUser)

	authorized := r.Group("/api/product")

	authorized.Use(authMiddleware())
	authorized.POST("addProduct", productHandler.CreatedProduct)
	authorized.GET("getProductAll", productHandler.GetProducts)
	authorized.GET("getProductById/:id", productHandler.GetProductById)
	authorized.PUT("updateProduct/:id", productHandler.UpdateProduct)
	authorized.PUT("updateStock/:id", productHandler.UpdateStockProduct)
	authorized.DELETE("deleteProduct/:id", productHandler.DeleteProduct)
	return r
}

func main() {
	initTimeZone()
	r := setupRouter()
	r.Run(":3000")
}

func initTimeZone() {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	time.Local = location
}

func initMongoDB() *mongo.Database {
	url, err := goDotEnvVariable("MONGODB_URL")
	if err != nil {
		panic(err)
	}
	dbName, err := goDotEnvVariable("DB_NAME")
	if err != nil {
		panic(err)
	}
	clientOptions := options.Client().ApplyURI(*url)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	mongoDB := client.Database(*dbName)
	fmt.Println("Connected to MongoDB!")
	return mongoDB
}

func goDotEnvVariable(key string) (*string, error) {

	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	value := os.Getenv(key)

	return &value, nil
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key, err := goDotEnvVariable("SECRET_KEY")
		if err != nil {
			panic(err)
		}
		tokenString := c.GetHeader("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(*key), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("username", claims["username"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func getDataAdim() gin.Accounts {
	usernameAdmin, err := goDotEnvVariable("USERNAME_REGIS")
	if err != nil {
		panic(err)
	}
	passwordAdim, err := goDotEnvVariable("PASSWORD_REGIS")
	if err != nil {
		panic(err)
	}

	return gin.Accounts{
		*usernameAdmin: *passwordAdim,
	}
}
