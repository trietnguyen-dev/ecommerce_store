package products

import (
	"github.com/example/golang-test/config"
	"github.com/example/golang-test/daos"
	"github.com/example/golang-test/models"
	"time"
)

type ProductServiceImpl struct {
	dao  *daos.DAO
	conf *config.Config
}

func NewProductService(dao *daos.DAO, conf *config.Config) *ProductServiceImpl {
	return &ProductServiceImpl{dao, conf}
}
func (ps *ProductServiceImpl) CreateProduct(product *models.Product) error {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	err := ps.dao.CreateProduct(product)
	if err != nil {
		return err
	}
	return nil
}
func (ps *ProductServiceImpl) DeleteProduct(productId string) error {

	err := ps.dao.DeleteProduct(productId)
	if err != nil {
		return err
	}
	return nil
}
func (ps *ProductServiceImpl) UpdateProduct(productId string, product *models.Product) error {
	product.UpdatedAt = time.Now()
	err := ps.dao.UpdateProduct(productId, product)
	if err != nil {
		return err
	}
	return nil
}
func (ps *ProductServiceImpl) GetProductById(productId string) (*models.Product, error) {

	value, err := ps.dao.GetProductById(productId)
	if err != nil {
		return nil, err
	}
	return value, nil
}
func (ps *ProductServiceImpl) GetListProducts(page int64) ([]*models.Product, int64, error) {
	products, count, err := ps.dao.GetListProducts(page)
	if err != nil {
		return nil, 0, err
	}
	return products, count, nil
}
