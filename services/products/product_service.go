package products

import "github.com/example/golang-test/models"

type ProductService interface {
	CreateProduct(product *models.Product) error
	DeleteProduct(productId string) error
	UpdateProduct(productId string, product *models.Product) error
	GetProductById(productId string) (*models.Product, error)
	GetListProducts(page int64) ([]*models.Product, int64, error)
}
