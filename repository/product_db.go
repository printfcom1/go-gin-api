package repository

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productrRepositoryDB struct {
	db *mongo.Database
}

func NewProductRepositoryDB(db *mongo.Database) productrRepositoryDB {
	return productrRepositoryDB{db: db}
}

var colNameProduct string = "product"

func (p productrRepositoryDB) GetProducts(page int, pageSize int, filter gin.H) ([]Product, error) {
	skip := (page - 1) * pageSize

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(pageSize))

	cursor, err := p.db.Collection(colNameProduct).Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}

	products := []Product{}
	for cursor.Next(context.Background()) {
		product := Product{}
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (p productrRepositoryDB) GetProduct(id primitive.ObjectID) (*Product, error) {

	filter := gin.H{"_id": id}
	product := &Product{}
	err := p.db.Collection(colNameProduct).FindOne(context.Background(), filter).Decode(product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p productrRepositoryDB) CreateProuct(product Product) (interface{}, error) {
	result, err := p.db.Collection(colNameProduct).InsertOne(context.Background(), product)
	if err != nil {
		return nil, err
	}
	id := result.InsertedID
	return id, nil
}

func (p productrRepositoryDB) UpdateProduct(id primitive.ObjectID, product ProductUpdate) (*Product, error) {
	filter := gin.H{"_id": id}
	update := gin.H{"$set": product}
	productRes := &Product{}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := p.db.Collection(colNameProduct).FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&productRes)
	if err != nil {
		return nil, err
	}
	return productRes, nil
}

func (p productrRepositoryDB) DeleteProduct(id primitive.ObjectID) (*Product, error) {
	filter := gin.H{"_id": id}
	productRes := &Product{}
	err := p.db.Collection(colNameProduct).FindOneAndDelete(context.Background(), filter).Decode(&productRes)
	if err != nil {
		return nil, err
	}
	return productRes, nil
}

func (p productrRepositoryDB) UpdateStockProduct(id primitive.ObjectID, stock Stock) (*Product, error) {
	filter := gin.H{"_id": id}
	update := gin.H{"$set": gin.H{"stock": stock}}
	product := &Product{}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := p.db.Collection(colNameProduct).FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&product)
	if err != nil {
		return nil, err
	}
	return product, nil
}
