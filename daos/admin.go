package daos

import (
	"errors"
	"github.com/example/golang-test/models"
	"github.com/example/golang-test/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

func (d *DAO) FindAdminByEmail(email string) (*models.DBResponse, error) {
	ctx, cancel := utils.NewCtx()
	defer cancel()

	var user models.DBResponse
	filter := bson.M{"email": strings.ToLower(email)}
	result := d.adminColl.FindOne(ctx, filter).Decode(&user)

	if errors.Is(result, mongo.ErrNoDocuments) {
		return &models.DBResponse{}, errors.New("email không tồn tại")
	}

	return &user, nil
}
func (d *DAO) GetListUsers(page int64) ([]*models.DBResponse, int64, error) {
	ctx, cancel := utils.NewCtx()
	defer cancel()
	var listUser []*models.DBResponse
	// Số lượng kết quả trên mỗi trang
	resultsPerPage := 10

	// Tính toán skip dựa trên trang
	skip := (page - 1) * int64(resultsPerPage)

	//optsCount := options.Count().SetSkip(skip).SetLimit(int64(resultsPerPage))

	count, err := d.userColl.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	optsPage := options.Find().SetLimit(int64(resultsPerPage)).SetSkip(skip)

	cursor, err := d.userColl.Find(ctx, bson.M{}, optsPage)
	if err != nil {
		return nil, 0, err
	}
	if err = cursor.All(ctx, &listUser); err != nil {
		panic(err)
	}

	return listUser, count, nil
}
func (d *DAO) FindAdminById(id string) (*models.DBResponse, error) {
	ctx, cancel := utils.NewCtx()
	defer cancel()
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.DBResponse
	query := bson.M{"_id": oid}
	err := d.adminColl.FindOne(ctx, query).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// This error means your query did not match any documents.
			return nil, err
		}
	}
	return user, nil
}
func (d *DAO) GetUserById(id string) (*models.UserResponse, error) {
	ctx, cancel := utils.NewCtx()
	defer cancel()
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.UserResponse
	query := bson.M{"_id": oid}

	err := d.userColl.FindOne(ctx, query).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// This error means your query did not match any documents.
			return nil, err
		}
	}
	return user, nil
}
func (d *DAO) DeleteUserById(id string) error {
	ctx, cancel := utils.NewCtx()
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = d.userColl.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}
