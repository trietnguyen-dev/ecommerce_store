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

func (d *DAO) SignUpUser(user *models.SignUpInput) (*models.DBResponse, error) {
	ctx, cancel := utils.NewCtx()
	defer cancel()

	res, err := d.userColl.InsertOne(ctx, &user)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("user with that email already exist")
		}
		return nil, err
	}

	// Create a unique index for the email field
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}

	if _, err := d.userColl.Indexes().CreateOne(ctx, index); err != nil {
		return nil, errors.New("could not create index for email")
	}

	var newUser *models.DBResponse
	query := bson.M{"_id": res.InsertedID}

	err = d.userColl.FindOne(ctx, query).Decode(&newUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (d *DAO) FindUserByEmail(email string) (*models.DBResponse, error) {
	ctx, cancel := utils.NewCtx()
	defer cancel()

	var user models.DBResponse
	filter := bson.M{"email": strings.ToLower(email)}
	result := d.userColl.FindOne(ctx, filter).Decode(&user)

	if errors.Is(result, mongo.ErrNoDocuments) {
		return &models.DBResponse{}, errors.New("email không tồn tại")
	}

	return &user, nil

}

func (d *DAO) FindUserById(id string) (*models.DBResponse, error) {
	ctx, cancel := utils.NewCtx()
	defer cancel()
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.DBResponse
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

func (d *DAO) UpdateUserById(id string, value *models.UserResponse) error {
	ctx, cancel := utils.NewCtx()
	defer cancel()

	userId, _ := primitive.ObjectIDFromHex(id)

	query := bson.M{"_id": userId}
	update := bson.M{"$set": value}
	result, err := d.userColl.UpdateOne(ctx, query, update)
	if err != nil {
		return errors.New("email update existing")
	}

	var updatedUser models.DBResponse
	if result.MatchedCount == 1 {
		err := d.userColl.FindOne(ctx, bson.M{"_id": userId}).Decode(&updatedUser)
		if err != nil {
			return errors.New("couldn't find user" + err.Error())
		}
	}

	return nil
}

//func (d *DAO) UpdateOne(field string, value interface{}) (*models.DBResponse, error) {
//	query := bson.D{{Key: field, Value: value}}
//	update := bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: value}}}}
//	result, err := uc.collection.UpdateOne(uc.ctx, query, update)
//
//	fmt.Print(result.ModifiedCount)
//	if err != nil {
//		fmt.Print(err)
//		return &models.DBResponse{}, err
//	}
//
//	return &models.DBResponse{}, nil
//}
