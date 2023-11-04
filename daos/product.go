package daos

import (
	"github.com/example/golang-test/models"
	"github.com/example/golang-test/utils"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *DAO) CreateProduct(product *models.Product) error {
	ctx, cancel := utils.NewCtx()
	defer cancel()

	filter := bson.M{"$or": []bson.M{
		{"id": product.Id},
	}}

	var existingProduct models.Product
	findOptions := options.FindOne()

	err := d.productColl.FindOne(
		ctx,
		filter,
		findOptions,
	).Decode(&existingProduct)

	if err == nil {
		return errors.New("product already exists")
	}

	_, err = d.productColl.InsertOne(ctx, &product)

	if err != nil {
		return err
	}

	return nil
}

func (d *DAO) DeleteProduct(productId string) error {
	ctx, cancel := utils.NewCtx()
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return err
	}
	_, err = d.productColl.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}

func (d *DAO) UpdateProduct(productId string, product *models.Product) error {
	ctx, cancel := utils.NewCtx()
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return err
	}

	condition := bson.M{
		"_id": bson.M{"$ne": oid},
		"$or": []bson.M{
			{"id": product.Id},
		},
	}
	query := bson.M{"_id": oid}
	update := bson.M{"$set": product}

	existingProduct := d.productColl.FindOne(ctx, condition)
	if existingProduct.Err() == nil {
		return errors.New("id product already exists")
	}

	_, err = d.productColl.UpdateOne(ctx, query, update)

	if err != nil {
		return errors.New("Update failed: " + err.Error())
	}

	return nil

}
func (d *DAO) GetProductById(productId string) (*models.Product, error) {
	ctx, cancel := utils.NewCtx()
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return nil, err
	}
	var product *models.Product
	err = d.productColl.FindOne(ctx, bson.M{"_id": oid}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return product, nil

}
func (d *DAO) GetListProducts(page int64) ([]*models.Product, int64, error) {
	ctx, cancel := utils.NewCtx()
	defer cancel()

	var products []*models.Product
	resultsPerPage := 5

	skip := (page - 1) * int64(resultsPerPage)

	//optsCount := options.Count().SetSkip(skip).SetLimit(int64(resultsPerPage))

	count, err := d.productColl.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	optsPage := options.Find().SetLimit(int64(resultsPerPage)).SetSkip(skip)

	cursor, err := d.productColl.Find(ctx, bson.M{}, optsPage)
	if err != nil {
		return nil, 0, err
	}
	if err = cursor.All(ctx, &products); err != nil {
		panic(err)
	}

	return products, count, nil
}
