package service

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/gin-api/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) productService {
	return productService{productRepo: productRepo}
}

func (s productService) GetProductsService(page int, pageSize int, queryStock []QueryStock, queryNameCategoryCode []QueryNameCategoryCode) ([]ProductRespose, error) {
	filter := gin.H{}
	generateQueryNameCategoryCodeFunc(&filter, queryNameCategoryCode)
	generateQueryStock(&filter, queryStock)

	products, err := s.productRepo.GetProducts(page, pageSize, filter)
	if err != nil {
		return nil, err
	}

	productsRespose := []ProductRespose{}
	for _, product := range products {
		ProductRes := ProductRespose{
			ID:          product.ID,
			ProductCode: product.ProductCode,
			Name:        product.Name,
			Price:       product.Price,
			Available:   product.Stock.Available,
			Category:    product.Category,
			Description: product.Description,
		}

		productsRespose = append(productsRespose, ProductRes)
	}

	return productsRespose, nil
}

func (s productService) CreateProductService(product repository.Product) (*string, error) {
	id, err := s.productRepo.CreateProuct(product)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors.New("data with the same unique value for name or productcode already exists")
		}
		return nil, err
	}
	objIDString := id.(primitive.ObjectID).Hex()

	message := "Product id " + objIDString + " created successfully"
	return &message, nil
}

func (s productService) GetProductService(id string) (*ProductRespose, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	product, err := s.productRepo.GetProduct(objectID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			message := "product item with id " + id + " not found"
			return nil, errors.New(message)
		}
		return nil, err
	}

	productRes := &ProductRespose{
		ID:          product.ID,
		ProductCode: product.ProductCode,
		Name:        product.Name,
		Price:       product.Price,
		Available:   product.Stock.Available,
		Category:    product.Category,
		Description: product.Description,
	}

	return productRes, nil
}

func (s productService) UpdateProductService(id string, product repository.ProductUpdate) (*ProductRespose, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	result, err := s.productRepo.UpdateProduct(objectID, product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			message := "product item with id " + id + " not found"
			return nil, errors.New(message)
		}

		if mongo.IsDuplicateKeyError(err) {
			return nil, errors.New("data with the same unique value for name or productcode already exists")
		}
		return nil, err
	}

	productRes := &ProductRespose{
		ID:          result.ID,
		ProductCode: result.ProductCode,
		Name:        result.Name,
		Price:       result.Price,
		Available:   result.Stock.Available,
		Category:    result.Category,
		Description: result.Description,
	}

	return productRes, nil
}

func (s productService) DeleteProductService(id string) (*ProductRespose, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	product, err := s.productRepo.DeleteProduct(objectID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			message := "product item with id " + id + " not found"
			return nil, errors.New(message)
		}
		return nil, err
	}

	productRes := &ProductRespose{
		ID:          product.ID,
		ProductCode: product.ProductCode,
		Name:        product.Name,
		Price:       product.Price,
		Available:   product.Stock.Available,
		Category:    product.Category,
		Description: product.Description,
	}

	return productRes, nil
}

func (s productService) UpdateStockProductService(id string, stock repository.Stock) (*repository.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	product, err := s.productRepo.UpdateStockProduct(objectID, stock)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			message := "product item with id " + id + " not found"
			return nil, errors.New(message)
		}
		return nil, err
	}
	return product, nil
}

func generateQueryStock(filter *gin.H, query []QueryStock) {
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

func generateQueryNameCategoryCodeFunc(filter *gin.H, query []QueryNameCategoryCode) {
	for _, item := range query {
		if item.Key != "" {
			regexPattern := "(?i)^" + regexp.QuoteMeta(item.Key)
			regex := gin.H{"$regex": regexPattern}
			(*filter)[item.Field] = regex
		}
	}

}
