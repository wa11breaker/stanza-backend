package controllers

import (
	"context"
	"stanza-api/src/configs"
	"stanza-api/src/models"
	"stanza-api/src/responses"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection("users")
var validate = validator.New()

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User

		// Call BindJSON to bind the received JSON to
		// user
		if err := c.BindJSON(&user); err != nil {
			res := responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}
			c.JSON(http.StatusBadRequest, res)
			return
		}

		// Validate required fields
		if err := validate.Struct(&user); err != nil {
			res := responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}
			c.JSON(http.StatusBadRequest, res)
			return
		}

		newUser := models.User{
			Id:        primitive.NewObjectID(),
			Name:      user.Name,
			AvatarUrl: user.AvatarUrl,
		}

		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			res := responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}
			c.JSON(http.StatusInternalServerError, res)
			return
		}
		res := responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: result}
		c.JSON(http.StatusCreated, res)
	}
}

func GetAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		var user models.User

		id, _ := primitive.ObjectIDFromHex(userId)
		filter := bson.M{"id": bson.M{"$eq": id}}

		err := userCollection.FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			res := responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		res := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: user}
		c.JSON(http.StatusOK, res)
	}
}

func EditAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()

		var user models.User

		// Call BindJSON to bind the received JSON to
		// user
		if err := c.BindJSON(&user); err != nil {
			res := responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}
			c.JSON(http.StatusBadRequest, res)
			return
		}

		// Validate required fields
		if err := validate.Struct(&user); err != nil {
			res := responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}
			c.JSON(http.StatusBadRequest, res)
			return
		}

		id, _ := primitive.ObjectIDFromHex(userId)
		filter := bson.M{"id": bson.M{"$eq": id}}
		update := bson.M{
			"$set": bson.M{
				"name":      user.Name,
				"avatarUrl": user.AvatarUrl,
			},
		}

		result, err := userCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			res := responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		//get updated user details
		var updatedUser models.User
		if result.MatchedCount == 1 {
			err := userCollection.FindOne(ctx, filter).Decode(&updatedUser)
			if err != nil {
				res := responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}
				c.JSON(http.StatusInternalServerError, res)
				return
			}
		}

		res := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: updatedUser}
		c.JSON(http.StatusOK, res)
	}
}

func DeleteAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()

		id, _ := primitive.ObjectIDFromHex(userId)
		filter := bson.M{"id": bson.M{"$eq": id}}

		result, err := userCollection.DeleteOne(ctx, filter)

		if err != nil {
			res := responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		if result.DeletedCount < 1 {
			res := responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		res := responses.UserResponse{Status: http.StatusOK, Message: "User successfully deleted", Data: nil}
		c.JSON(http.StatusOK, res)
	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var users []models.User

		results, err := userCollection.Find(ctx, bson.M{})

		if err != nil {
			res := responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		// Reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.User
			if err = results.Decode(&singleUser); err != nil {
				res := responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: nil}
				c.JSON(http.StatusInternalServerError, res)
			}

			users = append(users, singleUser)
		}

		res := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: users}
		c.JSON(http.StatusOK, res)
	}
}
