package ctrl

import (
	"context"
	"net/http"
	"strconv"
	"time"

	mongc "github.com/gin-api/db"
	strc "github.com/gin-api/struc"
	handler "github.com/gin-api/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Collection = mongc.InitMongoClient().Database("golang").Collection("product")

func Login(c *gin.Context) {

	auth := new(strc.AuthInput)

	if err := c.Bind(auth); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	expirationTime := time.Now().Add(1 * time.Hour)

	claims["username"] = auth.UserName
	claims["exp"] = expirationTime.Unix()

	key := handler.GoDotEnvVariable("SECRET_KEY")
	var secretKey = []byte(key)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
}

func AddProduct(c *gin.Context) {
	product := new(strc.Product)

	if err := c.Bind(product); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	} else {
		product.CreatedAt = time.Now()
		product.UpdatedAt = time.Now()
		res, err := db.InsertOne(context.Background(), product)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else {
			c.JSON(http.StatusCreated, map[string]interface{}{
				"message": "Product created successfully",
				"id":      res.InsertedID,
			})
		}
	}
}

func GetProductAll(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	skip := (page - 1) * pageSize

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(pageSize))

	filter := gin.H{}

	go handler.HandlerQuery(c, &filter)

	cursor, err := db.Find(context.Background(), filter, findOptions)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	} else {
		var productList []strc.Product
		for cursor.Next(context.Background()) {
			var product strc.Product
			if err := cursor.Decode(&product); err != nil {
				c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
				break
			}
			productList = append(productList, product)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, productList)
		}
	}
	defer cursor.Close(context.Background())

}

func GetProductById(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	} else {
		filter := gin.H{"_id": objectID}
		var result strc.Product
		err = db.FindOne(context.Background(), filter).Decode(&result)
		HandlerResponse(c, err, result, id)
	}

}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	} else {
		filter := gin.H{"_id": objectID}
		productUpdate := new(strc.ProductUpdate)

		if err := c.Bind(productUpdate); err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else {
			productUpdate.UpdatedAt = time.Now()
			update := gin.H{"$set": productUpdate}
			var result strc.Product
			opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
			err := db.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&result)
			HandlerResponse(c, err, result, id)
		}
	}

}

func UpdateStock(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	} else {
		filter := gin.H{"_id": objectID}
		stockProduct := new(strc.Stock)

		if err := c.Bind(stockProduct); err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else {
			update := gin.H{"$set": gin.H{"stock": stockProduct}}
			var result strc.Product
			opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
			err := db.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&result)
			HandlerResponse(c, err, result, id)
		}
	}
}

func DeleteStock(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	} else {
		filter := gin.H{"_id": objectID}
		var result strc.Product
		err = db.FindOneAndDelete(context.Background(), filter).Decode(&result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, result)
		}
	}

}

func HandlerResponse(c *gin.Context, err error, result strc.Product, id string) {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			response := map[string]interface{}{"message": "ToDo item with ID " + id + " not found."}
			c.JSON(http.StatusBadRequest, response)
		} else {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusOK, result)
	}
}
