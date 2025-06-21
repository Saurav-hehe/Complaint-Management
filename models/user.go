package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" validate:"required" json:"name"`
	Email    string             `bson:"email" validate:"required,email" json:"email"`
	Password string             `bson:"password" validate:"required,min=6" json:"password"`
	Role     string             `bson:"role" validate:"required,oneof=user warden" json:"role"`
	Hostel   string             `bson:"hostel,omitempty" json:"hostel"`
	Token    string             `bson:"token,omitempty" json:"token"`
}
