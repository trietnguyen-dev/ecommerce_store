package daos

import (
	"encoding/json"
	"errors"
	"fmt"
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

	opts := options.Find().SetLimit(2).SetSkip(page)
	cursor, err := d.userColl.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &listUser); err != nil {
		panic(err)
	}
	for _, result := range listUser {
		res, _ := json.Marshal(result)
		fmt.Println(string(res))
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
