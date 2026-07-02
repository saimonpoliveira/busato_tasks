package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/controllers"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/middlewares"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/utils"
)

type Handlers struct {
	Auth       *controllers.AuthController
	User       *controllers.UserController
	Project    *controllers.ProjectController
	Ticket     *controllers.TicketController
	Task       *controllers.TaskController
	Comment    *controllers.CommentController
	Attachment *controllers.AttachmentController
}

func Setup(router *gin.Engine, handlers *Handlers, jwtManager *utils.JWTManager, corsOrigins []string) {
	router.Use(middlewares.CORSMiddleware(corsOrigins))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/login", handlers.Auth.Login)
		auth.POST("/register", handlers.Auth.Register)
	}

	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware(jwtManager))
	{
		protected.GET("/me", handlers.User.Me)

		users := protected.Group("/users")
		{
			users.GET("", handlers.User.List)
			users.GET("/:id", handlers.User.GetByID)
			users.PUT("/:id", handlers.User.Update)
			users.DELETE("/:id", handlers.User.Delete)
			users.POST("", middlewares.AdminMiddleware(), handlers.User.Create)
		}

		projects := protected.Group("/projects")
		{
			projects.GET("", handlers.Project.List)
			projects.GET("/:id", handlers.Project.GetByID)
			projects.POST("", handlers.Project.Create)
			projects.PUT("/:id", handlers.Project.Update)
			projects.DELETE("/:id", handlers.Project.Delete)
		}

		tickets := protected.Group("/tickets")
		{
			tickets.GET("", handlers.Ticket.List)
			tickets.GET("/:id", handlers.Ticket.GetByID)
			tickets.POST("", handlers.Ticket.Create)
			tickets.PUT("/:id", handlers.Ticket.Update)
			tickets.DELETE("/:id", handlers.Ticket.Delete)
		}

		tasks := protected.Group("/tasks")
		{
			tasks.GET("", handlers.Task.List)
			tasks.GET("/:id", handlers.Task.GetByID)
			tasks.POST("", handlers.Task.Create)
			tasks.PUT("/:id", handlers.Task.Update)
			tasks.DELETE("/:id", handlers.Task.Delete)
		}

		comments := protected.Group("/comments")
		{
			comments.GET("", handlers.Comment.List)
			comments.GET("/:id", handlers.Comment.GetByID)
			comments.POST("", handlers.Comment.Create)
			comments.PUT("/:id", handlers.Comment.Update)
			comments.DELETE("/:id", handlers.Comment.Delete)
		}

		attachments := protected.Group("/attachments")
		{
			attachments.GET("", handlers.Attachment.List)
			attachments.GET("/:id", handlers.Attachment.GetByID)
			attachments.GET("/:id/download", handlers.Attachment.Download)
			attachments.POST("", handlers.Attachment.Upload)
			attachments.DELETE("/:id", handlers.Attachment.Delete)
		}
	}
}
