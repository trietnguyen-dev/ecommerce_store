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
func (d *DAO) GetListUsers(page int64) ([]*models.DBResponse, error) {
	ctx, cancel := utils.NewCtx()
	defer cancel()
	var listUser []*models.DBResponse
	// Số lượng kết quả trên mỗi trang
	resultsPerPage := 10

	// Tính toán skip dựa trên trang
	skip := (page - 1) * int64(resultsPerPage)

	opts := options.Find().SetLimit(int64(resultsPerPage)).SetSkip(skip)

	cursor, err := d.userColl.Find(ctx, bson.M{}, opts)

	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &listUser); err != nil {
		panic(err)
	}
	return listUser, nil
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
