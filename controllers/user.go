package controllers

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/barindebnath/gofiber/config"
	"github.com/barindebnath/gofiber/models"
	"github.com/barindebnath/gofiber/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var userCollection *mongo.Collection = config.GetCollection("users")
var validate = validator.New()

func SignUp(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": validationErr.Error()},
		})
	}

	_, err := getUserFromEmail(user.Email, ctx)
	if err == nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": "", "msg": "Email already exists."},
		})
	}

	newUser := models.User{
		ID:        primitive.NewObjectID(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		Email:     user.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	setSessionCookie(c, user.ID)

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &fiber.Map{"data": result, "msg": "User created successfully"},
	})
}

func SignIn(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type SigninData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	data := SigninData{}

	//validate the request body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	user, err := getUserFromEmail(data.Email, ctx)
	if err != nil {
		log.Println("\n\nEmail does not exist.")
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": "", "msg": "Email and password do not match."},
		})
	}

	if user.Password != data.Password {
		log.Println("\n\nPassword did not match.")
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": "", "msg": "Email and password do not match."},
		})
	}

	setSessionCookie(c, user.ID)

	return c.Status(http.StatusOK).JSON(responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": user, "msg": "Login successful"},
	})
}

func LogOut(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:    "userSession",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	return c.Status(http.StatusOK).JSON(responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": "", "msg": "User logged out"},
	})
}

// ----------------------- helper functions-----------------------------

func getUserFromEmail(email string, ctx context.Context) (models.User, error) {
	filter := bson.M{"email": email}
	result := models.User{}

	err := userCollection.FindOne(ctx, filter).Decode(&result)

	return result, err
}

func setSessionCookie(c *fiber.Ctx, userId primitive.ObjectID) {
	cookie := &fiber.Cookie{
		Name:     "userSession",
		Value:    base64.URLEncoding.EncodeToString([]byte(userId.Hex())),
		Path:     "/",
		MaxAge:   1 * 24 * 60 * 60,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}
	c.Cookie(cookie)
}
