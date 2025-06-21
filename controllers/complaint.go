package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Saurav-hehe/Complaint-Management/config"
	"github.com/Saurav-hehe/Complaint-Management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetComplaints(c *gin.Context) {
	email, _ := c.Get("email")
	var user models.User
	err := config.DB.Collection("users").FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	if user.Role == "warden" {

		filter := bson.A{
			bson.M{
				"type": models.HostelComplaint,
				"hostel": bson.M{
					"$regex":   "^" + strings.TrimSpace(strings.ToLower(user.Hostel)) + "$",
					"$options": "i",
				},
			},
			bson.M{"type": models.FacultyComplaint},
			bson.M{"type": models.AcademicsComplaint},
		}
		cursor, err := config.DB.Collection("complaints").Find(context.TODO(), bson.M{"$or": filter})
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch complaints"})
			return
		}
		var complaints []models.Complaint
		if err := cursor.All(context.TODO(), &complaints); err != nil {
			c.JSON(500, gin.H{"error": "Failed to decode complaints"})
			return
		}
		c.JSON(200, complaints)
	} else {

		cursor, err := config.DB.Collection("complaints").Find(context.TODO(), bson.M{"created_by": user.ID})
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch complaints"})
			return
		}
		var complaints []models.Complaint
		if err := cursor.All(context.TODO(), &complaints); err != nil {
			c.JSON(500, gin.H{"error": "Failed to decode complaints"})
			return
		}
		c.JSON(200, complaints)
	}
}

func ResolveComplaint(c *gin.Context) {
	complaintId := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(complaintId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid complaint ID"})
		return
	}

	result, err := config.DB.Collection("complaints").UpdateOne(
		context.TODO(),
		bson.M{"_id": objId},
		bson.M{"$set": bson.M{"status": "RESOLVED"}},
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update complaint"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(404, gin.H{"error": "Complaint not found"})
		return
	}
	c.JSON(200, gin.H{"message": "Complaint resolved"})
}

func CreateComplaint(c *gin.Context) {
	fmt.Println("CreateComplaint called")
	var complaint models.Complaint
	if err := c.BindJSON(&complaint); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	email, _ := c.Get("email")
	var user models.User
	err := config.DB.Collection("users").FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}
	complaint.Status = "PENDING"
	complaint.CreatedBy = user.ID
	complaint.CreatedAt = time.Now()

	result, err := config.DB.Collection("complaints").InsertOne(context.TODO(), complaint)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create complaint"})
		return
	}
	c.JSON(201, gin.H{"id": result.InsertedID})
}

func DeleteComplaint(c *gin.Context) {
	complaintID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(complaintID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid complaint ID"})
		return
	}

	email, _ := c.Get("email")
	var user models.User
	err = config.DB.Collection("users").FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	var complaint models.Complaint
	err = config.DB.Collection("complaints").FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&complaint)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(404, gin.H{"error": "Complaint not found"})
		} else {
			c.JSON(500, gin.H{"error": "Database error"})
		}
		return
	}

	if complaint.CreatedBy != user.ID {
		c.JSON(403, gin.H{"error": "You are not authorized to delete this complaint"})
		return
	}

	_, err = config.DB.Collection("complaints").DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete complaint"})
		return
	}

	c.JSON(200, gin.H{"message": "Complaint deleted successfully"})
}
func AssignComplaintToStaff(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Role string `json:"role"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	input.Role = strings.ToLower(input.Role)
	if input.Role != "electrician" && input.Role != "plumber" && input.Role != "carpenter" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid staff role"})
		return
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid complaint ID"})
		return
	}
	update := bson.M{
		"$set": bson.M{
			"assignedRole": input.Role,
			"staffStatus":  "pending",
		},
	}
	_, err = config.DB.Collection("complaints").UpdateByID(context.TODO(), objID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign complaint"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Complaint assigned to staff"})
}
