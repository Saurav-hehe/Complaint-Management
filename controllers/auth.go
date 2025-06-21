package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/Saurav-hehe/Complaint-Management/config"
	"github.com/Saurav-hehe/Complaint-Management/models"
	"github.com/Saurav-hehe/Complaint-Management/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		fmt.Println("BindJSON error:", err)
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	fmt.Printf("Register payload: %+v\n", user)

	if user.Role != "user" && user.Role != "warden" {
		c.JSON(400, gin.H{"error": "Invalid role"})
		return
	}
	if user.Role == "warden" && user.Hostel == "" {
		c.JSON(400, gin.H{"error": "Hostel is required for wardens"})
		return
	}

	existingUser := models.User{}
	err := config.DB.Collection("users").FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(400, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	user.Hostel = strings.ToLower(user.Hostel)
	result, err := config.DB.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	token := utils.GenerateToken(user.Email, user.Role)
	c.JSON(201, gin.H{"token": token, "userId": result.InsertedID})
}

func Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	err := config.DB.Collection("users").FindOne(context.TODO(), bson.M{"email": credentials.Email}).Decode(&user)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token := utils.GenerateToken(user.Email, user.Role)
	_, err = config.DB.Collection("users").UpdateOne(context.TODO(),
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"token": token}},
	)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
	}

	c.JSON(200, gin.H{"token": token, "userId": user.ID, "role": user.Role})
}
