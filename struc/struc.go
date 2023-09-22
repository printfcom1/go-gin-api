package struc

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthInput struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Product struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ProductCode string             `json:"productCode" bson:"productCode"`
	Name        string             `json:"name"`
	Price       float64            `json:"price"`
	Stock       Stock              `json:"stock"`
	Category    string             `json:"category"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type ProductUpdate struct {
	ProductCode string    `json:"productCode" bson:"productCode"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

type Stock struct {
	Available int `json:"available"`
	Reserved  int `json:"reserved"`
	Sold      int `json:"sold"`
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
