package controllers

import (
	"context"
	"net/http"
	"strings"

	"github.com/Saurav-hehe/Complaint-Management/config"
	"github.com/Saurav-hehe/Complaint-Management/models"
	"github.com/Saurav-hehe/Complaint-Management/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func StaffRegister(c *gin.Context) {
	var staff models.Staff
	if err := c.BindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	staff.Role = strings.ToLower(staff.Role)
	if staff.Role != "electrician" && staff.Role != "plumber" && staff.Role != "carpenter" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid staff role"})
		return
	}

	// Check if email already exists
	var existing models.Staff
	err := config.DB.Collection("staff").FindOne(context.TODO(), bson.M{"email": staff.Email}).Decode(&existing)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Staff already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(staff.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	staff.Password = string(hashedPassword)

	_, err = config.DB.Collection("staff").InsertOne(context.TODO(), staff)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register staff"})
		return
	}

	token := utils.GenerateToken(staff.Email, staff.Role)
	c.JSON(http.StatusCreated, gin.H{"token": token, "role": staff.Role})
}

func StaffLogin(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var staff models.Staff
	err := config.DB.Collection("staff").FindOne(context.TODO(), bson.M{"email": input.Email}).Decode(&staff)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token := utils.GenerateToken(staff.Email, staff.Role)
	c.JSON(http.StatusOK, gin.H{"token": token, "role": staff.Role})
}
func GetAssignedComplaints(c *gin.Context) {
	role := c.GetString("role") // Set by JWT middleware
	filter := bson.M{"assignedRole": role, "staffStatus": bson.M{"$ne": "resolved"}}
	cursor, err := config.DB.Collection("complaints").Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch complaints"})
		return
	}
	var complaints []models.Complaint
	if err := cursor.All(context.TODO(), &complaints); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse complaints"})
		return
	}
	c.JSON(http.StatusOK, complaints)
}
func UpdateComplaintStatus(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		StaffStatus string `json:"staffStatus"`
		StaffRemark string `json:"staffRemark"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid complaint ID"})
		return
	}
	update := bson.M{
		"$set": bson.M{
			"staffStatus": input.StaffStatus,
			"staffRemark": input.StaffRemark,
		},
	}
	_, err = config.DB.Collection("complaints").UpdateByID(context.TODO(), objID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update complaint"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Complaint updated"})
}
