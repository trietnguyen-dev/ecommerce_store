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

func (d *DAO) IsExistUser(user *models.SignUpInput) error {
	ctx, cancel := utils.NewCtx()
	defer cancel()

	filter := bson.M{"$or": []bson.M{
		{"email": user.Email},
		{"phone_number": user.PhoneNumber},
	}}
	//findOptions := options.FindOne().SetHint("email_1_phone_number_1")

	var existingUser models.SignUpInput
	err := d.userColl.FindOne(ctx, filter).Decode(&existingUser)

	if err == nil {
		return errors.New("user with that email or phone number already exists")
	}
	return nil
}
func (d *DAO) SignUpUser(user *models.SignUpInput) error {
	ctx, cancel := utils.NewCtx()
	defer cancel()

	filter := bson.M{"$or": []bson.M{
		{"email": user.Email},
		{"phone_number": user.PhoneNumber},
	}}

	var existingUser models.SignUpInput
	findOptions := options.FindOne().SetHint("email_1_phone_number_1")

	err := d.userColl.FindOne(
		ctx,
		filter,
		findOptions,
	).Decode(&existingUser)

	if err == nil {
		return errors.New("user with that email or phone number already exists")
	}

	_, err = d.userColl.InsertOne(ctx, user)

	if err != nil {
		return err
	}

	return nil
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
			return nil, err
		}
	}

	return user, nil
}

func (d *DAO) UpdateUserById(id string, value *models.UserResponse) error {
	ctx, cancel := utils.NewCtx()
	defer cancel()

	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	emailAndPhoneCheckCondition := bson.M{
		"_id": bson.M{"$ne": userId},
		"$or": []bson.M{
			{"email": value.Email},
			{"phone_number": value.PhoneNumber},
		},
	}
	query := bson.M{"_id": userId}
	update := bson.M{"$set": value}

	existingUser := d.userColl.FindOne(ctx, emailAndPhoneCheckCondition)
	if existingUser.Err() == nil {
		return errors.New("email or phone number already exists")
	}

	_, err = d.userColl.UpdateOne(ctx, query, update)

	if err != nil {
		return errors.New("Update failed: " + err.Error())
	}

	return nil
}
func (d *DAO) ChangePassword(id string, newPasswordHash string) error {
	ctx, cancel := utils.NewCtx()
	defer cancel()

	userId, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": userId}
	update := bson.M{"$set": bson.M{"password": newPasswordHash}}
	result, err := d.userColl.UpdateOne(ctx, query, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("change password failed")
	}

	return nil
}
