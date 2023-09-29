package service

import (
	"github.com/gin-api/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductRespose struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ProductCode string             `json:"productCode" bson:"productCode"`
	Name        string             `json:"name"`
	Price       float64            `json:"price"`
	Available   int                `json:"available"`
	Category    string             `json:"category"`
	Description string             `json:"description"`
}

type QueryStock struct {
	Field string
	Min   string
	Max   string
}

type QueryNameCategoryCode struct {
	Field string
	Key   string
}

type ProductService interface {
	GetProductsService(int, int, []QueryStock, []QueryNameCategoryCode) ([]ProductRespose, error)
	GetProductService(string) (*ProductRespose, error)
	CreateProductService(repository.Product) (*string, error)
	UpdateProductService(string, repository.ProductUpdate) (*ProductRespose, error)
	DeleteProductService(string) (*ProductRespose, error)
	UpdateStockProductService(string, repository.Stock) (*repository.Product, error)
}
