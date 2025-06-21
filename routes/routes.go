package routes

import (
	"github.com/Saurav-hehe/Complaint-Management/controllers"
	"github.com/Saurav-hehe/Complaint-Management/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(cors.Default())
	auth := r.Group("/auth")
	{

		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/complaints", controllers.CreateComplaint)
		api.GET("/complaints", controllers.GetComplaints)
		api.DELETE("/complaints/:id", controllers.DeleteComplaint)
		api.PUT("/complaints/:id", controllers.ResolveComplaint).Use(middleware.WardenMiddleware())

	}
	r.POST("/staff/register", controllers.StaffRegister)
	r.POST("/staff/login", controllers.StaffLogin)

	// Warden assigns complaint to staff
	r.PUT("/api/complaints/:id/assign", middleware.AuthMiddleware(), controllers.AssignComplaintToStaff)

	// Staff endpoints (protected)
	staff := r.Group("/api/staff")
	staff.Use(middleware.StaffAuthMiddleware())
	{
		staff.GET("/complaints", controllers.GetAssignedComplaints)
		staff.PUT("/complaints/:id", controllers.UpdateComplaintStatus)
	}
}
