package handler

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	strc "github.com/gin-api/struc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func GoDotEnvVariable(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println(map[string]string{"error": err.Error()})
	}

	return os.Getenv(key)
}

func HandlerQuery(c *gin.Context, filter *gin.H) {
	name := c.Query("name")
	category := c.Query("category")
	productCode := c.Query("productCode")

	queryKey := []strc.QueryNameCategoryCode{
		{Field: "name", Key: name},
		{Field: "category", Key: category},
		{Field: "productCode", Key: productCode},
	}

	QueryNameCategoryCode(filter, queryKey)

	priceMin := c.Query("priceMin")
	priceMax := c.Query("priceMax")

	availableMin := c.Query("availableMin")
	availableMax := c.Query("availableMax")

	reservedMin := c.Query("reservedMin")
	reservedMax := c.Query("reservedMax")

	soldMin := c.Query("soldMin")
	soldMax := c.Query("soldMax")

	queryStock := []strc.QueryStock{
		{Field: "price", Min: priceMin, Max: priceMax},
		{Field: "stock.available", Min: availableMin, Max: availableMax},
		{Field: "stock.reserved", Min: reservedMin, Max: reservedMax},
		{Field: "stock.sold", Min: soldMin, Max: soldMax},
	}

	QueryStock(filter, queryStock)
}

func QueryStock(filter *gin.H, query []strc.QueryStock) {
	for _, item := range query {
		if item.Min != "" {
			parsedVariableMin, _ := strconv.Atoi(item.Min)
			(*filter)[item.Field] = gin.H{"$gt": parsedVariableMin}
		}

		if item.Max != "" {
			parsedVariableMax, _ := strconv.Atoi(item.Max)
			if _, exists := (*filter)[item.Field]; exists {
				(*filter)[item.Field].(gin.H)["$lt"] = parsedVariableMax
			} else {
				(*filter)[item.Field] = gin.H{"$lt": parsedVariableMax}
			}
		}
	}
}

func QueryNameCategoryCode(filter *gin.H, query []strc.QueryNameCategoryCode) {
	for _, item := range query {
		if item.Key != "" {
			regexPattern := "(?i)^" + regexp.QuoteMeta(item.Key)
			regex := gin.H{"$regex": regexPattern}
			(*filter)[item.Field] = regex
		}
	}

}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := GoDotEnvVariable("SECRET_KEY")
		tokenString := c.GetHeader("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})

		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized", "message": err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("username", claims["username"])
		} else {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
