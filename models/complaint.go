package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	HostelComplaint    = "HOSTEL"
	FacultyComplaint   = "FACULTY"
	AcademicsComplaint = "ACADEMICS"
)

type Complaint struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Title        string             `bson:"title" validate:"required"`
	Description  string             `bson:"description" validate:"required"`
	Type         string             `bson:"type" validate:"required,oneof=HOSTEL FACULTY ACADEMICS"`
	Hostel       string             `bson:"hostel,omitempty"`
	Status       string             `bson:"status" validate:"required,oneof=PENDING IN_PROGRESS RESOLVED"`
	CreatedBy    primitive.ObjectID `bson:"created_by"`
	CreatedAt    time.Time          `bson:"created_at"`
	AssignedRole string             `bson:"assignedRole,omitempty" json:"assignedRole,omitempty"` // electrician/plumber/carpenter
	StaffStatus  string             `bson:"staffStatus,omitempty" json:"staffStatus,omitempty"`
	StaffRemark  string             `bson:"staffRemark,omitempty" json:"staffRemark,omitempty"`
}
